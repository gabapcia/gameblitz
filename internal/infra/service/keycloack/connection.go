package keycloack

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/lestrrat-go/jwx/jwk"
)

type service struct {
	keys jwk.Set
}

func getCerts(ctx context.Context, certsURL string) (jwk.Set, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, certsURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("KeyCloack Get Certs: [%d] - %s", res.StatusCode, body)
	}

	keys, err := jwk.Parse(body)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func New(ctx context.Context, certsURL string) (*service, error) {
	certs, err := getCerts(ctx, certsURL)
	if err != nil {
		return nil, err
	}

	svc := &service{
		keys: certs,
	}
	return svc, nil
}
