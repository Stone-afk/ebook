// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: search/v1/search.proto

package searchv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	SearchService_Search_FullMethodName = "/search.v1.SearchService/Search"
)

// SearchServiceClient is the client API for SearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchServiceClient interface {
	// 这个是最为模糊的搜索接口
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
}

type searchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchServiceClient(cc grpc.ClientConnInterface) SearchServiceClient {
	return &searchServiceClient{cc}
}

func (c *searchServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, SearchService_Search_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServiceServer is the server API for SearchService service.
// All implementations must embed UnimplementedSearchServiceServer
// for forward compatibility
type SearchServiceServer interface {
	// 这个是最为模糊的搜索接口
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	mustEmbedUnimplementedSearchServiceServer()
}

// UnimplementedSearchServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSearchServiceServer struct {
}

func (UnimplementedSearchServiceServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedSearchServiceServer) mustEmbedUnimplementedSearchServiceServer() {}

// UnsafeSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServiceServer will
// result in compilation errors.
type UnsafeSearchServiceServer interface {
	mustEmbedUnimplementedSearchServiceServer()
}

func RegisterSearchServiceServer(s grpc.ServiceRegistrar, srv SearchServiceServer) {
	s.RegisterService(&SearchService_ServiceDesc, srv)
}

func _SearchService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SearchService_ServiceDesc is the grpc.ServiceDesc for SearchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.v1.SearchService",
	HandlerType: (*SearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _SearchService_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search/v1/search.proto",
}

const (
	UserSearchService_SearchUser_FullMethodName = "/search.v1.UserSearchService/SearchUser"
)

// UserSearchServiceClient is the client API for UserSearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 业务专属接口
type UserSearchServiceClient interface {
	SearchUser(ctx context.Context, in *UserSearchRequest, opts ...grpc.CallOption) (*UserSearchResponse, error)
}

type userSearchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserSearchServiceClient(cc grpc.ClientConnInterface) UserSearchServiceClient {
	return &userSearchServiceClient{cc}
}

func (c *userSearchServiceClient) SearchUser(ctx context.Context, in *UserSearchRequest, opts ...grpc.CallOption) (*UserSearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserSearchResponse)
	err := c.cc.Invoke(ctx, UserSearchService_SearchUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserSearchServiceServer is the server API for UserSearchService service.
// All implementations must embed UnimplementedUserSearchServiceServer
// for forward compatibility
//
// 业务专属接口
type UserSearchServiceServer interface {
	SearchUser(context.Context, *UserSearchRequest) (*UserSearchResponse, error)
	mustEmbedUnimplementedUserSearchServiceServer()
}

// UnimplementedUserSearchServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserSearchServiceServer struct {
}

func (UnimplementedUserSearchServiceServer) SearchUser(context.Context, *UserSearchRequest) (*UserSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchUser not implemented")
}
func (UnimplementedUserSearchServiceServer) mustEmbedUnimplementedUserSearchServiceServer() {}

// UnsafeUserSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserSearchServiceServer will
// result in compilation errors.
type UnsafeUserSearchServiceServer interface {
	mustEmbedUnimplementedUserSearchServiceServer()
}

func RegisterUserSearchServiceServer(s grpc.ServiceRegistrar, srv UserSearchServiceServer) {
	s.RegisterService(&UserSearchService_ServiceDesc, srv)
}

