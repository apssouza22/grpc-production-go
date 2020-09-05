package grpc_server

import (
	"context"
	"github.com/apssouza22/grpc-production-go/grpcutils"
	"github.com/apssouza22/grpc-production-go/tlscert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	md, _ := metadata.FromIncomingContext(ctx)
	log.Print(md)
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
	err := s.Start("0.0.0.0:50051")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.AwaitTermination(func() {
		log.Print("Shutting down the server")
	})
}

func ServerInitializationWithTLS() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	builder := GrpcServerBuilder{}
	addInterceptors(&builder)
	builder.EnableReflection(true)

	builder.SetTlsCert(&tlscert.Cert)

	s := builder.Build()
	s.RegisterService(serviceRegister)
	err := s.Start("0.0.0.0:50051")
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
	s.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	s.SetStreamInterceptors(grpcutils.GetDefaultStreamServerInterceptors())
}
