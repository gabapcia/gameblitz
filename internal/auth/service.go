package auth

import "context"

type (
	// Validates the credential on the service provider
	ServiceValidateCredentialsFunc func(ctx context.Context, credential string) (Claims, error)
)
