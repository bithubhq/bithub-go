package wallet

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/bithubhq/bithub-go/wallet/address"
)

// defaultHTTPTimeout is the default timeout on the http.Client used by the library.
const defaultHTTPTimeout = 80 * time.Second

const (
	// APIURL is the BaseURL of the Wallet API service backend.
	APIURL string = "https://api.bithub.com/wallet"
)

var httpClient = &http.Client{
	Timeout: defaultHTTPTimeout,
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type API struct {
	Addresses *address.Client
}

func NewAPI(apiKey string, pinCode string) *API {
	return &API{
		Addresses: address.NewClient(&address.ClientConfig{
			BaseURL:    APIURL,
			HTTPClient: httpClient,
			APIKey:     apiKey,
			PINCode:    pinCode,
		}),
	}
}
