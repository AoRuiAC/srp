package providers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pigeonligh/srp/pkg/proxy"
)

type SocketHandler interface {
	ConvertHostPortToSocket(host, port string) (string, bool)
	SocketAlive(socket string) bool
}

type socketProvider struct {
	h            SocketHandler
	waitInterval time.Duration
}

func SocketProvider(h SocketHandler, waitInterval time.Duration) proxy.ProxyProvider {
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
