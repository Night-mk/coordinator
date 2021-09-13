/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 9/13/21$ 8:08 AM$
 **/
package hooks


import (
	pb "github.com/coordinator/protocol"
	"github.com/coordinator/server"
	"github.com/coordinator/hooks/src"
)

type ProposeShardHook func(req *pb.ProposeRequest) bool
type CommitShardHook func(req *pb.CommitRequest) bool

func GetHookF() ([]server.OptionShard, error) {
	proposeShardHook := func(f ProposeShardHook) func(*server.ServerShard) error {
		return func(server *server.ServerShard) error {
			server.ProposeShardHook = f
			return nil
		}
	}

	commitShardHook := func(f CommitShardHook) func(*server.ServerShard) error {
		return func(server *server.ServerShard) error {
			server.CommitShardHook = f
			return nil
		}
	}

	return []server.OptionShard{proposeShardHook(src.ProposeShard), commitShardHook(src.CommitShard)}, nil
}