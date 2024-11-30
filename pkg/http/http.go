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

func (s *HTTP) networkAddress() (string, string) {
	if s.Network == "unix" && s.Address == "" {
		dir, _ := os.MkdirTemp("", "srp-http")
		s.Address = filepath.Join(dir, "server.sock")
	}
	if s.Network == "" {
		s.Network = "tcp"
	}
	return s.Network, s.Address
}

func (s *HTTP) listen() (net.Listener, error) {
	if s.Listener != nil {
		return s.Listener, nil
	}
	network, address := s.networkAddress()
	return net.Listen(network, address)
}

func (s *HTTP) Run(ctx context.Context) error {
	l, err := s.listen()
	if err != nil {
		return err
	}
	server := &http.Server{
		Handler: s.Handler,
	}
	ctx = nets.ContextWithServerName(ctx, "HTTP["+s.Address+"]")
	return nets.RunNetServer(ctx, server, l)
}

func (s *HTTP) Provider() proxy.ProxyProvider {
	network, address := s.networkAddress()
	if network == "unix" {
		return providers.SocketProvider(providers.SocketFile(address), 0)
	}
	return proxy.ProxyProviderFunc(func(ctx context.Context, target string) (proxy.Proxy, error) {
		return proxy.Direct(network, address), nil
	})
}
