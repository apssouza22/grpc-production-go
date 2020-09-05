package grpc_client

import (
	"context"
	"github.com/apssouza22/grpc-production-go/grpcutils"
	grpc_server "github.com/apssouza22/grpc-production-go/server"
	"github.com/apssouza22/grpc-production-go/testdata"
	gtest "github.com/apssouza22/grpc-production-go/testing"
	"github.com/apssouza22/grpc-production-go/tlscert"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"testing"
)

var server gtest.GrpcInProcessingServer

func startServer() {
	builder := gtest.GrpcInProcessingServerBuilder{}
	builder.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	server.Start()
}
func startServerWithTLS() grpc_server.GrpcServer {
	builder := grpc_server.GrpcServerBuilder{}
	builder.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	builder.SetTlsCert(&tlscert.Cert)
	svr := builder.Build()
	svr.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	svr.Start("localhost", 8989)
	return svr
}

func TestSayHelloPassingContext(t *testing.T) {
	startServer()
	ctx := context.Background()
	clientBuilder := GrpcConnBuilder{}
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
	clientBuilder := GrpcConnBuilder{}
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

func TestTLSConnWithCert(t *testing.T) {
	serverWithTLS := startServerWithTLS()
	defer serverWithTLS.GetListener().Close()

	ctx := context.Background()
	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithContext(ctx)
	clientBuilder.WithBlock()
	clientBuilder.WithClientTransportCredentials(false, tlscert.CertPool)
	clientConn, _ := clientBuilder.GetTlsConn("localhost:8989")
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, resp.Message, "This is a mocked service test")
}

func TestTLSConnWithInsecure(t *testing.T) {
	serverWithTLS := startServerWithTLS()
	defer serverWithTLS.GetListener().Close()

	ctx := context.Background()
	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithContext(ctx)
	clientBuilder.WithBlock()
	clientBuilder.WithClientTransportCredentials(true, nil)
	clientConn, _ := clientBuilder.GetTlsConn("localhost:8989")
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, resp.Message, "This is a mocked service test")
}
