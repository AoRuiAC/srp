package proxy

import (
	"github.com/charmbracelet/ssh"
	"github.com/pigeonligh/srp/pkg/protocol"
)

type ProxyCallbacks struct {
	OnHandleProxyFunc     func(ctx ssh.Context)
	OnHandleProxyDoneFunc func(ctx ssh.Context)

	OnProxyCreatedFunc             func(ctx ssh.Context, payload protocol.DirectPayload)
	OnProxyCreateFailedFunc        func(ctx ssh.Context, payload protocol.DirectPayload, err error)
	OnProxyChannelAcceptedFunc     func(ctx ssh.Context, payload protocol.DirectPayload)
	OnProxyChannelAcceptFailedFunc func(ctx ssh.Context, payload protocol.DirectPayload, err error)
	OnProxyDialedFunc              func(ctx ssh.Context, payload protocol.DirectPayload)
	OnProxyDialFailedFunc          func(ctx ssh.Context, payload protocol.DirectPayload, err error)
	OnProxyConnectionDoneFunc      func(ctx ssh.Context, payload protocol.DirectPayload, err error)
}

func (c *ProxyCallbacks) OnHandleProxy(ctx ssh.Context) {
	if c == nil || c.OnHandleProxyFunc == nil {
		return
	}
	c.OnHandleProxyFunc(ctx)
}

func (c *ProxyCallbacks) OnHandleProxyDone(ctx ssh.Context) {
	if c == nil || c.OnHandleProxyDoneFunc == nil {
		return
	}
	c.OnHandleProxyDoneFunc(ctx)
}

func (c *ProxyCallbacks) OnProxyCreated(ctx ssh.Context, payload protocol.DirectPayload) {
	if c == nil || c.OnProxyCreatedFunc == nil {
		return
	}
	c.OnProxyCreatedFunc(ctx, payload)
}

func (c *ProxyCallbacks) OnProxyCreateFailed(ctx ssh.Context, payload protocol.DirectPayload, err error) {
	if c == nil || c.OnProxyCreateFailedFunc == nil {
		return
	}
	c.OnProxyCreateFailedFunc(ctx, payload, err)
}

func (c *ProxyCallbacks) OnProxyChannelAccepted(ctx ssh.Context, payload protocol.DirectPayload) {
	if c == nil || c.OnProxyChannelAcceptedFunc == nil {
		return
	}
	c.OnProxyChannelAcceptedFunc(ctx, payload)
}

func (c *ProxyCallbacks) OnProxyChannelAcceptFailed(ctx ssh.Context, payload protocol.DirectPayload, err error) {
	if c == nil || c.OnProxyChannelAcceptFailedFunc == nil {
		return
	}
	c.OnProxyChannelAcceptFailedFunc(ctx, payload, err)
}

func (c *ProxyCallbacks) OnProxyDialed(ctx ssh.Context, payload protocol.DirectPayload) {
	if c == nil || c.OnProxyDialedFunc == nil {
		return
	}
	c.OnProxyDialedFunc(ctx, payload)
}

func (c *ProxyCallbacks) OnProxyDialFailed(ctx ssh.Context, payload protocol.DirectPayload, err error) {
	if c == nil || c.OnProxyDialFailedFunc == nil {
		return
	}
	c.OnProxyDialFailedFunc(ctx, payload, err)
}

func (c *ProxyCallbacks) OnProxyConnectionDone(ctx ssh.Context, payload protocol.DirectPayload, err error) {
	if c == nil || c.OnProxyConnectionDoneFunc == nil {
		return
	}
	c.OnProxyConnectionDoneFunc(ctx, payload, err)
}
