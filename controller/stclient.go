/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 9/13/21$ 6:53 AM$
 **/
package main

import (
"context"
"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
"github.com/coordinator/peer"
pb "github.com/coordinator/protocol"
"sync"
"time"
)

const (
	TXNUM = 10000 // 消息总量
	ParallelChNum = 1500 // 并发协程数量
	CoordinatorAddr = "localhost:3000"
)

var (
	whitelist = []string{"127.0.0.1"}
	testTable = obtainTestData(shardList, stateList)
)

// 测试数据准备
type StateTransferMsg struct{
	source uint32
	target uint32
	keylist []byte
}
var shardList = []uint32{
	uint32(0), uint32(1), uint32(2),
}
// 需要rlp.encode和decode的结构体，所有元素首字母都要大写！！
type StateDataTest struct {
	Address string
	State string
}
// 这里直接使用harmony的地址, 传输的数据是包含state data的map数据
var stateList = []StateDataTest{
	{"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed", "balance=50"},
	{"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359", "balance=29.59"},
	{"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB", "balance=100"},
	{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", "balance=3000"},
}

// 设置测试数据
func obtainTestData(shard []uint32, stateList []StateDataTest) []StateTransferMsg {
	var enStateList = make([][]byte, len(stateList))
	fmt.Println("state ori: ", stateList)
	// 获取address的rlp编码
	for i, state := range stateList{
		enState, _ := rlp.EncodeToBytes(state)
		//fmt.Printf("%v → %X\n", state, enState)
		var deState StateDataTest
		err := rlp.DecodeBytes(enState, &deState)
		fmt.Printf("State:err=%+v,val=%+v\n", err, deState)
		enStateList[i] = enState
	}
	var testTable = []StateTransferMsg{
		{shard[0], shard[1], enStateList[0]},
		//{shard[0], shard[1], enStateList[1]},
		//{shard[0], shard[1], enStateList[2]},
		//{shard[1], shard[2], enStateList[3]},
	}

	// 构建更多交易
	TxNum := TXNUM
	for i:=0; i<TxNum; i++{
		testTable = append(testTable, testTable[0])
	}

	return testTable
}

// 客户端，准备大量交易，并行发给coordinator调用状态迁移协议
func ParallelTest(){
	client, err := peer.NewClient(CoordinatorAddr)
	if err != nil {
		panic(err)
	}
	time.Sleep(2* time.Second)
	fmt.Println("========= Create CLIENT END =========")

	// 控制coordinator并发的协程数量
	coordPCh := make(chan struct{}, ParallelChNum)

	startTime := time.Now() // 获取当前时间
	// 调用stateTransfer请求
	// 使用goroutine并行调用
	var wg sync.WaitGroup
	for index, state := range testTable {
		wg.Add(1)
		coordPCh <- struct{}{}
		go func(s StateTransferMsg, i int) {
			defer wg.Done()
			resp, err := client.StateTransfer(context.Background(), s.source, s.target, s.keylist, uint64(i))

			if err != nil {
				log.Error(err.Error())
			}
			if resp.Acktype != pb.AckType_ACK {
				fmt.Println("err response type: ", resp.Acktype)
				log.Error(err.Error())
			}
			<- coordPCh
		}(state, index+1)
	}

	wg.Wait()
	fmt.Println("=========== STATETRANSFER END ===========")
	elapsed := time.Since(startTime)

	fmt.Println("===== TX NUM:",len(testTable),"=====")
	// 计算纳秒
	duration := elapsed.Seconds()
	tps := TXNUM /duration
	fmt.Println("2PC time consumed: ", duration)
	fmt.Println("2PC TPS: ", tps)
}

func main(){
	ParallelTest()
}
