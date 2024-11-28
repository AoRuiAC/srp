package https

import (
	"context"
	"fmt"
	"net"
)

type SocketHandler interface {
	ConvertHostPortToSocket(host, port string) (string, bool)
}

func DialFromSocketHandler(h SocketHandler) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		socket, ok := h.ConvertHostPortToSocket(host, port)
		if !ok {
			return nil, fmt.Errorf("cannot resolve %v", addr)
		}
		d := net.Dialer{}
		return d.DialContext(ctx, "unix", socket)
	}
}
