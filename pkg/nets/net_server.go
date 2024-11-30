package nets

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

type NetServer interface {
	Serve(l net.Listener) error
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunNetServer(ctx context.Context, s NetServer, l net.Listener) error {
	var serverErr error
	done := make(chan struct{}, 1)

	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()

	go func() {
		var err error
		if l == nil {
			err = s.ListenAndServe()
		} else {
			err = s.Serve(l)
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr = err
		}
		done <- struct{}{}
	}()

	<-done

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(stopCtx); err != nil {
		return err
	}
	return serverErr
}
