package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

// 自定义消息状态结构体
type Status struct {
	Code    int32      `protobuf:"varint,1,opt,name=code,proto3"   json:"code,omitempty"`
	Message string     `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Details []*any.Any `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
}

func (m *Status) Reset() {
	*m = Status{}
}

func (m *Status) String() string {
	return proto.CompactTextString(m)
}

func (*Status) ProtoMessage() {}

func init() {
	// 注册GRPC消息类型，便于自动解析消息结构
	//proto.RegisterType((*Status)(nil), "MyCustomStatusType.rpc.Status")
}
