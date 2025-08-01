package query_test

import (
	"testing"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/model/request/query"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	model.Base
	Username  *string `json:"username"`
	Nickname  *string `json:"nickname"`
	Age       *int    `json:"age"`
	Height    *string `json:"height"`
	CarNumber *string `json:"carNumber"`
	Status    *uint   `json:"status"`
	Birth     *string `json:"birth"`
	Human     *bool   `json:"human"`
}

func (u TestUser) TableName() string {
	return "test_users"
}

func TestFieldColMapper(t *testing.T) {
	user := TestUser{}
	mapper := query.CamelToSnakeFromStruct(user)
	require.EqualValues(t, "id", mapper["id"])
	require.EqualValues(t, "created_at", mapper["createdAt"])
	require.EqualValues(t, "updated_at", mapper["updatedAt"])
	require.EqualValues(t, "created_by", mapper["createdBy"])
	require.EqualValues(t, "updated_by", mapper["updatedBy"])
	require.EqualValues(t, "version", mapper["version"])

	require.EqualValues(t, "username", mapper["username"])
	require.EqualValues(t, "nickname", mapper["nickname"])
	require.EqualValues(t, "age", mapper["age"])
	require.EqualValues(t, "height", mapper["height"])
	require.EqualValues(t, "car_number", mapper["carNumber"])
	require.EqualValues(t, "status", mapper["status"])
	require.EqualValues(t, "birth", mapper["birth"])
	require.EqualValues(t, "human", mapper["human"])
}

