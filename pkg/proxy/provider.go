package proxy

import (
	"context"
	"time"
)

type ProxyProvider interface {
	ProxyProvide(ctx context.Context, target string) (Proxy, error)
}

type ProxyProviderFunc func(ctx context.Context, target string) (Proxy, error)

func (f ProxyProviderFunc) ProxyProvide(ctx context.Context, target string) (Proxy, error) {
	return f(ctx, target)
}

func ProxyProviderWithTimeout(p ProxyProvider, timeout time.Duration) ProxyProvider {
	return ProxyProviderFunc(func(ctx context.Context, target string) (Proxy, error) {
		proxy, err := p.ProxyProvide(ctx, target)
		if err != nil {
			return nil, err
		}
		return ProxyWithTimeout(proxy, timeout), nil
	})
}

func ProxyProviderWithReadiness(p ProxyProvider, readiness func(context.Context, string) bool, interval time.Duration) ProxyProvider {
	return ProxyProviderFunc(func(ctx context.Context, target string) (Proxy, error) {
		proxy, err := p.ProxyProvide(ctx, target)
		if err != nil {
			return nil, err
		}
		return ProxyWithReadiness(proxy, func(ctx context.Context) bool {
			return readiness(ctx, target)
		}, interval), nil
	})
}
