package goveikkaus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type dummyResponseType struct {
	Foo string
}

var _ = Describe("goveikkaus", func() {
	var client *Client
	var mux *http.ServeMux
	var serverURL string
	var teardown func()

	BeforeEach(func() {
		client, mux, serverURL, teardown = setup()
	})

	AfterEach(func() {
		defer teardown()
	})

	Describe("Client", func() {
		It("should return HTTP Client with c.Client pointing to different client pointer", func() {
			newClient := NewClient(nil)
			client := newClient.Client()
			url, _ := url.Parse(serverURL)
			Expect(newClient).NotTo(Equal(client))
			Expect(newClient.client.Jar.Cookies(url)).To(HaveLen(0))
		})
	})
	Describe("NewClient", func() {
		It("should return new client with default constants as baseurl and user agent", func() {
			newClient := NewClient(nil)
			client := NewClient(nil)

			Expect(newClient.BaseURL.String()).To(Equal(api.BaseURL))
			Expect(newClient.UserAgent).To(Equal(api.UserAgent))
			Expect(newClient.client).NotTo(Equal(client))
		})
	})
	Describe("isContextOrUrlError", func() {
		It("returns context error when context.Done() detected", func() {
			canceledCtx, cancel := context.WithCancel(context.Background())
			// Immediately cancel the context
			cancel()
			err := isContextOrURLError(canceledCtx, nil)

			Expect(err).To(Equal(canceledCtx.Err()))
		})
		It("returns url.Parse error when url parsing fails", func() {
			ctx := context.Background()
			err := isContextOrURLError(ctx, &url.Error{
				URL: "invalid url",
				Err: errors.New("invalid url"),
			})

			Expect(err).To(BeAssignableToTypeOf(&url.Error{}))
		})
		It("returns nil when error is not cancelled/done context or url-parsing error, hinting towards error response from Veikkaus API", func() {
			ctx := context.Background()
			err := isContextOrURLError(ctx, nil)

			Expect(err).To(BeNil())
		})
	})
	DescribeTable("isAuthorizedCall",
		func(optionalBoolParams []bool, expectedResult bool) {
			Expect(isAuthorizedCall(optionalBoolParams)).To(Equal(expectedResult))
		},
		Entry("should return true when first element is true", []bool{true}, true),
		Entry("should return false when first element is false (not used / needed in practice but test-coverage is added)", []bool{false, true}, false),
		Entry("should return first boolean value when more values are passed", []bool{true, false, true}, true),
		Entry("should return false when boolean-array is nil / empty", nil, false),
	)
	Describe("do", func() {
		It(fmt.Sprintf("Can do bare HTTP-request to %s", serverURL), func() {
			expectedBody := "dummy response"

			mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, expectedBody)
			})

			req, err := api.GetRequest("hello", "GET", nil)
			Expect(err).To(BeNil())

			ctx := context.Background()
			resp, err := client.do(ctx, req)
			Expect(err).To(BeNil())
			Expect(resp).NotTo(BeNil())

			got, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(got)).To(Equal(expectedBody))

			if err := resp.Body.Close(); err != nil {
				log.Fatalf("could not close response body. ERR: %v", err)
			}
		})
		It("returns error when nil context is passed", func() {
			// defer teardown()
			req, _ := api.GetRequest("hello", "GET", nil)
			_, err := client.do(nil, req) //lint:ignore SA1012 ignoring this for unit-test purposes

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring(errNonNilContext.Error()))
		})
		It("returns error when context is getting cancelled in the middle", func() {
			canceledCtx, cancel := context.WithCancel(context.Background())
			// Immediately cancel the context
			cancel()

			req, _ := api.GetRequest("foobar", "GET", nil)
			_, err := client.do(canceledCtx, req)

			Expect(err).To(Equal(canceledCtx.Err()))
		})
		It("returns url.Parse error when url parsing fails", func() {
			expectedBody := "dummy response"

			mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, expectedBody)
			})

			req, _ := http.NewRequest(http.MethodGet, "invalid url here", nil)

			ctx := context.Background()
			resp, err := client.do(ctx, req)
			Expect(resp).To(BeNil())

			Expect(err).To(BeAssignableToTypeOf(&url.Error{}))
		})
		It("returns error when response status code was not in supported status codes", func() {
			expectedBody := "error with bar"

			mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMovedPermanently)
				fmt.Fprint(w, expectedBody)
			})

			req, err := api.GetRequest("foo", "GET", nil)
			Expect(err).To(BeNil())

			ctx := context.Background()
			resp, err := client.do(ctx, req)

			Expect(err).NotTo(BeNil())
			Expect(err).To(BeAssignableToTypeOf(&api.UnsupportedStatusCodeError{}))

			Expect(resp).To(BeNil())
		})
		It("returns error when api responds with a known API-error", func() {
			errorResponse := []byte(`{"code":"NOT_AUTHENTICATED", "fieldErrors":[]}`)

			mux.HandleFunc("/youshallnotpass", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				if _, err := w.Write(errorResponse); err != nil {
					log.Fatalf("could not write response-body in unit-test: %v", err)
				}
			})

			req, err := api.GetRequest("youshallnotpass", "GET", nil)
			Expect(err).To(BeNil())

			ctx := context.Background()
			resp, err := client.do(ctx, req)

			Expect(resp).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(&api.UnauthorizedError{}))
		})
		It("returns error when user is not logged in for authorized call", func() {
			errorResponse := []byte(`{"code":"NOT_AUTHENTICATED", "fieldErrors":[]}`)

			mux.HandleFunc("/youshallnotpass", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				if _, err := w.Write(errorResponse); err != nil {
					log.Fatalf("could not write response-body in unit-test: %v", err)
				}
			})

			req, err := api.GetRequest("youshallnotpass", "GET", nil)
			Expect(err).To(BeNil())

			ctx := context.Background()
			isAuthorized := true
			resp, err := client.do(ctx, req, isAuthorized)

			Expect(resp).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(&api.UserNotLoggedInError{}))
		})
	})
	Describe("Do", func() {
		It("should handle happy case fine", func() {
			expectedBody := dummyResponseType{
				Foo: "bar",
			}

			var responseInterface dummyResponseType

			mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
				bytes, err := json.Marshal(&expectedBody)
				if err != nil {
					log.Panicf("should have been able to marshal expectedBody to bytes. ERR: %v", err)
				}
				if _, err = w.Write(bytes); err != nil {
					log.Panicf("could not write response-body in unit-test %v", err)
				}
			})

			req, err := api.GetRequest("hello", "GET", nil)
			Expect(err).To(BeNil())

			ctx := context.Background()

			resp, err := client.Do(ctx, req, &responseInterface)
			Expect(err).To(BeNil())
			Expect(responseInterface.Foo).To(Equal(expectedBody.Foo))

			if err := resp.Body.Close(); err != nil {
				log.Fatalf("could not close response body. ERR: %v", err)
			}
		})
		It("returns error when api responds with a known API-error", func() {
			var v interface{}
			errorResponse := []byte(`{"code":"NOT_AUTHENTICATED", "fieldErrors":[]}`)

			mux.HandleFunc("/youshallnotpass", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				if _, err := w.Write(errorResponse); err != nil {
					log.Fatalf("could not write the response body in unit-test %v", err)
				}
			})

			req, err := api.GetRequest("youshallnotpass", "GET", nil)
			Expect(err).To(BeNil())

			ctx := context.Background()
			resp, err := client.Do(ctx, req, &v)

			Expect(resp).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(&api.UnauthorizedError{}))
		})
	})
	Describe("initialize", func() {
		It("should initialize client when client is nil", func() {
			client := &Client{}
			client.initialize()

			// Assert that client is not-nil and has been initialized
			Expect(client.client).NotTo(BeNil())
			// Assert that the client's Jar is of the correct type
			Expect(client.client.Jar).To(BeAssignableToTypeOf(&api.RequestCookies{}))
		})
		It("should not overwrite clients Jar when client is nil", func() {
			client := &http.Client{}
			veikkausClient := &Client{client: client}

			// initialize the client
			veikkausClient.initialize()

			Expect(veikkausClient.client).To(Equal(client))
		})
	})
})
