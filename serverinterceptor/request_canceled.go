package interceptors

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

// Log the request that has been cancelled by the client during the Unary request
// The request can be cancelled for many reasons, including timeout exceeded
func UnaryLogRequestCanceled() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (_ interface{}, err error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		if ctx.Err() == context.Canceled {
			logCanceledRequest(start, err, info.FullMethod)
		}
		return resp, err
	}
}

// Log the request that has been cancelled by the client during the Stream request
// The request can be cancelled for many reasons, including timeout exceeded
func StreamLogRequestCanceled() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		start := time.Now()
		err = handler(srv, stream)

		if stream.Context().Err() == context.Canceled {
			logCanceledRequest(start, err, info.FullMethod)
		}
		return err
	}
}

func logCanceledRequest(start time.Time, err error, method string) {
	auditEntry := log.Fields{
		"took_ns": time.Since(start),
		"status":  "Request Canceled",
		"err":     err,
	}
	log.WithFields(auditEntry).Warn(method)
}
