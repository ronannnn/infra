package models

import "github.com/shopspring/decimal"

type DecimalSafe struct {
	decimal.Decimal
}

func (d *DecimalSafe) UnmarshalJSON(decimalBytes []byte) error {
	if string(decimalBytes) == `""` {
		return nil
	}

	return d.Decimal.UnmarshalJSON(decimalBytes)
}