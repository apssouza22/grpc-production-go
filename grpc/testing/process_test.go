//This is a example how to use a in processing GRPC server that use memory instead of network
package testing

import (
	"context"
	interceptors "github.com/apssouza22/grpc-server-go/grpc/server/interceptor"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"testing"
)

var server GrpcServer

func init() {
	builder := GrpcServerBuilder{}
	builder.SetUnaryInterceptors(interceptors.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &mockedService{})
	})
	server.Start()
}

type mockedService struct{}

func (s *mockedService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "This is a mocked service " + in.Name}, nil
}

//TestSayHello will test the HelloWorld service using A in memory data transfer instead of network
func TestSayHello(t *testing.T) {
	ctx := context.Background()
	clientConn, err := GetInProcessingClientConn(ctx, server.GetListener())
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
