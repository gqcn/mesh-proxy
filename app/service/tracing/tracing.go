package tracing

import (
	"fmt"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
	"time"
)

const (
	SpanIdName    = "Span-Id"       // 层级Id，注意大小写敏感
	TraceIdName   = "Trace-Id"      // 链路Id，注意大小写敏感
	TraceMapName  = "Trace-Map"     // 应用层自定义Key-Value数据，注意大小写敏感
	NewTraceName  = "New-Trace"     // 客户端请求时没有传递TraceId的标记(一般为根节点)
	DefaultSpanId = "0"             // 默认的SpanId值，表示链路的根节点
	cacheTimeout  = 2 * time.Minute // SpanId在SideCar中的缓存时间
)

var (
	isRoot = genv.GetVar("ROOT", 0).Bool() // 是否当前节点为链路跟踪的根节点(入口节点，如定时任务节点)
	isGate = genv.GetVar("GATE", 0).Bool() // 是否当前节点为网关，注意整个集群的网关只有1个，并且必须为HTTP服务端
)

// 判断当前节点是否为链路跟踪的根节点
func IsRoot() bool {
	return isRoot
}

// 判断当前节点是否为网关节点
func IsGate() bool {
	return isGate
}

// 生成一个TraceId，32个字符
func NewTraceId() string {
	return gstr.ToUpper(guid.S())
}

// 根据TraceId递增链路层级。
// 如果是第一次创建链路信息那么不会记录到缓存中，即链路的根节点不会缓存。
func IncreaseSpanId(traceId, spanId string) string {
	cacheKey := "mesh-proxy-span-id-cache:" + traceId
	if v := gcache.GetVar(cacheKey); !v.IsNil() {
		spanId = v.String()
		if spanId == "" {
			spanId = DefaultSpanId
		}
		if spanId == DefaultSpanId {
			spanId = fmt.Sprintf(`%s.1`, spanId)
		} else {
			array := gstr.Split(spanId, ".")
			array[len(array)-1] = gconv.String(gconv.Int(array[len(array)-1]) + 1)
			spanId = gstr.Join(array, ".")
		}
	} else {
		if spanId == "" {
			spanId = DefaultSpanId
		}
		spanId = fmt.Sprintf(`%s.1`, spanId)
	}
	gcache.Set(cacheKey, spanId, cacheTimeout)
	return spanId
}
