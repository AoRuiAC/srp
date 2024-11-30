package providers

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pigeonligh/srp/pkg/nets"
	"github.com/pigeonligh/srp/pkg/proxy"
)

type socketProvider struct {
	h            nets.SocketHandler
	waitInterval time.Duration
}

func SocketProvider(h nets.SocketHandler, waitInterval time.Duration) proxy.ProxyProvider {
	return &socketProvider{h: h, waitInterval: waitInterval}
}

func (p *socketProvider) ProxyProvide(ctx context.Context, target string) (proxy.Proxy, error) {
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		return nil, err
	}
	socket, ok := p.h.ConvertHostPortToSocket(host, port)
	if !ok {
		return nil, fmt.Errorf("invalid target: %v", target)
	}

	ret := proxy.UnixSocket(socket)
	if p.waitInterval > 0 {
		ret = proxy.ProxyWithReadiness(ret, func(ctx context.Context) bool {
			return p.h.SocketAlive(socket)
		}, p.waitInterval)
	}
	return ret, nil
}

type SocketFile string

func (f SocketFile) ConvertHostPortToSocket(host, port string) (string, bool) {
	return string(f), true
}

func (SocketFile) SocketAlive(socket string) bool {
	stat, _ := os.Stat(socket)
	return stat != nil && !stat.IsDir()
}

type SocketNamer func(host, port string) string

func (n SocketNamer) ConvertHostPortToSocket(host, port string) (string, bool) {
	return n(host, port), true
}

func (SocketNamer) SocketAlive(socket string) bool {
	stat, _ := os.Stat(socket)
	return stat != nil && !stat.IsDir()
}
