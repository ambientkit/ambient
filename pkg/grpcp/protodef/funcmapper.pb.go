// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: funcmapper.proto

package protodef

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FuncMapperDoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requestid string       `protobuf:"bytes,1,opt,name=requestid,proto3" json:"requestid,omitempty"`
	Key       string       `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Params    []*anypb.Any `protobuf:"bytes,3,rep,name=params,proto3" json:"params,omitempty"`
}

func (x *FuncMapperDoRequest) Reset() {
	*x = FuncMapperDoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_funcmapper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FuncMapperDoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FuncMapperDoRequest) ProtoMessage() {}

func (x *FuncMapperDoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_funcmapper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FuncMapperDoRequest.ProtoReflect.Descriptor instead.
func (*FuncMapperDoRequest) Descriptor() ([]byte, []int) {
	return file_funcmapper_proto_rawDescGZIP(), []int{0}
}

func (x *FuncMapperDoRequest) GetRequestid() string {
	if x != nil {
		return x.Requestid
	}
	return ""
}

func (x *FuncMapperDoRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FuncMapperDoRequest) GetParams() []*anypb.Any {
	if x != nil {
		return x.Params
	}
	return nil
}

type FuncMapperDoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value *anypb.Any         `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Arr   []*structpb.Struct `protobuf:"bytes,2,rep,name=arr,proto3" json:"arr,omitempty"`
}

func (x *FuncMapperDoResponse) Reset() {
	*x = FuncMapperDoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_funcmapper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FuncMapperDoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FuncMapperDoResponse) ProtoMessage() {}

func (x *FuncMapperDoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_funcmapper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FuncMapperDoResponse.ProtoReflect.Descriptor instead.
func (*FuncMapperDoResponse) Descriptor() ([]byte, []int) {
	return file_funcmapper_proto_rawDescGZIP(), []int{1}
}

func (x *FuncMapperDoResponse) GetValue() *anypb.Any {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *FuncMapperDoResponse) GetArr() []*structpb.Struct {
	if x != nil {
		return x.Arr
	}
	return nil
}

var File_funcmapper_proto protoreflect.FileDescriptor

var file_funcmapper_proto_rawDesc = []byte{
	0x0a, 0x10, 0x66, 0x75, 0x6e, 0x63, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x64, 0x65, 0x66, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a,
	0x13, 0x46, 0x75, 0x6e, 0x63, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x44, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x22, 0x6d, 0x0a, 0x14, 0x46, 0x75, 0x6e, 0x63, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x44, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x29, 0x0a, 0x03, 0x61, 0x72, 0x72, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x03, 0x61, 0x72,
	0x72, 0x32, 0x63, 0x0a, 0x0a, 0x46, 0x75, 0x6e, 0x63, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12,
	0x55, 0x0a, 0x02, 0x44, 0x6f, 0x12, 0x25, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x46, 0x75, 0x6e, 0x63, 0x4d, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x44, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61,
	0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e,
	0x46, 0x75, 0x6e, 0x63, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x44, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x64, 0x65, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_funcmapper_proto_rawDescOnce sync.Once
	file_funcmapper_proto_rawDescData = file_funcmapper_proto_rawDesc
)

func file_funcmapper_proto_rawDescGZIP() []byte {
	file_funcmapper_proto_rawDescOnce.Do(func() {
		file_funcmapper_proto_rawDescData = protoimpl.X.CompressGZIP(file_funcmapper_proto_rawDescData)
	})
	return file_funcmapper_proto_rawDescData
}

var file_funcmapper_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_funcmapper_proto_goTypes = []interface{}{
	(*FuncMapperDoRequest)(nil),  // 0: ambient.protodef.FuncMapperDoRequest
	(*FuncMapperDoResponse)(nil), // 1: ambient.protodef.FuncMapperDoResponse
	(*anypb.Any)(nil),            // 2: google.protobuf.Any
	(*structpb.Struct)(nil),      // 3: google.protobuf.Struct
}
var file_funcmapper_proto_depIdxs = []int32{
	2, // 0: ambient.protodef.FuncMapperDoRequest.params:type_name -> google.protobuf.Any
	2, // 1: ambient.protodef.FuncMapperDoResponse.value:type_name -> google.protobuf.Any
	3, // 2: ambient.protodef.FuncMapperDoResponse.arr:type_name -> google.protobuf.Struct
	0, // 3: ambient.protodef.FuncMapper.Do:input_type -> ambient.protodef.FuncMapperDoRequest
	1, // 4: ambient.protodef.FuncMapper.Do:output_type -> ambient.protodef.FuncMapperDoResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_funcmapper_proto_init() }
func file_funcmapper_proto_init() {
	if File_funcmapper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_funcmapper_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FuncMapperDoRequest); i {
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
		file_funcmapper_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FuncMapperDoResponse); i {
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
			RawDescriptor: file_funcmapper_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_funcmapper_proto_goTypes,
		DependencyIndexes: file_funcmapper_proto_depIdxs,
		MessageInfos:      file_funcmapper_proto_msgTypes,
	}.Build()
	File_funcmapper_proto = out.File
	file_funcmapper_proto_rawDesc = nil
	file_funcmapper_proto_goTypes = nil
	file_funcmapper_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// FuncMapperClient is the client API for FuncMapper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FuncMapperClient interface {
	Do(ctx context.Context, in *FuncMapperDoRequest, opts ...grpc.CallOption) (*FuncMapperDoResponse, error)
}

type funcMapperClient struct {
	cc grpc.ClientConnInterface
}

func NewFuncMapperClient(cc grpc.ClientConnInterface) FuncMapperClient {
	return &funcMapperClient{cc}
}

func (c *funcMapperClient) Do(ctx context.Context, in *FuncMapperDoRequest, opts ...grpc.CallOption) (*FuncMapperDoResponse, error) {
	out := new(FuncMapperDoResponse)
	err := c.cc.Invoke(ctx, "/ambient.protodef.FuncMapper/Do", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FuncMapperServer is the server API for FuncMapper service.
type FuncMapperServer interface {
	Do(context.Context, *FuncMapperDoRequest) (*FuncMapperDoResponse, error)
}

// UnimplementedFuncMapperServer can be embedded to have forward compatible implementations.
type UnimplementedFuncMapperServer struct {
}

func (*UnimplementedFuncMapperServer) Do(context.Context, *FuncMapperDoRequest) (*FuncMapperDoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Do not implemented")
}

func RegisterFuncMapperServer(s *grpc.Server, srv FuncMapperServer) {
	s.RegisterService(&_FuncMapper_serviceDesc, srv)
}

func _FuncMapper_Do_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FuncMapperDoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FuncMapperServer).Do(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.FuncMapper/Do",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FuncMapperServer).Do(ctx, req.(*FuncMapperDoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _FuncMapper_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ambient.protodef.FuncMapper",
	HandlerType: (*FuncMapperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Do",
			Handler:    _FuncMapper_Do_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "funcmapper.proto",
}
