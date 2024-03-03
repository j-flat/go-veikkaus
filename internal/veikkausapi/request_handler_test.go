package veikkausapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	ExpectedLoginURL    = "https://www.veikkaus.fi/api/bff/v1/sessions"
	NoCookiesToday      = 0
	SingleCookieInSlice = 1
	DummyCookieCount    = 2
	AppendedCookieCount = 3
)

type DummyPayload struct {
	Foo string
	Bar string
	Baz string
}

var dummyURL, _ = url.Parse("http://localhost:8080")

var dummyCookie1 = &http.Cookie{
	Name:  "cookie1",
	Value: "value1",
}

var dummyCookie2 = &http.Cookie{
	Name:  "cookie2",
	Value: "value2",
}

var dummyCookies = []*http.Cookie{dummyCookie1, dummyCookie2}

var dummyPayload = DummyPayload{
	Foo: "foo",
	Bar: "bar",
	Baz: "baz",
}

var dummyPayloadBytes = []byte(`{"Foo":"foo","Bar":"bar","Baz":"baz"}`)
var invalidPayloadBytes = []byte(nil)

var unprocessableStruct = map[string]interface{}{
	"foo": make(chan int),
}

func mockNewRequestError(method, requestURL string, body io.Reader) (*http.Request, error) {
	return nil, errors.New("mocked error")
}

type TestStruct struct {
	Name  string
	Age   int
	Valid bool
}

type MockCookieJar struct{}

func (m *MockCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	// Implement SetCookies method for your mock

}

func (m *MockCookieJar) Cookies(u *url.URL) []*http.Cookie {
	// Implement Cookies method for your mock
	return []*http.Cookie{}
}

func getRequestBody(req *http.Request) []byte {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Panic("Error reading request body to a byte-array. ERR:", err)
	}
	return body
}

