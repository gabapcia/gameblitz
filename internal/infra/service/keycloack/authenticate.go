package keycloack

import (
	"context"
	"fmt"
	"strings"

	"github.com/gabarcia/gameblitz/internal/auth"
)

const splittedHeaderSize = 2

func (s service) Authenticate(ctx context.Context, header string) (auth.Claims, error) {
	splittedHeader := strings.Split(header, " ")
	if len(splittedHeader) != splittedHeaderSize {
		return auth.Claims{}, fmt.Errorf("%w: malformed token", auth.ErrInvalidCredentials)
	}

	var (
		tokenType = splittedHeader[0]
		rawToken  = splittedHeader[1]
	)

	switch strings.ToUpper(tokenType) {
	case "BEARER":
		return s.JWTAuthentication(ctx, rawToken)
	default:
		return auth.Claims{}, fmt.Errorf("%w: token type not supported", auth.ErrInvalidCredentials)
	}
}
