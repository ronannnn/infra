package validator_test

import (
	"testing"

	"github.com/shopspring/decimal"
)

func Test(t *testing.T) {
	a := decimal.RequireFromString("1.23450")
	println(a.String())
	println(a.Sub(a.Floor()).String())
}