var _ = Describe("internal/veikkausapi: request-handler code", func() {
	Describe("getRequestUrl", func() {
		It(fmt.Sprintf("should return request url with '%s' as base-url and given path", VeikkausAPIBaseURL), func() {
			Expect(getRequestURL("sessions")).To(Equal(ExpectedLoginURL))
		})
	})
	Describe("setRequestHeaders", func() {
		It("should set standard request headers for veikkaus api calls", func() {
			dummyRequest, _ := http.NewRequest("GET", "foobar.com", nil)

			setRequestHeaders(dummyRequest)

			Expect(dummyRequest.Header.Get("Content-Type")).To(Equal("application/json"))
			Expect(dummyRequest.Header.Get("Accept")).To(Equal("application/json"))
			Expect(dummyRequest.Header.Get("X-ESA-API-KEY")).To(Equal("ROBOT"))
		})
	})
	Describe("WithContext", func() {
		It("returns request pointer with the given context", func() {
			req := &http.Request{}
			origCtx := req.Context()
			Expect(origCtx).To(Equal(context.Background()))
			ctx := context.TODO()

			WithContext(ctx, req)
			Expect(req.Context()).NotTo(BeNil())
			Expect(req.Context()).To(BeEquivalentTo(context.TODO()))
		})
	})
	Describe("GetBaseURL", func() {
		It("should return overwritten base-url when needed", func() {
			overwrittenBaseURL := "http://localhost:80808"
			OverWriteBaseURL = true
			originalBaseURL := BaseURL
			BaseURL = overwrittenBaseURL

			Expect(GetBaseURL()).To(Equal(overwrittenBaseURL))
			OverWriteBaseURL = false
			BaseURL = originalBaseURL
		})
	})
	DescribeTable("validateRequestMethod",
		func(method string, allowedMethods []string, expectError bool, expectedError string) {
			isValid := validateRequestMethod(method, allowedMethods)
			if expectError {
				Expect(isValid).NotTo(BeNil())
				Expect(isValid.Error()).To(Equal(expectedError))
			} else {
				Expect(isValid).To(BeNil())
			}
		},
		Entry("should return 'nil' if given request method is in range of allowed methods", "GET", []string{"GET", "INFO"}, false, nil),
		Entry("should return error if given request method is not in range of allowed methods", "GET", []string{"PUT", "POST"}, true, "Expected to receive one of [|PUT||POST|] request-type, received 'GET'"),
	)
	DescribeTable("validateRequestUrl", func(requestUrl string, expectError bool, expectedError string) {
		isError := validateRequestURL(requestUrl)
		if expectError {
			Expect(isError).NotTo(BeNil())
			Expect(isError.Error()).To(Equal(expectedError))
		} else {
			Expect(isError).To(BeNil())
		}
	},
		Entry("should return 'nil' when request url is valid", "https://google.com", false, nil),
		Entry("should return error when request url is missing 'http(s)://'", "is.fi", true, "Request URL was malformatted. ERR: parse \"is.fi\": invalid URI for request"),
		Entry("should return error when request url has invalid structure", "https://[PLACEHOLDER].f", true, "Request URL was malformatted. ERR: parse \"https://[PLACEHOLDER].f\": invalid port \".f\" after host"),
	)
	DescribeTable("handlePutPost",
		func(url string, method string, payloadBytes []byte, expectError bool, mockNewRequest bool, expectedError string) {
			if expectError {
				originalNewRequest := newRequest
				newRequest = mockNewRequestError

				defer func() {
					newRequest = originalNewRequest
				}()
				req, err := handlePutPost(url, method, payloadBytes)
				Expect(err).NotTo(BeNil())
				Expect(req).To(BeNil())
				Expect(err.Error()).To(Equal(expectedError))
			} else {
				req, err := handlePutPost(url, method, payloadBytes)
				Expect(req).NotTo(BeNil())
				requestBody := getRequestBody(req)

				Expect(err).To(BeNil())
				Expect(req.Method).To(Equal(method))
				Expect(req.URL.String()).To(Equal(url))
				Expect(requestBody).To(Equal(payloadBytes))
			}
		},
		Entry("should return POST-request with given payload", "https://foobar.com/api/v1/hello", "POST", dummyPayloadBytes, false, false, nil),
		Entry("should return PUT-request with given payload", "http://localhost:8080", "PUT", dummyPayloadBytes, false, false, nil),
		Entry("should return error for POST-request without given payload", "https://amazingapi.com", "POST", nil, true, false, "payload bytes were expected, received nil"),
		Entry("should return error for PUT-request without given payload", "https://is.fi", "PUT", nil, true, false, "payload bytes were expected, received nil"),
		Entry("should return error for unsupported method", "https://is.fi", "GET", nil, true, false, "Expected to receive one of [|PUT||POST|] request-type, received 'GET'"),
		Entry("should return error for malformatted url", "˛∞é®§", "PUT", dummyPayloadBytes, true, false, "Request URL was malformatted. ERR: parse \"˛∞é®§\": invalid URI for request"),
		Entry("should return error when http.NewRequest unexpectedly fails", "https://localhost:8080", "PUT", dummyPayloadBytes, true, true, "Error creating 'PUT'-request for https://localhost:8080: mocked error"),
	)
	DescribeTable("handleGet",
		func(url string, expectError bool, mockNewRequest bool, expectedError string) {
			if expectError {
				originalNewRequest := newRequest
				newRequest = mockNewRequestError

				defer func() {
					newRequest = originalNewRequest
				}()

				req, err := handleGet(url)
				Expect(err).NotTo(BeNil())
				Expect(req).To(BeNil())
				Expect(err.Error()).To(Equal(expectedError))
			} else {
				req, err := handleGet(url)
				Expect(req).NotTo(BeNil())
				Expect(err).To(BeNil())
				Expect(req.Method).To(Equal("GET"))
				Expect(req.URL.String()).To(Equal(url))
			}
		},
		Entry("should return GET-request when all fields are valid", "https://foobar.com/api/v1/hello", false, false, nil),
		Entry("should return error for malformatted url", "˛∞é®§", true, false, "Request URL was malformatted. ERR: parse \"˛∞é®§\": invalid URI for request"),
		Entry("should return error when http.NewRequest unexpectedly returns error", "https://foobar.com/api/v2/baz", true, true, "Error creating 'GET'-request for https://foobar.com/api/v2/baz: mocked error"),
	)
	DescribeTable("requestHandler",
		func(requestUrl, requestMethod string, jsonPayload []byte, expectError bool, expectedError string) {
			if req, err := requestHandler(requestUrl, requestMethod, jsonPayload); expectError {
				Expect(err).ToNot(BeNil())
				Expect(req).To(BeNil())
				Expect(err.Error()).To(Equal(expectedError))
			} else {
				Expect(err).To(BeNil())
				Expect(req).NotTo((BeNil()))
			}
		},
		Entry("should return 'POST'-request with valid input", "https://localhost:8080", "POST", dummyPayloadBytes, false, nil),
		Entry("should return 'PUT'-request with valid input", "https://localhost:8080", "PUT", dummyPayloadBytes, false, nil),
		Entry("should return 'GET'-request with valid input", "https://localhost:8080", "GET", nil, false, nil),
		Entry("should return error when unsupported method is attempted", "https://localhost:8080", "DEL", nil, true, "unsupported method 'DEL' provided"),
		Entry("should return error when request url is invalid for 'POST'", "www.*****", "POST", dummyPayloadBytes, true, "Request URL was malformatted. ERR: parse \"www.*****\": invalid URI for request"),
		Entry("should return error when request url is invalid for 'PUT'", "**®é¸ƒ", "PUT", dummyPayloadBytes, true, "Request URL was malformatted. ERR: parse \"**®é¸ƒ\": invalid URI for request"),
		Entry("should return error when request url is invalid for 'GET'", "ww.notvalid.d", "GET", nil, true, "Request URL was malformatted. ERR: parse \"ww.notvalid.d\": invalid URI for request"),
	)
	DescribeTable("GetJsonPayload",
		func(testPayload interface{}, expectedBytes []byte, expectError bool, expectedError string) {
			if payloadBytes, err := GetJSONPayload(testPayload); expectError {
				Expect(err).NotTo(BeNil())
				Expect(payloadBytes).To(BeNil())
				Expect(err.Error()).To(ContainSubstring(expectedError))
			} else {
				Expect(err).To(BeNil())
				Expect(payloadBytes).NotTo(BeNil())
				Expect(payloadBytes).To(Equal(expectedBytes))
			}
		},
		Entry("should return payload as bytes-array when JSON-marshal is successful", dummyPayload, dummyPayloadBytes, false, nil),
		Entry("should return error when payload cannot be processed to byte-array", unprocessableStruct, nil, true, "could not parse request-payload -interface to bytes"),
	)
	DescribeTable("GetRequest",
		func(requestPath string, requestMethod string, requestPayload []byte, expectedRequestUrl string, expectError bool, expectedError string) {
			if req, err := GetRequest(requestPath, requestMethod, requestPayload); expectError {
				Expect(req).To(BeNil())
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal(expectedError))
			} else {
				Expect(err).To(BeNil())
				Expect(req).NotTo(BeNil())
				Expect(req.URL.String()).To(Equal(expectedRequestUrl))
			}
		},
		Entry("should return 'GET'-request with valid parameters", "hello", "GET", nil, "https://www.veikkaus.fi/api/bff/v1/hello", false, nil),
		Entry("should return 'POST'-request with valid parameters", "bar", "POST", dummyPayloadBytes, "https://www.veikkaus.fi/api/bff/v1/bar", false, nil),
		Entry("should return 'PUT'-request with valid parameters", "foo/bar/baz", "PUT", dummyPayloadBytes, "https://www.veikkaus.fi/api/bff/v1/foo/bar/baz", false, nil),
		Entry("should return error for invalid 'POST'-request", "irrelevant", "POST", invalidPayloadBytes, nil, true, "payload bytes were expected, received nil"),
		Entry("should return error for invalid 'PUT'-request", "/", "PUT", invalidPayloadBytes, nil, true, "payload bytes were expected, received nil"),
	)
	DescribeTable("SetCookies",
		func(requestCookies *RequestCookies, newCookies []*http.Cookie, originalCookieCount, expectedNumberOfCookies int) {
			Expect(len(requestCookies.CookieSlice)).To(Equal(originalCookieCount))
			requestCookies.SetCookies(dummyURL, newCookies)
			Expect(len(requestCookies.CookieSlice)).To(Equal(expectedNumberOfCookies))
		},
		Entry("should add new cookies to empty cookie-slice", &RequestCookies{}, dummyCookies, NoCookiesToday, DummyCookieCount),
		Entry(
			"should add new cookies to non-empty cookie-slice",
			&RequestCookies{
				CookieSlice: []*http.Cookie{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
			},
			dummyCookies,
			SingleCookieInSlice,
			AppendedCookieCount,
		),
	)
	DescribeTable("SetCookie",
		func(requestCookies *RequestCookies, newCookie *http.Cookie, originalCookieCount, expectedNumberOfCookies int) {
			Expect(len(requestCookies.CookieSlice)).To(Equal(originalCookieCount))
			requestCookies.SetCookie(dummyURL, newCookie)
			Expect(len(requestCookies.CookieSlice)).To(Equal(expectedNumberOfCookies))
		},
		Entry("should add new cookie to empty cookie-slice", &RequestCookies{}, dummyCookie2, NoCookiesToday, SingleCookieInSlice),
		Entry(
			"should add new cookie to non-empty cookie-slice",
			&RequestCookies{
				CookieSlice: []*http.Cookie{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
			},
			dummyCookie1,
			SingleCookieInSlice,
			DummyCookieCount,
		),
	)
	DescribeTable("Cookies",
		func(requestCookies *RequestCookies, expectedCookies []*http.Cookie) {
			cookies := requestCookies.Cookies(dummyURL)
			Expect(cookies).To(Equal(expectedCookies))
		},
		Entry("should return cookie-array from cookie-slice", &RequestCookies{CookieSlice: dummyCookies}, dummyCookies),
	)
	Describe("IsAuthenticated", func() {
		var (
			requestCookies *RequestCookies
		)

		BeforeEach(func() {
			requestCookies = &RequestCookies{
				CookieSlice: []*http.Cookie{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
			}
		})
		DescribeTable("IsAuthenticated",
			func(cookies []*http.Cookie, expected bool) {
				requestCookies.SetCookies(dummyURL, cookies)
				result := requestCookies.IsAuthenticated(dummyURL)
				Expect(result).To(Equal(expected))
			},
			Entry("has no cookies", []*http.Cookie{}, false),
			Entry("has authenticated cookie present", []*http.Cookie{{Name: AuthSessionCookie}}, true),
			Entry("has authenticated cookie not present", []*http.Cookie{{Name: "other_cookie"}}, false),
		)
	})

})
