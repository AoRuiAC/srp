package auth

import (
	"context"
)

type UserPasswordChecker interface {
	Check(ctx context.Context, user, password string) bool
}

type UserPasswordFunc func(ctx context.Context, user, password string) bool

func (f UserPasswordFunc) Check(ctx context.Context, user, password string) bool {
	return f(ctx, user, password)
}

type UserPasswordMap map[string]string

func (m UserPasswordMap) Check(ctx context.Context, user, password string) bool {
	want, ok := m[user]
	return ok && want == password
}

func UserPasswordAuthenticator(c UserPasswordChecker) Authenticator {
	return AuthenticateFunc(func(ctx context.Context, req AuthenticateRequest) bool {
		return c.Check(ctx, req.User, req.Password)
	})
}