func TestParseQuery(t *testing.T) {
	example := query.Query{
		Pagination: query.Pagination{
			PageNum:  1,
			PageSize: 10,
		},
		SelectQuery: []query.SelectQueryItem{
			{Field: "username"},
			{Field: "carNumber"},
			{Field: "status", Distinct: true},
		},
		WhereQuery: query.WhereQuery{
			Items: []query.WhereQueryItem{
				{AndOr: "or", Field: "username", Opr: query.TypeEq, Value: "ronan"},
				{AndOr: "and", Field: "username", Opr: query.TypeEq, Value: nil},
				{AndOr: "or", Field: "age", Opr: query.TypeRangeGtLte, Value: query.Range{Start: 18, End: 25}},
				{AndOr: "or", Field: "human", Opr: query.TypeIsNot, Value: true},
			},
			Groups: []query.WhereQueryItemGroup{
				{
					Items: []query.WhereQueryItem{
						{AndOr: "or", Field: "nickname", Opr: query.TypeNe, Value: "awe"},
					},
					Groups: []query.WhereQueryItemGroup{
						{
							Items: []query.WhereQueryItem{
								{Field: "carNumber", Opr: query.TypeLike, Value: "浙"},
								{Field: "carNumber", Opr: query.TypeStartLike, Value: "浙"},
								{Field: "carNumber", Opr: query.TypeEndLike, Value: "浙"},
							},
						},
						{
							AndOr: "or",
							Items: []query.WhereQueryItem{
								{Field: "status", Opr: query.TypeIn, Value: []uint{1, 2, 3}},
								{Field: "status", Opr: query.TypeNotIn, Value: []uint{4, 5}},
								{Field: "human", Opr: query.TypeIs, Value: true},
								{Field: "human", Opr: query.TypeIsNot, Value: false},
							},
						},
					},
				},
				{
					AndOr: "or",
					Items: []query.WhereQueryItem{
						{Field: "age", Opr: query.TypeGt, Value: 18},
						{Field: "age", Opr: query.TypeLt, Value: 25},
						{Field: "height", Opr: query.TypeGte, Value: 170.5},
						{Field: "height", Opr: query.TypeLte, Value: 185},
					},
				},
			},
		},
		OrderQuery: []query.OrderQueryItem{
			{Field: "createdAt", Order: "desc"},
			{Field: "nickname", Order: "asc"},
		},
	}
	var condition query.DbConditionImpl
	err := query.ResolveQuery(example, TestUser{}, &condition, query.GetAllConditionFilter())
	require.NoError(t, err)
	// distinct
	require.EqualValues(t, 1, len(condition.Distinct))
	require.EqualValues(t, "status", condition.Distinct[0])
	// select
	require.EqualValues(t, 2, len(condition.Select))
	require.EqualValues(t, "\"test_users\".\"username\"", condition.Select[0])
	require.EqualValues(t, "\"test_users\".\"car_number\"", condition.Select[1])
	// where
	require.EqualValues(t, 3, len(condition.Where))
	// where[0]
	require.EqualValues(t, "and", condition.Where[0].AndOr)
	require.EqualValues(t, 1, len(condition.Where[0].Items))
	require.EqualValues(t, "\"test_users\".\"nickname\" != ?", condition.Where[0].Items[0].Key)
	require.EqualValues(t, "awe", condition.Where[0].Items[0].Value)
	// where[0].Groups
	require.EqualValues(t, 2, len(condition.Where[0].Groups))
	// where[0].Groups[0]
	require.EqualValues(t, "and", condition.Where[0].Groups[0].AndOr)
	require.EqualValues(t, 3, len(condition.Where[0].Groups[0].Items))
	require.EqualValues(t, "\"test_users\".\"car_number\" like ?", condition.Where[0].Groups[0].Items[0].Key)
	require.EqualValues(t, "%浙%", condition.Where[0].Groups[0].Items[0].Value)
	require.EqualValues(t, "\"test_users\".\"car_number\" like ?", condition.Where[0].Groups[0].Items[1].Key)
	require.EqualValues(t, "浙%", condition.Where[0].Groups[0].Items[1].Value)
	require.EqualValues(t, "\"test_users\".\"car_number\" like ?", condition.Where[0].Groups[0].Items[2].Key)
	require.EqualValues(t, "%浙", condition.Where[0].Groups[0].Items[2].Value)
	// where[0].Groups[1]
	require.EqualValues(t, "or", condition.Where[0].Groups[1].AndOr)
	require.EqualValues(t, 4, len(condition.Where[0].Groups[1].Items))
	require.EqualValues(t, "\"test_users\".\"status\" in (?)", condition.Where[0].Groups[1].Items[0].Key)
	require.EqualValues(t, []uint{1, 2, 3}, condition.Where[0].Groups[1].Items[0].Value)
	require.EqualValues(t, "\"test_users\".\"status\" not in (?)", condition.Where[0].Groups[1].Items[1].Key)
	require.EqualValues(t, []uint{4, 5}, condition.Where[0].Groups[1].Items[1].Value)
	require.EqualValues(t, "\"test_users\".\"human\" is ?", condition.Where[0].Groups[1].Items[2].Key)
	require.EqualValues(t, true, condition.Where[0].Groups[1].Items[2].Value)
	require.EqualValues(t, "\"test_users\".\"human\" is not ?", condition.Where[0].Groups[1].Items[3].Key)
	require.EqualValues(t, false, condition.Where[0].Groups[1].Items[3].Value)
	// where[1]
	require.EqualValues(t, "or", condition.Where[1].AndOr)
	require.EqualValues(t, 4, len(condition.Where[1].Items))
	require.EqualValues(t, "\"test_users\".\"age\" > ?", condition.Where[1].Items[0].Key)
	require.EqualValues(t, 18, condition.Where[1].Items[0].Value)
	require.EqualValues(t, "\"test_users\".\"age\" < ?", condition.Where[1].Items[1].Key)
	require.EqualValues(t, 25, condition.Where[1].Items[1].Value)
	require.EqualValues(t, "\"test_users\".\"height\" >= ?", condition.Where[1].Items[2].Key)
	require.EqualValues(t, 170.5, condition.Where[1].Items[2].Value)
	require.EqualValues(t, "\"test_users\".\"height\" <= ?", condition.Where[1].Items[3].Key)
	require.EqualValues(t, 185, condition.Where[1].Items[3].Value)
	// where[2]
	require.EqualValues(t, "and", condition.Where[2].AndOr)
	require.EqualValues(t, 2, len(condition.Where[2].Items))
	require.EqualValues(t, "\"test_users\".\"username\" = ?", condition.Where[2].Items[0].Key)
	require.EqualValues(t, "ronan", condition.Where[2].Items[0].Value)
	require.EqualValues(t, "\"test_users\".\"human\" is not ?", condition.Where[2].Items[1].Key)
	require.EqualValues(t, true, condition.Where[2].Items[1].Value)
	// where[2].Groups
	require.EqualValues(t, 1, len(condition.Where[2].Groups))
	// where[2].Groups[0]
	require.EqualValues(t, "or", condition.Where[2].Groups[0].AndOr)
	require.EqualValues(t, 2, len(condition.Where[2].Groups[0].Items))
	require.EqualValues(t, "\"test_users\".\"age\" > ?", condition.Where[2].Groups[0].Items[0].Key)
	require.EqualValues(t, 18, condition.Where[2].Groups[0].Items[0].Value)
	require.EqualValues(t, "\"test_users\".\"age\" <= ?", condition.Where[2].Groups[0].Items[1].Key)
	require.EqualValues(t, 25, condition.Where[2].Groups[0].Items[1].Value)
	// order
	require.EqualValues(t, 2, len(condition.Order))
	require.EqualValues(t, "\"test_users\".\"created_at\" desc", condition.Order[0])
	require.EqualValues(t, "\"test_users\".\"nickname\" asc", condition.Order[1])
}
