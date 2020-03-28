package testing

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

func GetInProcessingClientConn(ctx context.Context, listener *bufconn.Listener, options []grpc.DialOption) (*grpc.ClientConn, error) {
	dialOptions := append(options, grpc.WithContextDialer(GetBufDialer(listener)))
	dialOptions = append(dialOptions, grpc.WithInsecure()) // Required to always set insecure for in-processing
	conn, err := grpc.DialContext(
		ctx,
		"bufconn",
		dialOptions...,
	)
	return conn, err
}

func GetInProcessingGRPCServer(options []grpc.ServerOption) (*grpc.Server, *bufconn.Listener) {
	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	srv := grpc.NewServer(options...)
	return srv, listener
}

func GetBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}
