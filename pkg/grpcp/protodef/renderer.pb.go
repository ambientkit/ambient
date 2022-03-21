// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: renderer.proto

package protodef

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type RendererPageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requestid    string           `protobuf:"bytes,1,opt,name=requestid,proto3" json:"requestid,omitempty"`
	Templatename string           `protobuf:"bytes,2,opt,name=templatename,proto3" json:"templatename,omitempty"`
	Vars         *structpb.Struct `protobuf:"bytes,3,opt,name=vars,proto3" json:"vars,omitempty"`
	Keys         []string         `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
	Files        []*EmbeddedFile  `protobuf:"bytes,5,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *RendererPageRequest) Reset() {
	*x = RendererPageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RendererPageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RendererPageRequest) ProtoMessage() {}

func (x *RendererPageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RendererPageRequest.ProtoReflect.Descriptor instead.
func (*RendererPageRequest) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{0}
}

func (x *RendererPageRequest) GetRequestid() string {
	if x != nil {
		return x.Requestid
	}
	return ""
}

func (x *RendererPageRequest) GetTemplatename() string {
	if x != nil {
		return x.Templatename
	}
	return ""
}

func (x *RendererPageRequest) GetVars() *structpb.Struct {
	if x != nil {
		return x.Vars
	}
	return nil
}

func (x *RendererPageRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *RendererPageRequest) GetFiles() []*EmbeddedFile {
	if x != nil {
		return x.Files
	}
	return nil
}

type RendererPageContentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requestid string           `protobuf:"bytes,1,opt,name=requestid,proto3" json:"requestid,omitempty"`
	Content   string           `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Vars      *structpb.Struct `protobuf:"bytes,3,opt,name=vars,proto3" json:"vars,omitempty"`
	Keys      []string         `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *RendererPageContentRequest) Reset() {
	*x = RendererPageContentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RendererPageContentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RendererPageContentRequest) ProtoMessage() {}

func (x *RendererPageContentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RendererPageContentRequest.ProtoReflect.Descriptor instead.
func (*RendererPageContentRequest) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{1}
}

func (x *RendererPageContentRequest) GetRequestid() string {
	if x != nil {
		return x.Requestid
	}
	return ""
}

func (x *RendererPageContentRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *RendererPageContentRequest) GetVars() *structpb.Struct {
	if x != nil {
		return x.Vars
	}
	return nil
}

