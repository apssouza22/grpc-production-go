package grpc_client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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
	GetConn(addr string) (*grpc.ClientConn, error)
}

//GRPC client builder
type GrpcConnBuilder struct {
	options              []grpc.DialOption
	enabledReflection    bool
	shutdownHook         func()
	enabledHealthCheck   bool
	ctx                  context.Context
	transportCredentials credentials.TransportCredentials
	err                  error
}

// WithContext set the context to be used in the dial
func (b *GrpcConnBuilder) WithContext(ctx context.Context) {
	b.ctx = ctx
}

// WithOptions set dial options
func (b *GrpcConnBuilder) WithOptions(opts ...grpc.DialOption) {
	b.options = append(b.options, opts...)
}

// WithInsecure set the connection as insecure
func (b *GrpcConnBuilder) WithInsecure() {
	b.options = append(b.options, grpc.WithInsecure())
}

// WithBlock the dialing blocks until the  underlying connection is up.
// Without this, Dial returns immediately and connecting the server happens in background.
func (b *GrpcConnBuilder) WithBlock() {
	b.options = append(b.options, grpc.WithBlock())
}

// WithKeepAliveParams set the keep alive params
// ClientParameters is used to set keepalive parameters on the client-side.
// These configure how the client will actively probe to notice when a
// connection is broken and send pings so intermediaries will be aware of the
// liveness of the connection. Make sure these parameters are set in
// coordination with the keepalive policy on the server, as incompatible
// settings can result in closing of connection.
func (b *GrpcConnBuilder) WithKeepAliveParams(params keepalive.ClientParameters) {
	keepAlive := grpc.WithKeepaliveParams(params)
	b.options = append(b.options, keepAlive)
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for unary connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcConnBuilder) WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor) {
	b.options = append(b.options, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)))
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for stream connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcConnBuilder) WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor) {
	b.options = append(b.options, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...)))
}

// ClientTransportCredentials builds transport credentials for a gRPC client using the given properties.
func (b *GrpcConnBuilder) WithClientTransportCredentials(insecureSkipVerify bool, certPool *x509.CertPool) {
	var tlsConf tls.Config

	if insecureSkipVerify {
		tlsConf.InsecureSkipVerify = true
		b.transportCredentials = credentials.NewTLS(&tlsConf)
		return
	}

	tlsConf.RootCAs = certPool
	b.transportCredentials = credentials.NewTLS(&tlsConf)
}

// GetConn returns the client connection to the server
func (b *GrpcConnBuilder) GetConn(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, fmt.Errorf("target connection parameter missing. address = %s", addr)
	}
	log.Debugf("Target to connect = %s", addr)
	cc, err := grpc.DialContext(b.getContext(), addr, b.options...)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to client. address = %s. error = %+v", addr, err)
	}
	return cc, nil
}

// GetTlsConn returns client connection to the server
func (b *GrpcConnBuilder) GetTlsConn(addr string) (*grpc.ClientConn, error) {
	b.options = append(b.options, grpc.WithTransportCredentials(b.transportCredentials))
	cc, err := grpc.DialContext(
		b.getContext(),
		addr,
		b.options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tls conn. Unable to connect to client. address = %s: %w", addr, err)
	}
	return cc, nil
}

func (b *GrpcConnBuilder) getContext() context.Context {
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx
}
