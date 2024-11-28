package https

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pigeonligh/srp/pkg/proxy/providers"
)

type Server struct {
	h        http.Handler
	certFile string
	keyFile  string

	socketFile string
}

func NewServer(h http.Handler, certFile, keyFile string) *Server {
	return &Server{
		h:        h,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (s *Server) Run(ctx context.Context) error {
	if s.socketFile == "" {
		dir, err := os.MkdirTemp("", "srp-https")
		if err != nil {
			return err
		}
		s.socketFile = filepath.Join(dir, "https.sock")
	}
	l, err := net.Listen("unix", s.socketFile)
	if err != nil {
		return err
	}
	return http.ServeTLS(l, s.h, s.certFile, s.keyFile)
}

func (s *Server) SocketHandler() providers.SocketHandler {
	return providers.SocketFile(s.socketFile)
}
