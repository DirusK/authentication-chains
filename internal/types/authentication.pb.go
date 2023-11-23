//
// Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: authentication.proto

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

// DeviceAuthenticationRequest is a request for authentication
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
		mi := &file_authentication_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceAuthenticationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceAuthenticationRequest) ProtoMessage() {}

func (x *DeviceAuthenticationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_proto_msgTypes[0]
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
	return file_authentication_proto_rawDescGZIP(), []int{0}
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

type DeviceAuthenticationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsAuthenticated bool `protobuf:"varint,1,opt,name=is_authenticated,json=isAuthenticated,proto3" json:"is_authenticated,omitempty"`
}

func (x *DeviceAuthenticationResponse) Reset() {
	*x = DeviceAuthenticationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceAuthenticationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceAuthenticationResponse) ProtoMessage() {}

func (x *DeviceAuthenticationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceAuthenticationResponse.ProtoReflect.Descriptor instead.
func (*DeviceAuthenticationResponse) Descriptor() ([]byte, []int) {
	return file_authentication_proto_rawDescGZIP(), []int{1}
}

func (x *DeviceAuthenticationResponse) GetIsAuthenticated() bool {
	if x != nil {
		return x.IsAuthenticated
	}
	return false
}

// AuthenticationEntry is a single record in authentication table
type AuthenticationEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId      []byte `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	ClusterHeadId []byte `protobuf:"bytes,2,opt,name=cluster_head_id,json=clusterHeadId,proto3" json:"cluster_head_id,omitempty"`
	BlockHash     []byte `protobuf:"bytes,3,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	BlockIndex    uint64 `protobuf:"varint,4,opt,name=block_index,json=blockIndex,proto3" json:"block_index,omitempty"`
}

func (x *AuthenticationEntry) Reset() {
	*x = AuthenticationEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticationEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticationEntry) ProtoMessage() {}

func (x *AuthenticationEntry) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticationEntry.ProtoReflect.Descriptor instead.
func (*AuthenticationEntry) Descriptor() ([]byte, []int) {
	return file_authentication_proto_rawDescGZIP(), []int{2}
}

func (x *AuthenticationEntry) GetDeviceId() []byte {
	if x != nil {
		return x.DeviceId
	}
	return nil
}

func (x *AuthenticationEntry) GetClusterHeadId() []byte {
	if x != nil {
		return x.ClusterHeadId
	}
	return nil
}

func (x *AuthenticationEntry) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *AuthenticationEntry) GetBlockIndex() uint64 {
	if x != nil {
		return x.BlockIndex
	}
	return 0
}

var File_authentication_proto protoreflect.FileDescriptor

var file_authentication_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x22, 0x80, 0x01, 0x0a, 0x1b, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12,
	0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x48, 0x65, 0x61, 0x64, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x49, 0x0a, 0x1c, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x41,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x69, 0x73, 0x5f, 0x61, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x69, 0x73, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64,
	0x22, 0x9a, 0x01, 0x0a, 0x13, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x48, 0x65, 0x61, 0x64, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1f, 0x0a, 0x0b,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x42, 0x10, 0x5a,
	0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_authentication_proto_rawDescOnce sync.Once
	file_authentication_proto_rawDescData = file_authentication_proto_rawDesc
)

func file_authentication_proto_rawDescGZIP() []byte {
	file_authentication_proto_rawDescOnce.Do(func() {
		file_authentication_proto_rawDescData = protoimpl.X.CompressGZIP(file_authentication_proto_rawDescData)
	})
	return file_authentication_proto_rawDescData
}

var file_authentication_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_authentication_proto_goTypes = []interface{}{
	(*DeviceAuthenticationRequest)(nil),  // 0: blockchain.DeviceAuthenticationRequest
	(*DeviceAuthenticationResponse)(nil), // 1: blockchain.DeviceAuthenticationResponse
	(*AuthenticationEntry)(nil),          // 2: blockchain.AuthenticationEntry
}
var file_authentication_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_authentication_proto_init() }
func file_authentication_proto_init() {
	if File_authentication_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_authentication_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_authentication_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceAuthenticationResponse); i {
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
		file_authentication_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticationEntry); i {
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
			RawDescriptor: file_authentication_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_authentication_proto_goTypes,
		DependencyIndexes: file_authentication_proto_depIdxs,
		MessageInfos:      file_authentication_proto_msgTypes,
	}.Build()
	File_authentication_proto = out.File
	file_authentication_proto_rawDesc = nil
	file_authentication_proto_goTypes = nil
	file_authentication_proto_depIdxs = nil
}