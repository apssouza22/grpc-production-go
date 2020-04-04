package clientinterceptor

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestUnaryPropagateHeaderInterceptor(t *testing.T) {
	fields := []string{"traceId", "clientID", "session-ID"}
	md := make(map[string]string)
	md["traceId"] = "123"
	md["clientID"] = "453"
	md["session-ID"] = "session"
	interceptor := UnaryPropagateHeaderInterceptor(fields)
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.New(md))
	rpc := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		mds, _ := metadata.FromOutgoingContext(ctx)
		assert.Equal(t, "123", mds["traceid"][0])
		assert.Equal(t, "453", mds["clientid"][0])
		assert.Equal(t, "session", mds["session-id"][0])
		return nil
	}
	interceptor(ctx, "test", "req", "reply", nil, rpc)
}

func TestUnaryPropagateAllHeaderInterceptor(t *testing.T) {
	fields := []string{}
	md := make(map[string]string)
	md["traceId"] = "123"
	md["clientID"] = "453"
	md["session-ID"] = "session"
	interceptor := UnaryPropagateHeaderInterceptor(fields)
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.New(md))
	rpc := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		mds, _ := metadata.FromOutgoingContext(ctx)
		assert.Equal(t, "123", mds["traceid"][0])
		assert.Equal(t, "453", mds["clientid"][0])
		assert.Equal(t, "session", mds["session-id"][0])
		return nil
	}
	interceptor(ctx, "test", "req", "reply", nil, rpc)
}

func TestStreamPropagateHeaderInterceptor(t *testing.T) {
	fields := []string{"traceId", "clientID", "session-ID"}
	md := make(map[string]string)
	md["traceId"] = "123"
	md["clientID"] = "453"
	md["session-ID"] = "session"
	interceptor := StreamPropagateHeaderInterceptor(fields)
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.New(md))

	rpc := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (stream grpc.ClientStream, e error) {
		mds, _ := metadata.FromOutgoingContext(ctx)
		assert.Equal(t, "123", mds["traceid"][0])
		assert.Equal(t, "453", mds["clientid"][0])
		assert.Equal(t, "session", mds["session-id"][0])
		return nil, nil
	}
	interceptor(ctx, &grpc.StreamDesc{StreamName: "test"}, nil, "test", rpc)
}
