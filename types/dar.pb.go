// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: dar.proto

package types

import (
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

type DeviceAuthenticationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId      []byte `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	ClusterHeadId []byte `protobuf:"bytes,2,opt,name=cluster_head_id,json=clusterHeadId,proto3" json:"cluster_head_id,omitempty"`
	Signature     []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *DeviceAuthenticationRequest) Reset() {
	*x = DeviceAuthenticationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dar_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceAuthenticationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceAuthenticationRequest) ProtoMessage() {}

func (x *DeviceAuthenticationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dar_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceAuthenticationRequest.ProtoReflect.Descriptor instead.
func (*DeviceAuthenticationRequest) Descriptor() ([]byte, []int) {
	return file_dar_proto_rawDescGZIP(), []int{0}
}

func (x *DeviceAuthenticationRequest) GetDeviceId() []byte {
	if x != nil {
		return x.DeviceId
	}
	return nil
}

func (x *DeviceAuthenticationRequest) GetClusterHeadId() []byte {
	if x != nil {
		return x.ClusterHeadId
	}
	return nil
}

func (x *DeviceAuthenticationRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

var File_dar_proto protoreflect.FileDescriptor

var file_dar_proto_rawDesc = []byte{
	0x0a, 0x09, 0x64, 0x61, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x22, 0x80, 0x01, 0x0a, 0x1b, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f,
	0x68, 0x65, 0x61, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x48, 0x65, 0x61, 0x64, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dar_proto_rawDescOnce sync.Once
	file_dar_proto_rawDescData = file_dar_proto_rawDesc
)

func file_dar_proto_rawDescGZIP() []byte {
	file_dar_proto_rawDescOnce.Do(func() {
		file_dar_proto_rawDescData = protoimpl.X.CompressGZIP(file_dar_proto_rawDescData)
	})
	return file_dar_proto_rawDescData
}

var file_dar_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_dar_proto_goTypes = []interface{}{
	(*DeviceAuthenticationRequest)(nil), // 0: blockchain.DeviceAuthenticationRequest
}
var file_dar_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_dar_proto_init() }
func file_dar_proto_init() {
	if File_dar_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dar_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceAuthenticationRequest); i {
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
			RawDescriptor: file_dar_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dar_proto_goTypes,
		DependencyIndexes: file_dar_proto_depIdxs,
		MessageInfos:      file_dar_proto_msgTypes,
	}.Build()
	File_dar_proto = out.File
	file_dar_proto_rawDesc = nil
	file_dar_proto_goTypes = nil
	file_dar_proto_depIdxs = nil
}
