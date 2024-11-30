package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pigeonligh/srp/pkg/nets"
	"github.com/pigeonligh/srp/pkg/proxy"
	"github.com/pigeonligh/srp/pkg/proxy/providers"
)

type HTTP struct {
	Network  string
	Address  string
	Listener net.Listener
	Handler  http.Handler
}

func (s *HTTP) listen() (net.Listener, error) {
	if s.Listener != nil {
		return s.Listener, nil
	}

	if s.Network == "unix" && s.Address == "" {
		dir, err := os.MkdirTemp("", "srp-http")
		if err != nil {
			return nil, err
		}
		s.Address = filepath.Join(dir, "server.sock")
	}

	return net.Listen(s.Network, s.Address)
}

func (s *HTTP) Run(ctx context.Context) error {
	l, err := s.listen()
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler: s.Handler,
	}

	return nets.RunNetServer(ctx, server, l)
}

func (s *HTTP) Provider() proxy.ProxyProvider {
	if s.Network == "unix" {
		return providers.SocketProvider(providers.SocketFile(s.Address), 0)
	}

	return proxy.ProxyProviderFunc(func(ctx context.Context, target string) (proxy.Proxy, error) {
		return proxy.Direct(s.Network, s.Address), nil
	})
}
