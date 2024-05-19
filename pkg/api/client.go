package api

import (
	"fmt"
	"io"
	"net/http"
)

// Client is an configured http client.
type Client struct {
	*http.Client
	ContentType string
	APIBaseURL  string
}

var (
	// UserAgent lets the API know where the call is being made from.
	// It's overridden from the root command so that we can set the version.
	UserAgent = "ipapi.co/#go-v1.5"

	// HTTPClient is the client used to make HTTP calls.
	HTTPClient = &http.Client{}
)

// NewClient returns an API client.
func NewClient(baseURL string) (*Client, error) {
	return &Client{
		Client:     HTTPClient,
		APIBaseURL: baseURL,
	}, nil
}

// NewRequest returns an http.Request with information for the API.
func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	if c.Client == nil {
		c.Client = HTTPClient
	}

	req, err := http.NewRequest(method, c.APIBaseURL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	if c.ContentType == "" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", c.ContentType)
	}

	return req, nil
}

// Do performs an http.Request and optionally parses the response body into the given interface.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		fmt.Print("err", err)
		return nil, err
	}

	return res, nil
}
