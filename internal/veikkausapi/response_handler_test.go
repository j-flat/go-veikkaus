package veikkausapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MockedReader is a custom io.Reader implementation that returns an error when its Read method is called.
type MockedReader struct {
	data []byte
	err  error
	idx  int
}

// Read always returns an error.
func (mr *MockedReader) Read(p []byte) (n int, err error) {
	errMsg := fmt.Sprintf("mocked errror: %v, %v, %d", mr.data, mr.err, mr.idx)
	return 0, errors.New(errMsg)
}

func (mr *MockedReader) ReadAll(p []byte) (n int, err error) {
	return 0, errors.New("mocked error")
}

func (mr *MockedReader) Close() error {
	return nil
}

type SomeCoolType struct {
	Cool  string
	Thing int
}

type SomeNotSoCoolType struct {
	NotCool string
	Thing   float32
}

var coolResponse = SomeCoolType{Cool: "stuff", Thing: 1}
var emptyResponseBody = &http.Response{
	StatusCode: http.StatusBadRequest,
	Body:       io.NopCloser(strings.NewReader("")),
}

var notCoolResponseBody = SomeNotSoCoolType{NotCool: "at all", Thing: 0.333}
var notCoolResponse = &http.Response{
	Body: io.NopCloser(strings.NewReader(`{"NotCool": "not so cool thing", "Thing": 3.333}`)),
}

// var responseBody =
var responseBody = getMockResponse(`{"key": "value"}`, 200)

func getMockResponse(body string, statusCode int) *http.Response {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write the sample response body to the response writer
		w.WriteHeader(statusCode)
		if _, err := w.Write([]byte(body)); err != nil {
			log.Panic("Could not write dummy response body to stream")
		}
	}))
	defer mockServer.Close()

	// Create a mock HTTP response using the mock server URL
	mockResponse := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	return mockResponse
}

var _ = Describe("internal/veikkausapi: response-handler code", func() {
	DescribeTable("ResponseCodeIsOk",
		func(statusCode int, expectedResult bool) {
			resp := getMockResponse(`{"key": "value"}`, statusCode)
			isOk := ResponseCodeIsOk(resp)

			Expect(isOk).To(Equal(expectedResult))
		},
		Entry("should return 'true' for 200 (OK)", 200, true),
		Entry("should return 'true' for 201 (Created)", 201, true),
		Entry("should return 'true' for 202 (Accepted)", 202, true),
		Entry("should return 'true' for 203 (Non-Authoritative Information)", 203, true),
		Entry("should return 'true' for 204 (No Content)", 204, true),
		Entry("should return 'true' for 205 (Reset Content)", 205, true),
		Entry("should return 'true' for 206 (Partial Content)", 206, true),
		Entry("should return 'true' for 207 (Multi-Status)", 207, true),
		Entry("should return 'true' for 208 (Already Reported)", 208, true),
		Entry("should return 'true' for 226 (IM Used)", 226, true),
		Entry("should return 'false' for 300 (Multiple Choices)", 300, false),
		Entry("should return 'false' for 301 (Moved Permanently)", 301, false),
		Entry("should return 'false' for 302 (Found)", 302, false),
		Entry("should return 'false' for 304 (Not Modified)", 304, false),
		Entry("should return 'false' for 307 (Temporary Redirect)", 307, false),
		Entry("should return 'false' for 308 (Permanent Redirect)", 308, false),
		Entry("should return 'false' for 400 (Bad Request)", 400, false),
		Entry("should return 'false' for 401 (Unauthorized)", 401, false),
		Entry("should return 'false' for 402 (Payment Required)", 402, false),
		Entry("should return 'false' for 403 (Forbidden)", 403, false),
		Entry("should return 'false' for 404 (Not Found)", 404, false),
		Entry("should return 'false' for 405 (Method Not Allowed)", 405, false),
		Entry("should return 'false' for 408 (Request Timeout)", 408, false),
		Entry("should return 'false' for 410 (Gone)", 410, false),
		Entry("should return 'false' for 500 (Internal Server Error)", 500, false),
		Entry("should return 'false' for 501 (Not Implemented)", 501, false),
		Entry("should return 'false' for 502 (Bad Gateway)", 502, false),
		Entry("should return 'false' for 503 (Service Unavailable)", 503, false),
		Entry("should return 'false' for 504 (Gateway Timeout)", 504, false),
	)
	DescribeTable("HandleResponse",
		func(response *http.Response, result interface{}, expectedResult interface{}, mockIOReader bool, expectError bool, expectedError string) {
			if mockIOReader {
				// Create a mock HTTP response with a mockedReader
				mockResponse := &http.Response{
					Body: &MockedReader{},
				}

				if err := HandleResponse(mockResponse, nil); err != nil {
					Expect(err.Error()).To(ContainSubstring(expectedError))
				}
			} else {
				if err := HandleResponse(response, result); expectError {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal(expectedError))
				} else {
					Expect(err).To(BeNil())
					if expectedResult != nil {
						Expect(result).To(Equal(expectedResult))
					} else {
						Expect(result).To(BeNil())
					}
				}
			}
		},
		Entry("should return parsed response when everything succeeds", responseBody, nil, nil, false, false, nil),
		Entry("should return error when response parsing fails for empty response body", emptyResponseBody, notCoolResponseBody, nil, false, true, "error unmarshaling response body: unexpected end of JSON input"),
		Entry("should return error when response parsing fails for incorrect result-interface", notCoolResponse, &coolResponse, nil, false, true, "error unmarshaling response body: json: cannot unmarshal number 3.333 into Go struct field SomeCoolType.Thing of type int"),
		Entry("should return error when io.ReadAll fails", responseBody, "value", "value", true, true, "error reading the response body"),
	)
	DescribeTable("HandleError",
		func(response *http.Response, expectedError error) {
			err := HandleError(response)
			if expectedError != nil {
				Expect(err).To(BeAssignableToTypeOf(expectedError))
			}
			Expect(err).To(HaveOccurred())
		},
		Entry("should return error for failed byte-stream conversion", emptyResponseBody, &json.SyntaxError{}),
		Entry("should return unauthorized-error when unauthorized", &http.Response{
			Status:     "401 Unauthorized",
			StatusCode: 401,
			Body:       io.NopCloser(strings.NewReader(`{"code":"NOT_AUTHENTICATED", "fieldErrors":[]}`)),
		}, &UnauthorizedError{}),
		Entry("should return input validation error when input validation fails", &http.Response{
			Status:     "400 Bad Request",
			StatusCode: 400,
			Body:       io.NopCloser(strings.NewReader(`{"code":"INPUT_VALIDATION_FAILED","fieldErrors":[{"field":"username","code":"INVALID","message":"Username is invalid"}]}`)),
		}, &ValidationError{}),
		Entry("should return error for unknown api-error", &http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader(`{"code":"UNKNOWN_ERROR"}`)),
		}, &APIErrorNotImplementedError{}),
		Entry("should return error for unprocessable bytes", &http.Response{
			Status:     "600",
			StatusCode: 600,
			Body:       io.NopCloser(strings.NewReader("")),
		}, &json.SyntaxError{}),
		Entry("should return error for empty response body", &http.Response{
			Body: nil,
		}, nil),
		Entry("should return error for failed byte-stream conversion", &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       &MockedReader{err: errors.New("read error")},
		}, nil),
		Entry("should return 'UnsupportedStatusCodeError' for 3xx status-codes", &http.Response{
			StatusCode: http.StatusMovedPermanently,
			Body:       io.NopCloser(strings.NewReader("otherwise fine response body")),
		}, &UnsupportedStatusCodeError{}),
	)
})