func _UserSearchService_SearchUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSearchServiceServer).SearchUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserSearchService_SearchUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSearchServiceServer).SearchUser(ctx, req.(*UserSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserSearchService_ServiceDesc is the grpc.ServiceDesc for UserSearchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserSearchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.v1.UserSearchService",
	HandlerType: (*UserSearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchUser",
			Handler:    _UserSearchService_SearchUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search/v1/search.proto",
}

const (
	ArticleSearchService_SearchArticle_FullMethodName = "/search.v1.ArticleSearchService/SearchArticle"
)

// ArticleSearchServiceClient is the client API for ArticleSearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArticleSearchServiceClient interface {
	SearchArticle(ctx context.Context, in *ArticleSearchRequest, opts ...grpc.CallOption) (*ArticleSearchResponse, error)
}

type articleSearchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewArticleSearchServiceClient(cc grpc.ClientConnInterface) ArticleSearchServiceClient {
	return &articleSearchServiceClient{cc}
}

func (c *articleSearchServiceClient) SearchArticle(ctx context.Context, in *ArticleSearchRequest, opts ...grpc.CallOption) (*ArticleSearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ArticleSearchResponse)
	err := c.cc.Invoke(ctx, ArticleSearchService_SearchArticle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArticleSearchServiceServer is the server API for ArticleSearchService service.
// All implementations must embed UnimplementedArticleSearchServiceServer
// for forward compatibility
type ArticleSearchServiceServer interface {
	SearchArticle(context.Context, *ArticleSearchRequest) (*ArticleSearchResponse, error)
	mustEmbedUnimplementedArticleSearchServiceServer()
}

// UnimplementedArticleSearchServiceServer must be embedded to have forward compatible implementations.
type UnimplementedArticleSearchServiceServer struct {
}

func (UnimplementedArticleSearchServiceServer) SearchArticle(context.Context, *ArticleSearchRequest) (*ArticleSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchArticle not implemented")
}
func (UnimplementedArticleSearchServiceServer) mustEmbedUnimplementedArticleSearchServiceServer() {}

// UnsafeArticleSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArticleSearchServiceServer will
// result in compilation errors.
type UnsafeArticleSearchServiceServer interface {
	mustEmbedUnimplementedArticleSearchServiceServer()
}

func RegisterArticleSearchServiceServer(s grpc.ServiceRegistrar, srv ArticleSearchServiceServer) {
	s.RegisterService(&ArticleSearchService_ServiceDesc, srv)
}

func _ArticleSearchService_SearchArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleSearchServiceServer).SearchArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleSearchService_SearchArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleSearchServiceServer).SearchArticle(ctx, req.(*ArticleSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ArticleSearchService_ServiceDesc is the grpc.ServiceDesc for ArticleSearchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ArticleSearchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.v1.ArticleSearchService",
	HandlerType: (*ArticleSearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchArticle",
			Handler:    _ArticleSearchService_SearchArticle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search/v1/search.proto",
}

const (
	TagSearchService_SearchBizTags_FullMethodName = "/search.v1.TagSearchService/SearchBizTags"
)

// TagSearchServiceClient is the client API for TagSearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TagSearchServiceClient interface {
	SearchBizTags(ctx context.Context, in *BizTagsSearchRequest, opts ...grpc.CallOption) (*BizTagsSearchResponse, error)
}

type tagSearchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTagSearchServiceClient(cc grpc.ClientConnInterface) TagSearchServiceClient {
	return &tagSearchServiceClient{cc}
}

func (c *tagSearchServiceClient) SearchBizTags(ctx context.Context, in *BizTagsSearchRequest, opts ...grpc.CallOption) (*BizTagsSearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BizTagsSearchResponse)
	err := c.cc.Invoke(ctx, TagSearchService_SearchBizTags_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TagSearchServiceServer is the server API for TagSearchService service.
// All implementations must embed UnimplementedTagSearchServiceServer
// for forward compatibility
type TagSearchServiceServer interface {
	SearchBizTags(context.Context, *BizTagsSearchRequest) (*BizTagsSearchResponse, error)
	mustEmbedUnimplementedTagSearchServiceServer()
}

// UnimplementedTagSearchServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTagSearchServiceServer struct {
}

func (UnimplementedTagSearchServiceServer) SearchBizTags(context.Context, *BizTagsSearchRequest) (*BizTagsSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchBizTags not implemented")
}
func (UnimplementedTagSearchServiceServer) mustEmbedUnimplementedTagSearchServiceServer() {}

// UnsafeTagSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TagSearchServiceServer will
// result in compilation errors.
type UnsafeTagSearchServiceServer interface {
	mustEmbedUnimplementedTagSearchServiceServer()
}

func RegisterTagSearchServiceServer(s grpc.ServiceRegistrar, srv TagSearchServiceServer) {
	s.RegisterService(&TagSearchService_ServiceDesc, srv)
}

func _TagSearchService_SearchBizTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BizTagsSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagSearchServiceServer).SearchBizTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TagSearchService_SearchBizTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagSearchServiceServer).SearchBizTags(ctx, req.(*BizTagsSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TagSearchService_ServiceDesc is the grpc.ServiceDesc for TagSearchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TagSearchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.v1.TagSearchService",
	HandlerType: (*TagSearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchBizTags",
			Handler:    _TagSearchService_SearchBizTags_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search/v1/search.proto",
}
