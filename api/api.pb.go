// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: github.com/na4ma4/rsca/api/api.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_github_com_na4ma4_rsca_api_api_proto protoreflect.FileDescriptor

var file_github_com_na4ma4_rsca_api_api_proto_rawDesc = []byte{
	0x0a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x34,
	0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69,
	0x1a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x34,
	0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x38, 0x0a, 0x04, 0x52, 0x53, 0x43,
	0x41, 0x12, 0x30, 0x0a, 0x04, 0x50, 0x69, 0x70, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x73, 0x63, 0x61,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x11, 0x2e, 0x72,
	0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x28,
	0x01, 0x30, 0x01, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6e, 0x61, 0x34, 0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_github_com_na4ma4_rsca_api_api_proto_goTypes = []any{
	(*Message)(nil), // 0: rsca.api.Message
}
var file_github_com_na4ma4_rsca_api_api_proto_depIdxs = []int32{
	0, // 0: rsca.api.RSCA.Pipe:input_type -> rsca.api.Message
	0, // 1: rsca.api.RSCA.Pipe:output_type -> rsca.api.Message
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_na4ma4_rsca_api_api_proto_init() }
func file_github_com_na4ma4_rsca_api_api_proto_init() {
	if File_github_com_na4ma4_rsca_api_api_proto != nil {
		return
	}
	file_github_com_na4ma4_rsca_api_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_na4ma4_rsca_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_na4ma4_rsca_api_api_proto_goTypes,
		DependencyIndexes: file_github_com_na4ma4_rsca_api_api_proto_depIdxs,
	}.Build()
	File_github_com_na4ma4_rsca_api_api_proto = out.File
	file_github_com_na4ma4_rsca_api_api_proto_rawDesc = nil
	file_github_com_na4ma4_rsca_api_api_proto_goTypes = nil
	file_github_com_na4ma4_rsca_api_api_proto_depIdxs = nil
}
