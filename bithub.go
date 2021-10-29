package bithub

import (
	"github.com/bithubhq/bithub-go/wallet"
)

type API struct {
	Wallet *wallet.API
}

func New(apiKey string, pinCode string) *API {
	return &API{
		Wallet: wallet.NewAPI(apiKey, pinCode),
	}
}
