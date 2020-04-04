package clientinterceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

//UnaryPropagateHeaderInterceptor copy given fields from Incoming request into Outgoing request
// Empty array will make the interceptor copy all metadata in the context
func UnaryPropagateHeaderInterceptor(fields []string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			pairs := transformMapToPairs(md, fields)
			ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

//StreamPropagateHeaderInterceptor copy given fields from Incoming request into Outgoing request
// Empty array will make the interceptor copy all metadata in the context
func StreamPropagateHeaderInterceptor(fields []string) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			pairs := transformMapToPairs(md, fields)
			ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
		}
		stream, err := streamer(ctx, desc, cc, method, opts...)
		return stream, err
	}
}

func transformMapToPairs(md map[string][]string, fields []string) []string {
	var kv []string
	for key, value := range md {
		if len(fields) > 0 && !contains(fields, key) {
			continue
		}
		for _, v := range value {
			kv = append(kv, key, v)
		}
	}
	return kv
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if strings.ToLower(x) == strings.ToLower(n) {
			return true
		}
	}
	return false
}
