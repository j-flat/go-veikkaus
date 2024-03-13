package goveikkaus

import (
	"fmt"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	errors "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

// MockCookieJar implements http.CookieJar for testing purposes.
type MockCookieJar struct {
	CookieSlice []*http.Cookie
}

// SetCookies implements the SetCookies method of the http.CookieJar interface.
func (jar *MockCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.CookieSlice = cookies
}

// Cookies implements the Cookies method of the http.CookieJar interface.
func (jar *MockCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return jar.CookieSlice
}

func getDummyCookies() []*http.Cookie {
	var dummyCokies []*http.Cookie

	for i := 1; i <= 5; i++ {
		cookie := &http.Cookie{
			Name:  fmt.Sprintf("cookie-%d", i),
			Value: fmt.Sprintf("value-%d", i),
		}
		dummyCokies = append(dummyCokies, cookie)
	}

	return dummyCokies
}

var _ = Describe("authservice: logout", func() {
	var url *url.URL

	BeforeEach(func() {
		url, _ = url.Parse("http://localhost")
	})

	Describe("Logout", func() {
		It("should logout from existing session", func() {
			var dummyJar = &MockCookieJar{}

			dummyJar.SetCookies(url, getDummyCookies())
			client := &http.Client{
				Jar: dummyJar,
			}

			veikkausClient := NewClient(client)
			Expect(len(veikkausClient.client.Jar.Cookies(url))).NotTo(BeZero())

			err := veikkausClient.Auth.Logout()
			Expect(veikkausClient.client.Jar).To(BeNil())
			Expect(err).To(BeNil())

		})
		It("should return error when no login session is active", func() {
			client := &http.Client{}
			veikkausClient := NewClient(client)

			Expect(veikkausClient.client.Jar).To(BeNil())

			err := veikkausClient.Auth.Logout()

			Expect(err).NotTo(BeNil())
			Expect(veikkausClient.client.Jar).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(&errors.UserNotLoggedInError{}))
		})
	})
})
