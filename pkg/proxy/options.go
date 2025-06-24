package proxy

import "github.com/pigeonligh/srp/pkg/auth"

type Option func(*handler)

func WithAuthenticator(authenticator auth.Authenticator) Option {
	return func(h *handler) {
		h.authenticator = authenticator
	}
}

func WithAuthorizer(authorizer auth.Authorizer) Option {
	return func(h *handler) {
		h.authorizer = authorizer
	}
}

func WithProxyProvider(provider ProxyProvider) Option {
	return func(h *handler) {
		h.provider = provider
	}
}

func WithCacheEnabled(enabled bool) Option {
	return func(h *handler) {
		h.cacheEnabled = enabled
	}
}

func WithProxyCallbacks(callbacks ProxyCallbacks) Option {
	return func(h *handler) {
		h.callbacks = callbacks
	}
}
