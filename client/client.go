package grpc_client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/apssouza22/grpc-server-go/cert"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"io/ioutil"
	"net"
)

//GrpcClientConnBuilder is a builder to create GRPC connection to the GRPC Server
type GrpcClientConnBuilder interface {
	WithContext(ctx context.Context)
	WithOptions(opts ...grpc.DialOption)
	WithInsecure()
	WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor)
	WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor)
	WithKeepAliveParams(params keepalive.ClientParameters)
	GetConn(addr string, port string) (*grpc.ClientConn, error)
}

//GRPC client builder
type GrpcClientBuilder struct {
	options              []grpc.DialOption
	enabledReflection    bool
	shutdownHook         func()
	enabledHealthCheck   bool
	ctx                  context.Context
	transportCredentials credentials.TransportCredentials
	err                  error
}

// WithContext set the context to be used in the dial
func (b *GrpcClientBuilder) WithContext(ctx context.Context) {
	b.ctx = ctx
}

// WithOptions set dial options
func (b *GrpcClientBuilder) WithOptions(opts ...grpc.DialOption) {
	b.options = append(b.options, opts...)
}

// WithInsecure set the connection as insecure
func (b *GrpcClientBuilder) WithInsecure() {
	b.options = append(b.options, grpc.WithInsecure())
}

// WithTLS set the connection with a self signed TLS certificate
func (b *GrpcClientBuilder) WithSelfSignedTLSCert() {
	b.options = append(b.options, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cert.CertPool, "")))
}

// WithKeepAliveParams set the keep alive params
// ClientParameters is used to set keepalive parameters on the client-side.
// These configure how the client will actively probe to notice when a
// connection is broken and send pings so intermediaries will be aware of the
// liveness of the connection. Make sure these parameters are set in
// coordination with the keepalive policy on the server, as incompatible
// settings can result in closing of connection.
func (b *GrpcClientBuilder) WithKeepAliveParams(params keepalive.ClientParameters) {
	keepAlive := grpc.WithKeepaliveParams(params)
	b.options = append(b.options, keepAlive)
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for unary connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcClientBuilder) WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor) {
	b.options = append(b.options, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)))
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for stream connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcClientBuilder) WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor) {
	b.options = append(b.options, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...)))
}

// ClientTransportCredentials builds transport credentials for a gRPC client using the
// given properties. If cacertFile is blank, only standard trusted certs are used to
// verify the server certs. If clientCertFile is blank, the client will not use a client
// certificate. If clientCertFile is not blank then clientKeyFile must not be blank.
func (b *GrpcClientBuilder) WithClientTransportCredentials(insecureSkipVerify bool, cacertFile, clientCertFile, clientKeyFile string) {
	var tlsConf tls.Config

	if clientCertFile != "" {
		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			b.err = fmt.Errorf("could not load client key pair: %v", err)
			return
		}
		tlsConf.Certificates = []tls.Certificate{certificate}
	}

	if insecureSkipVerify {
		tlsConf.InsecureSkipVerify = true
		b.transportCredentials = credentials.NewTLS(&tlsConf)
		return
	}
	if cacertFile != "" {
		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(cacertFile)
		if err != nil {
			b.err = fmt.Errorf("could not read ca certificate: %v", err)
			return
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			b.err = errors.New("failed to append ca certs")
			return
		}

		tlsConf.RootCAs = certPool
	}

	b.transportCredentials = credentials.NewTLS(&tlsConf)
}

// GetConn returns the client connection to the server
func (b *GrpcClientBuilder) GetConn(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, fmt.Errorf("target connection parameter missing. address = %s", addr)
	}
	log.Debugf("Target to connect = %s", addr)
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cc, err := grpc.DialContext(ctx, addr, b.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to client. address = %s. error = %+v", addr, err)
	}
	return cc, nil
}

// BlockingDial is a helper method to dial the given address, using optional TLS credentials,
// and blocking until the returned connection is ready. If the given credentials are nil, the
// connection will be insecure (plain-text).
func (b *GrpcClientBuilder) GetBlockingConn(addr string) (*grpc.ClientConn, error) {
	if b.err != nil {
		return nil, fmt.Errorf("get gRPC connection failed: %w", b.err)
	}
	network := "tcp"
	dial := newBlockingDial(b.ctx, network, addr, b.transportCredentials, b.options...)
	return dial.getConn()
}

func newBlockingDial(
	ctx context.Context,
	network,
	address string,
	creds credentials.TransportCredentials,
	opts ...grpc.DialOption,
) blockingDial {
	return blockingDial{
		result:  make(chan interface{}, 1),
		ctx:     ctx,
		network: network,
		address: address,
		creds:   creds,
		opts:    opts,
	}

}

// BlockingDial is a helper method to dial the given address, using optional TLS credentials,
// and blocking until the returned connection is ready. If the given credentials are nil, the
// connection will be insecure (plain-text).
type blockingDial struct {
	result  chan interface{}
	ctx     context.Context
	network string
	address string
	creds   credentials.TransportCredentials
	opts    []grpc.DialOption
}

func (d blockingDial) resultWriter(res interface{}) {
	// non-blocking write: we only need the first result
	select {
	case d.result <- res:
	default:
	}
}

// grpc.Dial doesn't provide any information on permanent connection errors (like
// TLS handshake failures). So in order to provide good error messages, we need a
// custom dialer that can provide that info. That means we manage the TLS handshake.
func (d blockingDial) customDialer(ctx context.Context, address string) (net.Conn, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, d.network, address)
	if err != nil {
		d.resultWriter(err)
		return nil, err
	}
	if d.creds != nil {
		conn, _, err = d.creds.ClientHandshake(ctx, address, conn)
		if err != nil {
			d.resultWriter(err)
			return nil, err
		}
	}
	return conn, nil
}

func (d blockingDial) getConn() (*grpc.ClientConn, error) {
	// Even with grpc.FailOnNonTempDialError, this call will usually timeout in
	// the face of TLS handshake errors. So we can't rely on grpc.WithBlock() to
	// know when we're done. So we run it in a goroutine and then use result
	// channel to either get the channel or fail-fast.
	go d.dial()

	select {
	case res := <-d.result:
		if conn, ok := res.(*grpc.ClientConn); ok {
			return conn, nil
		}
		return nil, res.(error)
	case <-d.ctx.Done():
		return nil, d.ctx.Err()
	}
}

func (d blockingDial) dial() {
	opts := append(d.opts,
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithContextDialer(d.customDialer),
		grpc.WithInsecure(), // we are handling TLS, so tell grpc not to
	)
	conn, err := grpc.DialContext(d.ctx, d.address, opts...)
	var res interface{}
	if err != nil {
		res = err
	} else {
		res = conn
	}
	d.resultWriter(res)
}
