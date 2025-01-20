package proxy

import (
	"context"
	"net"
	"time"

	"github.com/pigeonligh/srp/pkg/nets"
)

type Proxy interface {
	Dial(ctx context.Context) (net.Conn, error)
}

type directProxy struct {
	network string
	address string
	dialer  nets.NetDialer
}

func (p directProxy) Dial(ctx context.Context) (net.Conn, error) {
	return p.dialer.DialContext(ctx, p.network, p.address)
}

func Direct(network string, address string) Proxy {
	return directProxy{network: network, address: address, dialer: nets.DefaultNetDialer}
}

func DirectWithDialer(network string, address string, dialer nets.NetDialer) Proxy {
	return directProxy{network: network, address: address, dialer: dialer}
}

func UnixSocket(socket string) Proxy {
	return Direct("unix", socket)
}

type funcProxy func(ctx context.Context) (net.Conn, error)

func (f funcProxy) Dial(ctx context.Context) (net.Conn, error) {
	return f(ctx)
}

func ProxyWithTimeout(p Proxy, timeout time.Duration) Proxy {
	return funcProxy(func(ctx context.Context) (net.Conn, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return p.Dial(ctx)
	})
}

func ProxyWithReadiness(p Proxy, readiness func(context.Context) bool, interval time.Duration) Proxy {
	wait := func(ctx context.Context) error {
		if !readiness(ctx) {
			t := time.NewTicker(interval)
		LOOP:
			for {
				select {
				case <-ctx.Done():
					t.Stop()
					return context.Canceled

				case <-t.C:
					if readiness(ctx) {
						t.Stop()
						break LOOP
					}
				}

			}
		}
		return nil
	}

	return funcProxy(func(ctx context.Context) (net.Conn, error) {
		if err := wait(ctx); err != nil {
			return nil, err
		}
		return p.Dial(ctx)
	})
}
