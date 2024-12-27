//
// Copyright 2024 hopeio. All rights reserved.
// Licensed under the MIT License that can be found in the LICENSE file.
// @Created by jyb

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: hopeio/model/enum/enum.proto

package enum

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	deletedAt "github.com/hopeio/protobuf/time/deletedAt"
	timestamp "github.com/hopeio/protobuf/time/timestamp"
	_ "github.com/hopeio/protobuf/utils/patch"
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

type Enum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" gorm:"primaryKey"`
	Name      string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Group     uint64               `protobuf:"varint,3,opt,name=group,proto3" json:"group,omitempty" gorm:"size:20"`
	Type      uint32               `protobuf:"varint,6,opt,name=type,proto3" json:"type,omitempty" comment:"类型"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,16,opt,name=createdAt,proto3" json:"createdAt,omitempty" gorm:"type:timestamptz(6);default:now();index"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,26,opt,name=updatedAt,proto3" json:"updatedAt,omitempty" gorm:"type:timestamptz(6)"`
	DeletedAt *deletedAt.DeletedAt `protobuf:"bytes,28,opt,name=deletedAt,proto3" json:"deletedAt,omitempty" gorm:"<-:false;type:timestamptz(6);index"`
	Status    uint32               `protobuf:"varint,18,opt,name=status,proto3" json:"status,omitempty" gorm:"type:int2;default:0"`
}

func (x *Enum) Reset() {
	*x = Enum{}
	mi := &file_hopeio_model_enum_enum_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Enum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Enum) ProtoMessage() {}

func (x *Enum) ProtoReflect() protoreflect.Message {
	mi := &file_hopeio_model_enum_enum_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Enum.ProtoReflect.Descriptor instead.
func (*Enum) Descriptor() ([]byte, []int) {
	return file_hopeio_model_enum_enum_proto_rawDescGZIP(), []int{0}
}

func (x *Enum) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Enum) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Enum) GetGroup() uint64 {
	if x != nil {
		return x.Group
	}
	return 0
}

func (x *Enum) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Enum) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Enum) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Enum) GetDeletedAt() *deletedAt.DeletedAt {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *Enum) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type EnumValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" gorm:"primaryKey"`
	EnumId    uint64               `protobuf:"varint,2,opt,name=enumId,proto3" json:"enumId,omitempty"`
	Index     uint64               `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty" gorm:"comment:index"`
	Value     string               `protobuf:"bytes,5,opt,name=value,proto3" json:"value,omitempty" comment:"值"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,16,opt,name=createdAt,proto3" json:"createdAt,omitempty" gorm:"type:timestamptz(6);default:now();index"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,26,opt,name=updatedAt,proto3" json:"updatedAt,omitempty" gorm:"type:timestamptz(6)"`
	DeletedAt *deletedAt.DeletedAt `protobuf:"bytes,28,opt,name=deletedAt,proto3" json:"deletedAt,omitempty" gorm:"<-:false;type:timestamptz(6);index"`
	Status    uint32               `protobuf:"varint,18,opt,name=status,proto3" json:"status,omitempty" gorm:"type:int2;default:0"`
}

func (x *EnumValue) Reset() {
	*x = EnumValue{}
	mi := &file_hopeio_model_enum_enum_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnumValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumValue) ProtoMessage() {}

func (x *EnumValue) ProtoReflect() protoreflect.Message {
	mi := &file_hopeio_model_enum_enum_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumValue.ProtoReflect.Descriptor instead.
func (*EnumValue) Descriptor() ([]byte, []int) {
	return file_hopeio_model_enum_enum_proto_rawDescGZIP(), []int{1}
}

func (x *EnumValue) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *EnumValue) GetEnumId() uint64 {
	if x != nil {
		return x.EnumId
	}
	return 0
}

func (x *EnumValue) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *EnumValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *EnumValue) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *EnumValue) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *EnumValue) GetDeletedAt() *deletedAt.DeletedAt {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *EnumValue) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_hopeio_model_enum_enum_proto protoreflect.FileDescriptor

