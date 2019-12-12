package util

import (
	"github.com/apssouza22/grpc-server-go/clientinterceptor"
	interceptors "github.com/apssouza22/grpc-server-go/serverinterceptor"
	"google.golang.org/grpc"
)

// GetDefaultUnaryServerInterceptors returns the default interceptors server unary connections
func GetDefaultUnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		interceptors.UnaryAuditRequest(),
		interceptors.UnaryLogRequestCanceled(),
	}
}

// GetDefaultStreamServerInterceptors returns the default interceptors for server streams connections
func GetDefaultStreamServerInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		interceptors.StreamAuditRequest(),
		interceptors.StreamLogRequestCanceled(),
	}
}

//GetDefaultUnaryClientInterceptors returns the default interceptors for client unary connections
func GetDefaultUnaryClientInterceptors() []grpc.UnaryClientInterceptor {
	interceptors := []grpc.UnaryClientInterceptor{
		clientinterceptor.UnaryTimeoutInterceptor(),
	}
	return interceptors
}

//GetDefaultStreamClientInterceptors returns the default interceptors for client stream connections
func GetDefaultStreamClientInterceptors() []grpc.StreamClientInterceptor {
	var interceptors []grpc.StreamClientInterceptor
	return interceptors
}
