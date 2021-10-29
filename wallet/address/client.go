package address

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bithubhq/bithub-go/coin"
)

type Client struct {
	Config *ClientConfig
}

type ClientConfig struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	PINCode    string
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		Config: config,
	}
}

type CreateRequest struct {
	Label        string `json:"label"`
	CurrencyCode string `json:"currency_code"`
}

type Address struct {
	Label   string `json:"label"`
	Address string `json:"address"`
}

const (
	LabelMaxLength = 128
)

var (
	ErrLabelTooLong    = errors.New("label cannot exceed 128 symbols")
	ErrInvalidCoinType = errors.New("invalid coin type")
)

type AddressParams struct {
	Coin  coin.Type
	Label string
}

func (p AddressParams) Validate() error {
	if len(p.Label) > LabelMaxLength {
		return ErrLabelTooLong
	}

	if p.Coin.String() == "" {
		return ErrInvalidCoinType
	}

	return nil
}

type APIError struct {
	Msg string `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	ret, _ := json.Marshal(*e)
	return string(ret)
}

func (c *Client) Create(params AddressParams) (*Address, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	payload := CreateRequest{
		Label:        params.Label,
		CurrencyCode: params.Coin.CurrencyCode(),
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.createAddressEndpoint(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Api-Key", c.Config.APIKey)
	req.Header.Add("X-Pin-Code", c.Config.PINCode)

	res, err := c.Config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, c.responseToError(res)
	}

	var address Address
	if err := json.NewDecoder(res.Body).Decode(&address); err != nil {
		return nil, err
	}

	return &address, nil
}

func (c *Client) createAddressEndpoint() string {
	return fmt.Sprintf("%s/%s/%s", c.Config.BaseURL, "v1", "addresses")
}

func (c *Client) responseToError(res *http.Response) error {
	var errorResponse APIError
	if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
		return err
	}

	return &errorResponse
}

func (a *Address) equals(b *Address) bool {
	return a.Label == b.Label && a.Address == b.Address
}
