package goveikkaus

import (
	"context"
	"net/http"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

func (s *AuthService) AccountBalance(ctx context.Context) (*AccountBalance, *http.Response, error) {
	req, _ := api.GetRequest(api.AccountBalanceEndpoint, http.MethodGet, nil)

	var balance AccountBalance

	resp, err := s.apiClient.Do(ctx, req, &balance)

	if err != nil {
		return nil, resp, err
	}

	defer resp.Body.Close()

	return &balance, resp, nil
}
