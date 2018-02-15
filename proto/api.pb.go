// Code generated by protoc-gen-go.
// source: api.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	api.proto

It has these top-level messages:
	BalancesRequest
	Balance
	BalancesResponse
	LogRequest
	LogEntry
	LogResponse
	StatusRequest
	StatusResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type BalancesRequest struct {
	As          string `protobuf:"bytes,1,opt,name=as" json:"as,omitempty"`
	IgnoreSmall bool   `protobuf:"varint,2,opt,name=ignoreSmall" json:"ignoreSmall,omitempty"`
}

func (m *BalancesRequest) Reset()                    { *m = BalancesRequest{} }
func (m *BalancesRequest) String() string            { return proto1.CompactTextString(m) }
func (*BalancesRequest) ProtoMessage()               {}
func (*BalancesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *BalancesRequest) GetAs() string {
	if m != nil {
		return m.As
	}
	return ""
}

func (m *BalancesRequest) GetIgnoreSmall() bool {
	if m != nil {
		return m.IgnoreSmall
	}
	return false
}

type Balance struct {
	Symbol       string                     `protobuf:"bytes,1,opt,name=symbol" json:"symbol,omitempty"`
	Exchange     string                     `protobuf:"bytes,2,opt,name=exchange" json:"exchange,omitempty"`
	Free         float32                    `protobuf:"fixed32,3,opt,name=free" json:"free,omitempty"`
	Locked       float32                    `protobuf:"fixed32,4,opt,name=locked" json:"locked,omitempty"`
	Total        float32                    `protobuf:"fixed32,5,opt,name=total" json:"total,omitempty"`
	As           string                     `protobuf:"bytes,6,opt,name=as" json:"as,omitempty"`
	Price        float32                    `protobuf:"fixed32,7,opt,name=price" json:"price,omitempty"`
	Value        float32                    `protobuf:"fixed32,8,opt,name=value" json:"value,omitempty"`
	At           *google_protobuf.Timestamp `protobuf:"bytes,9,opt,name=at" json:"at,omitempty"`
	Price24H     float32                    `protobuf:"fixed32,10,opt,name=price24h" json:"price24h,omitempty"`
	Value24H     float32                    `protobuf:"fixed32,11,opt,name=value24h" json:"value24h,omitempty"`
	Change24H    float32                    `protobuf:"fixed32,12,opt,name=change24h" json:"change24h,omitempty"`
	ChangePct24H float32                    `protobuf:"fixed32,13,opt,name=changePct24h" json:"changePct24h,omitempty"`
}

func (m *Balance) Reset()                    { *m = Balance{} }
func (m *Balance) String() string            { return proto1.CompactTextString(m) }
func (*Balance) ProtoMessage()               {}
func (*Balance) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Balance) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *Balance) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Balance) GetFree() float32 {
	if m != nil {
		return m.Free
	}
	return 0
}

func (m *Balance) GetLocked() float32 {
	if m != nil {
		return m.Locked
	}
	return 0
}

func (m *Balance) GetTotal() float32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *Balance) GetAs() string {
	if m != nil {
		return m.As
	}
	return ""
}

func (m *Balance) GetPrice() float32 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Balance) GetValue() float32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Balance) GetAt() *google_protobuf.Timestamp {
	if m != nil {
		return m.At
	}
	return nil
}

func (m *Balance) GetPrice24H() float32 {
	if m != nil {
		return m.Price24H
	}
	return 0
}

func (m *Balance) GetValue24H() float32 {
	if m != nil {
		return m.Value24H
	}
	return 0
}

func (m *Balance) GetChange24H() float32 {
	if m != nil {
		return m.Change24H
	}
	return 0
}

func (m *Balance) GetChangePct24H() float32 {
	if m != nil {
		return m.ChangePct24H
	}
	return 0
}

type BalancesResponse struct {
	Balances []*Balance `protobuf:"bytes,1,rep,name=balances" json:"balances,omitempty"`
}

func (m *BalancesResponse) Reset()                    { *m = BalancesResponse{} }
func (m *BalancesResponse) String() string            { return proto1.CompactTextString(m) }
func (*BalancesResponse) ProtoMessage()               {}
func (*BalancesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *BalancesResponse) GetBalances() []*Balance {
	if m != nil {
		return m.Balances
	}
	return nil
}

