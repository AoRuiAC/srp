package providers

import (
	"context"

	"github.com/pigeonligh/srp/pkg/proxy"
)

type DirectProvider string

func (p DirectProvider) ProxyProvide(ctx context.Context, target string) (proxy.Proxy, error) {
	return proxy.Direct(string(p), target), nil
}

var (
	TCPProvider  = DirectProvider("tcp")
	UDPProvider  = DirectProvider("udp")
	IPProvider   = DirectProvider("ip")
	UnixProvider = DirectProvider("unix")
)
