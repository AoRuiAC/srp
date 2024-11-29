package client

import "context"

type Client interface {
	Run(ctx context.Context) error
}

type client struct {
}

func New(options ...Option) Client {
	c := &client{}
	for _, o := range options {
		o(c)
	}
	return c
}

func (c *client) Run(ctx context.Context) error {
	return nil
}
