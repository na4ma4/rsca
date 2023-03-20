// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.2
// source: github.com/na4ma4/rsca/api/admin.proto

package api

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

type RemoveHostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
}

func (x *RemoveHostRequest) Reset() {
	*x = RemoveHostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_na4ma4_rsca_api_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveHostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveHostRequest) ProtoMessage() {}

func (x *RemoveHostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_na4ma4_rsca_api_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveHostRequest.ProtoReflect.Descriptor instead.
func (*RemoveHostRequest) Descriptor() ([]byte, []int) {
	return file_github_com_na4ma4_rsca_api_admin_proto_rawDescGZIP(), []int{0}
}

func (x *RemoveHostRequest) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

type RemoveHostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
}

func (x *RemoveHostResponse) Reset() {
	*x = RemoveHostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_na4ma4_rsca_api_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveHostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveHostResponse) ProtoMessage() {}

func (x *RemoveHostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_na4ma4_rsca_api_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveHostResponse.ProtoReflect.Descriptor instead.
func (*RemoveHostResponse) Descriptor() ([]byte, []int) {
	return file_github_com_na4ma4_rsca_api_admin_proto_rawDescGZIP(), []int{1}
}

func (x *RemoveHostResponse) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

var File_github_com_na4ma4_rsca_api_admin_proto protoreflect.FileDescriptor

var file_github_com_na4ma4_rsca_api_admin_proto_rawDesc = []byte{
	0x0a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x34,
	0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61,
	0x70, 0x69, 0x1a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e,
	0x61, 0x34, 0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x29, 0x0a, 0x11, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x2a, 0x0a, 0x12, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x32, 0x82, 0x02, 0x0a, 0x05, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x30, 0x0a, 0x09,
	0x4c, 0x69, 0x73, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x0f, 0x2e, 0x72, 0x73, 0x63, 0x61,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x72, 0x73, 0x63,
	0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x30, 0x01, 0x12, 0x47,
	0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x1b, 0x2e, 0x72,
	0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x48, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x72, 0x73, 0x63, 0x61,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x41, 0x6c, 0x6c, 0x12, 0x11, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x1a, 0x1c, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x41, 0x6c, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0b, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x11, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x1a, 0x1d, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x34, 0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63,
	0x61, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_na4ma4_rsca_api_admin_proto_rawDescOnce sync.Once
	file_github_com_na4ma4_rsca_api_admin_proto_rawDescData = file_github_com_na4ma4_rsca_api_admin_proto_rawDesc
)

func file_github_com_na4ma4_rsca_api_admin_proto_rawDescGZIP() []byte {
	file_github_com_na4ma4_rsca_api_admin_proto_rawDescOnce.Do(func() {
		file_github_com_na4ma4_rsca_api_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_na4ma4_rsca_api_admin_proto_rawDescData)
	})
	return file_github_com_na4ma4_rsca_api_admin_proto_rawDescData
}

var file_github_com_na4ma4_rsca_api_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_na4ma4_rsca_api_admin_proto_goTypes = []interface{}{
	(*RemoveHostRequest)(nil),   // 0: rsca.api.RemoveHostRequest
	(*RemoveHostResponse)(nil),  // 1: rsca.api.RemoveHostResponse
	(*Empty)(nil),               // 2: rsca.api.Empty
	(*Members)(nil),             // 3: rsca.api.Members
	(*Member)(nil),              // 4: rsca.api.Member
	(*TriggerAllResponse)(nil),  // 5: rsca.api.TriggerAllResponse
	(*TriggerInfoResponse)(nil), // 6: rsca.api.TriggerInfoResponse
}
var file_github_com_na4ma4_rsca_api_admin_proto_depIdxs = []int32{
	2, // 0: rsca.api.Admin.ListHosts:input_type -> rsca.api.Empty
	0, // 1: rsca.api.Admin.RemoveHost:input_type -> rsca.api.RemoveHostRequest
	3, // 2: rsca.api.Admin.TriggerAll:input_type -> rsca.api.Members
	3, // 3: rsca.api.Admin.TriggerInfo:input_type -> rsca.api.Members
	4, // 4: rsca.api.Admin.ListHosts:output_type -> rsca.api.Member
	1, // 5: rsca.api.Admin.RemoveHost:output_type -> rsca.api.RemoveHostResponse
	5, // 6: rsca.api.Admin.TriggerAll:output_type -> rsca.api.TriggerAllResponse
	6, // 7: rsca.api.Admin.TriggerInfo:output_type -> rsca.api.TriggerInfoResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_na4ma4_rsca_api_admin_proto_init() }
func file_github_com_na4ma4_rsca_api_admin_proto_init() {
	if File_github_com_na4ma4_rsca_api_admin_proto != nil {
		return
	}
	file_github_com_na4ma4_rsca_api_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_github_com_na4ma4_rsca_api_admin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveHostRequest); i {
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
		file_github_com_na4ma4_rsca_api_admin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveHostResponse); i {
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
			RawDescriptor: file_github_com_na4ma4_rsca_api_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_na4ma4_rsca_api_admin_proto_goTypes,
		DependencyIndexes: file_github_com_na4ma4_rsca_api_admin_proto_depIdxs,
		MessageInfos:      file_github_com_na4ma4_rsca_api_admin_proto_msgTypes,
	}.Build()
	File_github_com_na4ma4_rsca_api_admin_proto = out.File
	file_github_com_na4ma4_rsca_api_admin_proto_rawDesc = nil
	file_github_com_na4ma4_rsca_api_admin_proto_goTypes = nil
	file_github_com_na4ma4_rsca_api_admin_proto_depIdxs = nil
}
