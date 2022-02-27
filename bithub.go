package bithub

import "net/http"

const (
	BaseURL string = "https://api.bithub.com"

	walletServiceEndpoint = "wallets"
	walletAddressEndpoint = "addresses"
	walletBalanceEndpoint = "balance"
	walletSendEndpoint    = "send"
)

type Config struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	PINCode    string
	IsMainNet  bool
}
