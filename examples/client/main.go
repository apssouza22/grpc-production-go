package main

import (
	grpcclient "github.com/apssouza22/grpc-production-go/client"
)

func main() {
	grpcclient.TimeoutLogExample()
	//grpcclient.TLSConnExample()
}
