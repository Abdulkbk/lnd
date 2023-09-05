// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: verrpc/verrpc.proto

package verrpc

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

type VersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *VersionRequest) Reset() {
	*x = VersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verrpc_verrpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionRequest) ProtoMessage() {}

func (x *VersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_verrpc_verrpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionRequest.ProtoReflect.Descriptor instead.
func (*VersionRequest) Descriptor() ([]byte, []int) {
	return file_verrpc_verrpc_proto_rawDescGZIP(), []int{0}
}

type Version struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A verbose description of the daemon's commit.
	Commit string `protobuf:"bytes,1,opt,name=commit,proto3" json:"commit,omitempty"`
	// The SHA1 commit hash that the daemon is compiled with.
	CommitHash string `protobuf:"bytes,2,opt,name=commit_hash,json=commitHash,proto3" json:"commit_hash,omitempty"`
	// The semantic version.
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	// The major application version.
	AppMajor uint32 `protobuf:"varint,4,opt,name=app_major,json=appMajor,proto3" json:"app_major,omitempty"`
	// The minor application version.
	AppMinor uint32 `protobuf:"varint,5,opt,name=app_minor,json=appMinor,proto3" json:"app_minor,omitempty"`
	// The application patch number.
	AppPatch uint32 `protobuf:"varint,6,opt,name=app_patch,json=appPatch,proto3" json:"app_patch,omitempty"`
	// The application pre-release modifier, possibly empty.
	AppPreRelease string `protobuf:"bytes,7,opt,name=app_pre_release,json=appPreRelease,proto3" json:"app_pre_release,omitempty"`
	// The list of build tags that were supplied during compilation.
	BuildTags []string `protobuf:"bytes,8,rep,name=build_tags,json=buildTags,proto3" json:"build_tags,omitempty"`
	// The version of go that compiled the executable.
	GoVersion string `protobuf:"bytes,9,opt,name=go_version,json=goVersion,proto3" json:"go_version,omitempty"`
}

func (x *Version) Reset() {
	*x = Version{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verrpc_verrpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Version) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Version) ProtoMessage() {}

func (x *Version) ProtoReflect() protoreflect.Message {
	mi := &file_verrpc_verrpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Version.ProtoReflect.Descriptor instead.
func (*Version) Descriptor() ([]byte, []int) {
	return file_verrpc_verrpc_proto_rawDescGZIP(), []int{1}
}

func (x *Version) GetCommit() string {
	if x != nil {
		return x.Commit
	}
	return ""
}

func (x *Version) GetCommitHash() string {
	if x != nil {
		return x.CommitHash
	}
	return ""
}

func (x *Version) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Version) GetAppMajor() uint32 {
	if x != nil {
		return x.AppMajor
	}
	return 0
}

func (x *Version) GetAppMinor() uint32 {
	if x != nil {
		return x.AppMinor
	}
	return 0
}

func (x *Version) GetAppPatch() uint32 {
	if x != nil {
		return x.AppPatch
	}
	return 0
}

func (x *Version) GetAppPreRelease() string {
	if x != nil {
		return x.AppPreRelease
	}
	return ""
}

func (x *Version) GetBuildTags() []string {
	if x != nil {
		return x.BuildTags
	}
	return nil
}

func (x *Version) GetGoVersion() string {
	if x != nil {
		return x.GoVersion
	}
	return ""
}

var File_verrpc_verrpc_proto protoreflect.FileDescriptor

var file_verrpc_verrpc_proto_rawDesc = []byte{
	0x0a, 0x13, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x22, 0x10, 0x0a,
	0x0e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x99, 0x02, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x63,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1b,
	0x0a, 0x09, 0x61, 0x70, 0x70, 0x5f, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x08, 0x61, 0x70, 0x70, 0x4d, 0x61, 0x6a, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x61,
	0x70, 0x70, 0x5f, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08,
	0x61, 0x70, 0x70, 0x4d, 0x69, 0x6e, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x70, 0x70, 0x5f,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x61, 0x70, 0x70,
	0x50, 0x61, 0x74, 0x63, 0x68, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x5f, 0x70, 0x72, 0x65,
	0x5f, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x61, 0x70, 0x70, 0x50, 0x72, 0x65, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x74, 0x61, 0x67, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x09, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x67, 0x73, 0x12, 0x1d, 0x0a, 0x0a,
	0x67, 0x6f, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x67, 0x6f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x32, 0x42, 0x0a, 0x09, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f,
	0x2e, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42,
	0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69,
	0x67, 0x68, 0x74, 0x6e, 0x69, 0x6e, 0x67, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x6c,
	0x6e, 0x64, 0x2f, 0x6c, 0x6e, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_verrpc_verrpc_proto_rawDescOnce sync.Once
	file_verrpc_verrpc_proto_rawDescData = file_verrpc_verrpc_proto_rawDesc
)

func file_verrpc_verrpc_proto_rawDescGZIP() []byte {
	file_verrpc_verrpc_proto_rawDescOnce.Do(func() {
		file_verrpc_verrpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_verrpc_verrpc_proto_rawDescData)
	})
	return file_verrpc_verrpc_proto_rawDescData
}

var file_verrpc_verrpc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_verrpc_verrpc_proto_goTypes = []interface{}{
	(*VersionRequest)(nil), // 0: verrpc.VersionRequest
	(*Version)(nil),        // 1: verrpc.Version
}
var file_verrpc_verrpc_proto_depIdxs = []int32{
	0, // 0: verrpc.Versioner.GetVersion:input_type -> verrpc.VersionRequest
	1, // 1: verrpc.Versioner.GetVersion:output_type -> verrpc.Version
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_verrpc_verrpc_proto_init() }
func file_verrpc_verrpc_proto_init() {
	if File_verrpc_verrpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_verrpc_verrpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionRequest); i {
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
		file_verrpc_verrpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Version); i {
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
			RawDescriptor: file_verrpc_verrpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_verrpc_verrpc_proto_goTypes,
		DependencyIndexes: file_verrpc_verrpc_proto_depIdxs,
		MessageInfos:      file_verrpc_verrpc_proto_msgTypes,
	}.Build()
	File_verrpc_verrpc_proto = out.File
	file_verrpc_verrpc_proto_rawDesc = nil
	file_verrpc_verrpc_proto_goTypes = nil
	file_verrpc_verrpc_proto_depIdxs = nil
}
