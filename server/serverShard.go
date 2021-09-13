/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 9/13/21$ 6:57 AM$
 **/
package server

import (
	"context"
	"fmt"
	"github.com/coordinator/cache"
	"github.com/coordinator/config"
	"github.com/coordinator/peer"
	pb "github.com/coordinator/protocol"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
)

type OptionShard func(server *ServerShard) error
const COORDINATOR_SHARDID = uint32(10000)
const (
	TWO_PHASE   = "two-phase"
	THREE_PHASE = "three-phase"
)

// 账户状态
type state struct {
	Address common.Address
	Balance *big.Int
	Nonce uint64 // 账户的nonce值也要设置
	Data []byte // 表示合约
}
// 状态证明
type stateProof struct {
	BlockHeight uint64
	RootHash *common.Hash
	Proof interface{}
}
// 状态数据
type StateData struct {
	State state
	StateProof stateProof
}

type ServerShard struct {
	pb.UnimplementedCommitServiceServer
	Addr                 string
	ShardID              uint32 // 初始化时需要整一个和node一样的shardID
	FollowerType         pb.FollowerType
	Followers            []*peer.CommitClientShard
	ShardFollowers       map[uint32][]*peer.CommitClientShard // 带shard的全体follower集合
	Config               *config.ConfigShard // 这个config改为node里的配置
	GRPCServer           *grpc.Server
	ProposeShardHook     func(req *pb.ProposeRequest) bool
	CommitShardHook      func(req *pb.CommitRequest) bool
	//CommitPool            *CommitPool
	CoordinatorCache     *cache.Cache
	cancelCommitOnHeight map[uint64]bool
	mu                   sync.RWMutex
}

// 2PC 第一阶段
func (s *ServerShard) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeResponse, error) {
	//fmt.Println("[server] ",s.Addr, " execute server propose")
	// 设置在某个Height是否需要cancel
	s.SetProgressForCommitPhase(req.Index, false)

	return s.ProposeHandler(ctx, req, s.ProposeShardHook)
}

// 2PC 第二阶段
func (s *ServerShard) Commit(ctx context.Context, req *pb.CommitRequest) (resp *pb.Response, err error) {
	//fmt.Println("[Server] execute server commit: [index=",req.Index,"]")
	// 两阶段commit
	resp, err = s.CommitHandler(ctx, req, s.CommitShardHook)

	if err != nil {
		return
	}
	//if resp.Acktype == c_pb.AckType_ACK {
	//	atomic.AddUint64(&s.Height, 1)
	//}
	return
}

/**
stateTransfer 并行化方案
*/
func (s *ServerShard) StateTransfer(ctx context.Context, req *pb.TransferRequest) (*pb.Response, error){
	var (
		response *pb.Response
		err      error
	)
	// 选择使用的协议，目前暂定只用2PC(propose, commit)
	var ctype pb.CommitType
	if s.Config.CommitType == THREE_PHASE {
		ctype = pb.CommitType_THREE_PHASE_COMMIT
	} else {
		ctype = pb.CommitType_TWO_PHASE_COMMIT
	}

	/**=======2PC协议流程=======*/
	// propose阶段
	// 打印address list
	// 反序列化req.Keylist为string slice, log输出, 调用ethereum的rlp编码进行解码
	var deState StateData
	err = rlp.DecodeBytes(req.MsgPayload, &deState)
	if err != nil{
		log.Info("decode error: ", err)
	}
	//fmt.Printf("[StateTransfer] Propose phase: index=[%d] source=[%d], target=[%d], state=[%+v] \n", req.Index, req.Source, req.Target, deState.Address)

	// 根据source和target选择followers组
	currentFollowers := s.ObtainFollowers(req.Source, req.Target)

	// 对每个follower并发处理
	var wg sync.WaitGroup
	proposeResChan := make(chan *pb.ProposeResponse, 1000)
	proposeErrChan := make(chan error, 1000)
	//proposeDone := make(chan bool)
	for _, follower := range currentFollowers {
		wg.Add(1)
		go func(follower *peer.CommitClientShard) {
			var (
				response *pb.ProposeResponse
				err      error
			)
			for response == nil || response != nil && response.Acktype == pb.AckType_NACK {
				// 每个follower（source, target）用grpc调用coordinator给他们发送消息
				response, err = follower.Propose(ctx, &pb.ProposeRequest{Source: req.Source, Target: req.Target, MsgPayload: req.MsgPayload, CommitType: ctype, Index: req.Index})
			}
			proposeResChan <- response
			proposeErrChan <- err

			if err != nil {
				return
			}
			// coordinator server设置cache
			// index设置成交易的绝对编号
			s.CoordinatorCache.Set(req.Index, strconv.FormatUint(req.Index, 10), []byte{})

			if response.Acktype != pb.AckType_ACK {
				log.Info("follower not acknowledged msg")
				return
			}
			wg.Done()
		}(follower)
	}

	// 等待所有follower反馈完成，再进入commit过程
	wg.Wait()
	close(proposeErrChan)
	close(proposeResChan)

	for err := range proposeErrChan{
		//fmt.Println("=============execute proposeErrChan=====================")
		if err != nil {
			log.Error("=============execute proposeErrChan=====================", err)
			return &pb.Response{Acktype: pb.AckType_NACK}, nil
		}
	}
	for res := range proposeResChan{
		//fmt.Println("=============execute proposeResChan=====================")
		if res.Acktype != pb.AckType_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}

	//fmt.Printf(s.Config.Nodeaddr,"PROPOSE END index=[", req.Index,"]")

	// the coordinator got all the answers, so it's time to persist msg and send commit command to followers
	_, _, ok := s.CoordinatorCache.Get(req.Index)
	//source, target, keylist, ok := s.NodeCache.Get(s.Height)
	if !ok {
		log.Info("[StateTransfer] can not find msg in cache for index=[", req.Index,"]")
		return nil, status.Error(codes.Internal, "can't to find msg in the coordinator's cache")
	}
	// coordinator server将事务写入
	//if err = s.DB.Put(deState.Address, statelist); err != nil {
	//	return &pb.Response{Type: pb.AckType_NACK}, status.Error(codes.Internal, "failed to save msg on coordinator")
	//}

	// commit阶段
	//fmt.Printf("[StateTransfer] Commit phase: index=[%d] source=[%d], target=[%d], state=[%s] \n", req.Index, req.Source, req.Target, deState)
	var wg1 sync.WaitGroup
	commitResChan := make(chan *pb.Response, 1000)
	commitErrChan := make(chan error, 1000)
	for _, follower := range currentFollowers{
		wg1.Add(1)
		go func(follower *peer.CommitClientShard) {
			// coordinator获取到commit的response
			response, err = follower.Commit(ctx, &pb.CommitRequest{Index: req.Index, Address: deState.State.Address.Bytes()})
			commitResChan <- response
			commitErrChan <- err
			if err != nil {
				log.Error("follower commit failure",err.Error())
			}
			if response.Acktype != pb.AckType_ACK {
				return
			}
			wg1.Done()
		}(follower)
	}

	wg1.Wait()
	close(commitResChan)
	close(commitErrChan)

	// 对所有follower节点的返回值进行处理
	for res := range commitResChan{
		if res.Acktype != pb.AckType_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}
	for err := range commitErrChan{
		if err != nil {
			log.Error("CommitResponse ERROR", err)
			return &pb.Response{Acktype: pb.AckType_NACK}, nil
		}
	}

	fmt.Printf("[Coordinator]: committed state for tx_index=[%d] \n", req.Index)
	//log.Infof("committed state for tx_index=[%d] \n", req.Index)
	// 返回调用端最终结果
	return &pb.Response{Acktype: pb.AckType_ACK}, nil
}