var file_hopeio_model_enum_enum_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x65,
	0x6e, 0x75, 0x6d, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04,
	0x65, 0x6e, 0x75, 0x6d, 0x1a, 0x1b, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x75, 0x74, 0x69,
	0x6c, 0x73, 0x2f, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x25, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2f,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e,
	0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x9e, 0x04, 0x0a, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x28, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x42, 0x18, 0xd2, 0xb5, 0x03, 0x14, 0xa2, 0x01, 0x11, 0x67, 0x6f, 0x72,
	0x6d, 0x3a, 0x22, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x22, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x2a, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x16, 0xd2, 0xb5, 0x03, 0x12, 0xa2, 0x01, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x3a, 0x22, 0xe5, 0x90, 0x8d, 0xe7, 0xa7, 0xb0, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b,
	0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x15, 0xd2,
	0xb5, 0x03, 0x11, 0xa2, 0x01, 0x0e, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x69, 0x7a, 0x65,
	0x3a, 0x32, 0x30, 0x22, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x2b, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x17, 0xd2, 0xb5, 0x03, 0x13, 0xa2,
	0x01, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe7, 0xb1, 0xbb, 0xe5, 0x9e,
	0x8b, 0x22, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x69, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x42, 0x35, 0xd2, 0xb5, 0x03, 0x31, 0xa2, 0x01, 0x2e, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22,
	0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a,
	0x28, 0x36, 0x29, 0x3b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x3a, 0x6e, 0x6f, 0x77, 0x28,
	0x29, 0x3b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x55, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x18, 0x1a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x21, 0xd2, 0xb5,
	0x03, 0x1d, 0xa2, 0x01, 0x1a, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x74, 0x79, 0x70, 0x65, 0x3a,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a, 0x28, 0x36, 0x29, 0x22, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x64, 0x0a, 0x09, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x42, 0x30, 0xd2, 0xb5, 0x03, 0x2c, 0xa2, 0x01, 0x29, 0x67, 0x6f, 0x72, 0x6d,
	0x3a, 0x22, 0x3c, 0x2d, 0x3a, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a, 0x28, 0x36, 0x29, 0x3b, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x22, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x3e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0d,
	0x42, 0x26, 0x92, 0x41, 0x02, 0x40, 0x01, 0xd2, 0xb5, 0x03, 0x1d, 0xa2, 0x01, 0x1a, 0x67, 0x6f,
	0x72, 0x6d, 0x3a, 0x22, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x69, 0x6e, 0x74, 0x32, 0x3b, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x3a, 0x30, 0x22, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x22, 0xae, 0x04, 0x0a, 0x09, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x28,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x18, 0xd2, 0xb5, 0x03, 0x14,
	0xa2, 0x01, 0x11, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79,
	0x4b, 0x65, 0x79, 0x22, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x06, 0x65, 0x6e, 0x75, 0x6d,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x42, 0x18, 0xd2, 0xb5, 0x03, 0x14, 0xa2, 0x01,
	0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe6, 0x9e, 0x9a, 0xe4, 0xb8, 0xbe,
	0x69, 0x64, 0x52, 0x06, 0x65, 0x6e, 0x75, 0x6d, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x1b, 0xd2, 0xb5, 0x03, 0x17, 0xa2,
	0x01, 0x14, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x2a, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x14, 0xd2, 0xb5,
	0x03, 0x10, 0xa2, 0x01, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe5, 0x80,
	0xbc, 0x22, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x69, 0x0a, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x42, 0x35, 0xd2, 0xb5, 0x03, 0x31, 0xa2, 0x01, 0x2e, 0x67, 0x6f, 0x72, 0x6d, 0x3a,
	0x22, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74,
	0x7a, 0x28, 0x36, 0x29, 0x3b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x3a, 0x6e, 0x6f, 0x77,
	0x28, 0x29, 0x3b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x55, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x21, 0xd2,
	0xb5, 0x03, 0x1d, 0xa2, 0x01, 0x1a, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x74, 0x79, 0x70, 0x65,
	0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a, 0x28, 0x36, 0x29, 0x22,
	0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x64, 0x0a, 0x09, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x42, 0x30, 0xd2, 0xb5, 0x03, 0x2c, 0xa2, 0x01, 0x29, 0x67, 0x6f, 0x72,
	0x6d, 0x3a, 0x22, 0x3c, 0x2d, 0x3a, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x3b, 0x74, 0x79, 0x70, 0x65,
	0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a, 0x28, 0x36, 0x29, 0x3b,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x3e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28,
	0x0d, 0x42, 0x26, 0x92, 0x41, 0x02, 0x40, 0x01, 0xd2, 0xb5, 0x03, 0x1d, 0xa2, 0x01, 0x1a, 0x67,
	0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x69, 0x6e, 0x74, 0x32, 0x3b, 0x64,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x3a, 0x30, 0x22, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x42, 0x48, 0x0a, 0x1d, 0x78, 0x79, 0x7a, 0x2e, 0x68, 0x6f, 0x70, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x65, 0x6e,
	0x75, 0x6d, 0x50, 0x01, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_hopeio_model_enum_enum_proto_rawDescOnce sync.Once
	file_hopeio_model_enum_enum_proto_rawDescData = file_hopeio_model_enum_enum_proto_rawDesc
)

func file_hopeio_model_enum_enum_proto_rawDescGZIP() []byte {
	file_hopeio_model_enum_enum_proto_rawDescOnce.Do(func() {
		file_hopeio_model_enum_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_hopeio_model_enum_enum_proto_rawDescData)
	})
	return file_hopeio_model_enum_enum_proto_rawDescData
}

var file_hopeio_model_enum_enum_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_hopeio_model_enum_enum_proto_goTypes = []any{
	(*Enum)(nil),                // 0: enum.Enum
	(*EnumValue)(nil),           // 1: enum.EnumValue
	(*timestamp.Timestamp)(nil), // 2: timestamp.Timestamp
	(*deletedAt.DeletedAt)(nil), // 3: deletedAt.DeletedAt
}
var file_hopeio_model_enum_enum_proto_depIdxs = []int32{
	2, // 0: enum.Enum.createdAt:type_name -> timestamp.Timestamp
	2, // 1: enum.Enum.updatedAt:type_name -> timestamp.Timestamp
	3, // 2: enum.Enum.deletedAt:type_name -> deletedAt.DeletedAt
	2, // 3: enum.EnumValue.createdAt:type_name -> timestamp.Timestamp
	2, // 4: enum.EnumValue.updatedAt:type_name -> timestamp.Timestamp
	3, // 5: enum.EnumValue.deletedAt:type_name -> deletedAt.DeletedAt
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_hopeio_model_enum_enum_proto_init() }
func file_hopeio_model_enum_enum_proto_init() {
	if File_hopeio_model_enum_enum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hopeio_model_enum_enum_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_hopeio_model_enum_enum_proto_goTypes,
		DependencyIndexes: file_hopeio_model_enum_enum_proto_depIdxs,
		MessageInfos:      file_hopeio_model_enum_enum_proto_msgTypes,
	}.Build()
	File_hopeio_model_enum_enum_proto = out.File
	file_hopeio_model_enum_enum_proto_rawDesc = nil
	file_hopeio_model_enum_enum_proto_goTypes = nil
	file_hopeio_model_enum_enum_proto_depIdxs = nil
}
