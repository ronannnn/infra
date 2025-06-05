package query

type DbCondition interface {
	SetWhere(k string, v []any)
	SetOr(k string, v []any)
	SetNot(k string, v []any)
	SetOrder(k string)
	SetSelect(k string)
	SetDistinct(k string)
}

type DbConditionImpl struct {
	Where    map[string][][]any
	Or       map[string][][]any
	Not      map[string][][]any
	Order    []string
	Select   []string
	Distinct []string
}

func (e *DbConditionImpl) SetWhere(k string, v []any) {
	if e.Where == nil {
		e.Where = make(map[string][][]any)
	}
	old, ok := e.Where[k]
	if ok {
		e.Where[k] = append(old, v)
	} else {
		e.Where[k] = [][]any{v}
	}
}

func (e *DbConditionImpl) SetOr(k string, v []any) {
	if e.Or == nil {
		e.Or = make(map[string][][]any)
	}
	old, ok := e.Or[k]
	if ok {
		e.Or[k] = append(old, v)
	} else {
		e.Or[k] = [][]any{v}
	}
}

func (e *DbConditionImpl) SetNot(k string, v []any) {
	if e.Not == nil {
		e.Not = make(map[string][][]any)
	}
	old, ok := e.Not[k]
	if ok {
		e.Not[k] = append(old, v)
	} else {
		e.Not[k] = [][]any{v}
	}
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
