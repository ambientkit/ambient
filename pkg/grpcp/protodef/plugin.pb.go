// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: plugin.proto

package protodef

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PluginNameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *PluginNameResponse) Reset() {
	*x = PluginNameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginNameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginNameResponse) ProtoMessage() {}

func (x *PluginNameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginNameResponse.ProtoReflect.Descriptor instead.
func (*PluginNameResponse) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *PluginNameResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PluginVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *PluginVersionResponse) Reset() {
	*x = PluginVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginVersionResponse) ProtoMessage() {}

func (x *PluginVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginVersionResponse.ProtoReflect.Descriptor instead.
func (*PluginVersionResponse) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *PluginVersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type GrantRequestsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GrantRequest []*GrantRequest `protobuf:"bytes,1,rep,name=GrantRequest,proto3" json:"GrantRequest,omitempty"`
}

func (x *GrantRequestsResponse) Reset() {
	*x = GrantRequestsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrantRequestsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantRequestsResponse) ProtoMessage() {}

func (x *GrantRequestsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantRequestsResponse.ProtoReflect.Descriptor instead.
func (*GrantRequestsResponse) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *GrantRequestsResponse) GetGrantRequest() []*GrantRequest {
	if x != nil {
		return x.GrantRequest
	}
	return nil
}

type GrantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Grant       string `protobuf:"bytes,1,opt,name=grant,proto3" json:"grant,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *GrantRequest) Reset() {
	*x = GrantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantRequest) ProtoMessage() {}

func (x *GrantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantRequest.ProtoReflect.Descriptor instead.
func (*GrantRequest) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{3}
}

func (x *GrantRequest) GetGrant() string {
	if x != nil {
		return x.Grant
	}
	return ""
}

func (x *GrantRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type Toolkit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid uint32 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *Toolkit) Reset() {
	*x = Toolkit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Toolkit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Toolkit) ProtoMessage() {}

func (x *Toolkit) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Toolkit.ProtoReflect.Descriptor instead.
func (*Toolkit) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{4}
}

func (x *Toolkit) GetUid() uint32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type EnableResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid uint32 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *EnableResponse) Reset() {
	*x = EnableResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnableResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnableResponse) ProtoMessage() {}

func (x *EnableResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnableResponse.ProtoReflect.Descriptor instead.
func (*EnableResponse) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{5}
}

func (x *EnableResponse) GetUid() uint32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

var File_plugin_proto protoreflect.FileDescriptor

var file_plugin_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10,
	0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66,
	0x1a, 0x0b, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x28, 0x0a,
	0x12, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x31, 0x0a, 0x15, 0x50, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x5b, 0x0a, 0x15, 0x47, 0x72,
	0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x0c, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x6d, 0x62, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x47, 0x72, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x0c, 0x47, 0x72, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x46, 0x0a, 0x0c, 0x47, 0x72, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x61, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x1b, 0x0a, 0x07, 0x54, 0x6f, 0x6f, 0x6c, 0x6b, 0x69, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x22, 0x0a, 0x0e,
	0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x69, 0x64,
	0x32, 0xf9, 0x02, 0x0a, 0x0d, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x12, 0x4d, 0x0a, 0x0a, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x24, 0x2e, 0x61, 0x6d, 0x62, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x50, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x53, 0x0a, 0x0d, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x27, 0x2e, 0x61, 0x6d,
	0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x06, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x19, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x64, 0x65, 0x66, 0x2e, 0x54, 0x6f, 0x6f, 0x6c, 0x6b, 0x69, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x6d,
	0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45,
	0x6e, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x3d, 0x0a, 0x07, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x17, 0x2e, 0x61, 0x6d, 0x62,
	0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c,
	0x0a, 0x06, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x12, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b,
	0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_plugin_proto_rawDescOnce sync.Once
	file_plugin_proto_rawDescData = file_plugin_proto_rawDesc
)

func file_plugin_proto_rawDescGZIP() []byte {
	file_plugin_proto_rawDescOnce.Do(func() {
		file_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_proto_rawDescData)
	})
	return file_plugin_proto_rawDescData
}

var file_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_plugin_proto_goTypes = []interface{}{
	(*PluginNameResponse)(nil),    // 0: ambient.protodef.PluginNameResponse
	(*PluginVersionResponse)(nil), // 1: ambient.protodef.PluginVersionResponse
	(*GrantRequestsResponse)(nil), // 2: ambient.protodef.GrantRequestsResponse
	(*GrantRequest)(nil),          // 3: ambient.protodef.GrantRequest
	(*Toolkit)(nil),               // 4: ambient.protodef.Toolkit
	(*EnableResponse)(nil),        // 5: ambient.protodef.EnableResponse
	(*Empty)(nil),                 // 6: ambient.protodef.Empty
}
var file_plugin_proto_depIdxs = []int32{
	3, // 0: ambient.protodef.GrantRequestsResponse.GrantRequest:type_name -> ambient.protodef.GrantRequest
	6, // 1: ambient.protodef.GenericPlugin.PluginName:input_type -> ambient.protodef.Empty
	6, // 2: ambient.protodef.GenericPlugin.PluginVersion:input_type -> ambient.protodef.Empty
	4, // 3: ambient.protodef.GenericPlugin.Enable:input_type -> ambient.protodef.Toolkit
	6, // 4: ambient.protodef.GenericPlugin.Disable:input_type -> ambient.protodef.Empty
	6, // 5: ambient.protodef.GenericPlugin.Routes:input_type -> ambient.protodef.Empty
	0, // 6: ambient.protodef.GenericPlugin.PluginName:output_type -> ambient.protodef.PluginNameResponse
	1, // 7: ambient.protodef.GenericPlugin.PluginVersion:output_type -> ambient.protodef.PluginVersionResponse
	5, // 8: ambient.protodef.GenericPlugin.Enable:output_type -> ambient.protodef.EnableResponse
	6, // 9: ambient.protodef.GenericPlugin.Disable:output_type -> ambient.protodef.Empty
	6, // 10: ambient.protodef.GenericPlugin.Routes:output_type -> ambient.protodef.Empty
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_plugin_proto_init() }
func file_plugin_proto_init() {
	if File_plugin_proto != nil {
		return
	}
	file_empty_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginNameResponse); i {
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
		file_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginVersionResponse); i {
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
		file_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrantRequestsResponse); i {
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
		file_plugin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrantRequest); i {
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
		file_plugin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Toolkit); i {
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
		file_plugin_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnableResponse); i {
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
			RawDescriptor: file_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_proto_goTypes,
		DependencyIndexes: file_plugin_proto_depIdxs,
		MessageInfos:      file_plugin_proto_msgTypes,
	}.Build()
	File_plugin_proto = out.File
	file_plugin_proto_rawDesc = nil
	file_plugin_proto_goTypes = nil
	file_plugin_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GenericPluginClient is the client API for GenericPlugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GenericPluginClient interface {
	PluginName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PluginNameResponse, error)
	PluginVersion(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PluginVersionResponse, error)
	// rpc GrantRequests(Empty) returns (GrantRequestsResponse) {}
	Enable(ctx context.Context, in *Toolkit, opts ...grpc.CallOption) (*EnableResponse, error)
	Disable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Routes(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type genericPluginClient struct {
	cc grpc.ClientConnInterface
}

func NewGenericPluginClient(cc grpc.ClientConnInterface) GenericPluginClient {
	return &genericPluginClient{cc}
}

func (c *genericPluginClient) PluginName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PluginNameResponse, error) {
	out := new(PluginNameResponse)
	err := c.cc.Invoke(ctx, "/ambient.protodef.GenericPlugin/PluginName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genericPluginClient) PluginVersion(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PluginVersionResponse, error) {
	out := new(PluginVersionResponse)
	err := c.cc.Invoke(ctx, "/ambient.protodef.GenericPlugin/PluginVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genericPluginClient) Enable(ctx context.Context, in *Toolkit, opts ...grpc.CallOption) (*EnableResponse, error) {
	out := new(EnableResponse)
	err := c.cc.Invoke(ctx, "/ambient.protodef.GenericPlugin/Enable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genericPluginClient) Disable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ambient.protodef.GenericPlugin/Disable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genericPluginClient) Routes(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ambient.protodef.GenericPlugin/Routes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GenericPluginServer is the server API for GenericPlugin service.
type GenericPluginServer interface {
	PluginName(context.Context, *Empty) (*PluginNameResponse, error)
	PluginVersion(context.Context, *Empty) (*PluginVersionResponse, error)
	// rpc GrantRequests(Empty) returns (GrantRequestsResponse) {}
	Enable(context.Context, *Toolkit) (*EnableResponse, error)
	Disable(context.Context, *Empty) (*Empty, error)
	Routes(context.Context, *Empty) (*Empty, error)
}

// UnimplementedGenericPluginServer can be embedded to have forward compatible implementations.
type UnimplementedGenericPluginServer struct {
}

func (*UnimplementedGenericPluginServer) PluginName(context.Context, *Empty) (*PluginNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PluginName not implemented")
}
func (*UnimplementedGenericPluginServer) PluginVersion(context.Context, *Empty) (*PluginVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PluginVersion not implemented")
}
func (*UnimplementedGenericPluginServer) Enable(context.Context, *Toolkit) (*EnableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enable not implemented")
}
func (*UnimplementedGenericPluginServer) Disable(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disable not implemented")
}
func (*UnimplementedGenericPluginServer) Routes(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Routes not implemented")
}

func RegisterGenericPluginServer(s *grpc.Server, srv GenericPluginServer) {
	s.RegisterService(&_GenericPlugin_serviceDesc, srv)
}

func _GenericPlugin_PluginName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenericPluginServer).PluginName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.GenericPlugin/PluginName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenericPluginServer).PluginName(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenericPlugin_PluginVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenericPluginServer).PluginVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.GenericPlugin/PluginVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenericPluginServer).PluginVersion(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenericPlugin_Enable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Toolkit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenericPluginServer).Enable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.GenericPlugin/Enable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenericPluginServer).Enable(ctx, req.(*Toolkit))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenericPlugin_Disable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenericPluginServer).Disable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.GenericPlugin/Disable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenericPluginServer).Disable(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenericPlugin_Routes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenericPluginServer).Routes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.GenericPlugin/Routes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenericPluginServer).Routes(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _GenericPlugin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ambient.protodef.GenericPlugin",
	HandlerType: (*GenericPluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PluginName",
			Handler:    _GenericPlugin_PluginName_Handler,
		},
		{
			MethodName: "PluginVersion",
			Handler:    _GenericPlugin_PluginVersion_Handler,
		},
		{
			MethodName: "Enable",
			Handler:    _GenericPlugin_Enable_Handler,
		},
		{
			MethodName: "Disable",
			Handler:    _GenericPlugin_Disable_Handler,
		},
		{
			MethodName: "Routes",
			Handler:    _GenericPlugin_Routes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin.proto",
}