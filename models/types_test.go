package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ronannnn/infra/models"
	"github.com/stretchr/testify/require"
)

type DecimalSafeModel struct {
	Age *models.DecimalSafe `json:"age"`
}

func TestNormalDecimalSafe(t *testing.T) {
	req := `{ "age": "1.23" }`
	var model DecimalSafeModel
	err := json.Unmarshal([]byte(req), &model)
	require.NoError(t, err)
	require.Equal(t, "1.23", model.Age.String())
}

func TestEmptyDecimalSafe(t *testing.T) {
	req := `{ "age": "" }`
	var model DecimalSafeModel
	err := json.Unmarshal([]byte(req), &model)
	require.NoError(t, err)
	fmt.Printf("model.Age: %v\n", model.Age)
	require.Nil(t, model.Age)
}

func TestUndefinedDecimalSafe(t *testing.T) {
	req := `{}`
	var model DecimalSafeModel
	err := json.Unmarshal([]byte(req), &model)
	require.NoError(t, err)
	fmt.Printf("model.Age: %v\n", model.Age)
	require.Nil(t, model.Age)
}
