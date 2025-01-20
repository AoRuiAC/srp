package nets

import (
	"io"

	"golang.org/x/sync/errgroup"
)

func HandleConnections(c1, c2 io.ReadWriteCloser) error {
	var pipes errgroup.Group
	pipes.Go(func() error {
		_, err := io.Copy(c1, c2)
		SafeCloseConn(c1)
		return err
	})
	pipes.Go(func() error {
		_, err := io.Copy(c2, c1)
		SafeCloseConn(c2)
		return err
	})

	return pipes.Wait()
}

func SafeCloseConn(c io.ReadWriteCloser) {
	if cw, ok := c.(interface {
		CloseWrite() error
	}); ok {
		_ = cw.CloseWrite()
	} else {
		_ = c.Close()
	}
}
