package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type Proxy interface {
	Proxy(ctx context.Context, r io.Reader, w io.Writer) error
}

type ProxyFunc func(ctx context.Context, r io.Reader, w io.Writer) error

func (f ProxyFunc) Proxy(ctx context.Context, r io.Reader, w io.Writer) error {
	return f(ctx, r, w)
}

func Direct(network string, address string) Proxy {
	return ProxyFunc(func(ctx context.Context, r io.Reader, w io.Writer) error {
		d := net.Dialer{}
		c, err := d.DialContext(ctx, network, address)
		if err != nil {
			return fmt.Errorf("connecting to %v://%v: %w", network, address, err)
		}
		defer c.Close()

		go func() {
			_, _ = io.Copy(c, r)
		}()
		_, _ = io.Copy(w, c)
		return nil
	})
}

func UnixSocket(socket string) Proxy {
	return Direct("unix", socket)
}

func ProxyWithTimeout(p Proxy, timeout time.Duration) Proxy {
	return ProxyFunc(func(ctx context.Context, r io.Reader, w io.Writer) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return p.Proxy(ctx, r, w)
	})
}

func ProxyWithReadiness(p Proxy, readiness func(context.Context) bool, interval time.Duration) Proxy {
	return ProxyFunc(func(ctx context.Context, r io.Reader, w io.Writer) error {
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

		return p.Proxy(ctx, r, w)
	})
}
