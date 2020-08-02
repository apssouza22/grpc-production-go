package grpc_server

import (
	"errors"
	"fmt"
	"github.com/apssouza22/grpc-server-go/cert"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

//Fiji GRPC server interface
type GrpcServer interface {
	Start(address string, port uint) error
	AwaitTermination(shutdownHook func())
	RegisterService(reg func(*grpc.Server))
}

//GRPC server builder
type GrpcServerBuilder struct {
	options                   []grpc.ServerOption
	enabledReflection         bool
	shutdownHook              func()
	enabledHealthCheck        bool
	disableDefaultHealthCheck bool
}

type grpcServer struct {
	server   *grpc.Server
	listener net.Listener
}

//DialOption configures how we set up the connection.
func (sb *GrpcServerBuilder) AddOption(o grpc.ServerOption) {
	sb.options = append(sb.options, o)
}

// EnableReflection enables the reflection
// gRPC Server Reflection provides information about publicly-accessible gRPC services on a server,
// and assists clients at runtime to construct RPC requests and responses without precompiled service information.
// It is used by gRPC CLI, which can be used to introspect server protos and send/receive test RPCs.
//Warning! We should not have this enabled in production
func (sb *GrpcServerBuilder) EnableReflection(e bool) {
	sb.enabledReflection = e
}

// DisableDefaultHealthCheck disables the default health check service
//Warning! if you disable the default health check you must provide a custom health check service
func (sb *GrpcServerBuilder) DisableDefaultHealthCheck(e bool) {
	sb.disableDefaultHealthCheck = e
}

// ServerParameters is used to set keepalive and max-age parameters on the server-side.
func (sb *GrpcServerBuilder) SetServerParameters(serverParams keepalive.ServerParameters) {
	keepAlive := grpc.KeepaliveParams(serverParams)
	sb.AddOption(keepAlive)
}

// SetStreamInterceptors set a list of interceptors to the Grpc server for stream connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (sb *GrpcServerBuilder) SetStreamInterceptors(interceptors []grpc.StreamServerInterceptor) {
	chain := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors...))
	sb.AddOption(chain)
}

// SetUnaryInterceptors set a list of interceptors to the Grpc server for unary connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (sb *GrpcServerBuilder) SetUnaryInterceptors(interceptors []grpc.UnaryServerInterceptor) {
	chain := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...))
	sb.AddOption(chain)
}

// SetSelfSignedTLS sets credentials for server connections
func (sb *GrpcServerBuilder) SetSelfSignedTLS() {
	sb.AddOption(grpc.Creds(credentials.NewServerTLSFromCert(&cert.Cert)))
}

//Build is responsible for building a Fiji GRPC server
func (sb *GrpcServerBuilder) Build() GrpcServer {
	srv := grpc.NewServer(sb.options...)
	if !sb.disableDefaultHealthCheck {
		grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	}

	if sb.enabledReflection {
		reflection.Register(srv)
	}
	return &grpcServer{srv, nil}
}

// RegisterService register the services to the server
func (s grpcServer) RegisterService(reg func(*grpc.Server)) {
	reg(s.server)
}

// Start the GRPC server
func (s *grpcServer) Start(address string, port uint) error {
	var err error
	add := fmt.Sprintf("%s:%d", address, port)
	s.listener, err = net.Listen("tcp", add)

	if err != nil {
		msg := fmt.Sprintf("Failed to listen: %v", err)
		return errors.New(msg)
	}

	go s.serv()

	log.Infof("grpcServer started on port: %d ", port)
	return nil
}

// AwaitTermination makes the program wait for the signal termination
// Valid signal termination (SIGINT, SIGTERM)
func (s *grpcServer) AwaitTermination(shutdownHook func()) {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	<-interruptSignal
	s.cleanup()
	if shutdownHook != nil {
		shutdownHook()
	}
}

func (s *grpcServer) cleanup() {
	log.Info("Stopping the server")
	s.server.GracefulStop()
	log.Info("Closing the listener")
	s.listener.Close()
	log.Info("End of Program")
}

func (s *grpcServer) serv() {
	if err := s.server.Serve(s.listener); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}
