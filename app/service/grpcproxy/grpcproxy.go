package grpcproxy

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/greuse"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"mesh-proxy/app/service/tracing"
	"mesh-proxy/library/utils"
	"os"
	"time"

	_ "mesh-proxy/app/service/grpcproxy/types"
)

const (
	// GRPC客户端链接创建超时时间
	defaultDialTimeout = 10 * time.Second
)

var (
	server      *grpc.Server
	hostname, _ = os.Hostname()
	address     = g.Cfg().GetString("proxy.grpc-address")
	clientMap   = gmap.NewStrAnyMap(true)
)

func init() {
	if address == "" {
		g.Log().Fatal("grpc proxy address cannot be empty")
	}
}

func Run() {
	// 注意，由于UnknownServiceHandler的存在，所有的GRPC请求都走的是Stream模式。
	server = grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(gRPCDirector)),
		grpc.UnaryInterceptor(gRPCUnaryInterceptor),
		grpc.StreamInterceptor(gRPCStreamInterceptor),
	)
	conn, err := greuse.Listen("tcp", address)
	if err != nil {
		g.Log().Fatal(err)
	}
	// 启动gRPC Server
	g.Log().Printf("%d: grpc proxy start running on %s", gproc.Pid(), address)
	if err := server.Serve(conn); err != nil {
		g.Log().Error(err)
	}
}

// 优雅关闭gRPC服务
func GracefulShutdown() {
	if server == nil {
		return
	}
	g.Log().Printf("%d: grpc proxy gracefully shutdown", gproc.Pid())
	server.GracefulStop()
}

// GRPCDirector is the gRPC routing handler.
func gRPCDirector(ctx context.Context, _ string) (context.Context, *grpc.ClientConn, error) {
	var host, spanId, traceId, traceMap string
	incomingHeaders, _ := metadata.FromIncomingContext(ctx)

	// 链路跟踪信息: TraceId
	if array := incomingHeaders.Get(tracing.TraceIdName); len(array) == 0 {
		traceId = tracing.NewTraceId()
		if !tracing.IsRoot() {
			incomingHeaders.Set(tracing.NewTraceName, "1")
		}
	} else {
		traceId = array[0]
	}

	// 链路跟踪信息: SpanId
	if array := incomingHeaders.Get(tracing.SpanIdName); len(array) == 0 {
		spanId = tracing.DefaultSpanId
		if !tracing.IsRoot() {
			incomingHeaders.Set(tracing.NewTraceName, "1")
		}
	} else {
		spanId = tracing.IncreaseSpanId(traceId, array[0])
	}

	// 链路跟踪信息: TraceMap
	if array := incomingHeaders.Get(tracing.TraceMapName); len(array) != 0 {
		traceMap = array[0]
	}

	incomingHeaders.Set(tracing.SpanIdName, spanId)
	incomingHeaders.Set(tracing.TraceIdName, traceId)
	if traceMap != "" {
		incomingHeaders.Set(tracing.TraceMapName, traceMap)
	}

	// 查找反向代理目标服务器
	if array := incomingHeaders.Get(":authority"); len(array) > 0 {
		host = utils.HandleServiceName(array[0])
	}

	// 需要注意，作为反向代理服务，这里接收到的incoming数据即需要发送给目标服务的outgoing数据
	ctx = metadata.NewOutgoingContext(ctx, incomingHeaders)

	// 获取/创建反向代理客户端，使用了Map保证客户端对象的复用。
	// 这里使用了Map而不是Pool，是因为gRPC底层使用了HTTP2实现连接的IO复用，并且实现了连接的失败重试。
	// 键值是一个*grpc.ClientConn类型，这是一个虚拟的连接对象。
	var err error
	result := clientMap.GetOrSetFuncLock(host, func() interface{} {
		var clientConn *grpc.ClientConn
		// 注意这里的ctx覆盖了外面的ctx变量，仅在当前作用域有效，用于控制连接创建的超时
		ctx, _ := context.WithTimeout(context.Background(), defaultDialTimeout)
		clientConn, err = grpc.DialContext(
			ctx,
			host,
			grpc.WithInsecure(),
			grpc.WithDefaultCallOptions(
				grpc.CustomCodecCallOption{
					Codec: proxy.Codec(),
				},
			),
		)
		return clientConn
	})
	if err != nil {
		return ctx, nil, err
	}
	return ctx, result.(*grpc.ClientConn), err
}
