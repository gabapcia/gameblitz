package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildAuthenticatorFunc(t *testing.T) {
	var (
		ctx = context.Background()

		expectedClaims = Claims{
			GameID: uuid.NewString(),
		}
	)

	t.Run("OK", func(t *testing.T) {
		authenticatorFunc := BuildAuthenticatorFunc(func(ctx context.Context, credential string) (Claims, error) {
			return expectedClaims, nil
		})

		claims, err := authenticatorFunc(ctx, uuid.NewString())
		assert.NoError(t, err)

		assert.Equal(t, expectedClaims, claims)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		authenticatorFunc := BuildAuthenticatorFunc(func(ctx context.Context, credential string) (Claims, error) {
			return Claims{}, ErrInvalidCredentials
		})

		claims, err := authenticatorFunc(ctx, uuid.NewString())
		assert.ErrorIs(t, err, ErrInvalidCredentials)

		assert.Empty(t, claims)
	})
}
