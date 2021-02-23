package grpcproxy

import (
	"context"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"mesh-proxy/app/service/tracing"
	"mesh-proxy/library/utils"
)

// GRPC中间件 - Unary
func gRPCUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var (
		array          []string
		host           string
		agent          string
		method         string
		extras         gset.StrSet
		spanId         string
		traceId        string
		traceMapStr    string
		startTimestamp = gtime.TimestampMilli()
	)

	// 执行请求
	resp, err = handler(ctx, req)

	incomingHeaders, _ := metadata.FromIncomingContext(ctx)
	// HTTP/2 Method，一般为POST
	if array = incomingHeaders.Get(":method"); len(array) > 0 {
		method = array[0]
	} else {
		method = "POST"
	}
	if array = incomingHeaders.Get(":authority"); len(array) > 0 {
		host = utils.HandleServiceName(array[0])
	}
	if array = incomingHeaders.Get("user-agent"); len(array) > 0 {
		agent = array[0]
	}
	if array = incomingHeaders.Get(tracing.NewTraceName); len(array) > 0 {
		extras.Add("NewCreatedTracing")
	}
	if array = incomingHeaders.Get(tracing.SpanIdName); len(array) > 0 {
		spanId = array[0]
	}
	if array = incomingHeaders.Get(tracing.TraceIdName); len(array) > 0 {
		traceId = array[0]
	}
	if array = incomingHeaders.Get(tracing.TraceMapName); len(array) > 0 {
		traceMapStr = array[0]
	}

	errCode, errMsg := extractError(err)
	// 日志格式(注意GRPC下Body日志记录为空)：
	// [类型] 状态码 HTTPMethod 客户端名称 服务名称 服务地址 调用方法 耗费时间(秒) 层级ID 链路ID "自定义信息" "错误信息" "客户端信息" "附加信息"#"请求Body"#"返回Body"
	costInSeconds := float64(gtime.TimestampMilli()-startTimestamp) / 1000
	g.Log().Async().Cat("proxy").Printf(
		`[GRPC] %d %s %s %s %s %s %.3f %s %s "%s" "%s" "%s" "%s"#"%s"#"%s"`,
		errCode, method, hostname, host, host, info.FullMethod,
		costInSeconds, spanId, traceId, traceMapStr, errMsg, agent, extras.Join(";"),
		marshalPbToJson(req),
		marshalPbToJson(resp),
	)
	return
}