func (x *RendererPageContentRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type RendererPostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requestid    string           `protobuf:"bytes,1,opt,name=requestid,proto3" json:"requestid,omitempty"`
	Templatename string           `protobuf:"bytes,2,opt,name=templatename,proto3" json:"templatename,omitempty"`
	Vars         *structpb.Struct `protobuf:"bytes,3,opt,name=vars,proto3" json:"vars,omitempty"`
	Keys         []string         `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *RendererPostRequest) Reset() {
	*x = RendererPostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RendererPostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RendererPostRequest) ProtoMessage() {}

func (x *RendererPostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RendererPostRequest.ProtoReflect.Descriptor instead.
func (*RendererPostRequest) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{2}
}

func (x *RendererPostRequest) GetRequestid() string {
	if x != nil {
		return x.Requestid
	}
	return ""
}

func (x *RendererPostRequest) GetTemplatename() string {
	if x != nil {
		return x.Templatename
	}
	return ""
}

func (x *RendererPostRequest) GetVars() *structpb.Struct {
	if x != nil {
		return x.Vars
	}
	return nil
}

func (x *RendererPostRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type RendererPostContentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requestid string           `protobuf:"bytes,1,opt,name=requestid,proto3" json:"requestid,omitempty"`
	Content   string           `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Vars      *structpb.Struct `protobuf:"bytes,3,opt,name=vars,proto3" json:"vars,omitempty"`
	Keys      []string         `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *RendererPostContentRequest) Reset() {
	*x = RendererPostContentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RendererPostContentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RendererPostContentRequest) ProtoMessage() {}

func (x *RendererPostContentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RendererPostContentRequest.ProtoReflect.Descriptor instead.
func (*RendererPostContentRequest) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{3}
}

func (x *RendererPostContentRequest) GetRequestid() string {
	if x != nil {
		return x.Requestid
	}
	return ""
}

func (x *RendererPostContentRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *RendererPostContentRequest) GetVars() *structpb.Struct {
	if x != nil {
		return x.Vars
	}
	return nil
}

func (x *RendererPostContentRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type RendererErrorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requestid  string           `protobuf:"bytes,1,opt,name=requestid,proto3" json:"requestid,omitempty"`
	Content    string           `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Vars       *structpb.Struct `protobuf:"bytes,3,opt,name=vars,proto3" json:"vars,omitempty"`
	Keys       []string         `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
	Statuscode uint32           `protobuf:"varint,5,opt,name=statuscode,proto3" json:"statuscode,omitempty"`
}

func (x *RendererErrorRequest) Reset() {
	*x = RendererErrorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RendererErrorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RendererErrorRequest) ProtoMessage() {}

func (x *RendererErrorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RendererErrorRequest.ProtoReflect.Descriptor instead.
func (*RendererErrorRequest) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{4}
}

func (x *RendererErrorRequest) GetRequestid() string {
	if x != nil {
		return x.Requestid
	}
	return ""
}

func (x *RendererErrorRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *RendererErrorRequest) GetVars() *structpb.Struct {
	if x != nil {
		return x.Vars
	}
	return nil
}

func (x *RendererErrorRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *RendererErrorRequest) GetStatuscode() uint32 {
	if x != nil {
		return x.Statuscode
	}
	return 0
}

type EmbeddedFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Body string `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *EmbeddedFile) Reset() {
	*x = EmbeddedFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmbeddedFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbeddedFile) ProtoMessage() {}

func (x *EmbeddedFile) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbeddedFile.ProtoReflect.Descriptor instead.
func (*EmbeddedFile) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{5}
}

func (x *EmbeddedFile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EmbeddedFile) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

var File_renderer_proto protoreflect.FileDescriptor

var file_renderer_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x10, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64,
	0x65, 0x66, 0x1a, 0x0b, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x01,
	0x0a, 0x13, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x76, 0x61, 0x72, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04,
	0x76, 0x61, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x34, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x62, 0x65, 0x64,
	0x64, 0x65, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x95,
	0x01, 0x0a, 0x1a, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x76, 0x61, 0x72, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x76, 0x61,
	0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0x98, 0x01, 0x0a, 0x13, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x65, 0x72, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64, 0x12, 0x22, 0x0a, 0x0c,
	0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x2b, 0x0a, 0x04, 0x76, 0x61, 0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x76, 0x61, 0x72, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79,
	0x73, 0x22, 0x95, 0x01, 0x0a, 0x1a, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x50, 0x6f,
	0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x76, 0x61, 0x72, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52,
	0x04, 0x76, 0x61, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0xaf, 0x01, 0x0a, 0x14, 0x52, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x76, 0x61,
	0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x52, 0x04, 0x76, 0x61, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x36, 0x0a, 0x0c, 0x45,
	0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x32, 0x9a, 0x03, 0x0a, 0x08, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72,
	0x12, 0x48, 0x0a, 0x04, 0x50, 0x61, 0x67, 0x65, 0x12, 0x25, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x52, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64,
	0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x0b, 0x50, 0x61,
	0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2c, 0x2e, 0x61, 0x6d, 0x62, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x48, 0x0a, 0x04, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x25, 0x2e, 0x61, 0x6d, 0x62,
	0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x52, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x0b,
	0x50, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2c, 0x2e, 0x61, 0x6d,
	0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x52,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x26, 0x2e,
	0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66,
	0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x66, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_renderer_proto_rawDescOnce sync.Once
	file_renderer_proto_rawDescData = file_renderer_proto_rawDesc
)

