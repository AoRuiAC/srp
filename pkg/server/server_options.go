package server

import (
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/pigeonligh/srp/pkg/proxy"
	"github.com/pigeonligh/srp/pkg/reverseproxy"
)

type Option func(s *server)

func WithReverseProxy(rp reverseproxy.Handler) Option {
	return func(s *server) {
		s.rp = rp
	}
}

func WithProxy(p proxy.Handler) Option {
	return func(s *server) {
		s.p = p
	}
}

func WithWishMiddleware(m wish.Middleware) Option {
	return func(s *server) {
		s.m = m
	}
}

func WithSSHHandler(h ssh.Handler) Option {
	return func(s *server) {
		s.h = h
	}
}

func WithSSHOptions(options ...ssh.Option) Option {
	return func(s *server) {
		s.sshOptions = append(s.sshOptions, options...)
	}
}
