package http

import (
	"net/http"
	"net/http/httputil"

	"github.com/pigeonligh/srp/pkg/nets"
)

type HTTPDirector func(hostname string) (scheme string, host string)

type handler struct {
	director HTTPDirector
	dialer   nets.NetDialer
}

func Handler(director HTTPDirector, netDialer nets.NetDialer) *handler {
	p := &handler{
		director: director,
		dialer:   netDialer,
	}
	if p.director == nil {
		p.director = func(hostname string) (string, string) {
			return "http", hostname + ":80"
		}
	}
	if p.dialer == nil {
		p.dialer = nets.DefaultNetDialer
	}
	return p
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Host = r.Host
	hostname := r.URL.Hostname()
	scheme, target := h.director(hostname)
	if scheme == "" || target == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.URL.Scheme = scheme
	r.URL.Host = target
	r.Host = target

	(&httputil.ReverseProxy{
		Director:  func(r *http.Request) {},
		Transport: &http.Transport{DialContext: h.dialer.DialContext},
	}).ServeHTTP(w, r)
}
