package httpproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/greuse"
	"io/ioutil"
	"log"
	"mesh-proxy/app/service/tracing"
	"mesh-proxy/library/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"
)

var (
	server      *http.Server
	hostname, _ = os.Hostname()
	address     = g.Cfg().GetString("proxy.http-address")
	errCodeReg  = regexp.MustCompile(`"(err){0,1}code":\s*(\d+)`)
	// 常见静态文件访问不做链路跟踪处理，提高请求转发性能
	staticFileExtSet = gset.NewStrSetFrom([]string{
		// 样式文件
		"js", "json", "css", "map", "less", "sass",
		// 网页文件
		"xml", "htm", "html", "xhtml", "shtml", "tpl",
		// 图片文件
		"png", "gif", "svg", "jpg", "jpeg", "bmp", "ico",
		// 字体文件
		"woff", "woff2", "ttf", "eot",
		// 压缩文件
		"zip", "rar", "7z", "gz",
		// 文档文件
		"doc", "docx", "pdf", "xls", "xlsx", "ppt", "txt", "log", "psd", "md",
	})
)

func init() {
	if address == "" {
		g.Log().Fatal("http proxy address cannot be empty")
	}
}

// 启动HTTP反向代理服务
func Run() {
	conn, err := greuse.Listen("tcp", address)
	if err != nil {
		g.Log().Fatal(err)
	}
	server = &http.Server{
		Handler:      http.HandlerFunc(httpHandler),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(&errorLogger{logger: g.Log()}, "", 0),
	}
	// 启动HTTP Server服务
	g.Log().Printf("%d: http proxy start running on %s", gproc.Pid(), address)
	if err := server.Serve(&Listener{conn}); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		g.Log().Error(err)
	}
}

// 优雅关闭HTTP服务
func GracefulShutdown() {
	if server == nil {
		return
	}
	g.Log().Printf("%d: http proxy gracefully shutdown", gproc.Pid())
	server.Shutdown(context.Background())
}

// 默认的HTTP反向代理处理方法
func httpHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		spanId      string
		traceId     string
		traceMapStr string
		upstream    string
		extras      gset.StrSet
		startTime   = gtime.TimestampMilli()
	)
	// 创建自定义的Writer，支持缓存控制
	writer := NewResponseWriter(w)
	// 判断静态文件请求
	isStaticRequest := false
	if ext := gfile.ExtName(r.URL.Path); ext != "" {
		if staticFileExtSet.Contains(ext) {
			isStaticRequest = true
		}
	}
	// 非静态文件请求才执行链路跟踪逻辑，以及日志写入逻辑
	if !isStaticRequest {
		requestBody, _ := ioutil.ReadAll(r.Body)
		r.Body = NewReadCloser(requestBody, false)
		requestBody = bytes.TrimSpace(requestBody)
		// 链路跟踪信息，如果请求Header中没有则创建并写入到请求Header中
		// TraceId
		if traceId = r.Header.Get(tracing.TraceIdName); traceId == "" {
			traceId = tracing.NewTraceId()
			if !tracing.IsRoot() {
				extras.Add("NewCreatedTracing")
			}
		}
		// SpanId
		if spanId = r.Header.Get(tracing.SpanIdName); spanId == "" {
			spanId = tracing.DefaultSpanId
			if !tracing.IsRoot() {
				extras.Add("NewCreatedTracing")
			}
		} else {
			spanId = tracing.IncreaseSpanId(traceId, spanId)
		}
		// TraceMap
		traceMapStr = r.Header.Get(tracing.TraceMapName)

		// 公网入口网关需要解析提交的SESSION并获取UserId传递到后续的链路中
		if tracing.IsGate() {
			// 网关不支持客户端提交TraceMap，需要清空，以保证安全性
			traceMapStr = ""
			if userId := parseUserId(r); userId != "" {
				data := make(map[string]interface{})
				if traceMapStr != "" {
					json.Unmarshal([]byte(traceMapStr), &data)
				}
				data["UserId"] = userId
				if b, err := json.Marshal(data); err == nil {
					traceMapStr = gconv.UnsafeBytesToStr(b)
					r.Header.Set(tracing.TraceMapName, traceMapStr)
				}
			}
		}

		r.Header.Set(tracing.SpanIdName, spanId)
		r.Header.Set(tracing.TraceIdName, traceId)

		// 返回的Header中增加链路跟踪信息，如果返回Header中已经存在则会覆盖
		writer.Header().Set(tracing.TraceIdName, traceId)

		// 反向代理日志记录
		defer func() {
			var (
				errMsg       = ""
				responseBody = writer.Buffer()
			)
			if gstr.Equal(writer.Header().Get("Content-Encoding"), "gzip") {
				responseBody, _ = gcompress.UnGzip(responseBody)
				if len(responseBody) == 0 {
					responseBody = writer.Buffer()
				}
			}
			if err != nil {
				errMsg = err.Error()
			}
			// 日志格式：
			// [类型] 状态码 HTTPMethod 客户端名称 服务名称 服务地址 URL 耗费时间(秒) 层级ID 链路ID "自定义信息" "错误信息" "客户端信息" "附加信息"#"请求Body"#"返回Body"
			costInSeconds := float64(gtime.TimestampMilli()-startTime) / 1000
			g.Log().Cat("proxy").Printf(
				`[HTTP] %d %s %s %s %s %s %.3f %s %s "%s" "%s" "%s" "%s"#"%s"#"%s"`,
				writer.Status(), r.Method, hostname,
				utils.HandleServiceName(r.Host), upstream, r.URL.String(),
				costInSeconds, spanId, traceId, traceMapStr, errMsg, r.UserAgent(), extras.Join(";"),
				gstr.Replace(gconv.UnsafeBytesToStr(requestBody), "\n", ""),
				gstr.Replace(gconv.UnsafeBytesToStr(responseBody), "\n", ""),
			)
			// 将缓存的返回内容输出到客户端
			writer.OutputBuffer()
		}()
	}
	// 检测反向代理配置，如果不存在则返回404
	upstream = GetOriginByRemoteAddr(r.RemoteAddr)
	if upstream == "" {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
	// 反向代理请求处理，后端HTTP目标服务统一使用HTTP
	var u *url.URL
	u, err = url.Parse("http://" + upstream)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		writer.WriteHeader(http.StatusBadGateway)
		err = e
	}
	if isStaticRequest {
		// 静态文件服务使用底层Writer支持Stream流式下载
		proxy.ServeHTTP(writer.RawWriter(), r)
	} else {
		// 非静态文件请求使用缓存Writer
		proxy.ServeHTTP(writer, r)
	}
}
