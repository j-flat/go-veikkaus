package goveikkaus

import (
	errors "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

func (s *AuthService) Logout() error {
	if s.apiClient.client.Jar == nil {
		return &errors.UserNotLoggedInError{}
	} else {
		s.apiClient.client.Jar = nil
	}

	return nil
}
