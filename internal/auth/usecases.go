package auth

import "context"

type (
	// Authenticate a game its credentials
	AuthenticateFunc func(ctx context.Context, credentials string) (Claims, error)
)
