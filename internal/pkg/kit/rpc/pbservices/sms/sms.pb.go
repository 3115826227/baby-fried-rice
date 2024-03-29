// proto 文件版本

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.14.0
// source: sms.proto

// 生成文件的包名

package sms

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 发送短信
type ReqSendMessageDao struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 用户id
	AccountId string `protobuf:"bytes,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	// 用户手机号码
	Phone string `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	// 发送手机验证码
	PhoneCode string `protobuf:"bytes,3,opt,name=phoneCode,proto3" json:"phoneCode,omitempty"`
	// 签名名称
	SignName string `protobuf:"bytes,4,opt,name=signName,proto3" json:"signName,omitempty"`
	// 发送模板代码
	Code int64 `protobuf:"varint,5,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *ReqSendMessageDao) Reset() {
	*x = ReqSendMessageDao{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sms_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqSendMessageDao) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqSendMessageDao) ProtoMessage() {}

func (x *ReqSendMessageDao) ProtoReflect() protoreflect.Message {
	mi := &file_sms_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqSendMessageDao.ProtoReflect.Descriptor instead.
func (*ReqSendMessageDao) Descriptor() ([]byte, []int) {
	return file_sms_proto_rawDescGZIP(), []int{0}
}

func (x *ReqSendMessageDao) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *ReqSendMessageDao) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *ReqSendMessageDao) GetPhoneCode() string {
	if x != nil {
		return x.PhoneCode
	}
	return ""
}

func (x *ReqSendMessageDao) GetSignName() string {
	if x != nil {
		return x.SignName
	}
	return ""
}

func (x *ReqSendMessageDao) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

var File_sms_proto protoreflect.FileDescriptor

var file_sms_proto_rawDesc = []byte{
	0x0a, 0x09, 0x73, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x73, 0x6d, 0x73,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01,
	0x0a, 0x11, 0x52, 0x65, 0x71, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x44, 0x61, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x68, 0x6f, 0x6e,
	0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69, 0x67, 0x6e, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x69, 0x67, 0x6e, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x32, 0x4a, 0x0a, 0x06, 0x44, 0x61, 0x6f, 0x53, 0x6d, 0x73, 0x12,
	0x40, 0x0a, 0x0e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x44, 0x61,
	0x6f, 0x12, 0x16, 0x2e, 0x73, 0x6d, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x53, 0x65, 0x6e, 0x64, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x44, 0x61, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x42, 0x35, 0x5a, 0x33, 0x62, 0x61, 0x62, 0x79, 0x2d, 0x66, 0x72, 0x69, 0x65, 0x64, 0x2d,
	0x72, 0x69, 0x63, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x6b, 0x69, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2f, 0x73, 0x6d, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sms_proto_rawDescOnce sync.Once
	file_sms_proto_rawDescData = file_sms_proto_rawDesc
)

func file_sms_proto_rawDescGZIP() []byte {
	file_sms_proto_rawDescOnce.Do(func() {
		file_sms_proto_rawDescData = protoimpl.X.CompressGZIP(file_sms_proto_rawDescData)
	})
	return file_sms_proto_rawDescData
}

var file_sms_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_sms_proto_goTypes = []interface{}{
	(*ReqSendMessageDao)(nil), // 0: sms.ReqSendMessageDao
	(*emptypb.Empty)(nil),     // 1: google.protobuf.Empty
}
var file_sms_proto_depIdxs = []int32{
	0, // 0: sms.DaoSms.SendMessageDao:input_type -> sms.ReqSendMessageDao
	1, // 1: sms.DaoSms.SendMessageDao:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sms_proto_init() }
func file_sms_proto_init() {
	if File_sms_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sms_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqSendMessageDao); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sms_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sms_proto_goTypes,
		DependencyIndexes: file_sms_proto_depIdxs,
		MessageInfos:      file_sms_proto_msgTypes,
	}.Build()
	File_sms_proto = out.File
	file_sms_proto_rawDesc = nil
	file_sms_proto_goTypes = nil
	file_sms_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DaoSmsClient is the client API for DaoSms service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DaoSmsClient interface {
	SendMessageDao(ctx context.Context, in *ReqSendMessageDao, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type daoSmsClient struct {
	cc grpc.ClientConnInterface
}

func NewDaoSmsClient(cc grpc.ClientConnInterface) DaoSmsClient {
	return &daoSmsClient{cc}
}

func (c *daoSmsClient) SendMessageDao(ctx context.Context, in *ReqSendMessageDao, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sms.DaoSms/SendMessageDao", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaoSmsServer is the server API for DaoSms service.
type DaoSmsServer interface {
	SendMessageDao(context.Context, *ReqSendMessageDao) (*emptypb.Empty, error)
}

// UnimplementedDaoSmsServer can be embedded to have forward compatible implementations.
type UnimplementedDaoSmsServer struct {
}

func (*UnimplementedDaoSmsServer) SendMessageDao(context.Context, *ReqSendMessageDao) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessageDao not implemented")
}

func RegisterDaoSmsServer(s *grpc.Server, srv DaoSmsServer) {
	s.RegisterService(&_DaoSms_serviceDesc, srv)
}

func _DaoSms_SendMessageDao_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqSendMessageDao)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaoSmsServer).SendMessageDao(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sms.DaoSms/SendMessageDao",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaoSmsServer).SendMessageDao(ctx, req.(*ReqSendMessageDao))
	}
	return interceptor(ctx, in, info, handler)
}

var _DaoSms_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sms.DaoSms",
	HandlerType: (*DaoSmsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessageDao",
			Handler:    _DaoSms_SendMessageDao_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sms.proto",
}
