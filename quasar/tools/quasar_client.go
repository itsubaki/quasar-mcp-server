package tools

import (
	"context"
	"fmt"

	"github.com/itsubaki/quasar/client"
	"google.golang.org/api/idtoken"
)

func NewQuasarClient(ctx context.Context, identityToken, targetURL string) (*client.Client, error) {
	if identityToken != "" {
		return client.New(targetURL, client.NewWithIdentityToken(identityToken)), nil
	}

	httpClient, err := idtoken.NewClient(ctx, targetURL)
	if err != nil {
		return nil, fmt.Errorf("new quasar client: %w", err)
	}

	return client.New(targetURL, httpClient), nil
}
