package veikkausapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var AuthSessionCookie = "JSESSIONID"

type RequestCookies struct {
	CookieSlice []*http.Cookie
}

func (rc *RequestCookies) SetCookies(u *url.URL, cookies []*http.Cookie) {
	rc.CookieSlice = append(rc.CookieSlice, cookies...)
}

func (rc *RequestCookies) SetCookie(u *url.URL, cookie *http.Cookie) {
	rc.CookieSlice = append(rc.CookieSlice, cookie)
}

func (rc *RequestCookies) Cookies(u *url.URL) []*http.Cookie {
	return rc.CookieSlice
}

func (rc *RequestCookies) IsAuthenticated(url *url.URL) bool {
	for _, cookie := range rc.Cookies(url) {
		if cookie.Name == AuthSessionCookie {
			return true
		}
	}

	return false
}

// Unit-Test purposes
var newRequest = http.NewRequest

func GetBaseURL() string {
	if !OverWriteBaseURL {
		return BaseURL + "/" + VeikkausAPIBaseURL + VeikkausAPIVersion + "/"
	} else {
		return BaseURL
	}
}

func getRequestURL(path string) string {
	return GetBaseURL() + path
}

func setRequestHeaders(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Accept", Accept)
	req.Header.Set(RobotIdentifierHeaderKey, RobotIdentifierHeaderValue)

	return req
}

func validateRequestMethod(method string, allowedMethods []string) error {
	var errorString = "Expected to receive one of ["
	for _, allowedMethod := range allowedMethods {
		if method == allowedMethod {
			return nil
		}
		errorString = fmt.Sprintf("%s|%s|", errorString, allowedMethod)
	}
	errorString = fmt.Sprintf("%s] request-type, received '%s'", errorString, method)
	return errors.New(errorString)
}

func validateRequestURL(requestURL string) error {
	_, err := url.ParseRequestURI(requestURL)
	if err != nil {
		errorString := fmt.Sprintf("Request URL was malformatted. ERR: %v", err)
		return errors.New(errorString)
	}
	return nil
}

func handlePutPost(requestURL, method string, jsonPayload []byte) (*http.Request, error) {
	if err := validateRequestMethod(method, []string{http.MethodPut, http.MethodPost}); err != nil {
		return nil, err
	}

	if err := validateRequestURL(requestURL); err != nil {
		return nil, err
	}

	if jsonPayload == nil {
		return nil, fmt.Errorf("payload bytes were expected, received nil")
	}

	req, err := newRequest(method, requestURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		errorString := fmt.Sprintf("Error creating '%s'-request for %s: %s", method, requestURL, err)
		return nil, errors.New(errorString)
	}

	return req, nil
}

func handleGet(requestURL string) (*http.Request, error) {
	if err := validateRequestURL(requestURL); err != nil {
		return nil, err
	}

	req, err := newRequest("GET", requestURL, nil)
	if err != nil {
		errorString := fmt.Sprintf("Error creating 'GET'-request for %s: %s", requestURL, err)
		return nil, errors.New(errorString)
	}
	return req, nil
}

func requestHandler(requestURL, requestMethod string, jsonPayload []byte) (*http.Request, error) {
	switch requestMethod {
	case http.MethodPut:
		fallthrough
	case http.MethodPost:
		return handlePutPost(requestURL, requestMethod, jsonPayload)
	case http.MethodGet:
		return handleGet(requestURL)
	default:
		return nil, fmt.Errorf("unsupported method '%s' provided", requestMethod)
	}
}

// JSONMarshal is stored to a variable for unit-test purposes
var JSONMarshal = json.Marshal

func GetJSONPayload(payload interface{}) ([]byte, error) {
	bytes, err := JSONMarshal(payload)
	if err != nil {
		return nil, &RequestPayloadError{Message: err.Error()}
	}

	return bytes, nil
}

func GetRequest(requestPath string, requestMethod string, requestPayloadBytes []byte) (*http.Request, error) {
	var req *http.Request
	var err error

	url := getRequestURL(requestPath)

	if req, err = requestHandler(url, requestMethod, requestPayloadBytes); err != nil {
		return nil, err
	}

	return setRequestHeaders(req), nil
}

func WithContext(ctx context.Context, req *http.Request) *http.Request {
	return req.WithContext(ctx)
}
