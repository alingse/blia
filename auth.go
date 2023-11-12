package blia

import (
	"context"
	"net/http"
)

type AuthParser[T any] interface {
	ParseFromRequest(r *http.Request) (T, error)
}

type Authenticator[T any] struct {
	parser     AuthParser[T]
	contextKey *struct{}
}

func NewAuthenticator[T any](parser AuthParser[T]) *Authenticator[T] {
	return &Authenticator[T]{
		parser:     parser,
		contextKey: new(struct{}),
	}
}

func (a *Authenticator[T]) WithLoginContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := a.parser.ParseFromRequest(r)
			if err == nil {
				ctx = context.WithValue(ctx, a.contextKey, user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (a *Authenticator[T]) LoginRequired() func(next http.Handler) http.Handler {
	return a.WithUserCheck(func(u T) bool { return true })
}

func (a *Authenticator[T]) WithUserCheck(f func(u T) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, ok := a.GetUserFromContext(r.Context())
			if !ok {
				WriteError(w, ErrUnauthorized)
			} else if !f(u) {
				WriteError(w, ErrForbidden)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func (a *Authenticator[T]) GetUserFromContext(ctx context.Context) (T, bool) {
	value := ctx.Value(a.contextKey)
	user, ok := value.(T)
	return user, ok
}
