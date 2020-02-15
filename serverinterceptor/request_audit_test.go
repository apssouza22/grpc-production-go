package interceptors

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"testing"
)

func TestStreamAuditServiceRequest(t *testing.T) {
	interceptor := StreamAuditServiceRequest()
	ctx := context.Background()
	addr := net.IPNet{}
	ctx = peer.NewContext(ctx, &peer.Peer{Addr: &addr})
	md := metadata.Pairs("user", "user", "pass", "123")
	ctx = metadata.NewIncomingContext(ctx, md)
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}
	err := interceptor(ctx, ServerStreamMock{}, &grpc.StreamServerInfo{
		FullMethod: "testMethod",
	}, handler)
	assert.NoError(t, err)
}

func TestUnaryAuditServiceRequest(t *testing.T) {
	interceptor := UnaryAuditServiceRequest()
	ctx := context.Background()
	addr := net.IPNet{}
	ctx = peer.NewContext(ctx, &peer.Peer{Addr: &addr})
	md := metadata.Pairs("user", "user", "pass", "123")
	ctx = metadata.NewIncomingContext(ctx, md)
	handler := func(ctx context.Context, req interface{}) (i interface{}, e error) {
		return nil, nil
	}
	info := &grpc.UnaryServerInfo{
		Server:     nil,
		FullMethod: "test",
	}
	_, e := interceptor(ctx, "test", info, handler)
	assert.NoError(t, e)
}

func Test_isHealthCheckRequest(t *testing.T) {
	SetHealthCheckMethodName("Health/Check")
	a := isHealthCheckRequest("Health/Check")
	assert.True(t, a)

	b := isHealthCheckRequest("Other/Method")
	assert.False(t, b)
}

type ServerStreamMock struct {
	grpc.ServerStream
}

func (s ServerStreamMock) Context() context.Context {
	ctx := context.Background()
	addr := net.IPNet{}
	ctx = peer.NewContext(ctx, &peer.Peer{Addr: &addr})
	md := metadata.Pairs("user", "user", "pass", "123")
	ctx = metadata.NewIncomingContext(ctx, md)
	return ctx
}
