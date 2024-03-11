package veikkausapi

import (
	"encoding/json"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func returnUnmarshalError() error {
	var i interface{}
	err := json.Unmarshal(invalidPayloadBytes, &i)
	if err != nil {
		return err
	}

	return nil
}

const ExampleUnsupportedHTTPStatusCode = 666

var errUnauthorizedError error = errors.New("user is not authorized to perform such action")
var unknownErrorBytes = []byte(`{"code": "UNKNOWN", "fieldErrors": []}`)
var unauthorizedErrorBytes = []byte(`{"code":"NOT_AUTHENTICATED", "fieldErrors":[]}`)

func getValidationErrors() ErrorResponse {
	validationErrors := &ErrorResponse{
		Code: "INPUT_VALIDATION_FAILED",
		FieldErrors: []FieldError{
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

	return *validationErrors
}

func getInputValidationErrorBytes() []byte {
	errors := getValidationErrors()
	bytes, err := json.Marshal(errors)
	Expect(err).To(BeNil())
	return bytes
}

func getErrorFieldsAsStringArray() []string {
	var errors []string
	errs := getValidationErrors()
	for _, fieldErr := range errs.FieldErrors {
		errors = append(errors, fieldErr.ToString())
	}

	return errors
}

const ExpectedValidationErrorLength = 2

func getSampleValidationErrors() []string {
	errSampleValidationErrors := []string{}
	errSampleValidationErrors = append(errSampleValidationErrors, "value should be 1, got 2")
	errSampleValidationErrors = append(errSampleValidationErrors, "'foo' is not valid value for 'bar'")

	return errSampleValidationErrors
}

func getValidationErrorMessage() string {
	errors := getSampleValidationErrors()

	return fmt.Sprintf("input validation errors: %v", errors)
}

var _ = Describe("internal/veikkausapi: errors", func() {
	DescribeTable("ToString",
		func(fieldError *FieldError, expectedErrorStr string) {
			actualErrorStr := fieldError.ToString()
			Expect(actualErrorStr).To(Equal(expectedErrorStr))
		},
		Entry("should convert field error to string", &FieldError{Field: "login", Code: "EMPTY", Message: "login cannot be empty"}, "Field 'login' had issue: 'login cannot be empty'"),
	)
	DescribeTable("Custom Errors return expected error strings",
		func(customError error, expectedErrorMsg string) {
			Expect(customError.Error()).To(Equal(expectedErrorMsg))
		},
		Entry("should return 'UnsupportedStatusCode' error with HTTP-status code that caused the error", &UnsupportedStatusCodeError{Code: ExampleUnsupportedHTTPStatusCode}, fmt.Sprintf("response status code was not in allowed range (200-299). Got %d", ExampleUnsupportedHTTPStatusCode)),
		Entry("should return 'UnauthorizedError' with original error message", &UnauthorizedError{Message: errUnauthorizedError.Error()}, "user is not authorized to perform such action"),
		Entry("should return 'ValidationError' with static text when original error had empty list for attribute 'errors'", &ValidationError{Errors: nil}, "input validation error"),
		Entry("should return 'ValidationError' with all validation errors in the error string, when attribute 'errors' is not empty list", &ValidationError{Errors: getSampleValidationErrors()}, getValidationErrorMessage()),
		Entry("should return 'UserNotLoggedInError' when user is not logged in", &UserNotLoggedInError{}, "No Authenticated session active, user not logged in"),
		Entry("should return 'APIErrorNotImplementedError' when the API error is not known", &APIErrorNotImplementedError{Code: "TOO_JUICY"}, "API Returned error that has not been implemented in this library. Error code was 'TOO_JUICY'"),
	)
	DescribeTable("ParseAPIError",
		func(inputBytes []byte, expectedError error) {
			err := ParseAPIError(inputBytes)
			Expect(err).To(BeAssignableToTypeOf(expectedError))
		},
		Entry("should parse 'NOT_AUTHENTICATED' code to 'UnauthorizedError'", unauthorizedErrorBytes, &UnauthorizedError{Message: "User not authenticated or login failed"}),
		Entry("should parse complex validation error to 'ValidationError'", getInputValidationErrorBytes(), &ValidationError{Errors: getErrorFieldsAsStringArray()}),
		Entry("should return error for unprocessable bytes", invalidPayloadBytes, returnUnmarshalError()),
		Entry("should return generidc error for unknown API error", unknownErrorBytes, &APIErrorNotImplementedError{Code: "UNKNOWN"}),
	)
})
