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
	Label   string
	Address string
	Balance Balance
}

type Balance struct {
	Amount    float64
	AmountUSD float64
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

	type createAddressResponse struct {
		Label   string `json:"label"`
		Address string `json:"address"`
	}

	var response createAddressResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	address := Address{
		Label:   response.Label,
		Address: response.Address,
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

	type listAddressesResponse struct {
		Data []struct {
			Label   string `json:"label"`
			Address string `json:"address"`
			Balance struct {
				Amount    float64 `json:"amount"`
				AmountUSD float64 `json:"amount_usd"`
			} `json:"balance"`
		} `json:"data"`
	}

	var response listAddressesResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	var addresses []*Address
	for _, r := range response.Data {
		addresses = append(addresses, &Address{
			Label:   r.Label,
			Address: r.Address,
			Balance: Balance{
				Amount:    r.Balance.Amount,
				AmountUSD: r.Balance.AmountUSD,
			},
		})
	}

	return addresses, nil
}

type GetBalanceParams struct {
	Coin *CoinType
}

func (p GetBalanceParams) validate() error {
	if p.Coin == nil || p.Coin.String() == "" {
		return ErrInvalidCoinType
	}

	return nil
}

func (c *Client) GetBalance(params GetBalanceParams) (*Balance, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	endpoint := c.getBalanceEndpoint(&params)
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

	type getBalanceResponse struct {
		Data struct {
			Amount    float64 `json:"amount"`
			AmountUSD float64 `json:"amount_usd"`
		} `json:"data"`
	}

	var response getBalanceResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	balance := Balance{
		Amount:    response.Data.Amount,
		AmountUSD: response.Data.AmountUSD,
	}

	return &balance, nil
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

func (c *Client) getBalanceEndpoint(params *GetBalanceParams) string {
	wallet := params.Coin.CurrencyCode()

	return fmt.Sprintf("%s/wallets/%s/%s", c.BaseURL, wallet, "balance")
}

func (c *Client) responseToError(res *http.Response) error {
	var errorResponse APIError
	if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
		return err
	}

	return &errorResponse
}
