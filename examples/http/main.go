package main

import (
	"context"
	"fmt"
	gohttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pigeonligh/srp/pkg/http"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s := http.HTTP{
		Address: "127.0.0.1:8000",
		Handler: gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
			fmt.Fprintln(w, r.Host, r.URL)
			w.WriteHeader(gohttp.StatusOK)
		}),
	}
	_ = s.Run(ctx)
}
