package grpc_server

import (
	"context"
	interceptors "github.com/apssouza22/grpc-server-go/grpc/server/interceptor"
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

	builder := GrpcServerBuilder{}
	addInterceptors(&builder)
	builder.EnableReflection(true)
	s := builder.Build()
	s.RegisterService(serviceRegister)
	err := s.Start("0.0.0.0", 50051)
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.AwaitTermination(func() {
		log.Print("Shutting down the server")
	})
}

func serviceRegister(sv *grpc.Server) {
	helloworld.RegisterGreeterServer(sv, &server{})
}

func addInterceptors(s *GrpcServerBuilder) {
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
