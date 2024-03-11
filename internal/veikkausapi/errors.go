package veikkausapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	InputValidationFailed ErrorCode = "INPUT_VALIDATION_FAILED"
	NotAuthenticated      ErrorCode = "NOT_AUTHENTICATED"
	Unknown               ErrorCode = "UNKNOWN"
)

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (f *FieldError) ToString() string {
	return fmt.Sprintf("Field '%s' had issue: '%s'", f.Field, f.Message)
}

func getFieldErrors(fieldErrors []FieldError) []string {
	var errors []string

	for _, fieldErr := range fieldErrors {
		errors = append(errors, fieldErr.ToString())
	}
	return errors
}

type ErrorResponse struct {
	Response    *http.Response `json:"-"`
	FieldErrors []FieldError   `json:"fieldErrors"`
	Code        ErrorCode      `json:"code"`
}

type UserNotLoggedInError struct{}

func (e *UserNotLoggedInError) Error() string {
	return "No Authenticated session active, user not logged in"
}

type UnsupportedStatusCodeError struct {
	Code int
}

func (err *UnsupportedStatusCodeError) Error() string {
	return fmt.Sprintf("response status code was not in allowed range (200-299). Got %d", err.Code)
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf(e.Message)
}

type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	if fmt.Sprint(e.Errors) == "[]" {
		return "input validation error"
	}
	return fmt.Sprintf("input validation errors: %v", e.Errors)
}

type APIErrorNotImplementedError struct {
	Code        ErrorCode
	Message     string
	FieldErrors []string
}

func (e *APIErrorNotImplementedError) Error() string {
	return fmt.Sprintf("API Returned error that has not been implemented in this library. Error code was '%s'", e.Code)
}

type RequestPayloadError struct {
	Message string
}

func (e *RequestPayloadError) Error() string {
	return "could not parse request-payload -interface to bytes"
}

func ParseAPIError(rawJSON []byte) error {
	var response ErrorResponse

	if err := json.Unmarshal(rawJSON, &response); err != nil {
		return err // Failed to parse JSON
	}

	switch response.Code {
	case NotAuthenticated:
		return &UnauthorizedError{Message: "User not authenticated or login failed"}
	case InputValidationFailed:
		return &ValidationError{Errors: getFieldErrors(response.FieldErrors)}
	default:
		return &APIErrorNotImplementedError{
			Code:        response.Code,
			FieldErrors: getFieldErrors(response.FieldErrors),
			Message:     "Unsupported API Error",
		}
	}
}
