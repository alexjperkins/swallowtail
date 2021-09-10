// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: s.satoshi/proto/satoshi.proto

package satoshiproto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type PublishStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PublishStatusRequest) Reset() {
	*x = PublishStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_satoshi_proto_satoshi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishStatusRequest) ProtoMessage() {}

func (x *PublishStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_s_satoshi_proto_satoshi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishStatusRequest.ProtoReflect.Descriptor instead.
func (*PublishStatusRequest) Descriptor() ([]byte, []int) {
	return file_s_satoshi_proto_satoshi_proto_rawDescGZIP(), []int{0}
}

type PublishStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alive   bool   `protobuf:"varint,1,opt,name=alive,proto3" json:"alive,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *PublishStatusResponse) Reset() {
	*x = PublishStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_satoshi_proto_satoshi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishStatusResponse) ProtoMessage() {}

func (x *PublishStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_s_satoshi_proto_satoshi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishStatusResponse.ProtoReflect.Descriptor instead.
func (*PublishStatusResponse) Descriptor() ([]byte, []int) {
	return file_s_satoshi_proto_satoshi_proto_rawDescGZIP(), []int{1}
}

func (x *PublishStatusResponse) GetAlive() bool {
	if x != nil {
		return x.Alive
	}
	return false
}

func (x *PublishStatusResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_s_satoshi_proto_satoshi_proto protoreflect.FileDescriptor

var file_s_satoshi_proto_satoshi_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x2e, 0x73, 0x61, 0x74, 0x6f, 0x73, 0x68, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x61, 0x74, 0x6f, 0x73, 0x68, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x16, 0x0a, 0x14, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x47, 0x0a, 0x15, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x05, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x32, 0x4b, 0x0a, 0x07, 0x73, 0x61, 0x74, 0x6f, 0x73, 0x68, 0x69, 0x12, 0x40, 0x0a, 0x0d, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x15, 0x2e, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x11, 0x5a,
	0x0f, 0x2e, 0x2f, 0x3b, 0x73, 0x61, 0x74, 0x6f, 0x73, 0x68, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_s_satoshi_proto_satoshi_proto_rawDescOnce sync.Once
	file_s_satoshi_proto_satoshi_proto_rawDescData = file_s_satoshi_proto_satoshi_proto_rawDesc
)

func file_s_satoshi_proto_satoshi_proto_rawDescGZIP() []byte {
	file_s_satoshi_proto_satoshi_proto_rawDescOnce.Do(func() {
		file_s_satoshi_proto_satoshi_proto_rawDescData = protoimpl.X.CompressGZIP(file_s_satoshi_proto_satoshi_proto_rawDescData)
	})
	return file_s_satoshi_proto_satoshi_proto_rawDescData
}

var file_s_satoshi_proto_satoshi_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_s_satoshi_proto_satoshi_proto_goTypes = []interface{}{
	(*PublishStatusRequest)(nil),  // 0: PublishStatusRequest
	(*PublishStatusResponse)(nil), // 1: PublishStatusResponse
}
var file_s_satoshi_proto_satoshi_proto_depIdxs = []int32{
	0, // 0: satoshi.PublishStatus:input_type -> PublishStatusRequest
	1, // 1: satoshi.PublishStatus:output_type -> PublishStatusResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_s_satoshi_proto_satoshi_proto_init() }
func file_s_satoshi_proto_satoshi_proto_init() {
	if File_s_satoshi_proto_satoshi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_s_satoshi_proto_satoshi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishStatusRequest); i {
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
		file_s_satoshi_proto_satoshi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishStatusResponse); i {
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
			RawDescriptor: file_s_satoshi_proto_satoshi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_s_satoshi_proto_satoshi_proto_goTypes,
		DependencyIndexes: file_s_satoshi_proto_satoshi_proto_depIdxs,
		MessageInfos:      file_s_satoshi_proto_satoshi_proto_msgTypes,
	}.Build()
	File_s_satoshi_proto_satoshi_proto = out.File
	file_s_satoshi_proto_satoshi_proto_rawDesc = nil
	file_s_satoshi_proto_satoshi_proto_goTypes = nil
	file_s_satoshi_proto_satoshi_proto_depIdxs = nil
}
