package main

import (
	grpcserver "github.com/apssouza22/grpc-production-go/server"
)

func main() {
	grpcserver.ServerInitialization()
	//grpcserver.ServerInitializationWithTLS()
}