type LogRequest struct {
}

func (m *LogRequest) Reset()                    { *m = LogRequest{} }
func (m *LogRequest) String() string            { return proto1.CompactTextString(m) }
func (*LogRequest) ProtoMessage()               {}
func (*LogRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type LogEntry struct {
	Time *google_protobuf.Timestamp `protobuf:"bytes,1,opt,name=time" json:"time,omitempty"`
	Text string                     `protobuf:"bytes,2,opt,name=text" json:"text,omitempty"`
}

func (m *LogEntry) Reset()                    { *m = LogEntry{} }
func (m *LogEntry) String() string            { return proto1.CompactTextString(m) }
func (*LogEntry) ProtoMessage()               {}
func (*LogEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *LogEntry) GetTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *LogEntry) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type LogResponse struct {
	Entries []*LogEntry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *LogResponse) Reset()                    { *m = LogResponse{} }
func (m *LogResponse) String() string            { return proto1.CompactTextString(m) }
func (*LogResponse) ProtoMessage()               {}
func (*LogResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *LogResponse) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto1.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type StatusResponse struct {
	ServerStarted *google_protobuf.Timestamp `protobuf:"bytes,1,opt,name=serverStarted" json:"serverStarted,omitempty"`
	LastUpdate    *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=lastUpdate" json:"lastUpdate,omitempty"`
	UpdateCount   int32                      `protobuf:"varint,3,opt,name=updateCount" json:"updateCount,omitempty"`
	TotalSymbols  int32                      `protobuf:"varint,4,opt,name=totalSymbols" json:"totalSymbols,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto1.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *StatusResponse) GetServerStarted() *google_protobuf.Timestamp {
	if m != nil {
		return m.ServerStarted
	}
	return nil
}

func (m *StatusResponse) GetLastUpdate() *google_protobuf.Timestamp {
	if m != nil {
		return m.LastUpdate
	}
	return nil
}

func (m *StatusResponse) GetUpdateCount() int32 {
	if m != nil {
		return m.UpdateCount
	}
	return 0
}

func (m *StatusResponse) GetTotalSymbols() int32 {
	if m != nil {
		return m.TotalSymbols
	}
	return 0
}

func init() {
	proto1.RegisterType((*BalancesRequest)(nil), "proto.BalancesRequest")
	proto1.RegisterType((*Balance)(nil), "proto.Balance")
	proto1.RegisterType((*BalancesResponse)(nil), "proto.BalancesResponse")
	proto1.RegisterType((*LogRequest)(nil), "proto.LogRequest")
	proto1.RegisterType((*LogEntry)(nil), "proto.LogEntry")
	proto1.RegisterType((*LogResponse)(nil), "proto.LogResponse")
	proto1.RegisterType((*StatusRequest)(nil), "proto.StatusRequest")
	proto1.RegisterType((*StatusResponse)(nil), "proto.StatusResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Teletrada service

type TeletradaClient interface {
	// Get current balances
	GetBalances(ctx context.Context, in *BalancesRequest, opts ...grpc.CallOption) (*BalancesResponse, error)
	GetLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error)
	GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type teletradaClient struct {
	cc *grpc.ClientConn
}

func NewTeletradaClient(cc *grpc.ClientConn) TeletradaClient {
	return &teletradaClient{cc}
}

func (c *teletradaClient) GetBalances(ctx context.Context, in *BalancesRequest, opts ...grpc.CallOption) (*BalancesResponse, error) {
	out := new(BalancesResponse)
	err := grpc.Invoke(ctx, "/proto.teletrada/GetBalances", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teletradaClient) GetLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error) {
	out := new(LogResponse)
	err := grpc.Invoke(ctx, "/proto.teletrada/GetLog", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teletradaClient) GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := grpc.Invoke(ctx, "/proto.teletrada/GetStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Teletrada service

type TeletradaServer interface {
	// Get current balances
	GetBalances(context.Context, *BalancesRequest) (*BalancesResponse, error)
	GetLog(context.Context, *LogRequest) (*LogResponse, error)
	GetStatus(context.Context, *StatusRequest) (*StatusResponse, error)
}

func RegisterTeletradaServer(s *grpc.Server, srv TeletradaServer) {
	s.RegisterService(&_Teletrada_serviceDesc, srv)
}

func _Teletrada_GetBalances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeletradaServer).GetBalances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.teletrada/GetBalances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeletradaServer).GetBalances(ctx, req.(*BalancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Teletrada_GetLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeletradaServer).GetLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.teletrada/GetLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeletradaServer).GetLog(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Teletrada_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeletradaServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.teletrada/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeletradaServer).GetStatus(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Teletrada_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.teletrada",
	HandlerType: (*TeletradaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBalances",
			Handler:    _Teletrada_GetBalances_Handler,
		},
		{
			MethodName: "GetLog",
			Handler:    _Teletrada_GetLog_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _Teletrada_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto1.RegisterFile("api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 570 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x25, 0x69, 0x93, 0xc6, 0x93, 0x26, 0x81, 0x55, 0x29, 0x56, 0x40, 0xa2, 0xf2, 0x09, 0x38,
	0xb8, 0x10, 0x38, 0xa0, 0x1e, 0x10, 0x4a, 0x85, 0xe0, 0x80, 0x50, 0xe4, 0x80, 0x38, 0x6f, 0x9c,
	0xa9, 0x63, 0xb1, 0xf1, 0x1a, 0x7b, 0x5d, 0xd2, 0x7f, 0xe3, 0x33, 0xf8, 0x01, 0xfe, 0x84, 0xdd,
	0xd9, 0xb5, 0x93, 0x80, 0x54, 0x4e, 0x9e, 0xf7, 0x66, 0xe6, 0x79, 0xf6, 0xcd, 0x80, 0xc7, 0xf3,
	0x34, 0xcc, 0x0b, 0xa9, 0x24, 0xeb, 0xd0, 0x67, 0xfc, 0x38, 0x91, 0x32, 0x11, 0x78, 0x4e, 0x68,
	0x51, 0x5d, 0x9d, 0xab, 0x74, 0x8d, 0xa5, 0xe2, 0xeb, 0xdc, 0xd6, 0x05, 0x97, 0x30, 0x9a, 0x72,
	0xc1, 0xb3, 0x18, 0xcb, 0x08, 0xbf, 0x57, 0x3a, 0xc7, 0x86, 0xd0, 0xe6, 0xa5, 0xdf, 0x3a, 0x6b,
	0x3d, 0xf1, 0x22, 0x1d, 0xb1, 0x33, 0xe8, 0xa7, 0x49, 0x26, 0x0b, 0x9c, 0xaf, 0xb9, 0x10, 0x7e,
	0x5b, 0x27, 0x7a, 0xd1, 0x2e, 0x15, 0xfc, 0x6e, 0xc3, 0x91, 0x53, 0x61, 0xa7, 0xd0, 0x2d, 0x6f,
	0xd6, 0x0b, 0x29, 0x9c, 0x82, 0x43, 0x6c, 0x0c, 0x3d, 0xdc, 0xc4, 0x2b, 0x9e, 0x25, 0x48, 0x12,
	0x5e, 0xd4, 0x60, 0xc6, 0xe0, 0xf0, 0xaa, 0x40, 0xf4, 0x0f, 0x34, 0xdf, 0x8e, 0x28, 0x36, 0x3a,
	0x42, 0xc6, 0xdf, 0x70, 0xe9, 0x1f, 0x12, 0xeb, 0x10, 0x3b, 0x81, 0x8e, 0x92, 0x8a, 0x0b, 0xbf,
	0x43, 0xb4, 0x05, 0x6e, 0xe6, 0x6e, 0x33, 0xb3, 0xae, 0xca, 0x8b, 0x34, 0x46, 0xff, 0xc8, 0x56,
	0x11, 0x30, 0xec, 0x35, 0x17, 0x15, 0xfa, 0x3d, 0xcb, 0x12, 0x60, 0xcf, 0x74, 0xaf, 0xf2, 0x3d,
	0x4d, 0xf5, 0x27, 0xe3, 0xd0, 0x1a, 0x16, 0xd6, 0x86, 0x85, 0x9f, 0x6b, 0xc3, 0xb4, 0xae, 0x32,
	0xaf, 0x20, 0xa9, 0xc9, 0xab, 0x95, 0x0f, 0x24, 0xd2, 0x60, 0x93, 0x23, 0x41, 0x93, 0xeb, 0xdb,
	0x5c, 0x8d, 0xd9, 0x23, 0xf0, 0xec, 0x5b, 0x4d, 0xf2, 0x98, 0x92, 0x5b, 0x82, 0x05, 0x70, 0x6c,
	0xc1, 0x2c, 0x56, 0xa6, 0x60, 0x40, 0x05, 0x7b, 0x5c, 0xf0, 0x06, 0xee, 0x6e, 0x17, 0x55, 0xe6,
	0x32, 0x2b, 0xcd, 0xe4, 0xbd, 0x85, 0xe3, 0xb4, 0xdb, 0x07, 0x7a, 0xfe, 0xa1, 0x1d, 0x3c, 0x74,
	0xa5, 0x51, 0x93, 0x0f, 0x8e, 0x01, 0x3e, 0xca, 0xc4, 0xed, 0x38, 0xf8, 0x04, 0x3d, 0x8d, 0xde,
	0x65, 0xaa, 0xb8, 0x61, 0x21, 0x1c, 0x9a, 0xab, 0xa0, 0x7d, 0xdd, 0xee, 0x00, 0xd5, 0x99, 0x6d,
	0x29, 0xdc, 0x28, 0xb7, 0x45, 0x8a, 0x83, 0xd7, 0xd0, 0x27, 0x75, 0x37, 0xd8, 0x53, 0x38, 0x42,
	0xad, 0x9d, 0x36, 0x73, 0x8d, 0xdc, 0x5c, 0xf5, 0x4f, 0xa3, 0x3a, 0x1f, 0x8c, 0x60, 0x30, 0x57,
	0x5c, 0x55, 0xf5, 0xf9, 0x05, 0xbf, 0x5a, 0x30, 0xac, 0x19, 0x27, 0xf7, 0x16, 0x06, 0x25, 0x16,
	0xd7, 0x58, 0x68, 0xbe, 0x50, 0xfa, 0x24, 0xfe, 0x3f, 0xea, 0x7e, 0x03, 0xbb, 0x00, 0x10, 0xbc,
	0x54, 0x5f, 0xf2, 0x25, 0x57, 0xf6, 0xfe, 0x6e, 0x6f, 0xdf, 0xa9, 0x36, 0xf7, 0x5f, 0x51, 0x74,
	0x29, 0xab, 0x4c, 0xd1, 0x91, 0x76, 0xa2, 0x5d, 0xca, 0xec, 0x8f, 0xce, 0x70, 0x4e, 0xa7, 0x5e,
	0xd2, 0xc5, 0x76, 0xa2, 0x3d, 0x6e, 0xf2, 0xb3, 0x05, 0x9e, 0x42, 0x81, 0xaa, 0xe0, 0x4b, 0xae,
	0x5f, 0xd4, 0x7f, 0x8f, 0xaa, 0x5e, 0x28, 0x3b, 0xdd, 0x5f, 0x5b, 0xed, 0xc5, 0xf8, 0xc1, 0x3f,
	0xbc, 0x75, 0x24, 0xb8, 0xc3, 0x5e, 0x40, 0x57, 0x2b, 0x68, 0x3f, 0xd9, 0xbd, 0xad, 0xb7, 0x75,
	0x1f, 0xdb, 0xa5, 0x9a, 0x96, 0x0b, 0xf0, 0x74, 0x8b, 0xf5, 0x96, 0x9d, 0xb8, 0x92, 0x3d, 0xf3,
	0xc7, 0xf7, 0xff, 0x62, 0xeb, 0xde, 0xe9, 0x73, 0x78, 0x98, 0xca, 0x30, 0x29, 0xf2, 0x38, 0xc4,
	0x8d, 0x76, 0x48, 0x60, 0x19, 0xae, 0x50, 0x08, 0xf9, 0x43, 0x16, 0x62, 0x39, 0x1d, 0x7d, 0x30,
	0xf1, 0x57, 0x13, 0xcf, 0x8c, 0xc0, 0xac, 0xb5, 0xe8, 0x92, 0xd2, 0xcb, 0x3f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x5c, 0xa4, 0x83, 0x21, 0x95, 0x04, 0x00, 0x00,
}
