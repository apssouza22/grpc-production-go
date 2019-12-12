package grpc_client

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"strings"
)

//GrpcClientConnBuilder is a builder to create GRPC connection to the GRPC Server
type GrpcClientConnBuilder interface {
	WithContext(ctx context.Context)
	WithOptions(opts ...grpc.DialOption)
	WithInsecure()
	WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor)
	WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor)
	GetConn(addr string, port string) (*grpc.ClientConn, error)
}

//GRPC client builder
type GrpcClientBuilder struct {
	options            []grpc.DialOption
	enabledReflection  bool
	shutdownHook       func()
	enabledHealthCheck bool
	ctx                context.Context
}

// WithContext set the context to be used in the dial
func (b *GrpcClientBuilder) WithContext(ctx context.Context) {
	b.ctx = ctx
}

// WithOptions set dial options
func (b *GrpcClientBuilder) WithOptions(opts ...grpc.DialOption) {
	b.options = append(b.options, opts...)
}

// WithInsecure set the connection as insecure
func (b *GrpcClientBuilder) WithInsecure() {
	b.options = append(b.options, grpc.WithInsecure())
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for unary connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcClientBuilder) WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor) {
	b.options = append(b.options, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)))
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for stream connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcClientBuilder) WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor) {
	b.options = append(b.options, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...)))
}

// GetConn returns the client connection to the server
func (b *GrpcClientBuilder) GetConn(addr string, port string) (*grpc.ClientConn, error) {
	if addr == "" || port == "" {
		return nil, fmt.Errorf("target connection parameter missing. address = %s, port = %s", addr, port)
	}
	target := strings.Join([]string{addr, port}, ":")
	log.Debugf("Target to connect = %s", target)
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cc, err := grpc.DialContext(ctx, target, b.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to client. address = %s, port = %s. error = %+v", addr, port, err)
	}
	return cc, nil
}
