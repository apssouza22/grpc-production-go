package testing

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

func GetInProcessingClientConn(ctx context.Context, listener *bufconn.Listener) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(
		ctx,
		"bufconn",
		grpc.WithContextDialer(GetBufDialer(listener)),
		grpc.WithInsecure(),
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
