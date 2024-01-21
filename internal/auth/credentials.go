package auth

import (
	"context"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Claims struct {
	GameID string
}

func BuildAuthenticatorFunc(serviceValidateCredentialsFunc ServiceValidateCredentialsFunc) AuthenticateFunc {
	return func(ctx context.Context, credentials string) (Claims, error) {
		return serviceValidateCredentialsFunc(ctx, credentials)
	}
}
