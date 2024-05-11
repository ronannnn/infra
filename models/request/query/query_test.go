package query_test

import (
	"testing"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Height    string `json:"height"`
	CarNumber string `json:"carNumber"`
	Status    uint   `json:"status"`
	Birth     string `json:"birth"`
}

func (u TestUser) TableName() string {
	return "test_users"
}

func (u TestUser) FieldColMapper() map[string]string {
	return models.CamelToSnakeFromStruct(u)
}

func TestParseSearch(t *testing.T) {
	example := query.Query{
		Pagination: query.Pagination{
			PageNum:  1,
			PageSize: 10,
		},
		SelectQuery: []query.SelectQueryItem{
			{Field: "username"},
			{Field: "carNumber"},
		},
		WhereQuery: []query.WhereQueryItem{
			{Field: "username", Opr: query.TypeEq, Value: "ronan"},
			{Field: "nickname", Opr: query.TypeNe, Value: "awe"},
			{Field: "age", Opr: query.TypeGt, Value: 18},
			{Field: "age", Opr: query.TypeLt, Value: 25},
			{Field: "height", Opr: query.TypeGte, Value: 170.5},
			{Field: "height", Opr: query.TypeLte, Value: 185},
			{Field: "carNumber", Opr: query.TypeLike, Value: "浙"},
			{Field: "carNumber", Opr: query.TypeStartLike, Value: "浙"},
			{Field: "carNumber", Opr: query.TypeEndLike, Value: "浙"},
			{Field: "status", Opr: query.TypeIn, Value: []uint{1, 2, 3}},
			{Field: "status", Opr: query.TypeNotIn, Value: []uint{4, 5}},
			{Field: "birth", Opr: query.TypeRange, Value: query.Range{Start: "2000-01-01", End: "2000-12-31"}},
		},
		OrderQuery: []query.OrderQueryItem{
			{Field: "birth", Order: "desc"},
			{Field: "nickname", Order: "asc"},
		},
	}
	var condition query.DbConditionImpl
	err := query.ResolveQuery(example, TestUser{}, &condition)
	require.NoError(t, err)
	// select
	require.EqualValues(t, 2, len(condition.Select))
	require.EqualValues(t, "`test_users`.`username`", condition.Select[0])
	require.EqualValues(t, "`test_users`.`car_number`", condition.Select[1])
	// where
	require.EqualValues(t, "ronan", condition.Where["`test_users`.`username` = ?"][0][0])
	require.EqualValues(t, "awe", condition.Not["`test_users`.`nickname` != ?"][0][0])
	require.EqualValues(t, 18, condition.Where["`test_users`.`age` > ?"][0][0])
	require.EqualValues(t, 25, condition.Where["`test_users`.`age` < ?"][0][0])
	require.EqualValues(t, 170.5, condition.Where["`test_users`.`height` >= ?"][0][0])
	require.EqualValues(t, 185, condition.Where["`test_users`.`height` <= ?"][0][0])
	require.EqualValues(t, "%浙%", condition.Where["`test_users`.`car_number` like ?"][0][0])
	require.EqualValues(t, "浙%", condition.Where["`test_users`.`car_number` like ?"][1][0])
	require.EqualValues(t, "%浙", condition.Where["`test_users`.`car_number` like ?"][2][0])
	require.EqualValues(t, 3, len(condition.Where["`test_users`.`status` in (?)"][0][0].([]uint)))
	require.EqualValues(t, 2, len(condition.Where["`test_users`.`status` not in (?)"][0][0].([]uint)))
	require.EqualValues(t, "2000-01-01", condition.Where["`test_users`.`birth` >= ?"][0][0])
	require.EqualValues(t, "2000-12-31", condition.Where["`test_users`.`birth` <= ?"][0][0])
	// order
	require.EqualValues(t, 2, len(condition.Order))
	require.EqualValues(t, "`test_users`.`birth` desc", condition.Order[0])
	require.EqualValues(t, "`test_users`.`nickname` asc", condition.Order[1])
}
