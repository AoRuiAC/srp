package server

import (
	"github.com/charmbracelet/ssh"
	"github.com/pigeonligh/srp/pkg/protocol"
)

func (s *server) channelOption(srv *ssh.Server) error {
	if srv.ChannelHandlers == nil {
		srv.ChannelHandlers = make(map[string]ssh.ChannelHandler)
	}

	srv.ChannelHandlers["direct-tcpip"] = s.p.HandleProxyfunc
	srv.ChannelHandlers["session"] = ssh.DefaultSessionHandler
	return nil
}

func (s *server) requestOption(srv *ssh.Server) error {
	srv.RequestHandlers = map[string]ssh.RequestHandler{
		protocol.ForwardRequestType: s.rp.HandleSSHRequest,
		protocol.CancelRequestType:  s.rp.HandleSSHRequest,
	}
	return nil
}

func (s *server) passwordOption(srv *ssh.Server) error {
	return ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		rp := s.rp.PasswordHandler()(ctx, password)
		p := s.p.PasswordHandler()(ctx, password)
		return rp || p
	})(srv)
}

func (s *server) publickeyOption(srv *ssh.Server) error {
	return ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		rp := s.rp.PublicKeyHandler()(ctx, key)
		p := s.p.PublicKeyHandler()(ctx, key)
		return rp || p
	})(srv)
}
