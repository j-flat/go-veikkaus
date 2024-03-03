package goveikkaus

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("auth login: authservice utils", func() {
	Describe("getSessionTimeout", func() {
		It("should return session timeout that is 30 minutes from now", func() {
			currentTime := time.Now()
			sessionTimeOutTime := getSessionTimeout()
			Expect(sessionTimeOutTime.Truncate(time.Millisecond)).To(Equal(currentTime.Add(HalfHourInSeconds * time.Second).Truncate(time.Millisecond)))
		})
	})
})
