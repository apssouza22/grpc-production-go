package interceptors

import (
	"context"
	logrus "github.com/sirupsen/logrus"
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
	sts := status.Convert(err)
	auditEntry := logrus.Fields{
		"user-agent":  userAgents,
		"peer":        ip,
		"took_ns":     time.Since(start),
		"status":      sts.Code().String(),
		"err":         sts.Message(),
		"err-details": sts.Details(),
	}
	log := logrus.WithFields(auditEntry)

	switch sts.Code() {

	case codes.OK:
		log.Info("gRPC call succeeded")

	// Caused by invalid client requests (http 4xx equiv.)
	case codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented, // usually caused by client requesting invalid operation (even though it matches 501)
		codes.Unauthenticated:

		log.Warn("gRPC call failed")

	// Server errors (http 5xx equiv.):
	// Unknown, DeadlineExceeded, ResourceExhausted, Internal, Unavailable, DataLoss
	// (ResourceExhausted is somewhere in between, from user quota exhausted to OOM, rather have it on error for now)
	default:
		log.Error("gRPC call errored")
	}
}

func isHealthCheckRequest(requestMethod string) bool {
	if strings.Contains(requestMethod, "Health/Check") {
		return true
	}
	return false
}
