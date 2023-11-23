//
// Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: genesis.proto

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

// GenesisHashRequest is the request for getting genesis block hash.
type GenesisHashRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GenesisHashRequest) Reset() {
	*x = GenesisHashRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenesisHashRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenesisHashRequest) ProtoMessage() {}

func (x *GenesisHashRequest) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenesisHashRequest.ProtoReflect.Descriptor instead.
func (*GenesisHashRequest) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{0}
}

// GenesisHashResponse is the response for getting genesis block hash.
type GenesisHashResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *GenesisHashResponse) Reset() {
	*x = GenesisHashResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenesisHashResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenesisHashResponse) ProtoMessage() {}

func (x *GenesisHashResponse) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenesisHashResponse.ProtoReflect.Descriptor instead.
func (*GenesisHashResponse) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{1}
}

func (x *GenesisHashResponse) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

var File_genesis_proto protoreflect.FileDescriptor

var file_genesis_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x22, 0x14, 0x0a, 0x12, 0x47,
	0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x48, 0x61, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x29, 0x0a, 0x13, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x48, 0x61, 0x73, 0x68,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x42, 0x10, 0x5a, 0x0e,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_genesis_proto_rawDescOnce sync.Once
	file_genesis_proto_rawDescData = file_genesis_proto_rawDesc
)

func file_genesis_proto_rawDescGZIP() []byte {
	file_genesis_proto_rawDescOnce.Do(func() {
		file_genesis_proto_rawDescData = protoimpl.X.CompressGZIP(file_genesis_proto_rawDescData)
	})
	return file_genesis_proto_rawDescData
}

var file_genesis_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_genesis_proto_goTypes = []interface{}{
	(*GenesisHashRequest)(nil),  // 0: blockchain.GenesisHashRequest
	(*GenesisHashResponse)(nil), // 1: blockchain.GenesisHashResponse
}
var file_genesis_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_genesis_proto_init() }
func file_genesis_proto_init() {
	if File_genesis_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_genesis_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenesisHashRequest); i {
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
		file_genesis_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenesisHashResponse); i {
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
			RawDescriptor: file_genesis_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_genesis_proto_goTypes,
		DependencyIndexes: file_genesis_proto_depIdxs,
		MessageInfos:      file_genesis_proto_msgTypes,
	}.Build()
	File_genesis_proto = out.File
	file_genesis_proto_rawDesc = nil
	file_genesis_proto_goTypes = nil
	file_genesis_proto_depIdxs = nil
}
