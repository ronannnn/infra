package query

type DbConditionWhereItem struct {
	AndOr string // and/or
	Key   string
	Value any
}

type DbConditionWhereGroup struct {
	AndOr  string // and/or
	Items  []DbConditionWhereItem
	Groups []DbConditionWhereGroup
}

type DbCondition interface {
	SetWhere(group DbConditionWhereGroup)
	SetOrder(k string)
	SetSelect(k string)
	SetDistinct(k string)
}

type DbConditionImpl struct {
	Where    []DbConditionWhereGroup
	Order    []string
	Select   []string
	Distinct []string
}

func (e *DbConditionImpl) SetWhere(group DbConditionWhereGroup) {
	if e.Where == nil {
		e.Where = make([]DbConditionWhereGroup, 0)
	}
	e.Where = append(e.Where, group)
}

func (e *DbConditionImpl) SetOrder(k string) {
	if e.Order == nil {
		e.Order = make([]string, 0)
	}
	e.Order = append(e.Order, k)
}

func (e *DbConditionImpl) SetSelect(k string) {
	if e.Select == nil {
		e.Select = make([]string, 0)
	}
	e.Select = append(e.Select, k)
}

func (e *DbConditionImpl) SetDistinct(k string) {
	if e.Distinct == nil {
		e.Distinct = make([]string, 0)
	}
	e.Distinct = append(e.Distinct, k)
}

func (e *DbConditionImpl) UnquotedSelect() []string {
	quotedSelect := make([]string, len(e.Select))
	for i, s := range e.Select {
		quotedSelect[i] = s[1 : len(s)-1]
	}
	return quotedSelect
}
