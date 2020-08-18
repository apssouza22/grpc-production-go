// In Processing server uses memory to transfer data between the server and the client
// This is ideal for testing propose as it not require networking to run integration testes
package testing

import (
	"crypto/tls"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//GRPC server interface
type GrpcInProcessingServer interface {
	Start() error
	AwaitTermination(shutdownHook func())
	RegisterService(reg func(*grpc.Server))
	Cleanup()
	GetListener() *bufconn.Listener
}

//GRPC in-processing server builder
type GrpcInProcessingServerBuilder struct {
	options []grpc.ServerOption
}

//DialOption configures how we set up the connection.
func (sb *GrpcInProcessingServerBuilder) AddOption(o grpc.ServerOption) {
	sb.options = append(sb.options, o)
}

// SetStreamInterceptors set a list of interceptors to the Grpc server for stream connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (sb *GrpcInProcessingServerBuilder) SetStreamInterceptors(interceptors []grpc.StreamServerInterceptor) {
	chain := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors...))
	sb.AddOption(chain)
}

// SetUnaryInterceptors set a list of interceptors to the Grpc server for unary connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (sb *GrpcInProcessingServerBuilder) SetUnaryInterceptors(interceptors []grpc.UnaryServerInterceptor) {
	chain := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...))
	sb.AddOption(chain)
}

// SetTlsCert sets credentials for server connections
func (sb *GrpcInProcessingServerBuilder) SetTlsCert(cert *tls.Certificate) {
	sb.AddOption(grpc.Creds(credentials.NewServerTLSFromCert(cert)))
}

//Build is responsible for building a Fiji GRPC server
func (sb *GrpcInProcessingServerBuilder) Build() GrpcInProcessingServer {
	server, listener := GetInProcessingGRPCServer(sb.options)
	return &grpcServer{server, listener}
}

type grpcServer struct {
	server   *grpc.Server
	listener *bufconn.Listener
}

// GetListener register the services to the server
func (s *grpcServer) GetListener() *bufconn.Listener {
	return s.listener
}

// RegisterService register the services to the server
func (s *grpcServer) RegisterService(reg func(*grpc.Server)) {
	reg(s.server)
}

// Start the GRPC server
func (s *grpcServer) Start() error {
	go s.serv()
	log.Printf("In processing server started")
	return nil
}

// AwaitTermination makes the program wait for the signal termination
// Valid signal termination (SIGINT, SIGTERM)
func (s *grpcServer) AwaitTermination(shutdownHook func()) {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	<-interruptSignal
	s.Cleanup()
	if shutdownHook != nil {
		shutdownHook()
	}
}

func (s *grpcServer) Cleanup() {
	s.server.Stop()
	s.listener.Close()
	log.Println("Server stopped")
}

func (s *grpcServer) serv() {
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}
