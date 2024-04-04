package query

import (
	"reflect"
	"strconv"
	"strings"
)

type Exceler interface {
	TableName() string
	SheetName() string
	Model() any
}

const (
	ExcelTag = "excel"
)

// ExcelColumnMeta, SheetMeta, ExcelMeta are used to store excel data
// including columns and rows
type ExcelColumnMeta struct {
	HeaderName    string
	HeaderWidth   float64
	HeaderStyleId int // generate by excelize with `file.NewStyle()`
}

type SheetMeta struct {
	ColumnsMeta []ExcelColumnMeta
	Data        [][]any
	Name        string // sheet name
}

type ExcelMeta struct {
	SheetsMeta []SheetMeta
}

// ExcelCondition is an interface to store conditions for generating excel
type ExcelCondition interface {
	SetExcelColumnMeta(meta ExcelColumnMeta)
}

type ExcelConditionImpl struct {
	ExcelColumnMeta []ExcelColumnMeta
}

func (e *ExcelConditionImpl) SetExcelColumnMeta(meta ExcelColumnMeta) {
	if e.ExcelColumnMeta == nil {
		e.ExcelColumnMeta = make([]ExcelColumnMeta, 0)
	}
	e.ExcelColumnMeta = append(e.ExcelColumnMeta, meta)
}

func (e *ExcelConditionImpl) GetColumnNames() []string {
	columnNames := make([]string, len(e.ExcelColumnMeta))
	for i, columnMeta := range e.ExcelColumnMeta {
		columnNames[i] = columnMeta.HeaderName
	}
	return columnNames
}

type excelTag struct {
	Column string
	Width  float64
}

func parseExcelTag(tagStr string) *excelTag {
	tag := &excelTag{}
	tags := strings.Split(tagStr, ";")
	var ts []string
	for _, t := range tags {
		ts = strings.Split(t, ":")
		if len(ts) != 2 {
			continue
		}
		switch ts[0] {
		case "column":
			tag.Column = ts[1]
		case "width":
			width, _ := strconv.ParseFloat(ts[1], 64)
			tag.Width = width
		}
	}
	return tag
}

func ResolveExcel(excelModel any, condition ExcelCondition) {
	excelType := reflect.TypeOf(excelModel)
	excelValue := reflect.ValueOf(excelModel)
	for i := 0; i < excelType.NumField(); i++ {
		tagStr, ok := excelType.Field(i).Tag.Lookup(ExcelTag)
		if !ok {
			// 递归调用
			ResolveExcel(excelValue.Field(i).Interface(), condition)
			continue
		}
		if tagStr == TypeSkip || excelValue.Field(i).IsZero() {
			continue
		}
		tag := parseExcelTag(tagStr)
		condition.SetExcelColumnMeta(ExcelColumnMeta{
			HeaderName:  tag.Column,
			HeaderWidth: tag.Width,
		})
	}
}
