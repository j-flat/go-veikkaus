package goveikkaus

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

var HalfHourInSeconds = 1800

func timeTruncated(timeToTruncate time.Time) time.Time {
	return timeToTruncate.Truncate(time.Second)
}

func pauseForNSeconds(seconds int) {
	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)
}

var _ = Describe("auth login: authservice utils", func() {
	var currentTime time.Time

	BeforeEach(func() {
		currentTime = time.Now()
	})

	Describe("getSessionTimeout", func() {
		It("should return session timeout that is 30 minutes from now", func() {
			sessionTimeOutTime := getSessionTimeout()

			expectedTime := currentTime.Add(time.Duration(HalfHourInSeconds) * time.Second)

			Expect(timeTruncated(sessionTimeOutTime)).To(Equal(timeTruncated(expectedTime)))
		})
	})
	DescribeTable("AuthSessionIsActive",
		func(sessionTimeoutSeconds, sleepForNSeconds int, authSessionShouldBeActive bool) {
			api.SessionTimeoutSeconds = sessionTimeoutSeconds

			veikkausClient := NewClient(nil)

			sessionTimeOutTime := getSessionTimeout()
			veikkausClient.Auth.apiClient.SessionTimeout = sessionTimeOutTime

			Expect(veikkausClient.Auth.AuthSessionIsActive()).To(BeTrue())

			pauseForNSeconds(sleepForNSeconds)

			Expect(veikkausClient.Auth.AuthSessionIsActive()).To(Equal(authSessionShouldBeActive))
		},
		Entry("should return 'true' after sleep time when auth-session is still valid", 5, 1, true),
		Entry("should return 'false' after sleep time when auth-session is no longer valid", 1, 2, false),
	)
})
