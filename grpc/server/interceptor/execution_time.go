package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func UnaryLogExecutionTime() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (_ interface{}, err error) {
		start := time.Now()
		i, err := handler(ctx, req)
		log.Printf(
			"Request - Method:%s\tDuration:%s\tError:%v\n",
			info.FullMethod,
			time.Since(start),
			err,
		)
		return i, err
	}
}
func StreamLogExecutionTime() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		start := time.Now()
		err = handler(srv, stream)

		log.Printf(
			"Request - Method:%s\tDuration:%s\tError:%v\n",
			info.FullMethod,
			time.Since(start),
			err,
		)
		return err
	}
}
