package veikkausapi

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MockedReader is a custom io.Reader implementation that returns an error when its Read method is called.
type MockedReader struct{}

// Read always returns an error.
func (mr *MockedReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mocked error")
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
	Body: io.NopCloser(strings.NewReader("")),
}

var notCoolResponseBody = SomeNotSoCoolType{NotCool: "at all", Thing: 0.333}
var notCoolResponse = &http.Response{
	Body: io.NopCloser(strings.NewReader(`{"NotCool": "not so cool thing", "Thing": 3.333}`)),
}

// var responseBody =
var responseBody = getMockResponse(`{"key": "value"}`)

func getMockResponse(body string) *http.Response {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write the sample response body to the response writer
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(body)); err != nil {
			log.Panic("Could not write dummy response body to stream")
		}
	}))
	defer mockServer.Close()

	// Create a mock HTTP response using the mock server URL
	mockResponse := &http.Response{
		Body: io.NopCloser(strings.NewReader(body)),
	}

	return mockResponse
}

var _ = Describe("internal/veikkausapi: response-handler code", func() {
	DescribeTable("ResponseCodeIsOk",
		func(statusCode int, expectedResult bool) {
			actualResult := ResponseCodeIsOk(statusCode)
			Expect(actualResult).To(Equal(expectedResult))
		},
		Entry("should return 'true' for 200 (OK)", 200, true),
		Entry("should return 'true' for 201 (Created)", 201, true),
		Entry("should return 'true' for 202 (Accepted)", 202, true),
		Entry("should return 'true' for 203 (Non-Authorative Information)", 203, true),
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
					Expect(err.Error()).To(Equal(expectedError))
				}
			} else {
				if err := HandleResponse(response, result); expectError {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal(expectedError))
				} else {
					Expect(err).To(BeNil())
					Expect(result).To(Equal(expectedResult))
				}
			}
		},
		Entry("should return parsed response when everything succeeds", responseBody, "value", "value", false, false, nil),
		Entry("should return error when response parsing fails for empty response body", emptyResponseBody, notCoolResponseBody, nil, false, true, "Error unmarshaling response body: unexpected end of JSON input"),
		Entry("should return error when response parsing fails for incorrect result-interface", notCoolResponse, &coolResponse, nil, false, true, "Error unmarshaling response body: json: cannot unmarshal number 3.333 into Go struct field SomeCoolType.Thing of type int"),
		Entry("should return error when io.ReadAll fails", responseBody, "value", "value", true, true, "Error reading the response body: mocked error"),
	)
})