func file_renderer_proto_rawDescGZIP() []byte {
	file_renderer_proto_rawDescOnce.Do(func() {
		file_renderer_proto_rawDescData = protoimpl.X.CompressGZIP(file_renderer_proto_rawDescData)
	})
	return file_renderer_proto_rawDescData
}

var file_renderer_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_renderer_proto_goTypes = []interface{}{
	(*RendererPageRequest)(nil),        // 0: ambient.protodef.RendererPageRequest
	(*RendererPageContentRequest)(nil), // 1: ambient.protodef.RendererPageContentRequest
	(*RendererPostRequest)(nil),        // 2: ambient.protodef.RendererPostRequest
	(*RendererPostContentRequest)(nil), // 3: ambient.protodef.RendererPostContentRequest
	(*RendererErrorRequest)(nil),       // 4: ambient.protodef.RendererErrorRequest
	(*EmbeddedFile)(nil),               // 5: ambient.protodef.EmbeddedFile
	(*structpb.Struct)(nil),            // 6: google.protobuf.Struct
	(*Empty)(nil),                      // 7: ambient.protodef.Empty
}
var file_renderer_proto_depIdxs = []int32{
	6,  // 0: ambient.protodef.RendererPageRequest.vars:type_name -> google.protobuf.Struct
	5,  // 1: ambient.protodef.RendererPageRequest.files:type_name -> ambient.protodef.EmbeddedFile
	6,  // 2: ambient.protodef.RendererPageContentRequest.vars:type_name -> google.protobuf.Struct
	6,  // 3: ambient.protodef.RendererPostRequest.vars:type_name -> google.protobuf.Struct
	6,  // 4: ambient.protodef.RendererPostContentRequest.vars:type_name -> google.protobuf.Struct
	6,  // 5: ambient.protodef.RendererErrorRequest.vars:type_name -> google.protobuf.Struct
	0,  // 6: ambient.protodef.Renderer.Page:input_type -> ambient.protodef.RendererPageRequest
	1,  // 7: ambient.protodef.Renderer.PageContent:input_type -> ambient.protodef.RendererPageContentRequest
	2,  // 8: ambient.protodef.Renderer.Post:input_type -> ambient.protodef.RendererPostRequest
	3,  // 9: ambient.protodef.Renderer.PostContent:input_type -> ambient.protodef.RendererPostContentRequest
	4,  // 10: ambient.protodef.Renderer.Error:input_type -> ambient.protodef.RendererErrorRequest
	7,  // 11: ambient.protodef.Renderer.Page:output_type -> ambient.protodef.Empty
	7,  // 12: ambient.protodef.Renderer.PageContent:output_type -> ambient.protodef.Empty
	7,  // 13: ambient.protodef.Renderer.Post:output_type -> ambient.protodef.Empty
	7,  // 14: ambient.protodef.Renderer.PostContent:output_type -> ambient.protodef.Empty
	7,  // 15: ambient.protodef.Renderer.Error:output_type -> ambient.protodef.Empty
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_renderer_proto_init() }
func file_renderer_proto_init() {
	if File_renderer_proto != nil {
		return
	}
	file_empty_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_renderer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RendererPageRequest); i {
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
		file_renderer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RendererPageContentRequest); i {
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
		file_renderer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RendererPostRequest); i {
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
		file_renderer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RendererPostContentRequest); i {
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
		file_renderer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RendererErrorRequest); i {
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
		file_renderer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmbeddedFile); i {
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
			RawDescriptor: file_renderer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_renderer_proto_goTypes,
		DependencyIndexes: file_renderer_proto_depIdxs,
		MessageInfos:      file_renderer_proto_msgTypes,
	}.Build()
	File_renderer_proto = out.File
	file_renderer_proto_rawDesc = nil
	file_renderer_proto_goTypes = nil
	file_renderer_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RendererClient is the client API for Renderer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RendererClient interface {
	Page(ctx context.Context, in *RendererPageRequest, opts ...grpc.CallOption) (*Empty, error)
	PageContent(ctx context.Context, in *RendererPageContentRequest, opts ...grpc.CallOption) (*Empty, error)
	Post(ctx context.Context, in *RendererPostRequest, opts ...grpc.CallOption) (*Empty, error)
	PostContent(ctx context.Context, in *RendererPostContentRequest, opts ...grpc.CallOption) (*Empty, error)
	Error(ctx context.Context, in *RendererErrorRequest, opts ...grpc.CallOption) (*Empty, error)
}

