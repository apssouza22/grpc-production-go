package main

import (
	"context"
	grpc_server "github.com/apssouza22/grpc-server/grpc"
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
