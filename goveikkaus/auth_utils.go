package goveikkaus

import (
	"time"
)

const HalfHourInSeconds = 1800

func getSessionTimeout() time.Time {
	currentTime := time.Now()

	sessionTimeoutTime := currentTime.Add(HalfHourInSeconds * time.Second)

	return sessionTimeoutTime
}