type rendererClient struct {
	cc grpc.ClientConnInterface
}

func NewRendererClient(cc grpc.ClientConnInterface) RendererClient {
	return &rendererClient{cc}
}

func (c *rendererClient) Page(ctx context.Context, in *RendererPageRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ambient.protodef.Renderer/Page", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendererClient) PageContent(ctx context.Context, in *RendererPageContentRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ambient.protodef.Renderer/PageContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendererClient) Post(ctx context.Context, in *RendererPostRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ambient.protodef.Renderer/Post", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendererClient) PostContent(ctx context.Context, in *RendererPostContentRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ambient.protodef.Renderer/PostContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendererClient) Error(ctx context.Context, in *RendererErrorRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ambient.protodef.Renderer/Error", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RendererServer is the server API for Renderer service.
type RendererServer interface {
	Page(context.Context, *RendererPageRequest) (*Empty, error)
	PageContent(context.Context, *RendererPageContentRequest) (*Empty, error)
	Post(context.Context, *RendererPostRequest) (*Empty, error)
	PostContent(context.Context, *RendererPostContentRequest) (*Empty, error)
	Error(context.Context, *RendererErrorRequest) (*Empty, error)
}

// UnimplementedRendererServer can be embedded to have forward compatible implementations.
type UnimplementedRendererServer struct {
}

func (*UnimplementedRendererServer) Page(context.Context, *RendererPageRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Page not implemented")
}
func (*UnimplementedRendererServer) PageContent(context.Context, *RendererPageContentRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PageContent not implemented")
}
func (*UnimplementedRendererServer) Post(context.Context, *RendererPostRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Post not implemented")
}
func (*UnimplementedRendererServer) PostContent(context.Context, *RendererPostContentRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostContent not implemented")
}
func (*UnimplementedRendererServer) Error(context.Context, *RendererErrorRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Error not implemented")
}

func RegisterRendererServer(s *grpc.Server, srv RendererServer) {
	s.RegisterService(&_Renderer_serviceDesc, srv)
}

func _Renderer_Page_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RendererPageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendererServer).Page(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.Renderer/Page",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendererServer).Page(ctx, req.(*RendererPageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Renderer_PageContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RendererPageContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendererServer).PageContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.Renderer/PageContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendererServer).PageContent(ctx, req.(*RendererPageContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Renderer_Post_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RendererPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendererServer).Post(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.Renderer/Post",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendererServer).Post(ctx, req.(*RendererPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Renderer_PostContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RendererPostContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendererServer).PostContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.Renderer/PostContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendererServer).PostContent(ctx, req.(*RendererPostContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Renderer_Error_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RendererErrorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendererServer).Error(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ambient.protodef.Renderer/Error",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendererServer).Error(ctx, req.(*RendererErrorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Renderer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ambient.protodef.Renderer",
	HandlerType: (*RendererServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Page",
			Handler:    _Renderer_Page_Handler,
		},
		{
			MethodName: "PageContent",
			Handler:    _Renderer_PageContent_Handler,
		},
		{
			MethodName: "Post",
			Handler:    _Renderer_Post_Handler,
		},
		{
			MethodName: "PostContent",
			Handler:    _Renderer_PostContent_Handler,
		},
		{
			MethodName: "Error",
			Handler:    _Renderer_Error_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "renderer.proto",
}
