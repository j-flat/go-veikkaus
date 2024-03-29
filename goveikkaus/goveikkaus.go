package goveikkaus

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

const (
	Version = "v0.4.0"
)

var errNonNilContext = errors.New("context must be non-nil")

type service struct {
	apiClient *Client
}

type Client struct {
	clientMu sync.Mutex
	client   *http.Client

	// Base URL for API requests.
	// Base URL should end with the trailing slash
	BaseURL *url.URL

	// User Agent to use when communicating with Veikkaus JSON API
	UserAgent string

	common         service
	SessionTimeout time.Time

	// Services used for interacting with different endpoints on Veikkaus API
	Auth     *AuthService
	Glossary *GlossaryService
}

func (veikkausClient *Client) UserIsLoggedIn() bool {
	return veikkausClient.Auth.AuthSessionIsActive()
}

func (veikkausClient *Client) Client() *http.Client {
	veikkausClient.clientMu.Lock()
	defer veikkausClient.clientMu.Unlock()
	clientCopy := *veikkausClient.client
	return &clientCopy
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Jar: &api.RequestCookies{},
		}
	}

	httpClient2 := *httpClient

	veikkausClient := &Client{client: &httpClient2}
	veikkausClient.initialize()

	return veikkausClient
}

func (veikkausClient *Client) initialize() {
	if veikkausClient.client == nil {
		veikkausClient.client = &http.Client{
			Jar: &api.RequestCookies{},
		}
	}
	if veikkausClient.BaseURL == nil {
		veikkausClient.BaseURL, _ = url.Parse(api.BaseURL)
	}
	if veikkausClient.UserAgent == "" {
		veikkausClient.UserAgent = api.UserAgent
	}

	veikkausClient.common.apiClient = veikkausClient
	veikkausClient.Auth = (*AuthService)(&veikkausClient.common)
	veikkausClient.Glossary = (*GlossaryService)(&veikkausClient.common)
}

func isContextOrURLError(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if urlErr, ok := err.(*url.Error); ok {
		return urlErr
	}

	return nil
}

func isAuthorizedCall(isAuthorizedCall []bool) bool {
	authorizedCall := false
	if len(isAuthorizedCall) > 0 {
		authorizedCall = isAuthorizedCall[0]
	}

	return authorizedCall
}

func (veikkausClient *Client) do(ctx context.Context, req *http.Request, authorizedCall ...bool) (*http.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	if isAuthorizedCall(authorizedCall) && !veikkausClient.UserIsLoggedIn() {
		return nil, &api.UserNotLoggedInError{}
	}

	req = api.WithContext(ctx, req)

	resp, err := veikkausClient.client.Do(req)
	if err = isContextOrURLError(ctx, err); err != nil {
		return nil, err
	}

	if !api.ResponseCodeIsOk(resp) {
		defer resp.Body.Close()
		return nil, api.HandleError(resp)
	}

	return resp, err
}

func (veikkausClient *Client) Do(ctx context.Context, req *http.Request, responseInterface interface{}) (*http.Response, error) {
	resp, err := veikkausClient.do(ctx, req)
	if err != nil {
		return resp, err
	}

	defer resp.Body.Close()

	err = api.HandleResponse(resp, &responseInterface)

	return resp, err
}
