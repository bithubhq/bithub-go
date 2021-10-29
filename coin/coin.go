package coin

import (
	"errors"
	"math"
	"strings"
)

var ErrUnknownCurrencyCode = errors.New("unknown currency code")

type Type uint32

const (
	Bitcoin     Type = 0
	Litecoin    Type = 1
	Zcash       Type = 133
	BitcoinCash Type = 145
	Ethereum    Type = 60

	TestnetBitcoin     Type = 1000000
	TestnetLitecoin    Type = 1000001
	TestnetZcash       Type = 1000133
	TestnetBitcoinCash Type = 1000145
	TestnetEthereum    Type = 1000060
)

func (c Type) String() string {
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

func (c Type) CurrencyCode() string {
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

func (c Type) ID() int {
	return int(c)
}

func ParseCurrencyCode(code string) (Type, error) {
	switch strings.ToUpper(code) {
	case "BTC":
		return Bitcoin, nil
	case "TBTC":
		return TestnetBitcoin, nil
	}
	return math.MaxUint32, ErrUnknownCurrencyCode
}
