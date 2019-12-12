package clientinterceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

//UnaryTimeoutInterceptor monitor the DeadlineExceeded error and log it
func UnaryTimeoutInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		return handleError(err, method, start)
	}
}

//StreamTimeoutInterceptor monitor the DeadlineExceeded error and log it
func StreamTimeoutInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()
		stream, err := streamer(ctx, desc, cc, method, opts...)
		err = handleError(err, method, start)
		return stream, err
	}
}

func handleError(err error, method string, start time.Time) error {
	if err == nil {
		return err
	}
	statusErr, ok := status.FromError(err)
	if !ok {
		return err
	}
	if statusErr.Code() != codes.DeadlineExceeded {
		return err
	}
	log.Printf(
		"Timeout - Invoked RPC method=%s; Duration=%s; Error=%+v",
		method,
		time.Since(start), err,
	)
	return err
}
