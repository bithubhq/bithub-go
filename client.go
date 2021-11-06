package bithub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	PINCode    string
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

type CreateAddressParams struct {
	Coin  *CoinType
	Label *string
}

func (p CreateAddressParams) validate() error {
	if p.Label != nil && len(*p.Label) > LabelMaxLength {
		return ErrLabelTooLong
	}

	if p.Coin == nil || p.Coin.String() == "" {
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

func (c *Client) CreateAddress(params CreateAddressParams) (*Address, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	type createRequest struct {
		Label        string `json:"label"`
		CurrencyCode string `json:"currency_code"`
	}

	var payload createRequest
	if params.Label != nil {
		payload.Label = *params.Label
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	endpoint := c.createAddressEndpoint(&params)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.APIKey)
	req.Header.Add("X-Pin-Code", c.PINCode)

	res, err := c.HTTPClient.Do(req)
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

type ListAddressesParams struct {
	Coin  *CoinType
	Label *string
}

func (p ListAddressesParams) validate() error {
	if p.Label != nil && len(*p.Label) > LabelMaxLength {
		return ErrLabelTooLong
	}

	if p.Coin == nil || p.Coin.String() == "" {
		return ErrInvalidCoinType
	}

	return nil
}

func (c *Client) ListAddresses(params ListAddressesParams) ([]*Address, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	endpoint := c.listAddressesEndpoint(&params)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.APIKey)
	req.Header.Add("X-Pin-Code", c.PINCode)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, c.responseToError(res)
	}

	type listResponse struct {
		Data []*Address `json:"data"`
	}

	var response listResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *Client) createAddressEndpoint(params *CreateAddressParams) string {
	wallet := params.Coin.CurrencyCode()

	return fmt.Sprintf("%s/wallets/%s/%s", c.BaseURL, wallet, "addresses")
}

func (c *Client) listAddressesEndpoint(params *ListAddressesParams) string {
	wallet := params.Coin.CurrencyCode()

	var label string
	if params.Label != nil {
		label = *params.Label
	}

	label = strings.TrimSpace(label)
	if label != "" {
		return fmt.Sprintf("%s/wallets/%s/%s?label=%s",
			c.BaseURL,
			wallet,
			"addresses",
			url.QueryEscape(label),
		)
	}

	return fmt.Sprintf("%s/wallets/%s/%s",
		c.BaseURL,
		wallet,
		"addresses",
	)
}

func (c *Client) responseToError(res *http.Response) error {
	var errorResponse APIError
	if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
		return err
	}

	return &errorResponse
}
