package main

import (
	"../../pkg/grpc"
	"context"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
)
type server struct {}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gs := grpc.Server{}
	gs.EnableReflection(true)
	s := gs.NewServer()
	helloworld.RegisterGreeterServer(s,  &server{})
	err := gs.ListenAndServe("0.0.0.0", 50051)
	if err != nil {
		log.Fatalf("%v")
	}
	gs.AddShutdownHook(func() {
		log.Print("Shutdown call")
	})
	gs.AwaitTermination()
}
