package util

import (
	"github.com/apssouza22/grpc-server-go/clientinterceptor"
	interceptors "github.com/apssouza22/grpc-server-go/serverinterceptor"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func requestErrorHandler(p interface{}) (err error) {
	logrus.Error(p)
	return status.Errorf(codes.Internal, "Something went wrong :( ")
}

// GetDefaultUnaryServerInterceptors returns the default interceptors server unary connections
func GetDefaultUnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		interceptors.UnaryAuditRequest(),
		interceptors.UnaryLogRequestCanceled(),
		//Recovery handlers should typically be last in the chain so that other middleware
		// (e.g. logging) can operate on the recovered state instead of being directly affected by any panic
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(requestErrorHandler)),
	}
}

// GetDefaultStreamServerInterceptors returns the default interceptors for server streams connections
func GetDefaultStreamServerInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		interceptors.StreamAuditRequest(),
		interceptors.StreamLogRequestCanceled(),
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(requestErrorHandler)),
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
