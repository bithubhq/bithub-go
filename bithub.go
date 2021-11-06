package bithub

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"
)

const BaseURL string = "https://api.bithub.com"

// defaultHTTPTimeout is the default timeout on the http.Client used by the library.
const defaultHTTPTimeout = 80 * time.Second

var defaultHTTPClient = &http.Client{
	Timeout: defaultHTTPTimeout,
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type Config struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	PINCode    string
}

func New(apiKey string, pinCode string) *Client {
	return newClient(apiKey, pinCode, BaseURL, defaultHTTPClient)
}

func (c *Config) New() *Client {
	if c.HTTPClient == nil {
		c.HTTPClient = defaultHTTPClient
	}

	return newClient(c.APIKey, c.PINCode, c.BaseURL, c.HTTPClient)
}

func newClient(apiKey string, pinCode string, baseURL string, httpClient *http.Client) *Client {
	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		APIKey:     apiKey,
		PINCode:    pinCode,
	}
}
