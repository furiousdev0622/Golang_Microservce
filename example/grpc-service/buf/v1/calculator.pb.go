// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.2
// source: buf/v1/calculator.proto

package v1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/httpbody"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/anypb"
	_ "google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddIntRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// param a
	A int32 `protobuf:"varint,1,opt,name=a,proto3" json:"a,omitempty"`
	// param b
	B int32 `protobuf:"varint,2,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *AddIntRequest) Reset() {
	*x = AddIntRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_v1_calculator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddIntRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddIntRequest) ProtoMessage() {}

func (x *AddIntRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_v1_calculator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddIntRequest.ProtoReflect.Descriptor instead.
func (*AddIntRequest) Descriptor() ([]byte, []int) {
	return file_buf_v1_calculator_proto_rawDescGZIP(), []int{0}
}

func (x *AddIntRequest) GetA() int32 {
	if x != nil {
		return x.A
	}
	return 0
}

func (x *AddIntRequest) GetB() int32 {
	if x != nil {
		return x.B
	}
	return 0
}

type AddIntResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result int32 `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *AddIntResponse) Reset() {
	*x = AddIntResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_v1_calculator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddIntResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddIntResponse) ProtoMessage() {}

func (x *AddIntResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_v1_calculator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddIntResponse.ProtoReflect.Descriptor instead.
func (*AddIntResponse) Descriptor() ([]byte, []int) {
	return file_buf_v1_calculator_proto_rawDescGZIP(), []int{1}
}

func (x *AddIntResponse) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

var File_buf_v1_calculator_proto protoreflect.FileDescriptor

var file_buf_v1_calculator_proto_rawDesc = []byte{
	0x0a, 0x17, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61,
	0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x62, 0x75, 0x66, 0x2e, 0x76,
	0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x68, 0x74,
	0x74, 0x70, 0x62, 0x6f, 0x64, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x0d, 0x41, 0x64, 0x64,
	0x49, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x01, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x1a, 0x04, 0x10, 0x64, 0x28, 0x00,
	0x52, 0x01, 0x61, 0x12, 0x0c, 0x0a, 0x01, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01,
	0x62, 0x22, 0x28, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x6b, 0x0a, 0x11, 0x43,
	0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x56, 0x0a, 0x06, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x74, 0x12, 0x15, 0x2e, 0x62, 0x75, 0x66,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x49, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x17, 0x22, 0x12, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x6c, 0x63, 0x75,
	0x6c, 0x61, 0x74, 0x6f, 0x72, 0x3a, 0x01, 0x2a, 0x42, 0x11, 0x5a, 0x0f, 0x67, 0x72, 0x70, 0x63,
	0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_buf_v1_calculator_proto_rawDescOnce sync.Once
	file_buf_v1_calculator_proto_rawDescData = file_buf_v1_calculator_proto_rawDesc
)

func file_buf_v1_calculator_proto_rawDescGZIP() []byte {
	file_buf_v1_calculator_proto_rawDescOnce.Do(func() {
		file_buf_v1_calculator_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_v1_calculator_proto_rawDescData)
	})
	return file_buf_v1_calculator_proto_rawDescData
}

var file_buf_v1_calculator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_buf_v1_calculator_proto_goTypes = []interface{}{
	(*AddIntRequest)(nil),  // 0: buf.v1.AddIntRequest
	(*AddIntResponse)(nil), // 1: buf.v1.AddIntResponse
}
var file_buf_v1_calculator_proto_depIdxs = []int32{
	0, // 0: buf.v1.CalculatorService.AddInt:input_type -> buf.v1.AddIntRequest
	1, // 1: buf.v1.CalculatorService.AddInt:output_type -> buf.v1.AddIntResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_buf_v1_calculator_proto_init() }
func file_buf_v1_calculator_proto_init() {
	if File_buf_v1_calculator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_buf_v1_calculator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddIntRequest); i {
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
		file_buf_v1_calculator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddIntResponse); i {
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
			RawDescriptor: file_buf_v1_calculator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_buf_v1_calculator_proto_goTypes,
		DependencyIndexes: file_buf_v1_calculator_proto_depIdxs,
		MessageInfos:      file_buf_v1_calculator_proto_msgTypes,
	}.Build()
	File_buf_v1_calculator_proto = out.File
	file_buf_v1_calculator_proto_rawDesc = nil
	file_buf_v1_calculator_proto_goTypes = nil
	file_buf_v1_calculator_proto_depIdxs = nil
}
