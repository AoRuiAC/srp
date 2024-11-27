package auth

import (
	"context"

	gossh "golang.org/x/crypto/ssh"
)

// req

type AuthenticateRequest struct {
	User      string
	Password  string
	PublicKey gossh.PublicKey
}

// def

type Authenticator interface {
	Authenticate(context.Context, AuthenticateRequest) bool
}

type AuthenticateFunc func(context.Context, AuthenticateRequest) bool

func (f AuthenticateFunc) Authenticate(ctx context.Context, req AuthenticateRequest) bool {
	return f(ctx, req)
}

var _ Authenticator = AuthenticateFunc(nil)

// slice

type Authenticators []Authenticator

func (slice Authenticators) Authenticate(ctx context.Context, req AuthenticateRequest) bool {
	for _, a := range slice {
		if a.Authenticate(ctx, req) {
			return true
		}
	}
	return false
}

func MergeAuthenticators(slice ...Authenticator) Authenticator {
	return Authenticators(slice)
}
