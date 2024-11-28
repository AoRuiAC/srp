package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type Proxy interface {
	Dial(ctx context.Context) (net.Conn, error)
	Proxy(ctx context.Context, r io.Reader, w io.Writer) error
}

type directProxy struct {
	network string
	address string
}

func (p directProxy) Dial(ctx context.Context) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, p.network, p.address)
}

func (p directProxy) Proxy(ctx context.Context, r io.Reader, w io.Writer) error {
	c, err := p.Dial(ctx)
	if err != nil {
		return fmt.Errorf("connecting to %v://%v: %w", p.network, p.address, err)
	}
	defer c.Close()

	go func() {
		_, _ = io.Copy(c, r)
	}()
	_, _ = io.Copy(w, c)
	return nil
}

func Direct(network string, address string) Proxy {
	return directProxy{network: network, address: address}
}

func UnixSocket(socket string) Proxy {
	return Direct("unix", socket)
}

type funcProxy struct {
	dial  func(ctx context.Context) (net.Conn, error)
	proxy func(ctx context.Context, r io.Reader, w io.Writer) error
}

func (p funcProxy) Dial(ctx context.Context) (net.Conn, error) {
	return p.dial(ctx)
}

func (p funcProxy) Proxy(ctx context.Context, r io.Reader, w io.Writer) error {
	return p.proxy(ctx, r, w)
}

func ProxyWithTimeout(p Proxy, timeout time.Duration) Proxy {
	return funcProxy{
		dial: func(ctx context.Context) (net.Conn, error) {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return p.Dial(ctx)
		},
		proxy: func(ctx context.Context, r io.Reader, w io.Writer) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return p.Proxy(ctx, r, w)
		},
	}
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

	return funcProxy{
		dial: func(ctx context.Context) (net.Conn, error) {
			if err := wait(ctx); err != nil {
				return nil, err
			}
			return p.Dial(ctx)
		},
		proxy: func(ctx context.Context, r io.Reader, w io.Writer) error {
			if err := wait(ctx); err != nil {
				return err
			}
			return p.Proxy(ctx, r, w)
		},
	}
}
