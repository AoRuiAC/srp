package nets

import (
	"context"
	"fmt"
	"net"
)

type SocketHandler interface {
	ConvertHostPortToSocket(host, port string) (string, bool)
	SocketAlive(socket string) bool
}

func SocketsDialer(h SocketHandler) NetDialer {
	return NetDialerFunc(func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		socket, ok := h.ConvertHostPortToSocket(host, port)
		if !ok {
			return nil, fmt.Errorf("cannot resolve %v", addr)
		}
		if !h.SocketAlive(socket) {
			return nil, fmt.Errorf("target %v is not alive", addr)
		}
		d := net.Dialer{}
		return d.DialContext(ctx, "unix", socket)
	})
}
