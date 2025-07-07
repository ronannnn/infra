package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

// Scan 实现 Gorm Scanner 接口
func (d *DecimalSafe) Scan(value any) error {
	if value == nil {
		*d = DecimalSafe{}
		return nil
	}

	var dec decimal.Decimal
	var err error

	switch v := value.(type) {
	case []byte:
		dec, err = decimal.NewFromString(string(v))
	case string:
		dec, err = decimal.NewFromString(v)
	default:
		return fmt.Errorf("unsupported type for gorm decimal: %T", value)
	}

	if err != nil {
		return err
	}
	*d = DecimalSafe{
		Decimal: &dec,
	}
	return nil
}

// Value 实现 Gorm Valuer 接口
func (d DecimalSafe) Value() (driver.Value, error) {
	return d.Decimal.String(), nil
}

// GormDataType 实现 GormDataTypeInterface 接口
func (DecimalSafe) GormDataType() string {
	return "decimal" // 基础类型声明
}

// GormDBDataType 实现 GormDBDataTypeInterface 接口
func (DecimalSafe) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "DECIMAL(19, 4)" // PostgreSQL 使用字符串类型
	case "mysql":
		return "VARCHAR(191)" // MySQL 使用较短的字符串
	default:
		return "VARCHAR(255)" // 其他数据库默认
	}
}
