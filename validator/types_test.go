package validator_test

import (
	"encoding/json"
	"os"
	"sort"
	"testing"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/validator"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	models.Base
	Name   *string      `json:"name" validate:"required"`
	Age    *int         `json:"age" validate:"required"`
	Houses []*TestHouse `json:"houses" validate:"required"`
}

type TestHouse struct {
	Price   *int        `json:"price" validate:"required"`
	Address *string     `json:"address" validate:"required"`
	Rooms   []*TestRoom `json:"rooms" validate:"required"`
}

type TestRoom struct {
	Area  *int    `json:"area" validate:"required"`
	Floor *string `json:"floor" validate:"required"`
}

func TestGetNonZeroFields(t *testing.T) {
	var user1 TestUser
	jsonData1, err := os.ReadFile("./testdata/user1.json")
	require.NoError(t, err)
	err = json.Unmarshal(jsonData1, &user1)
	require.NoError(t, err)
	expectedA := []string{
		"Name",
		"Houses[0].Price",
		"Houses[0].Address",
		"Houses[0].Rooms[0].Floor",
		"Houses[0].Rooms[1].Area",
		"Houses[1].Price",
		"Houses[1].Address",
		"Houses[1].Rooms[0].Area",
		"Houses[1].Rooms[0].Floor",
	}
	sort.Strings(expectedA)
	actualA := validator.GetNonZeroFields(&user1)
	sort.Strings(actualA)
	require.Equal(t, expectedA, actualA)

	var user2 TestUser
	jsonData2, err := os.ReadFile("./testdata/user2.json")
	require.NoError(t, err)
	err = json.Unmarshal(jsonData2, &user2)
	require.NoError(t, err)
	expectedB := []string{"Age"}
	sort.Strings(expectedB)
	actualB := validator.GetNonZeroFields(&user2)
	sort.Strings(actualB)
	require.Equal(t, []string{"Age"}, actualB)
}
