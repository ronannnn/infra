package models_test

import (
	"sort"
	"testing"

	"github.com/ronannnn/infra/models"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	Name   *string
	Age    *int
	Houses []*TestHouse
}

type TestHouse struct {
	Price   *int
	Address *string
	Rooms   []*TestRoom
}

type TestRoom struct {
	Area  *int
	Floor *string
}

func TestGetNonZeroFields(t *testing.T) {
	address := "BB"
	area := 100
	floor := ""
	age := 0
	house1Rooms := []*TestRoom{
		{Floor: &floor},
		{Area: &area},
	}
	house2Rooms := []*TestRoom{
		{Area: &area, Floor: &floor},
	}
	houses := []*TestHouse{
		{Price: &area, Address: &address, Rooms: house1Rooms},
		{Price: &area, Address: &address, Rooms: house2Rooms},
	}
	name := "Alice"
	user1 := TestUser{
		Name:   &name,
		Houses: houses,
	}
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
	actualA := models.GetNonZeroFields(&user1)
	sort.Strings(actualA)
	require.Equal(t, expectedA, actualA)

	user2 := TestUser{
		Age: &age,
	}
	expectedB := []string{"Age"}
	sort.Strings(expectedB)
	actualB := models.GetNonZeroFields(&user2)
	sort.Strings(actualB)
	require.Equal(t, []string{"Age"}, actualB)
}
