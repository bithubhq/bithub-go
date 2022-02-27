package bithub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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

func (s *walletService) ListAddresses(params ListAddressesParams) ([]*Address, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	path := s.listAddressesEndpointPath(params.Coin.CurrencyCode(), params.Label)
	b, err := s.client.sendHTTPRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
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
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	var addresses []*Address
	for _, r := range response.Data {
		addresses = append(addresses, &Address{
			Label:   r.Label,
			Address: r.Address,
			Balance: AddressBalance{
				Amount:    r.Balance.Amount,
				AmountUSD: r.Balance.AmountUSD,
			},
		})
	}

	return addresses, nil
}

func (s *walletService) listAddressesEndpointPath(currencyCode string, label *string) string {
	path := fmt.Sprintf("%s/%s/%s/%s",
		walletServiceEndpoint,
		s.netType.String(),
		currencyCode,
		walletAddressEndpoint,
	)

	if label != nil {
		trimmedLabel := strings.TrimSpace(*label)
		if trimmedLabel != "" {
			path = fmt.Sprintf("%s?label=%s", path, url.QueryEscape(trimmedLabel))
		}
	}

	return path
}
