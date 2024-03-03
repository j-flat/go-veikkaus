package goveikkaus

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

func ReturnError(v any) ([]byte, error) {
	return nil, errors.New("mocked error during marshaling")
}

var validationErrors = api.ErrorResponse{
	Code: "INPUT_VALIDATION_FAILED",
	FieldErrors: []api.FieldError{
		{
			Field:   "login",
			Code:    "EMPTY",
			Message: "may not be empty",
		},
		{
			Field:   "password",
			Code:    "EMPTY",
			Message: "may not be empty",
		},
	},
}

func getInputValidationErrorBytes() []byte {
	bytes, err := json.Marshal(&validationErrors)
	Expect(err).To(BeNil())
	return bytes
}

var _ = Describe("auth login: authservice login", func() {
	var username = "johndoe"
	var password = "verysecret"
	var sucessfulRequestBytes = []byte(`{}`)
	var unknownErrorBytes = []byte(`{"code": "UNKNOWN", "fieldErrors": []}`)
	var unauthorizedErrorBytes = []byte(`{"code":"NOT_AUTHENTICATED", "fieldErrors":[]}`)
	// var unsupportedErrorResposeBytes = []byte(`{"errorcode":"NOT_AUTHENTICATED", "fieldError":""}`)

	var client *Client
	var mux *http.ServeMux
	// var serverURL string
	var teardown func()

	BeforeEach(func() {
		client, mux, _, teardown = setup()
	})

	AfterEach(func() {
		defer teardown()
	})

	DescribeTable("Login",
		func(shouldSucceed bool, expectedStatusCode int, expectedErr error, expectedResponseBody []byte) {
			mux.HandleFunc("/"+api.LoginEndpoint, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(expectedStatusCode)
				if _, err := w.Write(expectedResponseBody); err != nil {
					log.Fatalf("Error while writing the response body in test: %v", err)
				}
			})

			ctx := context.Background()
			data, _, err := client.Auth.Login(ctx, username, password)

			if shouldSucceed {
				Expect(data).To(BeAssignableToTypeOf(&LoginSuccessful{}))
				Expect(err).To(BeNil())

			} else {
				Expect(err).To(BeAssignableToTypeOf(expectedErr))
				Expect(data).To(BeNil())
			}
		},
		Entry("should login successfully in the happy case", true, http.StatusOK, nil, sucessfulRequestBytes),
		Entry("should return error if login is unsuccessful due to wrong credentials", false, http.StatusUnauthorized, &api.UnauthorizedError{}, unauthorizedErrorBytes),
		Entry("should return error when input validation fails", false, http.StatusBadRequest, &api.ValidationError{}, getInputValidationErrorBytes()),
		Entry("should return error when response status code is unsupported, but response is otherwise successful", false, http.StatusMovedPermanently, &api.UnsupportedStatusCodeError{}, sucessfulRequestBytes),
		Entry("should return error when response errored and code is unknown", false, http.StatusBadRequest, &api.APIErrorNotImplementedError{}, unknownErrorBytes),
		Entry("should return error for nil response body", false, http.StatusServiceUnavailable, &json.SyntaxError{}, nil),
		Entry("should return error when error-response from API is not in supported format", false, http.StatusInternalServerError, &json.SyntaxError{}, nil),
	)
	Describe("Login", func() {
		It("should return error when request-payload byte conversion fails", func() {
			api.JSONMarshal = ReturnError
			mux.HandleFunc("/"+api.LoginEndpoint, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				if _, err := w.Write(nil); err != nil {
					log.Fatalf("could not write response body in unit-test: %v", err)
				}
			})

			ctx := context.Background()
			data, _, err := client.Auth.Login(ctx, username, password)

			Expect(err).To(BeAssignableToTypeOf(&api.RequestPayloadError{}))
			Expect(data).To(BeNil())
		})
	})
})
