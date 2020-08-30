package grpcutils

import (
	"github.com/apssouza22/grpc-production-go/clientinterceptor"
	interceptors "github.com/apssouza22/grpc-production-go/serverinterceptor"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
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
		interceptors.UnaryAuditServiceRequest(),
		interceptors.UnaryLogRequestCanceled(),
		//Recovery handlers should typically be last in the chain so that other middleware
		// (e.g. logging) can operate on the recovered state instead of being directly affected by any panic
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(requestErrorHandler)),
	}
}

// GetDefaultStreamServerInterceptors returns the default interceptors for server streams connections
func GetDefaultStreamServerInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		interceptors.StreamAuditServiceRequest(),
		interceptors.StreamLogRequestCanceled(),
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(requestErrorHandler)),
	}
}

//GetDefaultUnaryClientInterceptors returns the default interceptors for client unary connections
func GetDefaultUnaryClientInterceptors() []grpc.UnaryClientInterceptor {
	tracing := grpc_opentracing.UnaryClientInterceptor(
		grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
	)
	interceptors := []grpc.UnaryClientInterceptor{
		clientinterceptor.UnaryTimeoutInterceptor(),
		tracing,
	}
	return interceptors
}

//GetDefaultStreamClientInterceptors returns the default interceptors for client stream connections
func GetDefaultStreamClientInterceptors() []grpc.StreamClientInterceptor {
	tracing := grpc_opentracing.StreamClientInterceptor(
		grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
	)
	interceptors := []grpc.StreamClientInterceptor{
		clientinterceptor.StreamTimeoutInterceptor(),
		tracing,
	}
	return interceptors
}
