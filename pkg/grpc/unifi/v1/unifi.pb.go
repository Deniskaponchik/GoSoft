// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: unifi.proto

//определяет просто пространтсво имён, чтобы не было конфликтов
//package unifi_v1;

package v1

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

type ClientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
}

func (x *ClientRequest) Reset() {
	*x = ClientRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unifi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientRequest) ProtoMessage() {}

func (x *ClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_unifi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientRequest.ProtoReflect.Descriptor instead.
func (*ClientRequest) Descriptor() ([]byte, []int) {
	return file_unifi_proto_rawDescGZIP(), []int{0}
}

func (x *ClientRequest) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

type ClientResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Error    string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	// ClientInfo info = 2;
	Anomalies []*Anomaly `protobuf:"bytes,3,rep,name=anomalies,proto3" json:"anomalies,omitempty"`
}

func (x *ClientResponse) Reset() {
	*x = ClientResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unifi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientResponse) ProtoMessage() {}

func (x *ClientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_unifi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientResponse.ProtoReflect.Descriptor instead.
func (*ClientResponse) Descriptor() ([]byte, []int) {
	return file_unifi_proto_rawDescGZIP(), []int{1}
}

func (x *ClientResponse) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *ClientResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *ClientResponse) GetAnomalies() []*Anomaly {
	if x != nil {
		return x.Anomalies
	}
	return nil
}

type Anomaly struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApName   string `protobuf:"bytes,1,opt,name=apName,proto3" json:"apName,omitempty"`
	DateHour string `protobuf:"bytes,2,opt,name=dateHour,proto3" json:"dateHour,omitempty"`
	// google.protobuf.Timestamp date_hour = 2;
	// repeated AnomalyString sliceAnomStr = 3;
	AnomStr []string `protobuf:"bytes,3,rep,name=anomStr,proto3" json:"anomStr,omitempty"`
}

func (x *Anomaly) Reset() {
	*x = Anomaly{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unifi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Anomaly) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Anomaly) ProtoMessage() {}

func (x *Anomaly) ProtoReflect() protoreflect.Message {
	mi := &file_unifi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Anomaly.ProtoReflect.Descriptor instead.
func (*Anomaly) Descriptor() ([]byte, []int) {
	return file_unifi_proto_rawDescGZIP(), []int{2}
}

func (x *Anomaly) GetApName() string {
	if x != nil {
		return x.ApName
	}
	return ""
}

func (x *Anomaly) GetDateHour() string {
	if x != nil {
		return x.DateHour
	}
	return ""
}

func (x *Anomaly) GetAnomStr() []string {
	if x != nil {
		return x.AnomStr
	}
	return nil
}

var File_unifi_proto protoreflect.FileDescriptor

var file_unifi_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x61,
	0x70, 0x69, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x5f, 0x76, 0x31, 0x22, 0x2b, 0x0a, 0x0d, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x77, 0x0a, 0x0e, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f,
	0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f,
	0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x33, 0x0a, 0x09,
	0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x69, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x5f, 0x76, 0x31, 0x2e, 0x41,
	0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x52, 0x09, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x69, 0x65,
	0x73, 0x22, 0x57, 0x0a, 0x07, 0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x65, 0x48, 0x6f, 0x75, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x65, 0x48, 0x6f, 0x75, 0x72,
	0x12, 0x18, 0x0a, 0x07, 0x61, 0x6e, 0x6f, 0x6d, 0x53, 0x74, 0x72, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x6e, 0x6f, 0x6d, 0x53, 0x74, 0x72, 0x32, 0x56, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x69, 0x65, 0x73, 0x12, 0x46, 0x0a, 0x09, 0x47, 0x65,
	0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75, 0x6e,
	0x69, 0x66, 0x69, 0x5f, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69,
	0x5f, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x64, 0x65, 0x6e, 0x69, 0x73, 0x6b, 0x61, 0x70, 0x6f, 0x6e, 0x63, 0x68, 0x69, 0x6b, 0x2f,
	0x47, 0x6f, 0x53, 0x6f, 0x66, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f,
	0x75, 0x6e, 0x69, 0x66, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_unifi_proto_rawDescOnce sync.Once
	file_unifi_proto_rawDescData = file_unifi_proto_rawDesc
)

func file_unifi_proto_rawDescGZIP() []byte {
	file_unifi_proto_rawDescOnce.Do(func() {
		file_unifi_proto_rawDescData = protoimpl.X.CompressGZIP(file_unifi_proto_rawDescData)
	})
	return file_unifi_proto_rawDescData
}

var file_unifi_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_unifi_proto_goTypes = []interface{}{
	(*ClientRequest)(nil),  // 0: api.unifi_v1.ClientRequest
	(*ClientResponse)(nil), // 1: api.unifi_v1.ClientResponse
	(*Anomaly)(nil),        // 2: api.unifi_v1.Anomaly
}
var file_unifi_proto_depIdxs = []int32{
	2, // 0: api.unifi_v1.ClientResponse.anomalies:type_name -> api.unifi_v1.Anomaly
	0, // 1: api.unifi_v1.GetAnomalies.GetClient:input_type -> api.unifi_v1.ClientRequest
	1, // 2: api.unifi_v1.GetAnomalies.GetClient:output_type -> api.unifi_v1.ClientResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_unifi_proto_init() }
func file_unifi_proto_init() {
	if File_unifi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_unifi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientRequest); i {
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
		file_unifi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientResponse); i {
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
		file_unifi_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Anomaly); i {
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
			RawDescriptor: file_unifi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_unifi_proto_goTypes,
		DependencyIndexes: file_unifi_proto_depIdxs,
		MessageInfos:      file_unifi_proto_msgTypes,
	}.Build()
	File_unifi_proto = out.File
	file_unifi_proto_rawDesc = nil
	file_unifi_proto_goTypes = nil
	file_unifi_proto_depIdxs = nil
}