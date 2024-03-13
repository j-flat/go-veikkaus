package goveikkaus

import (
	"time"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

func (s *AuthService) AuthSessionIsActive() bool {
	return time.Now().Before(s.apiClient.SessionTimeout)
}

func getSessionTimeout() time.Time {
	currentTime := time.Now()
	duration := time.Duration(api.SessionTimeoutSeconds)
	sessionTimeoutTime := currentTime.Add(duration * time.Second)

	return sessionTimeoutTime
}
