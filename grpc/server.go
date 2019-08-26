package grpc

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

type Server struct {
	server            *grpc.Server
	listener          net.Listener
	options           []grpc.ServerOption
	enabledReflection bool
	shutdownHook      func()
}

func (s *Server) AddOption(o grpc.ServerOption) {
	s.options = append(s.options, o)
}

func (s *Server) EnableReflection(e bool) {
	s.enabledReflection = e
}

func (s *Server) NewServer() *grpc.Server {
	s.server = grpc.NewServer(s.options...)
	return s.server
}

func (s *Server) ListenAndServe(address string, port uint) error {
	var err error
	add := fmt.Sprintf("%s:%d", address, port)
	s.listener, err = net.Listen("tcp", add)

	if err != nil {
		msg := fmt.Sprintf("Failed to listen: %v", err)
		return errors.New(msg)
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
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
