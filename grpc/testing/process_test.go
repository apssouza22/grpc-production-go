//This is a example how to use a in processing GRPC server that use memory instead of network
package testing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"testing"
)

var lis *bufconn.Listener
var srv *grpc.Server

func init() {
	srv, lis = GetInProcessingGRPCServer()
	helloworld.RegisterGreeterServer(srv, &mockedService{})
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

type mockedService struct{}

func (s *mockedService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "This is a mocked service " + in.Name}, nil
}

//TestSayHello will test the HelloWorld service using A in memory data transfer instead of network
func TestSayHello(t *testing.T) {
	ctx := context.Background()
	clientConn, err := GetInProcessingClientConn(ctx, lis)
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
	srv.Stop()
	clientConn.Close()
	assert.Equal(t, resp.Message, "This is a mocked service test")
}
