package bithub

import (
	"errors"
	"math"
	"strings"
)

var ErrUnknownCurrencyCode = errors.New("unknown currency code")

type CoinType uint32

const (
	Bitcoin     CoinType = 0
	Litecoin    CoinType = 1
	Zcash       CoinType = 133
	BitcoinCash CoinType = 145
	Ethereum    CoinType = 60
)

func (c CoinType) String() string {
	switch c {
	case Bitcoin:
		return "Bitcoin"
	case BitcoinCash:
		return "Bitcoin Cash"
	case Zcash:
		return "Zcash"
	case Litecoin:
		return "Litecoin"
	case Ethereum:
		return "Ethereum"
	default:
		return ""
	}
}

func (c CoinType) CurrencyCode() string {
	switch c {
	case Bitcoin:
		return "BTC"
	case BitcoinCash:
		return "BCH"
	case Zcash:
		return "ZEC"
	case Litecoin:
		return "LTC"
	case Ethereum:
		return "ETH"
	default:
		return ""
	}
}

func (c CoinType) ID() int {
	return int(c)
}

func ParseCurrencyCode(code string) (CoinType, error) {
	switch strings.ToUpper(code) {
	case Bitcoin.CurrencyCode():
		return Bitcoin, nil
	case BitcoinCash.CurrencyCode():
		return BitcoinCash, nil
	case Zcash.CurrencyCode():
		return Zcash, nil
	case Litecoin.CurrencyCode():
		return Litecoin, nil
	case Ethereum.CurrencyCode():
		return Ethereum, nil
	}
	return math.MaxUint32, ErrUnknownCurrencyCode
}
