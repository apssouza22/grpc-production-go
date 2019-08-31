package grpc_server

import (
	"context"
	interceptors "github.com/apssouza22/grpc-server/grpc/server/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func ServerInitialization() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gs := Server{}
	addInterceptors(&gs)
	gs.EnableReflection(true)
	gs.EnableHealthCheck(true)
	s := gs.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	err := gs.ListenAndServe("0.0.0.0", 50051)
	if err != nil {
		log.Fatalf("%v")
	}
	gs.AddShutdownHook(func() {
		log.Print("Shutdown call")
	})
	gs.AwaitTermination()
}

func addInterceptors(s *Server) {
	ui := []grpc.UnaryServerInterceptor{
		interceptors.UnaryAuthentication(),
		interceptors.UnaryLogExecutionTime(),
		interceptors.UnaryLogRequestCanceled(),
	}
	si := []grpc.StreamServerInterceptor{
		interceptors.StreamAuthentication(),
		interceptors.StreamLogExecutionTime(),
		interceptors.StreamLogRequestCanceled(),
	}
	s.SetUnaryInterceptors(ui)
	s.SetStreamInterceptors(si)
}
