package utils

import (
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
)

var (
	isDebug = genv.GetVar("DEBUG", 0).Bool() // 是否开启调试模式，调试模式下会打印出详细的处理日志信息
)

// 是否打开调试模式
func IsDebug() bool {
	return isDebug
}

// 统一规范化处理服务名称。
// 统一使用K8S的短连接地址。
// 请求HostName中不带端口则不带。
func HandleServiceName(name string) string {
	return gstr.Replace(name, ".svc.cluster.local", "")
}
