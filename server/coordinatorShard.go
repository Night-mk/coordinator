/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 9/13/21$ 6:58 AM$
 **/
package server

import (
	"context"
	pb "github.com/coordinator/protocol"
)

/*
	Propose消息处理
*/
// 这里只要实现类似接口的东西就好，不会在这个程序里调用。。大概
// follower server处理propose逻辑
func (s *ServerShard) ProposeHandler(ctx context.Context, req *pb.ProposeRequest, hook func(req *pb.ProposeRequest) bool) (*pb.ProposeResponse, error) {
	var (
		response *pb.ProposeResponse
	)
	return response, nil
}

// follower server 处理commit逻辑
func (s *ServerShard) CommitHandler(ctx context.Context, req *pb.CommitRequest, hook func(req *pb.CommitRequest) bool) (*pb.Response, error) {
	var (
		response *pb.Response
	)
	return response, nil
}