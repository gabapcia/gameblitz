package keycloack

import (
	"context"
	"errors"
	"fmt"

	"github.com/gabarcia/gameblitz/internal/auth"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	jwt.RegisteredClaims
	GameID string `json:"client_id"`
}

func (c claims) toDomain() auth.Claims {
	return auth.Claims{
		GameID: c.GameID,
	}
}

func (s service) JWTAuthentication(ctx context.Context, rawToken string) (auth.Claims, error) {
	var c claims
	token, err := jwt.ParseWithClaims(rawToken, &c, func(t *jwt.Token) (any, error) {
		keyID, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: kid header not found", auth.ErrInvalidCredentials)
		}

		key, ok := s.keys.LookupKeyID(keyID)
		if !ok {
			return nil, fmt.Errorf("%w: unable to find a key that matches with the kid header provided", auth.ErrInvalidCredentials)
		}

		var publicKey any
		return publicKey, key.Raw(&publicKey)
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = fmt.Errorf("%w: expired credentials", auth.ErrInvalidCredentials)
		}

		return auth.Claims{}, err
	}

	if !token.Valid {
		return auth.Claims{}, auth.ErrInvalidCredentials
	}

	return c.toDomain(), nil
}
