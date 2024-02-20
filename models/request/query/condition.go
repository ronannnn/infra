package query

type DbCondition interface {
	SetWhere(k string, v []any)
	SetOr(k string, v []any)
	SetNot(k string, v []any)
	SetOrder(k string)
	SetSelect(k string)
}

type DbConditionImpl struct {
	Where  map[string][]any
	Or     map[string][]any
	Not    map[string][]any
	Order  []string
	Select []string
}

func (e *DbConditionImpl) SetWhere(k string, v []any) {
	if e.Where == nil {
		e.Where = make(map[string][]any)
	}
	e.Where[k] = v
}

func (e *DbConditionImpl) SetOr(k string, v []any) {
	if e.Or == nil {
		e.Or = make(map[string][]any)
	}
	e.Or[k] = v
}

func (e *DbConditionImpl) SetNot(k string, v []any) {
	if e.Not == nil {
		e.Not = make(map[string][]any)
	}
	e.Not[k] = v
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

func (e *DbConditionImpl) UnquotedSelect() []string {
	quotedSelect := make([]string, len(e.Select))
	for i, s := range e.Select {
		quotedSelect[i] = s[1 : len(s)-1]
	}
	return quotedSelect
}
