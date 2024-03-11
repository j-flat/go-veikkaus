package goveikkaus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	baseURLPath = "/api-v1"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)
	url, _ := url.Parse(server.URL + baseURLPath + "/")

	// Overwrite BaseURL variable on internal/veikkausapi
	api.OverWriteBaseURL = true
	api.BaseURL = url.String()

	client = NewClient(nil)
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}
