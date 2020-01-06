// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/hyparview.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_hyparview_d3c1d53ca2be5b46, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type FromRequest struct {
	From                 string   `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FromRequest) Reset()         { *m = FromRequest{} }
func (m *FromRequest) String() string { return proto.CompactTextString(m) }
func (*FromRequest) ProtoMessage()    {}
func (*FromRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_hyparview_d3c1d53ca2be5b46, []int{1}
}
func (m *FromRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FromRequest.Unmarshal(m, b)
}
func (m *FromRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FromRequest.Marshal(b, m, deterministic)
}
func (dst *FromRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FromRequest.Merge(dst, src)
}
func (m *FromRequest) XXX_Size() int {
	return xxx_messageInfo_FromRequest.Size(m)
}
func (m *FromRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FromRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FromRequest proto.InternalMessageInfo

func (m *FromRequest) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

type ForwardJoinRequest struct {
	Ttl                  int32    `protobuf:"varint,1,opt,name=ttl,proto3" json:"ttl,omitempty"`
	Join                 string   `protobuf:"bytes,2,opt,name=join,proto3" json:"join,omitempty"`
	From                 string   `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardJoinRequest) Reset()         { *m = ForwardJoinRequest{} }
func (m *ForwardJoinRequest) String() string { return proto.CompactTextString(m) }
func (*ForwardJoinRequest) ProtoMessage()    {}
func (*ForwardJoinRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_hyparview_d3c1d53ca2be5b46, []int{2}
}
func (m *ForwardJoinRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardJoinRequest.Unmarshal(m, b)
}
func (m *ForwardJoinRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardJoinRequest.Marshal(b, m, deterministic)
}
func (dst *ForwardJoinRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardJoinRequest.Merge(dst, src)
}
func (m *ForwardJoinRequest) XXX_Size() int {
	return xxx_messageInfo_ForwardJoinRequest.Size(m)
}
func (m *ForwardJoinRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardJoinRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardJoinRequest proto.InternalMessageInfo

func (m *ForwardJoinRequest) GetTtl() int32 {
	if m != nil {
		return m.Ttl
	}
	return 0
}

func (m *ForwardJoinRequest) GetJoin() string {
	if m != nil {
		return m.Join
	}
	return ""
}

func (m *ForwardJoinRequest) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

type NeighborRequest struct {
	Priority             bool     `protobuf:"varint,1,opt,name=priority,proto3" json:"priority,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NeighborRequest) Reset()         { *m = NeighborRequest{} }
func (m *NeighborRequest) String() string { return proto.CompactTextString(m) }
func (*NeighborRequest) ProtoMessage()    {}
func (*NeighborRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_hyparview_d3c1d53ca2be5b46, []int{3}
}
func (m *NeighborRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NeighborRequest.Unmarshal(m, b)
}
func (m *NeighborRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NeighborRequest.Marshal(b, m, deterministic)
}
func (dst *NeighborRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NeighborRequest.Merge(dst, src)
}
func (m *NeighborRequest) XXX_Size() int {
	return xxx_messageInfo_NeighborRequest.Size(m)
}
func (m *NeighborRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NeighborRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NeighborRequest proto.InternalMessageInfo

func (m *NeighborRequest) GetPriority() bool {
	if m != nil {
		return m.Priority
	}
	return false
}

func (m *NeighborRequest) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

type NeighborResponse struct {
	Accept               bool     `protobuf:"varint,1,opt,name=accept,proto3" json:"accept,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NeighborResponse) Reset()         { *m = NeighborResponse{} }
func (m *NeighborResponse) String() string { return proto.CompactTextString(m) }
func (*NeighborResponse) ProtoMessage()    {}
func (*NeighborResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_hyparview_d3c1d53ca2be5b46, []int{4}
}
func (m *NeighborResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NeighborResponse.Unmarshal(m, b)
}
func (m *NeighborResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NeighborResponse.Marshal(b, m, deterministic)
}
func (dst *NeighborResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NeighborResponse.Merge(dst, src)
}
func (m *NeighborResponse) XXX_Size() int {
	return xxx_messageInfo_NeighborResponse.Size(m)
}
func (m *NeighborResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NeighborResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NeighborResponse proto.InternalMessageInfo

func (m *NeighborResponse) GetAccept() bool {
	if m != nil {
		return m.Accept
	}
	return false
}

func (m *NeighborResponse) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

type ShuffleRequest struct {
	Ttl                  int32    `protobuf:"varint,1,opt,name=ttl,proto3" json:"ttl,omitempty"`
	Active               []string `protobuf:"bytes,2,rep,name=active,proto3" json:"active,omitempty"`
	Passive              []string `protobuf:"bytes,3,rep,name=passive,proto3" json:"passive,omitempty"`
	From                 string   `protobuf:"bytes,4,opt,name=from,proto3" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShuffleRequest) Reset()         { *m = ShuffleRequest{} }
func (m *ShuffleRequest) String() string { return proto.CompactTextString(m) }
func (*ShuffleRequest) ProtoMessage()    {}
func (*ShuffleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_hyparview_d3c1d53ca2be5b46, []int{5}
}
func (m *ShuffleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShuffleRequest.Unmarshal(m, b)
}
func (m *ShuffleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShuffleRequest.Marshal(b, m, deterministic)
}
func (dst *ShuffleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShuffleRequest.Merge(dst, src)
}
func (m *ShuffleRequest) XXX_Size() int {
	return xxx_messageInfo_ShuffleRequest.Size(m)
}
func (m *ShuffleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShuffleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShuffleRequest proto.InternalMessageInfo

func (m *ShuffleRequest) GetTtl() int32 {
	if m != nil {
		return m.Ttl
	}
	return 0
}

func (m *ShuffleRequest) GetActive() []string {
	if m != nil {
		return m.Active
	}
	return nil
}

func (m *ShuffleRequest) GetPassive() []string {
	if m != nil {
		return m.Passive
	}
	return nil
}

func (m *ShuffleRequest) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

type ShuffleReply struct {
	Passive              []string `protobuf:"bytes,1,rep,name=passive,proto3" json:"passive,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShuffleReply) Reset()         { *m = ShuffleReply{} }
func (m *ShuffleReply) String() string { return proto.CompactTextString(m) }
func (*ShuffleReply) ProtoMessage()    {}
func (*ShuffleReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_hyparview_d3c1d53ca2be5b46, []int{6}
}
func (m *ShuffleReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShuffleReply.Unmarshal(m, b)
}
func (m *ShuffleReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShuffleReply.Marshal(b, m, deterministic)
}
func (dst *ShuffleReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShuffleReply.Merge(dst, src)
}
func (m *ShuffleReply) XXX_Size() int {
	return xxx_messageInfo_ShuffleReply.Size(m)
}
func (m *ShuffleReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ShuffleReply.DiscardUnknown(m)
}

var xxx_messageInfo_ShuffleReply proto.InternalMessageInfo

func (m *ShuffleReply) GetPassive() []string {
	if m != nil {
		return m.Passive
	}
	return nil
}

func (m *ShuffleReply) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func init() {
	proto.RegisterType((*Empty)(nil), "hashicorp.hyparview.example.hyparview.Empty")
	proto.RegisterType((*FromRequest)(nil), "hashicorp.hyparview.example.hyparview.FromRequest")
	proto.RegisterType((*ForwardJoinRequest)(nil), "hashicorp.hyparview.example.hyparview.ForwardJoinRequest")
	proto.RegisterType((*NeighborRequest)(nil), "hashicorp.hyparview.example.hyparview.NeighborRequest")
	proto.RegisterType((*NeighborResponse)(nil), "hashicorp.hyparview.example.hyparview.NeighborResponse")
	proto.RegisterType((*ShuffleRequest)(nil), "hashicorp.hyparview.example.hyparview.ShuffleRequest")
	proto.RegisterType((*ShuffleReply)(nil), "hashicorp.hyparview.example.hyparview.ShuffleReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HyparviewClient is the client API for Hyparview service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HyparviewClient interface {
	Join(ctx context.Context, in *FromRequest, opts ...grpc.CallOption) (*Empty, error)
	ForwardJoin(ctx context.Context, in *ForwardJoinRequest, opts ...grpc.CallOption) (*Empty, error)
	Disconnect(ctx context.Context, in *FromRequest, opts ...grpc.CallOption) (*Empty, error)
	Neighbor(ctx context.Context, in *NeighborRequest, opts ...grpc.CallOption) (*NeighborResponse, error)
	Shuffle(ctx context.Context, in *ShuffleRequest, opts ...grpc.CallOption) (*ShuffleReply, error)
}

type hyparviewClient struct {
	cc *grpc.ClientConn
}

func NewHyparviewClient(cc *grpc.ClientConn) HyparviewClient {
	return &hyparviewClient{cc}
}

func (c *hyparviewClient) Join(ctx context.Context, in *FromRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/hashicorp.hyparview.example.hyparview.Hyparview/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hyparviewClient) ForwardJoin(ctx context.Context, in *ForwardJoinRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/hashicorp.hyparview.example.hyparview.Hyparview/ForwardJoin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hyparviewClient) Disconnect(ctx context.Context, in *FromRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/hashicorp.hyparview.example.hyparview.Hyparview/Disconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hyparviewClient) Neighbor(ctx context.Context, in *NeighborRequest, opts ...grpc.CallOption) (*NeighborResponse, error) {
	out := new(NeighborResponse)
	err := c.cc.Invoke(ctx, "/hashicorp.hyparview.example.hyparview.Hyparview/Neighbor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hyparviewClient) Shuffle(ctx context.Context, in *ShuffleRequest, opts ...grpc.CallOption) (*ShuffleReply, error) {
	out := new(ShuffleReply)
	err := c.cc.Invoke(ctx, "/hashicorp.hyparview.example.hyparview.Hyparview/Shuffle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HyparviewServer is the server API for Hyparview service.
type HyparviewServer interface {
	Join(context.Context, *FromRequest) (*Empty, error)
	ForwardJoin(context.Context, *ForwardJoinRequest) (*Empty, error)
	Disconnect(context.Context, *FromRequest) (*Empty, error)
	Neighbor(context.Context, *NeighborRequest) (*NeighborResponse, error)
	Shuffle(context.Context, *ShuffleRequest) (*ShuffleReply, error)
}

func RegisterHyparviewServer(s *grpc.Server, srv HyparviewServer) {
	s.RegisterService(&_Hyparview_serviceDesc, srv)
}

func _Hyparview_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FromRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HyparviewServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hashicorp.hyparview.example.hyparview.Hyparview/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HyparviewServer).Join(ctx, req.(*FromRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hyparview_ForwardJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardJoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HyparviewServer).ForwardJoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hashicorp.hyparview.example.hyparview.Hyparview/ForwardJoin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HyparviewServer).ForwardJoin(ctx, req.(*ForwardJoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hyparview_Disconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FromRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HyparviewServer).Disconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hashicorp.hyparview.example.hyparview.Hyparview/Disconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HyparviewServer).Disconnect(ctx, req.(*FromRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hyparview_Neighbor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NeighborRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HyparviewServer).Neighbor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hashicorp.hyparview.example.hyparview.Hyparview/Neighbor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HyparviewServer).Neighbor(ctx, req.(*NeighborRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hyparview_Shuffle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShuffleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HyparviewServer).Shuffle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hashicorp.hyparview.example.hyparview.Hyparview/Shuffle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HyparviewServer).Shuffle(ctx, req.(*ShuffleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Hyparview_serviceDesc = grpc.ServiceDesc{
	ServiceName: "hashicorp.hyparview.example.hyparview.Hyparview",
	HandlerType: (*HyparviewServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _Hyparview_Join_Handler,
		},
		{
			MethodName: "ForwardJoin",
			Handler:    _Hyparview_ForwardJoin_Handler,
		},
		{
			MethodName: "Disconnect",
			Handler:    _Hyparview_Disconnect_Handler,
		},
		{
			MethodName: "Neighbor",
			Handler:    _Hyparview_Neighbor_Handler,
		},
		{
			MethodName: "Shuffle",
			Handler:    _Hyparview_Shuffle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/hyparview.proto",
}

func init() { proto.RegisterFile("proto/hyparview.proto", fileDescriptor_hyparview_d3c1d53ca2be5b46) }

var fileDescriptor_hyparview_d3c1d53ca2be5b46 = []byte{
	// 385 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x94, 0x4f, 0x4b, 0xeb, 0x40,
	0x14, 0xc5, 0xdb, 0x26, 0x6d, 0xda, 0xdb, 0xc7, 0x7b, 0x65, 0xe0, 0x49, 0xc8, 0xaa, 0x06, 0x84,
	0x2e, 0x24, 0x42, 0x8b, 0x8a, 0x20, 0x82, 0xa2, 0x45, 0x5c, 0x74, 0x11, 0x77, 0xee, 0xd2, 0x38,
	0x35, 0x53, 0x92, 0xcc, 0x38, 0x99, 0xfe, 0xc9, 0xc2, 0x4f, 0xea, 0x97, 0x91, 0x4c, 0x93, 0x74,
	0xd4, 0x0a, 0xe9, 0xc6, 0x55, 0xee, 0xcd, 0x9d, 0xf3, 0xbb, 0x19, 0xce, 0x21, 0xf0, 0x9f, 0x71,
	0x2a, 0xe8, 0x49, 0x90, 0x32, 0x8f, 0x2f, 0x09, 0x5e, 0x39, 0xb2, 0x47, 0x47, 0x81, 0x97, 0x04,
	0xc4, 0xa7, 0x9c, 0x39, 0xdb, 0x11, 0x5e, 0x7b, 0x11, 0x0b, 0xf1, 0xf6, 0x8d, 0x6d, 0x40, 0xf3,
	0x2e, 0x62, 0x22, 0xb5, 0x0f, 0xa1, 0x3b, 0xe6, 0x34, 0x72, 0xf1, 0xeb, 0x02, 0x27, 0x02, 0x21,
	0xd0, 0x67, 0x9c, 0x46, 0x66, 0xbd, 0x5f, 0x1f, 0x74, 0x5c, 0x59, 0xdb, 0x13, 0x40, 0x63, 0xca,
	0x57, 0x1e, 0x7f, 0x7e, 0xa0, 0x24, 0x2e, 0x4e, 0xf6, 0x40, 0x13, 0x22, 0x94, 0x07, 0x9b, 0x6e,
	0x56, 0x66, 0xda, 0x39, 0x25, 0xb1, 0xd9, 0xd8, 0x68, 0xb3, 0xba, 0xe4, 0x69, 0x0a, 0xef, 0x1a,
	0xfe, 0x4d, 0x30, 0x79, 0x09, 0xa6, 0x94, 0x17, 0x30, 0x0b, 0xda, 0x8c, 0x13, 0xca, 0x89, 0x48,
	0x25, 0xb1, 0xed, 0x96, 0x7d, 0x89, 0x68, 0x28, 0x88, 0x2b, 0xe8, 0x6d, 0x11, 0x09, 0xa3, 0x71,
	0x82, 0xd1, 0x01, 0xb4, 0x3c, 0xdf, 0xc7, 0x4c, 0xe4, 0x84, 0xbc, 0xdb, 0xa9, 0x0f, 0xe0, 0xef,
	0x63, 0xb0, 0x98, 0xcd, 0x42, 0xfc, 0xf3, 0x75, 0x24, 0x4f, 0x90, 0x25, 0x36, 0x1b, 0x7d, 0x6d,
	0xd0, 0x71, 0xf3, 0x0e, 0x99, 0x60, 0x30, 0x2f, 0x49, 0xb2, 0x81, 0x26, 0x07, 0x45, 0x5b, 0x6e,
	0xd2, 0x95, 0x4d, 0x97, 0xf0, 0xa7, 0xdc, 0xc4, 0xc2, 0x54, 0x55, 0xd7, 0x77, 0xab, 0x95, 0xef,
	0x1c, 0xbe, 0xeb, 0xd0, 0xb9, 0x2f, 0x4c, 0x43, 0x73, 0xd0, 0x33, 0x07, 0xd0, 0xd0, 0xa9, 0x64,
	0xb2, 0xa3, 0x18, 0x6b, 0x1d, 0x57, 0xd4, 0x6c, 0x52, 0x51, 0x43, 0x6b, 0xe8, 0x2a, 0xa6, 0xa3,
	0x8b, 0xaa, 0x2b, 0xbf, 0x05, 0x65, 0xef, 0xcd, 0x0c, 0xe0, 0x96, 0x24, 0x3e, 0x8d, 0x63, 0xec,
	0x8b, 0x5f, 0xb9, 0xeb, 0x1b, 0xb4, 0x8b, 0x34, 0xa1, 0xb3, 0x8a, 0xda, 0x2f, 0x09, 0xb6, 0xce,
	0xf7, 0xd6, 0x6d, 0x62, 0x6b, 0xd7, 0xd0, 0x0a, 0x8c, 0x3c, 0x22, 0xe8, 0xb4, 0x22, 0xe5, 0x73,
	0x78, 0xad, 0xd1, 0xbe, 0x32, 0x16, 0xa6, 0x76, 0xed, 0xc6, 0x78, 0x6a, 0xca, 0x9f, 0xc6, 0xb4,
	0x25, 0x1f, 0xa3, 0x8f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xce, 0x78, 0x8c, 0x69, 0x54, 0x04, 0x00,
	0x00,
}