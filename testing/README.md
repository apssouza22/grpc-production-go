This package is bilt using  the google.golang.org/grpc/test/bufconn package to help you avoid starting up a service 
with a real port number, but still allowing testing of streaming RPCs.

The benefit of this approach is that you're still getting network behavior, but over an in-memory connection without 
using OS-level resources like ports that may or may not clean up quickly. And it allows you to test it the way 
it's actually used, and it gives you proper streaming behavior.

We provide an In memory communication between client and server, helpful to write unit and integration tests. 
When writing integration tests we should avoid having the networking element from your test as it is slow to assign and release ports

```
var server GrpcInProcessingServer

func serverStart() {
	builder := GrpcInProcessingServerBuilder{}
	builder.SetUnaryInterceptors(util.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	server.Start()
}

//TestSayHello will test the HelloWorld service using A in memory data transfer instead of the normal networking
func TestSayHello(t *testing.T) {
	serverStart()
	ctx := context.Background()
	clientConn, err := GetInProcessingClientConn(ctx, server.GetListener(), []grpc.DialOption{})
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	server.Cleanup()
	clientConn.Close()
	assert.Equal(t, resp.Message, "This is a mocked service test")
}
```