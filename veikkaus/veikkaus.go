package goveikkaus

import (
	"net/http"
	"net/url"
	"sync"

	"go-veikkaus/internal/veikkausapi"
)

const (
	Version = "v0.0.1"

	// API Basic Configuration
	defaultAPIVersion = veikkausapi.VeikkausAPIVersion
	defaultBaseURL    = veikkausapi.VeikkausApiBaseUrl + "/" + defaultAPIVersion
	defaultUserAgent  = "go-veikkaus" + Version

	// Header configurations
	// headerRobotIdentifierKey   = veikkausapi.RobotIdentifierHeaderKey
	// headerRobotIdentifierValue = veikkausapi.RobotIdentifierHeaderValue
	// headerAccept               = veikkausapi.Accept
	// headerContentType          = veikkausapi.ContentType

	// // API Endpoints
	// loginEndpoint = "sessions"
)

// var errNonNilContext = errors.New("context must be non-nil")

type RequestOption func(req *http.Request)

// type requestContext uint8

type Client struct {
	clientMu sync.Mutex
	client   *http.Client

	// Base URL for API requests.
	// Base URL should end with the trailing slash
	BaseURL *url.URL

	// User Agent to use when communicating with Veikkaus JSON API
	UserAgent string

	common service

	// Services used for interacting with different endpoints on Veikkaus API
	Login *LoginService
}

type service struct {
	client *Client
}

func (c *Client) Client() *http.Client {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	clientCopy := *c.client
	return &clientCopy
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient
	c := &Client{client: &httpClient2}
	c.initialize()
	return c
}

func (c *Client) initialize() {
	if c.client == nil {
		c.client = &http.Client{}
	}
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}
	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent
	}
	c.common.client = c
	c.Login = (*LoginService)(&c.common)
}

// func (c *Client) copy() *Client {
// 	c.clientMu.Lock()
// 	clone := Client{
// 		client:    &http.Client{},
// 		UserAgent: c.UserAgent,
// 		BaseURL:   c.BaseURL,
// 	}
// 	c.clientMu.Unlock()
// 	if c.client != nil {
// 		clone.client.Transport = c.client.Transport
// 		clone.client.CheckRedirect = c.client.CheckRedirect
// 		clone.client.Jar = c.client.Jar
// 		clone.client.Timeout = c.client.Timeout
// 	}
// 	return &clone
// }

// func (c *Client) NewRequest(method, endpointPath string, body interface{}, opts ...RequestOption) (*http.Request, error) {
// 	if !strings.HasSuffix(c.BaseURL.Path, "/") {
// 		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
// 	}

// 	u, err := c.BaseURL.Parse(endpointPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var buf io.ReadWriter
// 	if body != nil {
// 		buf = &bytes.Buffer{}
// 		enc := json.NewEncoder(buf)
// 		enc.SetEscapeHTML(false)
// 		err := enc.Encode(body)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	req, err := http.NewRequest(method, u.String(), buf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if body != nil {
// 		req.Header.Set("Content-Type", headerContentType)
// 	}

// 	req.Header.Set("Accept", headerAccept)
// 	if c.UserAgent != "" {
// 		req.Header.Set("User-Agent", c.UserAgent)
// 	}
// 	req.Header.Set(headerRobotIdentifierKey, headerRobotIdentifierValue)

// 	for _, opt := range opts {
// 		opt(req)
// 	}

// 	return req, nil
// }
