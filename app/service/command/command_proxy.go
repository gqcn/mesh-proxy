package command

import (
	"mesh-proxy/app/service/grpcproxy"
	"mesh-proxy/app/service/httpproxy"
	"os"
	"os/signal"
	"syscall"
)

// 启动反向代理网关服务
func Proxy() {
	ch := make(chan struct{})
	// 子进程同时开启HTTP及gRPC网关，
	// 任意服务挂掉则该进程销毁。
	go func() {
		httpproxy.Run()
		close(ch)
	}()
	go func() {
		grpcproxy.Run()
		close(ch)
	}()
	// 信号量监听，优雅关闭
	go func() {
		var procSignalChan = make(chan os.Signal)
		var sig os.Signal
		signal.Notify(
			procSignalChan,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
		)
		for {
			sig = <-procSignalChan
			switch sig {
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM:
				httpproxy.GracefulShutdown()
				grpcproxy.GracefulShutdown()
				close(ch)
			}
		}
	}()
	<-ch
}
