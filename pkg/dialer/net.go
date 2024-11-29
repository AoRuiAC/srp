package dialer

import (
	"context"
	"net"
)

type NetDialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

type NetDialerFunc func(ctx context.Context, network, addr string) (net.Conn, error)

func (f NetDialerFunc) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return f(ctx, network, addr)
}

var DefaultNetDialer NetDialer = NetDialerFunc(func(ctx context.Context, network, addr string) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, addr, network)
})

func NetDialerWithConnModifier(d NetDialer, m func(net.Conn) net.Conn) NetDialer {
	return NetDialerFunc(func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := d.DialContext(ctx, network, addr)
		if err == nil {
			conn = m(conn)
		}
		return conn, err
	})
}
