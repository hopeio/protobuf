// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: hopeio/oauth/oauth.proto

package oauth

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

type OauthReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResponseType   string `protobuf:"bytes,1,opt,name=responseType,proto3" json:"responseType,omitempty"`
	ClientID       string `protobuf:"bytes,2,opt,name=clientID,proto3" json:"clientID,omitempty"`
	Scope          string `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	RedirectURI    string `protobuf:"bytes,4,opt,name=redirectURI,proto3" json:"redirectURI,omitempty"`
	State          string `protobuf:"bytes,5,opt,name=state,proto3" json:"state,omitempty"`
	UserID         string `protobuf:"bytes,6,opt,name=userID,proto3" json:"userID,omitempty"`
	AccessTokenExp int64  `protobuf:"varint,7,opt,name=accessTokenExp,proto3" json:"accessTokenExp,omitempty"`
	ClientSecret   string `protobuf:"bytes,11,opt,name=clientSecret,proto3" json:"clientSecret,omitempty"`
	Code           string `protobuf:"bytes,12,opt,name=code,proto3" json:"code,omitempty"`
	RefreshToken   string `protobuf:"bytes,13,opt,name=refreshToken,proto3" json:"refreshToken,omitempty"`
	GrantType      string `protobuf:"bytes,14,opt,name=grantType,proto3" json:"grantType,omitempty"`
	AccessType     string `protobuf:"bytes,15,opt,name=accessType,proto3" json:"accessType,omitempty"`
	LoginURI       string `protobuf:"bytes,16,opt,name=loginURI,proto3" json:"loginURI,omitempty"`
}

func (x *OauthReq) Reset() {
	*x = OauthReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hopeio_oauth_oauth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OauthReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OauthReq) ProtoMessage() {}

func (x *OauthReq) ProtoReflect() protoreflect.Message {
	mi := &file_hopeio_oauth_oauth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OauthReq.ProtoReflect.Descriptor instead.
func (*OauthReq) Descriptor() ([]byte, []int) {
	return file_hopeio_oauth_oauth_proto_rawDescGZIP(), []int{0}
}

func (x *OauthReq) GetResponseType() string {
	if x != nil {
		return x.ResponseType
	}
	return ""
}

func (x *OauthReq) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *OauthReq) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *OauthReq) GetRedirectURI() string {
	if x != nil {
		return x.RedirectURI
	}
	return ""
}

func (x *OauthReq) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *OauthReq) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *OauthReq) GetAccessTokenExp() int64 {
	if x != nil {
		return x.AccessTokenExp
	}
	return 0
}

func (x *OauthReq) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *OauthReq) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *OauthReq) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *OauthReq) GetGrantType() string {
	if x != nil {
		return x.GrantType
	}
	return ""
}

func (x *OauthReq) GetAccessType() string {
	if x != nil {
		return x.AccessType
	}
	return ""
}

func (x *OauthReq) GetLoginURI() string {
	if x != nil {
		return x.LoginURI
	}
	return ""
}

type Client struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID     string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Secret string `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	Domain string `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	UserID string `protobuf:"bytes,4,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *Client) Reset() {
	*x = Client{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hopeio_oauth_oauth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Client) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Client) ProtoMessage() {}

func (x *Client) ProtoReflect() protoreflect.Message {
	mi := &file_hopeio_oauth_oauth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Client.ProtoReflect.Descriptor instead.
func (*Client) Descriptor() ([]byte, []int) {
	return file_hopeio_oauth_oauth_proto_rawDescGZIP(), []int{1}
}

func (x *Client) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Client) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

func (x *Client) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Client) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

var File_hopeio_oauth_oauth_proto protoreflect.FileDescriptor

var file_hopeio_oauth_oauth_proto_rawDesc = []byte{
	0x0a, 0x18, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x6f,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6f, 0x61, 0x75, 0x74,
	0x68, 0x22, 0x8e, 0x03, 0x0a, 0x08, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x12, 0x22,
	0x0a, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x55, 0x52, 0x49, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x55, 0x52, 0x49, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x12, 0x26, 0x0a, 0x0e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x45, 0x78, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x45, 0x78, 0x70, 0x12, 0x22, 0x0a, 0x0c,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x67, 0x72, 0x61, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x61,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x55,
	0x52, 0x49, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x55,
	0x52, 0x49, 0x22, 0x60, 0x0a, 0x06, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x42, 0x3e, 0x0a, 0x18, 0x78, 0x79, 0x7a, 0x2e, 0x68, 0x6f, 0x70, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x6f, 0x61, 0x75, 0x74, 0x68,
	0x50, 0x01, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68,
	0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6f,
	0x61, 0x75, 0x74, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hopeio_oauth_oauth_proto_rawDescOnce sync.Once
	file_hopeio_oauth_oauth_proto_rawDescData = file_hopeio_oauth_oauth_proto_rawDesc
)

func file_hopeio_oauth_oauth_proto_rawDescGZIP() []byte {
	file_hopeio_oauth_oauth_proto_rawDescOnce.Do(func() {
		file_hopeio_oauth_oauth_proto_rawDescData = protoimpl.X.CompressGZIP(file_hopeio_oauth_oauth_proto_rawDescData)
	})
	return file_hopeio_oauth_oauth_proto_rawDescData
}

var file_hopeio_oauth_oauth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_hopeio_oauth_oauth_proto_goTypes = []interface{}{
	(*OauthReq)(nil), // 0: oauth.OauthReq
	(*Client)(nil),   // 1: oauth.Client
}
var file_hopeio_oauth_oauth_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_hopeio_oauth_oauth_proto_init() }
func file_hopeio_oauth_oauth_proto_init() {
	if File_hopeio_oauth_oauth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hopeio_oauth_oauth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OauthReq); i {
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
		file_hopeio_oauth_oauth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Client); i {
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
			RawDescriptor: file_hopeio_oauth_oauth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_hopeio_oauth_oauth_proto_goTypes,
		DependencyIndexes: file_hopeio_oauth_oauth_proto_depIdxs,
		MessageInfos:      file_hopeio_oauth_oauth_proto_msgTypes,
	}.Build()
	File_hopeio_oauth_oauth_proto = out.File
	file_hopeio_oauth_oauth_proto_rawDesc = nil
	file_hopeio_oauth_oauth_proto_goTypes = nil
	file_hopeio_oauth_oauth_proto_depIdxs = nil
}
