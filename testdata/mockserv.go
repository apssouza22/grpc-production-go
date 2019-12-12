package testdata

import (
	"context"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type MockedService struct{}

func (s *MockedService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "This is a mocked service " + in.Name}, nil
}
