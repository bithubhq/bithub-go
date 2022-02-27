package bithub

import (
	"encoding/json"
	"errors"
)

var (
	ErrLabelTooLong    = errors.New("label cannot exceed 128 symbols")
	ErrInvalidCoinType = errors.New("invalid coin type")
	ErrInvalidAddress  = errors.New("invalid address value")
	ErrInvalidAmount   = errors.New("invalid amount value")
)

type walletService struct {
	client  *Client
	netType netType
}

type netType int

const (
	mainNet netType = iota
	testNet
)

func (t netType) String() string {
	switch t {
	case mainNet:
		return "mainnet"
	case testNet:
		return "testnet"
	}

	return ""
}

type APIError struct {
	Msg string `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	ret, _ := json.Marshal(*e)
	return string(ret)
}

func newMainNetWalletService(client *Client) *walletService {
	return &walletService{
		client:  client,
		netType: mainNet,
	}
}

func newTestNetWalletService(client *Client) *walletService {
	return &walletService{
		client:  client,
		netType: testNet,
	}
}
