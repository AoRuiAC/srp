package https

import (
	"net/http"
	"net/http/httputil"

	"github.com/pigeonligh/srp/pkg/dialer"
)

type HandleOption func(*handler)

func WithDirector(d func(string) string) HandleOption {
	return func(p *handler) {
		p.director = d
	}
}

func WithDialer(d dialer.NetDialer) HandleOption {
	return func(p *handler) {
		p.dialer = d
	}
}

type handler struct {
	director func(hostname string) string
	dialer   dialer.NetDialer
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
	if p.dialer == nil {
		p.dialer = dialer.DefaultNetDialer
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
		Transport: &http.Transport{DialContext: p.dialer.DialContext},
	}).ServeHTTP(w, r)
}
