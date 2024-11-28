package server

import (
	"cmp"

	"github.com/charmbracelet/ssh"
	"github.com/pigeonligh/srp/pkg/protocol"
)

func (s *server) channelOption(srv *ssh.Server) error {
	if s.p == nil {
		return nil
	}

	if srv.ChannelHandlers == nil {
		srv.ChannelHandlers = make(map[string]ssh.ChannelHandler)
	}
	srv.ChannelHandlers["direct-tcpip"] = s.p.HandleProxyfunc
	srv.ChannelHandlers["session"] = ssh.DefaultSessionHandler
	return nil
}

func (s *server) requestOption(srv *ssh.Server) error {
	if s.rp == nil {
		return nil
	}

	srv.RequestHandlers = map[string]ssh.RequestHandler{
		protocol.ForwardRequestType: s.rp.HandleSSHRequest,
		protocol.CancelRequestType:  s.rp.HandleSSHRequest,
	}
	return nil
}

func (s *server) passwordOption(srv *ssh.Server) error {
	return ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		ret := make([]bool, 0)
		if s.rp != nil {
			ret = append(ret, s.rp.PasswordHandler()(ctx, password))
		}
		if s.p != nil {
			ret = append(ret, s.p.PasswordHandler()(ctx, password))
		}
		return cmp.Or(ret...) || len(ret) == 0
	})(srv)
}

func (s *server) publickeyOption(srv *ssh.Server) error {
	return ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		ret := make([]bool, 0)
		if s.rp != nil {
			ret = append(ret, s.rp.PublicKeyHandler()(ctx, key))
		}
		if s.p != nil {
			ret = append(ret, s.p.PublicKeyHandler()(ctx, key))
		}
		return cmp.Or(ret...) || len(ret) == 0
	})(srv)
}
