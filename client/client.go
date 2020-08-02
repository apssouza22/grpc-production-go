package grpc_client

import (
	"context"
	"fmt"
	"github.com/apssouza22/grpc-server-go/cert"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

//GrpcClientConnBuilder is a builder to create GRPC connection to the GRPC Server
type GrpcClientConnBuilder interface {
	WithContext(ctx context.Context)
	WithOptions(opts ...grpc.DialOption)
	WithInsecure()
	WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor)
	WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor)
	WithKeepAliveParams(params keepalive.ClientParameters)
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

// WithTLS set the connection with a self signed TLS certificate
func (b *GrpcClientBuilder) WithSelfSignedTLSCert() {
	b.options = append(b.options, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cert.CertPool, "")))
}

// WithKeepAliveParams set the keep alive params
// ClientParameters is used to set keepalive parameters on the client-side.
// These configure how the client will actively probe to notice when a
// connection is broken and send pings so intermediaries will be aware of the
// liveness of the connection. Make sure these parameters are set in
// coordination with the keepalive policy on the server, as incompatible
// settings can result in closing of connection.
func (b *GrpcClientBuilder) WithKeepAliveParams(params keepalive.ClientParameters) {
	keepAlive := grpc.WithKeepaliveParams(params)
	b.options = append(b.options, keepAlive)
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
func (b *GrpcClientBuilder) GetConn(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, fmt.Errorf("target connection parameter missing. address = %s", addr)
	}
	log.Debugf("Target to connect = %s", addr)
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cc, err := grpc.DialContext(ctx, addr, b.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to client. address = %s. error = %+v", addr, err)
	}
	return cc, nil
}