func (s *ServerShard)ObtainFollowers(source uint32, target uint32) (follower []*peer.CommitClientShard){
	for _, node := range s.ShardFollowers[source]{
		follower = append(follower, node)
	}
	for _, node := range s.ShardFollowers[target]{
		follower = append(follower, node)
	}
	return follower
}

// 根据config，创建server端instance
func NewCommitServerShard(conf *config.ConfigShard, opts ...OptionShard) (*ServerShard, error){
	var wg sync.WaitGroup
	wg.Add(1)
	//fmt.Println("============2PC CONFIG============", conf)

	// 初始化serverShard
	//server := &ServerShard{Addr: conf.Nodeaddr, DBPath: conf.DBPath}
	server := &ServerShard{}
	// 对coordinator和follower分开初始化
	if conf.Role == "coordinator"{
		server.Addr = conf.Nodeaddr
		server.ShardID = COORDINATOR_SHARDID
	}else{
		// follower分解nodeaddr
		addrWithShard := strings.SplitN(conf.Nodeaddr, "+", 2)
		server.Addr = addrWithShard[1]
		intNum, _ := strconv.Atoi(addrWithShard[0])
		server.ShardID = uint32(intNum)
	}

	var err error
	// 查看option是否有问题
	for _, option := range opts {
		err = option(server)
		if err != nil {
			return nil, err
		}
	}

	nodeaddr := conf.Nodeaddr
	if conf.Role != "coordinator"{
		nodeaddr = strings.SplitN(conf.Nodeaddr, "+", 2)[1]
	}
	//utils.Logger().Info().Str("nodeAddress: ", nodeaddr)
	log.Info("[2PC] NodeAddress: ", nodeaddr)

	// 初始化配置里的follower client（coordinator中的follower）
	server.ShardFollowers = make(map[uint32][]*peer.CommitClientShard)
	for key, nodes := range conf.Followers {
		for _, node := range nodes{
			//fmt.Println("config init node: ", node)
			cli, err := peer.NewClient(node)
			if err != nil {
				return nil, err
			}
			server.Followers = append(server.Followers, cli)
			server.ShardFollowers[key] = append(server.ShardFollowers[key], cli)
		}
	}

	// 在server.config中标记coordinator
	if conf.Role == "coordinator" {
		server.Config.Coordinator = server.Addr
	}

	// 初始化DB, cache, rollback参数
	//server.DB, err = db.New(conf.DBPath)
	server.CoordinatorCache = cache.NewCoorCache()
	server.cancelCommitOnHeight = map[uint64]bool{}

	if server.Config.CommitType == TWO_PHASE {
		log.Info("two-phase-commit mode enabled")
	} else {
		log.Info("three-phase-commit mode enabled")
	}

	wg.Done()
	wg.Wait()

	return server, err
}

// Run starts non-blocking GRPC server
func (s *ServerShard) Run(opts ...grpc.UnaryServerInterceptor) {
	var err error

	s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))
	pb.RegisterCommitServiceServer(s.GRPCServer, s)

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Error("failed to listen: %v", err)
	}
	log.Info("listening on tcp://%s", s.Addr)

	go s.GRPCServer.Serve(l)
}

// Stop stops server
func (s *ServerShard) Stop() {
	log.Info("stopping server")
	s.GRPCServer.GracefulStop()
	log.Info("server stopped")
}

//func (s *ServerShard) rollback() {
//	s.NodeCache.Delete(s.Height)
//}

func (s *ServerShard) SetProgressForCommitPhase(height uint64, docancel bool) {
	s.mu.Lock()
	s.cancelCommitOnHeight[height] = docancel
	s.mu.Unlock()
}

func (s *ServerShard) HasProgressForCommitPhase(height uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cancelCommitOnHeight[height]
}
