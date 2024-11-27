package proxy

import (
	"fmt"
	"net"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/pigeonligh/srp/pkg/auth"
	"github.com/pigeonligh/srp/pkg/protocol"
	gossh "golang.org/x/crypto/ssh"
)

type Handler interface {
	PasswordHandler() ssh.PasswordHandler
	PublicKeyHandler() ssh.PublicKeyHandler

	HandleProxyfunc(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context)
}

type handler struct {
	authenticator auth.Authenticator
	authorizer    auth.Authorizer
	provider      ProxyProvider
	cacheEnabled  bool
}

func New(authenticator auth.Authenticator, authorizer auth.Authorizer, provider ProxyProvider, cacheEnabled bool) Handler {
	return &handler{
		authenticator: authenticator,
		authorizer:    authorizer,
		provider:      provider,
		cacheEnabled:  cacheEnabled,
	}
}

func (h *handler) PasswordHandler() ssh.PasswordHandler {
	return func(ctx ssh.Context, password string) bool {
		var ret bool
		if h.authenticator == nil {
			ret = true
		} else {
			ret = h.authenticator.Authenticate(ctx, auth.AuthenticateRequest{
				User:     ctx.User(),
				Password: password,
			})
		}

		ctx.SetValue(protocol.ContextKeyProxyAuthed, ret)
		return ret
	}
}

func (h *handler) PublicKeyHandler() ssh.PublicKeyHandler {
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		var ret bool
		if h.authenticator == nil {
			ret = true
		} else {
			ret = h.authenticator.Authenticate(ctx, auth.AuthenticateRequest{
				User:      ctx.User(),
				PublicKey: key,
			})
		}

		ctx.SetValue(protocol.ContextKeyProxyAuthed, ret)
		return ret
	}
}

func (h *handler) GetProxy(ctx ssh.Context, target string) (Proxy, error) {
	authed, _ := ctx.Value(protocol.ContextKeyProxyAuthed).(bool)
	if !authed {
		return nil, fmt.Errorf("unauthenticated for proxy")
	}

	var cachedResult any
	if h.cacheEnabled {
		cacheKey := protocol.CachedProxyKey{Target: target}
		cachedResult = ctx.Value(cacheKey)
		if cachedResult != nil {
			if proxy, ok := cachedResult.(Proxy); ok {
				return proxy, nil
			}
			if err, ok := cachedResult.(error); ok {
				return nil, err
			}
		}
		defer func() {
			if cachedResult != nil {
				ctx.SetValue(cacheKey, cachedResult)
			}
		}()
	}

	if h.authorizer != nil {
		if !h.authorizer.Authorize(ctx, auth.AuthorizeRequest{
			User:   ctx.User(),
			Target: target,
		}) {
			err := fmt.Errorf("access denied")
			cachedResult = err
			return nil, err
		}
	}

	proxy, err := h.provider.ProxyProvide(ctx, target)
	if err != nil {
		cachedResult = err
		return nil, err
	}
	cachedResult = proxy
	return proxy, nil
}

func (h *handler) HandleProxyfunc(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context) {
	log.Infof("Handle direct-tcpip for user %v in %v", ctx.User(), ctx.SessionID())

	ch, _, err := newChan.Accept()
	if err != nil {
		log.Errorf("Cannot accept channel for %v: %v", ctx.SessionID(), err)
		return
	}
	defer ch.Close()

	var payload protocol.DirectPayload
	err = gossh.Unmarshal(newChan.ExtraData(), &payload)
	if err != nil {
		log.Errorf("Cannot accept extra data for %v: %v", ctx.SessionID(), err)
		return
	}
	log.Infof("Payload for session %v: %v", ctx.SessionID(), payload)

	proxy, err := h.GetProxy(ctx, net.JoinHostPort(payload.Host, fmt.Sprint(payload.Port)))
	if err != nil {
		log.Errorf("Cannot create proxy for %v: %v", ctx.SessionID(), err)
		return
	}

	log.Infof("Proxy created for session %v.", ctx.SessionID())
	err = proxy.Proxy(ctx, ch, ch)
	if err != nil {
		log.Errorf("Cannot handle proxy for %v: %v", ctx.SessionID(), err)
		return
	}

	log.Infof("Proxy done for session %v.", ctx.SessionID())
	<-ctx.Done()
}
