package dialer

import (
	"context"
	"errors"
	"fmt"
	"net"
)

type listenDialer chan net.Conn

func (ld listenDialer) Accept() (net.Conn, error) {
	c, ok := <-ld
	if ok {
		return c, nil
	}
	return nil, net.ErrClosed
}

func (ld listenDialer) Close() error {
	close(ld)
	return nil
}

func (ld listenDialer) Addr() net.Addr {
	return &net.UnixAddr{
		Net:  "channel",
		Name: "listendialer",
	}
}

func (ld listenDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	c1, c2 := net.Pipe()
	select {
	case ld <- c1:
		return c2, nil

	default:
		return nil, net.ErrClosed
	}
}

func ListenDialer() (net.Listener, NetDialer) {
	ld := make(listenDialer)
	return ld, ld
}

func ListenDialerWithBuffer(size int) (net.Listener, NetDialer) {
	ld := make(listenDialer, size)
	return ld, ld
}

func HandleListener(l net.Listener, h func(net.Conn)) error {
	for {
		c, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("listener accept: %w", err)
		}
		go func() {
			h(c)
			_ = c.Close()
		}()
	}
}
