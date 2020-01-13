package grpc_server

import (
	"context"
	interceptors "github.com/apssouza22/grpc-server-go/serverinterceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to get hostname %v", err)
	}
	if hostname != "" {
		grpc.SendHeader(ctx, metadata.Pairs("hostname", hostname))
	}
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
		interceptors.UnaryAuditRequest(),
		interceptors.UnaryLogRequestCanceled(),
		interceptors.UnaryAuthentication(),
	}
	si := []grpc.StreamServerInterceptor{
		interceptors.StreamAuditRequest(),
		interceptors.StreamLogRequestCanceled(),
		interceptors.StreamAuthentication(),
	}
	s.SetUnaryInterceptors(ui)
	s.SetStreamInterceptors(si)
}
