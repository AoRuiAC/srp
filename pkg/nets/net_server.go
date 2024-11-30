package nets

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type contextServerName struct{}

func ContextWithServerName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, contextServerName{}, name)
}

func GetServerNameFromContext(ctx context.Context) (string, bool) {
	obj := ctx.Value(contextServerName{})
	name, ok := obj.(string)
	return name, ok
}

type NetServer interface {
	Serve(l net.Listener) error
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunNetServer(ctx context.Context, s NetServer, l net.Listener) error {
	name, _ := GetServerNameFromContext(ctx)
	logger := logrus.WithField("netserver", name)

	var serverErr error
	done := make(chan struct{}, 1)

	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()

	go func() {
		logger.Info("Server start")

		var err error
		if l == nil {
			err = s.ListenAndServe()
		} else {
			err = s.Serve(l)
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Infof("Server run error: %v", err)
			serverErr = err
		}
		done <- struct{}{}
	}()

	<-done
	logger.Info("Server stopping")

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(stopCtx); err != nil {
		logger.Infof("Server stop error: %v", err)
		return err
	}
	logger.Info("Server stopped")
	return serverErr
}
