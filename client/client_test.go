package grpc_client

import (
	"context"
	"github.com/apssouza22/grpc-server-go/testdata"
	gtest "github.com/apssouza22/grpc-server-go/testing"
	"github.com/apssouza22/grpc-server-go/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"testing"
)

var server gtest.GrpcInProcessingServer

func startServer() {
	builder := gtest.GrpcInProcessingServerBuilder{}
	builder.SetUnaryInterceptors(util.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	server.Start()
}

func TestSayHelloPassingContext(t *testing.T) {
	startServer()
	ctx := context.Background()
	clientBuilder := GrpcClientBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithContext(ctx)
	clientBuilder.WithOptions(grpc.WithContextDialer(gtest.GetBufDialer(server.GetListener())))
	clientConn, err := clientBuilder.GetConn("localhost:50051")

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	server.Cleanup()
	clientConn.Close()
	assert.Equal(t, resp.Message, "This is a mocked service test")
}

func TestSayHelloNotPassingContext(t *testing.T) {
	startServer()
	ctx := context.Background()
	clientBuilder := GrpcClientBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithOptions(grpc.WithContextDialer(gtest.GetBufDialer(server.GetListener())))
	clientConn, err := clientBuilder.GetConn("localhost:50051")

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	server.Cleanup()
	clientConn.Close()
	assert.Equal(t, resp.Message, "This is a mocked service test")
}
