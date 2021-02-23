package grpcproxy

import (
	"fmt"
	"github.com/gogf/gf/text/gstr"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CustomUnknownErrorCode = 800
)

var (
	pbJsonMarshaller = &jsonpb.Marshaler{
		EmitDefaults: true,
	}
)

// 将二进制PB内容转换为JSON字符串，便于日志记录
func marshalPbToJson(v interface{}) (msg string) {
	var err error
	pb, ok := v.(proto.Message)
	if ok {
		msg, err = pbJsonMarshaller.MarshalToString(pb)
	}
	if err != nil || !ok {
		msg = fmt.Sprintf("%s", v)
	}
	return msg
}

// 解析GRPC error.
func extractError(err error) (code codes.Code, msg string) {
	if err == nil {
		return codes.OK, ""
	}
	// 自定义错误状态码
	code = CustomUnknownErrorCode
	if rpcErr, ok := status.FromError(err); ok {
		code = rpcErr.Code()
		msg = fmt.Sprintf(`%+v`, rpcErr.Details())
		if len(msg) < 3 {
			msg = rpcErr.Message()
		}
	} else {
		msg = fmt.Sprintf(`%+v`, err)
	}
	if msg != "" {
		msg = gstr.OctStr(msg)
	}
	return
}
