/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 9/13/21$ 6:55 AM$
 **/
package peer

import (
	"context"
	pb "github.com/coordinator/protocol"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)

type CommitClientShard struct {
	Connection pb.CommitServiceClient
}

// New creates instance of peer client.
// 'addr' is a coordinator network address (host + port).
func NewClient(addr string) (*CommitClientShard, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	// 增加最小超时时间MinConnectTimeout
	connParams := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay: 100 * time.Millisecond,
			MaxDelay:  5 * time.Second,
		},
		MinConnectTimeout: 200 * time.Millisecond,
	}
	conn, err = grpc.Dial(addr, grpc.WithConnectParams(connParams), grpc.WithInsecure())
	if err != nil {
		time.Sleep(10*time.Millisecond)
		return nil, errors.Wrap(err, "failed to connect")
	}
	return &CommitClientShard{Connection: pb.NewCommitServiceClient(conn)}, nil
}

func (client *CommitClientShard) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeResponse, error) {
	//fmt.Println("[Client] execute client propose: [index=",req.Index,"]")
	return client.Connection.Propose(ctx, req)
}


func (client *CommitClientShard) Commit(ctx context.Context, req *pb.CommitRequest) (*pb.Response, error) {
	//fmt.Println("[Client] execute client commit: [index=",req.Index,"]")
	return client.Connection.Commit(ctx, req)
}

// Get queries value of specific key
func (client *CommitClientShard) Get(ctx context.Context, key string) (*pb.Value, error) {
	return client.Connection.Get(ctx, &pb.Msg{Key: key})
}

/**
	dynamic sharding修改
**/
// StateTransfer将一组state从 source followers转到 target followers
// 2PC StateTransfer client端协议
func (client *CommitClientShard) StateTransfer(ctx context.Context, source uint32, target uint32, msgPayload []byte, index uint64) (*pb.Response, error) {
	//fmt.Println("execute client stateTransfer")
	return client.Connection.StateTransfer(ctx, &pb.TransferRequest{Source: source, Target: target, MsgPayload: msgPayload, Index: index})
}
