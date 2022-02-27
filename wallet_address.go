package bithub

const (
	LabelMaxLength = 128
)

type Address struct {
	Label   string
	Address string
	Balance AddressBalance
}

type AddressBalance struct {
	Amount    float64
	AmountUSD float64
}
