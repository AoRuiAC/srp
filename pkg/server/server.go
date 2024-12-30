package server

import (
	"context"
	"fmt"
	"net"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/pigeonligh/srp/pkg/nets"
	"github.com/pigeonligh/srp/pkg/proxy"
	"github.com/pigeonligh/srp/pkg/reverseproxy"
)

type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	name string

	rp reverseproxy.Handler
	p  proxy.Handler
	m  wish.Middleware
	h  ssh.Handler
	l  net.Listener

	sshOptions []ssh.Option
}

func New(name string, options ...Option) Server {
	s := &server{
		name: name,
	}
	for _, o := range options {
		o(s)
	}
	return s
}

func (s *server) HandleSession(_ ssh.Handler) ssh.Handler {
	defaultHandler := func(sess ssh.Session) {
		if s.h != nil {
			s.h(sess)
			return
		}

		if len(sess.Command()) > 0 {
			fmt.Fprintln(sess, "Disallowed command")
		} else {
			_, _, isPty := sess.Pty()
			if !isPty {
				fmt.Fprintf(sess, "Welcome to %v, @%v!\n", s.name, sess.User())
			} else {
				fmt.Fprintln(sess, "PTY allocation request failed")
			}
		}
	}

	if s.m != nil {
		return s.m(defaultHandler)
	}
	return defaultHandler
}

func (s *server) Run(ctx context.Context) error {
	options := make([]ssh.Option, 0)
	options = append(options, s.sshOptions...)
	options = append(options,
		s.channelOption,
		s.requestOption,
		s.passwordOption,
		s.publickeyOption,
		wish.WithMiddleware(
			s.HandleSession,
			logging.Middleware(),
		),
	)

	srv, err := wish.NewServer(options...)
	if err != nil {
		return fmt.Errorf("create SSH server: %w", err)
	}

	ctx = nets.ContextWithServerName(ctx, s.name)
	return nets.RunNetServer(ctx, srv, s.l)
}
