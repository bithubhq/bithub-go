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

	TestnetBitcoin     CoinType = 1000000
	TestnetLitecoin    CoinType = 1000001
	TestnetZcash       CoinType = 1000133
	TestnetBitcoinCash CoinType = 1000145
	TestnetEthereum    CoinType = 1000060
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
	case TestnetBitcoin:
		return "Testnet Bitcoin"
	case TestnetBitcoinCash:
		return "Testnet Bitcoin Cash"
	case TestnetZcash:
		return "Testnet Zcash"
	case TestnetLitecoin:
		return "Testnet Litecoin"
	case TestnetEthereum:
		return "Testnet Ethereum"
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
	case TestnetBitcoin:
		return "TBTC"
	case TestnetBitcoinCash:
		return "TBCH"
	case TestnetZcash:
		return "TZEC"
	case TestnetLitecoin:
		return "TLTC"
	case TestnetEthereum:
		return "TETH"
	default:
		return ""
	}
}

func (c CoinType) ID() int {
	return int(c)
}

func ParseCurrencyCode(code string) (CoinType, error) {
	switch strings.ToUpper(code) {
	case "BTC":
		return Bitcoin, nil
	case "TBTC":
		return TestnetBitcoin, nil
	}
	return math.MaxUint32, ErrUnknownCurrencyCode
}
