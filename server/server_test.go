package grpc_server

import (
	"github.com/apssouza22/grpc-server-go/tlscert"
	"github.com/apssouza22/grpc-server-go/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildGrpcServer(t *testing.T) {
	builder := &GrpcServerBuilder{}
	builder.SetTlsCert(&tlscert.Cert)
	builder.DisableDefaultHealthCheck(true)
	builder.EnableReflection(true)
	builder.SetStreamInterceptors(util.GetDefaultStreamServerInterceptors())
	builder.SetUnaryInterceptors(util.GetDefaultUnaryServerInterceptors())
	server := builder.Build()
	assert.NotNil(t, server)
}
