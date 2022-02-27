package bithub

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (s *walletService) CreateAddress(params CreateAddressParams) (*Address, error) {
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

	path := s.createAddressEndpointPath(params.Coin.CurrencyCode())
	b, err := s.client.sendHTTPRequest(http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}

	type createAddressResponse struct {
		Label   string `json:"label"`
		Address string `json:"address"`
	}

	var response createAddressResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	address := Address{
		Label:   response.Label,
		Address: response.Address,
	}

	return &address, nil
}

func (s *walletService) createAddressEndpointPath(currencyCode string) string {
	return fmt.Sprintf("%s/%s/%s/%s",
		walletServiceEndpoint,
		s.netType.String(),
		currencyCode,
		walletAddressEndpoint,
	)
}
