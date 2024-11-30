package nets

import (
	"context"
	"net"
)

type TLSServer interface {
	ServeTLS(l net.Listener, certFile, keyFile string) error
	ListenAndServeTLS(certFile, keyFile string) error
	Shutdown(ctx context.Context) error
}

type tlsServer struct {
	s        TLSServer
	certFile string
	keyFile  string
}

func (s *tlsServer) Serve(l net.Listener) error {
	return s.s.ServeTLS(l, s.certFile, s.keyFile)
}

func (s *tlsServer) ListenAndServe() error {
	return s.s.ListenAndServeTLS(s.certFile, s.keyFile)
}

func (s *tlsServer) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}

func WrapTLSServer(s TLSServer, certFile, keyFile string) NetServer {
	return &tlsServer{s, certFile, keyFile}
}
