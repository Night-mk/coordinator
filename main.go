/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 9/13/21$ 6:07 AM$
 **/
package main

import (
	"github.com/coordinator/config"
	"github.com/coordinator/hooks"
	"github.com/coordinator/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	conf := config.GetConfig()
	hooks, err := hooks.GetHookF()
	if err != nil {
		panic(err)
	}
	s, err := server.NewCommitServerShard(conf, hooks...)
	if err != nil {
		panic(err)
	}

	// Run函数启动一个non-blocking GRPC server
	s.Run(server.WhiteListCheckerShard)
	<-ch // 如果一直没有系统信号，就一直等待？
	s.Stop()
}