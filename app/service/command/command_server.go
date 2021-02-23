package command

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 启动控制面服务
func Server() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {

	})
	s.Run()
}
