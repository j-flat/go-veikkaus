package goveikkaus

import (
	"context"
	"net/http"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginSuccessful, *http.Response, error) {
	payloadStruct := LoginPayload{
		Type:     "STANDARD_LOGIN",
		User:     username,
		Password: password,
	}

	body, err := api.GetJSONPayload(payloadStruct)

	if err != nil {
		return nil, nil, err
	}

	req, _ := api.GetRequest(api.LoginEndpoint, http.MethodPost, body)

	var loginSuccessful LoginSuccessful

	resp, err := s.apiClient.Do(ctx, req, &loginSuccessful)

	if err != nil {
		return nil, resp, err
	}

	defer resp.Body.Close()

	// Store session timeout information and logged in state to client
	s.apiClient.SessionTimeout = getSessionTimeout()

	// Veikkaus API Returns empty JSON-response, no need to parse it to an empty object
	return &loginSuccessful, resp, nil
}
