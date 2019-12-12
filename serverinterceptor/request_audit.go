package interceptors

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

// Logging request information for Unary requests
func UnaryAuditRequest() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Extract all needed info for audit from the RPC call
		peer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing peer info")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		mkAuditEntry(
			start,
			md["authority"],
			md["content-type"],
			md["user-agent"],
			peer.Addr,
			info.FullMethod,
			err,
		)

		return resp, err // passing up the chain the response and the err
	}
}

// Logging request information for Stream requests
func StreamAuditRequest() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) (err error) {
		peer, ok := peer.FromContext(stream.Context())
		if !ok {
			return status.Errorf(codes.InvalidArgument, "missing peer info")
		}
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		start := time.Now()
		err = handler(srv, stream)
		mkAuditEntry(
			start,
			md["authority"],
			md["content-type"],
			md["user-agent"],
			peer.Addr,
			info.FullMethod,
			err,
		)
		return err
	}
}

func mkAuditEntry(
	start time.Time,
	authorities []string,
	contentTypes []string,
	userAgents []string,
	ip net.Addr,
	fullMethod string,
	err error,
) {
	status := "OK"
	if err != nil {
		status = "KO"
	}
	auditEntry := log.Fields{
		"authority":    authorities,
		"content-type": contentTypes,
		"user-agent":   userAgents,
		"peer":         ip,
		"took_ns":      time.Since(start),
		"status":       status,
		"err":          err,
	}
	if err != nil {
		log.WithFields(auditEntry).Error(fullMethod)
	} else {
		log.WithFields(auditEntry).Info(fullMethod)
	}
}
