package main

import (
	"context"
	grpc_server "git.deem.com/fijigroup/shared/fiji-grpc-core-library/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	grpc_server.ServerInitialization()
}
