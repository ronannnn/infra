package validator_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	a := decimal.RequireFromString("1.23450")
	require.Equal(t, "1.2345", a.String())
	require.Equal(t, "0.2345", a.Sub(a.Floor()).String())
}
