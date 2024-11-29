package dialer

import (
	"context"

	gossh "golang.org/x/crypto/ssh"
)

type SSHDialer interface {
	DialContext(ctx context.Context, network, addr string, config *gossh.ClientConfig) (*gossh.Client, error)
}

type SSHDialerFunc func(ctx context.Context, network, addr string, config *gossh.ClientConfig) (*gossh.Client, error)

func (f SSHDialerFunc) DialContext(ctx context.Context, network, addr string, config *gossh.ClientConfig) (*gossh.Client, error) {
	return f(ctx, network, addr, config)
}

func NetSSHDialer(netDialer NetDialer) SSHDialer {
	if netDialer == nil {
		netDialer = DefaultNetDialer
	}
	return SSHDialerFunc(func(ctx context.Context, network, addr string, config *gossh.ClientConfig) (*gossh.Client, error) {
		conn, err := netDialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		sshConn, chans, reqs, err := gossh.NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return gossh.NewClient(sshConn, chans, reqs), nil
	})
}
