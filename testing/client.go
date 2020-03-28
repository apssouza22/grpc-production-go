package testing

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type InProcessingClientBuilder struct {
	Server  GrpcInProcessingServer
	options []grpc.DialOption
	ctx     context.Context
}

// WithContext set the context to be used in the dial
func (b *InProcessingClientBuilder) WithContext(ctx context.Context) {
	b.ctx = ctx
}

// WithOptions set dial options
func (b *InProcessingClientBuilder) WithOptions(opts ...grpc.DialOption) {
	b.options = append(b.options, opts...)
}

// WithInsecure set the connection as insecure
func (b *InProcessingClientBuilder) WithInsecure() {
	b.options = append(b.options, grpc.WithInsecure())
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for unary connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the Server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *InProcessingClientBuilder) WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor) {
	b.options = append(b.options, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)))
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for stream connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the Server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *InProcessingClientBuilder) WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor) {
	b.options = append(b.options, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...)))
}

// GetConn returns the client connection to the Server
func (b *InProcessingClientBuilder) GetConn(addr string, port string) (*grpc.ClientConn, error) {
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return GetInProcessingClientConn(context.Background(), b.Server.GetListener(), b.options)
}
