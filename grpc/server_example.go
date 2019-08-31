package grpc_server

import (
	"context"
	execution_time "git.deem.com/fijigroup/shared/fiji-grpc-core-library/grpc/interceptor"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
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
		grpc_auth.UnaryServerInterceptor(exampleAuthFunc),
		execution_time.UnaryLogExecutionTime(),
	}
	si := []grpc.StreamServerInterceptor{
		grpc_auth.StreamServerInterceptor(exampleAuthFunc),
		execution_time.StreamLogExecutionTime(),
	}
	s.setUnaryInterceptors(ui)
	s.setStreamInterceptors(si)
}

func exampleAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	if token != "123" {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}
	grpc_ctxtags.Extract(ctx).Set("auth.sub", "info")

	type authInfo struct {
		name string
	}

	newCtx := context.WithValue(ctx, "tokenInfo", authInfo{"foo"})
	return newCtx, nil
}
