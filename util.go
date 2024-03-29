package infra

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { // path exists
		return true, nil
	} else if os.IsNotExist(err) { // error is 'not exist'
		return false, nil
	}
	return false, err // other error
}

func createDirsIfNotExist(dirs ...string) (err error) {
	for _, dir := range dirs {
		if existing, pathExistsErr := pathExists(dir); !existing && pathExistsErr == nil {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return
			}
		}
	}
	return
}

func IsZeroValue(v interface{}) bool {
	return v == nil || v == reflect.Zero(reflect.TypeOf(v)).Interface()
}

// LeftJustifyingPrint 打印的log内容中，第二列之后的内容会左对齐
// e.g.:
// DELETE /api/v1/tasks/out/subtasks/{id}
// GET    /api/v1/tasks/out/{id}
func LeftJustifyingPrint(rows [][]string, log *zap.SugaredLogger) {
	maxFieldLength := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, field := range row {
			if len(field) > maxFieldLength[i] {
				maxFieldLength[i] = len(field)
			}
		}
	}

	for _, row := range rows {
		paddedRow := make([]string, len(row))
		for i, field := range row {
			paddedRow[i] = fmt.Sprintf("%-*s", maxFieldLength[i], field)
		}
		log.Infof("%s", strings.Join(paddedRow, " "))
	}
}
