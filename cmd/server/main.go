package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/wish"
	"github.com/pigeonligh/srp/pkg/proxy"
	"github.com/pigeonligh/srp/pkg/proxy/providers"
	"github.com/pigeonligh/srp/pkg/reverseproxy"
	"github.com/pigeonligh/srp/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	var name string
	var address string
	var socketDir string
	var hostKey string

	cmd := &cobra.Command{
		Use: "srp-server",
		Run: func(cmd *cobra.Command, args []string) {
			rp, err := reverseproxy.New(nil, nil, socketDir)
			if err != nil {
				log.Fatalln("Error:", err)
			}
			p := proxy.New(nil, nil, providers.SocketProvider(rp, 0), true)

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
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "SRP", "SRP Server Name")
	cmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1:22", "SRP listen address")
	cmd.Flags().StringVarP(&socketDir, "socket-dir", "d", "", "Path for unix socket files")
	cmd.Flags().StringVarP(&hostKey, "host-key", "k", "ssh_host_ed25519_key", "Host Key File for SSH Server")

	_ = cmd.Execute()
}
