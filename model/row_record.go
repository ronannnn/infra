package model

import (
	"fmt"
	"net/http"
	"reflect"
)

type RowRecord struct {
	Base
	Table    string `json:"-"`
	RowId    uint   `json:"rowId"`
	Key      string `json:"key"`
	OldValue string `json:"oldValue" gorm:"type:text"`
	NewValue string `json:"newValue" gorm:"type:text"`
}

func (RowRecord) TableName() string {
	return "row_records"
}

func (model RowRecord) WithOprFromReq(r *http.Request) Crudable {
	model.OprBy = GetOprFromReq(r)
	return model
}

func (model RowRecord) WithUpdaterFromReq(r *http.Request) Crudable {
	model.OprBy = GetUpdaterFromReq(r)
	return model
}

type rowRecordHelper struct {
	TableName string
	RowId     uint
	Records   []RowRecord
}

func (trh *rowRecordHelper) record(key string, oldValue any, newValue any, oprId uint) {
	trh.Records = append(trh.Records, RowRecord{
		Base:     Base{OprBy: OprBy{CreatedBy: &oprId, UpdatedBy: &oprId}},
		Table:    trh.TableName,
		RowId:    trh.RowId,
		Key:      key,
		OldValue: fmt.Sprintf("%v", oldValue),
		NewValue: fmt.Sprintf("%v", newValue),
	})
}

func newRowRecordHelper(tableName string, rowId uint) *rowRecordHelper {
	return &rowRecordHelper{
		TableName: tableName,
		RowId:     rowId,
	}
}

func resolveRecord(trh *rowRecordHelper, oldModel any, newModel any, oprId uint) {
	types := reflect.TypeOf(newModel)
	newValues := reflect.ValueOf(newModel)
	oldValues := reflect.ValueOf(oldModel)
	for i := 0; i < types.NumField(); i++ {
		if newValues.Field(i).IsZero() {
			continue
		}
		tagStr, ok := types.Field(i).Tag.Lookup("json")
		if !ok {
			// 递归调用
			resolveRecord(trh, oldValues.Field(i).Interface(), newValues.Field(i).Interface(), oprId)
			continue
		}
		switch tagStr {
		case "id", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "version":
			continue
		}
		if oldValues.Field(i).IsZero() {
			trh.record(tagStr, "", newValues.Field(i).Elem(), oprId)
		} else {
			trh.record(tagStr, oldValues.Field(i).Elem(), newValues.Field(i).Elem(), oprId)
		}
	}
}

func RecordRowModifications(tableName string, rowId uint, oldModel any, newModel any, oprId uint) []RowRecord {
	trh := newRowRecordHelper(tableName, rowId)
	resolveRecord(trh, oldModel, newModel, oprId)
	return trh.Records
}
