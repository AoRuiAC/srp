package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pigeonligh/srp/pkg/nets"
)

type HTTPS struct {
	HTTP
	CertFile    string
	KeyFile     string
	TLSProvider func() (string, string, error)
}

func (s *HTTPS) tls() (string, string, error) {
	if s.TLSProvider != nil {
		return s.TLSProvider()
	}
	if s.CertFile == "" || s.KeyFile == "" {
		return "", "", fmt.Errorf("cert not found")
	}
	return s.CertFile, s.KeyFile, nil
}

func (s *HTTPS) Run(ctx context.Context) error {
	l, err := s.listen()
	if err != nil {
		return err
	}

	cert, key, err := s.tls()
	if err != nil {
		return err
	}

	server := nets.WrapTLSServer(&http.Server{
		Handler: s.Handler,
	}, cert, key)

	return nets.RunNetServer(ctx, server, l)
}
