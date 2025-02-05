package validator_test

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	a := decimal.RequireFromString("1.23450")
	require.Equal(t, "1.2345", a.String())
	require.Equal(t, "0.2345", a.Sub(a.Floor()).String())
}

func TestDecimalConvertion(t *testing.T) {
	a := decimal.RequireFromString("120.0")
	strA := a.String()
	require.Equal(t, "120", strA)
	require.True(t, strings.HasSuffix(strA, "0"))

	b := decimal.RequireFromString("120")
	strB := b.String()
	require.Equal(t, "120", strB)
	require.True(t, strings.HasSuffix(strB, "0"))
}
