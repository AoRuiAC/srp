package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/wish"
	"github.com/pigeonligh/srp/pkg/auth"
	"github.com/pigeonligh/srp/pkg/proxy"
	"github.com/pigeonligh/srp/pkg/proxy/providers"
	"github.com/pigeonligh/srp/pkg/reverseproxy"
	"github.com/pigeonligh/srp/pkg/server"
)

var (
	name    = "SRP Auth Example"
	address = "127.0.0.1:8022"
	hostKey = "examples/auth/keys/host"
)

func main() {
	rp, err := reverseproxy.New(
		auth.UserPublicKeysAuthenticator(auth.PublicKeysDir("examples/auth/reverseproxy_auth")),
		auth.UserGlobsAuthorizer(auth.UserGlobsDir("examples/auth/reverseproxy_rules")),
		"",
	)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	p := proxy.New(
		auth.UserPublicKeysAuthenticator(auth.PublicKeysDir("examples/auth/proxy_auth")),
		auth.UserGlobsAuthorizer(auth.UserGlobsDir("examples/auth/proxy_rules")),
		providers.SocketProvider(rp, 0),
		true,
	)

	s, err := server.New(
		name,
		server.WithReverseProxy(rp),
		server.WithProxy(p),
		server.WithSSHOptions(
			wish.WithHostKeyPath(hostKey),
			wish.WithAddress(address),
		),
	)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		log.Fatalln("Error:", err)
	}
}
