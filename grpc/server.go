package grpc

import (
	"errors"
	"fmt"
	_ "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server             *grpc.Server
	listener           net.Listener
	options            []grpc.ServerOption
	enabledReflection  bool
	shutdownHook       func()
	enabledHealthcheck bool
}

func (s *Server) AddOption(o grpc.ServerOption) {
	s.options = append(s.options, o)
}

func (s *Server) EnableReflection(e bool) {
	s.enabledReflection = e
}

func (s *Server) EnableHealthcheck(e bool) {
	s.enabledHealthcheck = e
}

func (s *Server) NewServer() *grpc.Server {
	s.server = grpc.NewServer(s.options...)
	return s.server
}

func (s *Server) setKeepaliveParams(duration time.Duration) {
	// MaxConnectionAge is just to avoid long connection, to facilitate load balancing
	// MaxConnectionAgeGrace will torn them, default to infinity
	keepAlive := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: duration})
	s.options = append(s.options, keepAlive)
}

func (s *Server) setStreamInterceptors(interceptors []grpc.StreamServerInterceptor) {
	chain := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors...))
	s.options = append(s.options, chain)
}

func (s *Server) setUnaryInterceptors(interceptors []grpc.UnaryServerInterceptor) {
	chain := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...))
	s.options = append(s.options, chain)
}

func (s *Server) ListenAndServe(address string, port uint) error {
	var err error
	add := fmt.Sprintf("%s:%d", address, port)
	s.listener, err = net.Listen("tcp", add)

	if err != nil {
		msg := fmt.Sprintf("Failed to listen: %v", err)
		return errors.New(msg)
	}

	if s.enabledHealthcheck {
		grpc_health_v1.RegisterHealthServer(s.server, health.NewServer())
	}

	if s.enabledReflection {
		reflection.Register(s.server)
	}
	go s.serv()

	log.Printf("Server started on port: %d \n", port)
	return nil
}

func (s *Server) AddShutdownHook(f func()) {
	s.shutdownHook = f
}

func (s *Server) AwaitTermination() {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGTERM)
	<-interruptSignal
	s.cleanup()
	if s.shutdownHook != nil {
		s.shutdownHook()
	}
}

func (s *Server) cleanup() {
	log.Println("Stopping the server")
	s.server.GracefulStop()
	log.Println("Closing the listener")
	s.listener.Close()
	log.Println("End of Program")
}

func (s *Server) serv() {
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
