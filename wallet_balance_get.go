package bithub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetBalanceParams struct {
	Coin *CoinType
}

func (p GetBalanceParams) validate() error {
	if p.Coin == nil || p.Coin.String() == "" {
		return ErrInvalidCoinType
	}

	return nil
}

func (s *walletService) GetBalance(params GetBalanceParams) (*WalletBalance, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	path := s.getBalanceEndpointPath(params.Coin.CurrencyCode())
	b, err := s.client.sendHTTPRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	type getBalanceResponse struct {
		Data struct {
			Amount    float64 `json:"amount"`
			AmountUSD float64 `json:"amount_usd"`
		} `json:"data"`
	}

	var response getBalanceResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	balance := WalletBalance{
		Amount:    response.Data.Amount,
		AmountUSD: response.Data.AmountUSD,
	}

	return &balance, nil
}

func (s *walletService) getBalanceEndpointPath(currencyCode string) string {
	return fmt.Sprintf("%s/%s/%s/%s",
		walletServiceEndpoint,
		s.netType.String(),
		currencyCode,
		walletBalanceEndpoint,
	)
}
