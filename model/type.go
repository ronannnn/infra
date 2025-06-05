package model

import (
	"github.com/shopspring/decimal"
)

type DecimalSafe struct {
	*decimal.Decimal
}

func (d *DecimalSafe) UnmarshalJSON(decimalBytes []byte) error {
	if string(decimalBytes) == `""` {
		return nil // Handle empty string case
	}

	d.Decimal = &decimal.Decimal{}
	return d.Decimal.UnmarshalJSON(decimalBytes)
}
