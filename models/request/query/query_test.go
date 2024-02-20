package query_test

import (
	"testing"

	"github.com/ronannnn/infra/models/request/query"
	"github.com/stretchr/testify/require"
)

type Example struct {
	EqField    string `search:"type:eq;column:eq_field"`
	NeField    string `search:"type:ne;column:ne_field"`
	GtField    string `search:"type:gt;column:gt_field"`
	GteField   string `search:"type:gte;column:gte_field"`
	LtField    string `search:"type:lt;column:lt_field"`
	LteField   string `search:"type:lte;column:lte_field"`
	LikeField  string `search:"type:like;column:like_field"`
	InField    []uint `search:"type:in;column:in_field"`
	OrderField string `search:"type:order;column:order_field"`
}

func TestParseSearch(t *testing.T) {
	example := Example{
		EqField:    "exact",
		NeField:    "not exact",
		GtField:    "2023",
		GteField:   "2024",
		LtField:    "2025",
		LteField:   "2026",
		LikeField:  "like this",
		InField:    []uint{1, 2, 3},
		OrderField: "asc",
	}
	var condition query.DbConditionImpl
	query.ResolveQuery(example, &condition)
	require.EqualValues(t, "exact", condition.Where["`eq_field` = ?"][0])
	require.EqualValues(t, "not exact", condition.Not["`ne_field` = ?"][0])
	require.EqualValues(t, "2023", condition.Where["`gt_field` > ?"][0])
	require.EqualValues(t, "2024", condition.Where["`gte_field` >= ?"][0])
	require.EqualValues(t, "2025", condition.Where["`lt_field` < ?"][0])
	require.EqualValues(t, "2026", condition.Where["`lte_field` <= ?"][0])
	require.EqualValues(t, "%like this%", condition.Where["`like_field` like ?"][0])
	require.EqualValues(t, 3, len(condition.Where["`in_field` in (?)"][0].([]uint)))
	require.EqualValues(t, "`order_field` asc", condition.Order[0])
}
