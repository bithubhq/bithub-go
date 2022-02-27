package bithub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Transaction struct {
	TxID string
}

type SendParams struct {
	Coin    *CoinType
	Address string
	Amount  float64
}

func (p SendParams) validate() error {
	if p.Coin == nil || p.Coin.String() == "" {
		return ErrInvalidCoinType
	}

	if strings.TrimSpace(p.Address) == "" {
		return ErrInvalidAddress
	}

	if p.Amount <= 0 {
		return ErrInvalidAmount
	}

	return nil
}

func (s *walletService) Send(params SendParams) (*Transaction, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	type sendRequest struct {
		CurrencyCode string  `json:"currency_code"`
		Address      string  `json:"address"`
		Amount       float64 `json:"amount"`
	}

	payload := sendRequest{
		CurrencyCode: params.Coin.CurrencyCode(),
		Address:      params.Address,
		Amount:       params.Amount,
	}

	path := s.sendEndpointPath(params.Coin.CurrencyCode())
	b, err := s.client.sendHTTPRequest(http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}

	type sendToAddressResponse struct {
		TxID string `json:"txid"`
	}

	var response sendToAddressResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	transaction := Transaction{TxID: response.TxID}

	return &transaction, nil
}

func (s *walletService) sendEndpointPath(currencyCode string) string {
	return fmt.Sprintf("%s/%s/%s/%s",
		walletServiceEndpoint,
		s.netType.String(),
		currencyCode,
		walletSendEndpoint,
	)
}
