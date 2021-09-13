/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 9/13/21$ 6:54 AM$
 **/
package src

import (
	pb "github.com/coordinator/protocol"
	"github.com/ethereum/go-ethereum/log"
)

func ProposeShard(req *pb.ProposeRequest) bool {
	log.Info("propose hook on height %d is OK", req.Index)
	return true
}

func CommitShard(req *pb.CommitRequest) bool {
	log.Info("commit hook on height %d is OK", req.Index)
	return true
}