package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
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

	sshOptions []ssh.Option
}

func New(name string, options ...Option) (Server, error) {
	s := &server{
		name: name,
	}
	for _, o := range options {
		o(s)
	}
	return s, nil
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

	var serverErr error
	done := make(chan struct{}, 1)

	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()

	go func() {
		log.Info("Starting SSH server")
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Failed to run SSH server", "error", err)
			serverErr = err
			done <- struct{}{}
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Stopping SSH server")
	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
	return serverErr
}
