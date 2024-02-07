package veikkausapi

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	ExpectedLoginUrl = "https://www.veikkaus.fi/api/bff/v1/sessions"
)

type DummyPayload struct {
	Foo string
	Bar string
	Baz string
}

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

func mockNewRequestError(method, requestUrl string, body io.Reader) (*http.Request, error) {
	return nil, errors.New("mocked error")
}

type TestStruct struct {
	Name  string
	Age   int
	Valid bool
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
		It(fmt.Sprintf("should return request url with '%s' as base-url and given path", VeikkausApiBaseUrl), func() {
			Expect(getRequestUrl("sessions")).To(Equal(ExpectedLoginUrl))
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
	DescribeTable("getJsonPayload",
		func(payloadStruct interface{}, expectedBytes []byte, expectError bool, expectedError string) {
			if bytes, err := getJsonPayload(payloadStruct); expectError {
				Expect(bytes).To(BeNil())
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal(expectedError))
			} else {
				Expect(err).To(BeNil())
				Expect(bytes).NotTo(BeNil())
				Expect(bytes).To(Equal(expectedBytes))
			}
		},
		Entry("should format struct-type object into byte-array payload", dummyPayload, dummyPayloadBytes, false, nil),
		Entry("should return error when JSON Marshalling fails", unprocessableStruct, nil, true, "Unable to marshal the request body to JSON"),
	)
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
		isError := validateRequestUrl(requestUrl)
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
		Entry("should return error for POST-request without given payload", "https://amazingapi.com", "POST", nil, true, false, "Payload bytes were expected, received nil"),
		Entry("should return error for PUT-request without given payload", "https://is.fi", "PUT", nil, true, false, "Payload bytes were expected, received nil"),
		Entry("should return error for unsupported method", "https://is.fi", "GET", nil, true, false, "Expected to receive one of [|PUT||POST|] request-type, received 'GET'"),
		Entry("should return error for malformatted url", "˛∞é®§", "PUT", dummyPayloadBytes, true, false, "Request URL was malformatted. ERR: parse \"˛∞é®§\": invalid URI for request"),
		Entry("should return error when http.NewRequest unexpectetly fails", "https://localhost:8080", "PUT", dummyPayloadBytes, true, true, "Error creating 'PUT'-request for https://localhost:8080: mocked error"),
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
		Entry("should return error when http.NewRequest unexpectetly returns error", "https://foobar.com/api/v2/baz", true, true, "Error creating 'GET'-request for https://foobar.com/api/v2/baz: mocked error"),
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
		Entry("should return error when unsupported method is attempted", "https://localhost:8080", "DEL", nil, true, "Unsupported method 'DEL' provided."),
		Entry("should return error when request url is invalid for 'POST'", "www.*****", "POST", dummyPayloadBytes, true, "Request URL was malformatted. ERR: parse \"www.*****\": invalid URI for request"),
		Entry("should return error when request url is invalid for 'PUT'", "**®é¸ƒ", "PUT", dummyPayloadBytes, true, "Request URL was malformatted. ERR: parse \"**®é¸ƒ\": invalid URI for request"),
		Entry("should return error when request url is invalid for 'GET'", "ww.notvalid.d", "GET", nil, true, "Request URL was malformatted. ERR: parse \"ww.notvalid.d\": invalid URI for request"),
	)
	DescribeTable("GetJsonPayload",
		func(testPayload interface{}, expectedBytes []byte, expectError bool, expectedError string) {
			if payloadBytes, err := GetJsonPayload(testPayload); expectError {
				Expect(err).NotTo(BeNil())
				Expect(payloadBytes).To(BeNil())
				Expect(err.Error()).To(Equal(expectedError))
			} else {
				Expect(err).To(BeNil())
				Expect(payloadBytes).NotTo(BeNil())
				Expect(payloadBytes).To(Equal(expectedBytes))
			}
		},
		Entry("should return payload as bytes-array when JSON-marshal is successful", dummyPayload, dummyPayloadBytes, false, nil),
		Entry("should return error when payload cannot be processed to byte-array", unprocessableStruct, nil, true, "Could not parse struct to bytes. ERR: Unable to marshal the request body to JSON"),
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
		Entry("should return error for invalid 'POST'-request", "irrelevant", "POST", invalidPayloadBytes, nil, true, "Payload bytes were expected, received nil"),
		Entry("should return error for invalid 'PUT'-request", "/", "PUT", invalidPayloadBytes, nil, true, "Payload bytes were expected, received nil"),
	)
})
