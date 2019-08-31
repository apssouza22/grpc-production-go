package interceptors

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryAuthentication() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(exampleAuthFunc)
}

func StreamAuthentication() grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(exampleAuthFunc)
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
