package command

import (
	"github.com/gogf/gf/frame/g"
	"net"
	"os"
	"time"
)

// 执行网关健康检查
func Health() {
	grpcAddress := g.Cfg().GetString("proxy.http-address")
	httpAddress := g.Cfg().GetString("proxy.grpc-address")
	// 本地 GRPC 端口检查
	grpcConn, err := net.DialTimeout("tcp", "127.0.0.1"+grpcAddress, 3*time.Second)
	if err == nil {
		grpcConn.Close()
	} else {
		os.Exit(1)
	}
	// 本地 HTTP 端口检查
	httpConn, err := net.DialTimeout("tcp", "127.0.0.1"+httpAddress, 3*time.Second)
	if err == nil {
		httpConn.Close()
	} else {
		os.Exit(1)
	}
	os.Exit(0)
}
