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
	"strings"
	"time"
)

var healthCheckMethodName = "/grpc.health.v1.Health/Check"

// SetHealthCheckMethodName changes the default health check method name
func SetHealthCheckMethodName(methodName string) {
	healthCheckMethodName = methodName
}

// Logging request information for Unary requests
func UnaryAuditServiceRequest() grpc.UnaryServerInterceptor {
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
		logRequest(
			start,
			info.FullMethod,
			md["user-agent"],
			peer.Addr,
			info.FullMethod,
			err,
		)

		return resp, err // passing up the chain the response and the err
	}
}

// Logging request information for Stream requests
func StreamAuditServiceRequest() grpc.StreamServerInterceptor {
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
		logRequest(
			start,
			info.FullMethod,
			md["user-agent"],
			peer.Addr,
			info.FullMethod,
			err,
		)
		return err
	}
}

func logRequest(start time.Time, requestMethod string, userAgents []string, ip net.Addr, fullMethod string, err error) {
	if isHealthCheckRequest(requestMethod) {
		return
	}
	status := "OK"
	if err != nil {
		status = "KO"
	}
	auditEntry := log.Fields{
		"user-agent": userAgents,
		"peer":       ip,
		"took_ns":    time.Since(start),
		"status":     status,
		"err":        err,
	}
	if err != nil {
		log.WithFields(auditEntry).Error(fullMethod)
	} else {
		log.WithFields(auditEntry).Info(fullMethod)
	}
}

func isHealthCheckRequest(requestMethod string) bool {
	if strings.Contains(requestMethod, "Health/Check") {
		return true
	}
	return false
}
