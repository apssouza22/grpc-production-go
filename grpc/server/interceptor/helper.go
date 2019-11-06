package interceptors

import (
	"github.com/apssouza22/grpc-server-go/grpc/client/interceptor"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func GetDefaultUnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		UnaryLogExecutionTime(),
		UnaryLogRequestCanceled(),
	}
}

func GetDefaultStreamServerInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		StreamLogExecutionTime(),
		StreamLogRequestCanceled(),
	}
}

func GetDefaultUnaryClientInterceptors() grpc.DialOption {
	interceptors := []grpc.UnaryClientInterceptor{
		interceptor.UnaryTimeoutInterceptor(),
	}
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...))
}

func GetDefaultStreamClientInterceptors() grpc.DialOption {
	interceptors := []grpc.StreamClientInterceptor{}
	return grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...))
}
