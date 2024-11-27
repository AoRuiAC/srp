package auth

import "context"

// req

type AuthorizeRequest struct {
	User   string
	Target string
}

// def

type Authorizer interface {
	Authorize(context.Context, AuthorizeRequest) bool
}

type AuthorizeFunc func(context.Context, AuthorizeRequest) bool

func (f AuthorizeFunc) Authorize(ctx context.Context, req AuthorizeRequest) bool {
	return f(ctx, req)
}

var _ Authorizer = AuthorizeFunc(nil)

// slice

type Authorizers []Authorizer

func (slice Authorizers) Authorize(ctx context.Context, req AuthorizeRequest) bool {
	for _, a := range slice {
		if a.Authorize(ctx, req) {
			return true
		}
	}
	return false
}

func MergeAuthorizers(slice ...Authorizer) Authorizer {
	return Authorizers(slice)
}
