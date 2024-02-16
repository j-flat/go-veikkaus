package veikkausapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Unit-Test purposes
var newRequest = http.NewRequest

func getRequestUrl(path string) string {
	return fmt.Sprintf("%s%s/%s", VeikkausApiBaseUrl, VeikkausAPIVersion, path)
}

func getJsonPayload(payloadStruct interface{}) ([]byte, error) {
	jsonPayload, err := json.Marshal(payloadStruct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, errors.New("unable to marshal the request body to JSON")
	}

	return jsonPayload, nil
}

func setRequestHeaders(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Accept", Accept)
	req.Header.Set(RobotIdentifierHeaderKey, RobotIdentifierHeaderValue)

	return req
}

func validateRequestMethod(method string, allowedMethods []string) error {
	var errorString string = "Expected to receive one of ["
	for _, allowedMethod := range allowedMethods {
		if method == allowedMethod {
			return nil
		}
		errorString = fmt.Sprintf("%s|%s|", errorString, allowedMethod)
	}
	errorString = fmt.Sprintf("%s] request-type, received '%s'", errorString, method)
	return errors.New(errorString)
}

func validateRequestUrl(requestUrl string) error {
	_, err := url.ParseRequestURI(requestUrl)
	if err != nil {
		errorString := fmt.Sprintf("Request URL was malformatted. ERR: %v", err)
		return errors.New(errorString)
	}
	return nil
}

func handlePutPost(requestUrl, method string, jsonPayload []byte) (*http.Request, error) {
	if err := validateRequestMethod(method, []string{Put, Post}); err != nil {
		return nil, err
	}

	if err := validateRequestUrl(requestUrl); err != nil {
		return nil, err
	}

	if jsonPayload == nil {
		return nil, errors.New("Payload bytes were expected, received nil")
	}

	req, err := newRequest(method, requestUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		errorString := fmt.Sprintf("Error creating '%s'-request for %s: %s", method, requestUrl, err)
		return nil, errors.New(errorString)
	}

	return req, nil
}

func handleGet(requestUrl string) (*http.Request, error) {
	if err := validateRequestUrl(requestUrl); err != nil {
		return nil, err
	}

	req, err := newRequest("GET", requestUrl, nil)
	if err != nil {
		errorString := fmt.Sprintf("Error creating 'GET'-request for %s: %s", requestUrl, err)
		return nil, errors.New(errorString)
	}
	return req, nil
}

func requestHandler(requestUrl, requestMethod string, jsonPayload []byte) (*http.Request, error) {
	switch requestMethod {
	case Put:
		fallthrough
	case Post:
		return handlePutPost(requestUrl, requestMethod, jsonPayload)
	case Get:
		return handleGet(requestUrl)
	default:
		return nil, fmt.Errorf("Unsupported method '%s' provided.", requestMethod)
	}
}

func GetJsonPayload(payload interface{}) ([]byte, error) {
	payloadBytes, err := getJsonPayload(payload)

	if err != nil {
		return nil, fmt.Errorf("Could not parse struct to bytes. ERR: %v", err)
	}

	return payloadBytes, nil
}

func GetRequest(requestPath string, requestMethod string, requestPayloadBytes []byte) (*http.Request, error) {
	var req *http.Request
	var err error

	url := getRequestUrl(requestPath)

	if req, err = requestHandler(url, requestMethod, requestPayloadBytes); err != nil {
		return nil, err
	}

	return setRequestHeaders(req), nil
}
