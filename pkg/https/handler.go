package https

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
)

type HandleOption func(*handler)

func WithDirector(d func(string) string) HandleOption {
	return func(p *handler) {
		p.director = d
	}
}

func WithDial(d func(ctx context.Context, network, addr string) (net.Conn, error)) HandleOption {
	return func(p *handler) {
		p.dial = d
	}
}

type handler struct {
	director func(hostname string) string
	dial     func(ctx context.Context, network, addr string) (net.Conn, error)
}

func Handler(options ...HandleOption) *handler {
	p := &handler{}
	for _, o := range options {
		o(p)
	}
	if p.director == nil {
		p.director = func(hostname string) string {
			return hostname + ":80"
		}
	}
	if p.dial == nil {
		p.dial = func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, addr)
		}
	}
	return p
}

func (p *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Host = r.Host
	hostname := r.URL.Hostname()
	target := p.director(hostname)
	if target == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.URL.Scheme = "http"
	r.URL.Host = target
	r.Host = target

	(&httputil.ReverseProxy{
		Director:  func(r *http.Request) {},
		Transport: &http.Transport{DialContext: p.dial},
	}).ServeHTTP(w, r)
}
