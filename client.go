package bithub

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// defaultHTTPTimeout is the default timeout on the http.Client used by the library.
const defaultHTTPTimeout = 80 * time.Second

var defaultHTTPClient = &http.Client{
	Timeout: defaultHTTPTimeout,
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type Client struct {
	baseURL     string
	httpClient  *http.Client
	credentials Credentials

	// Services
	Wallets *walletService
}

type Credentials struct {
	APIKey  string
	PINCode string
}

func NewMainNetClient(credentials Credentials) *Client {
	config := Config{
		BaseURL:    BaseURL,
		HTTPClient: defaultHTTPClient,
		APIKey:     credentials.APIKey,
		PINCode:    credentials.PINCode,
		IsMainNet:  true,
	}

	return newClientFromConfig(config)
}

func NewTestNetClient(credentials Credentials) *Client {
	config := Config{
		BaseURL:    BaseURL,
		HTTPClient: defaultHTTPClient,
		APIKey:     credentials.APIKey,
		PINCode:    credentials.PINCode,
		IsMainNet:  false,
	}

	return newClientFromConfig(config)
}

func (c *Config) NewMainNetClient() *Client {
	if c.HTTPClient == nil {
		c.HTTPClient = defaultHTTPClient
	}

	c.IsMainNet = true
	return newClientFromConfig(*c)
}

func (c *Config) NewTestNetClient() *Client {
	if c.HTTPClient == nil {
		c.HTTPClient = defaultHTTPClient
	}

	c.IsMainNet = false
	return newClientFromConfig(*c)
}

func newClientFromConfig(c Config) *Client {
	baseURL := c.BaseURL
	if strings.HasSuffix(c.BaseURL, "/") {
		baseURL = strings.TrimSuffix(c.BaseURL, "/")
	}

	client := &Client{
		baseURL:    baseURL,
		httpClient: defaultHTTPClient,
		credentials: Credentials{
			APIKey:  c.APIKey,
			PINCode: c.PINCode,
		},
	}

	if c.IsMainNet {
		client.Wallets = newMainNetWalletService(client)
	} else {
		client.Wallets = newTestNetWalletService(client)
	}

	return client
}

func (c *Client) responseToError(res *http.Response) error {
	var errorResponse APIError
	if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
		return err
	}

	return &errorResponse
}

func (c *Client) sendHTTPRequest(httpMethod string, path string, body interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/%s", c.baseURL, path)

	var buffer io.ReadWriter
	if body != nil {
		buffer = new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(httpMethod, endpoint, buffer)
	if err != nil {
		return nil, err
	}

	if httpMethod == http.MethodPost {
		req.Header.Add("Content-type", "application/json")
	}
	req.Header.Add("Authorization", "Bearer "+c.credentials.APIKey)
	req.Header.Add("X-Pin-Code", c.credentials.PINCode)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, c.responseToError(res)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
