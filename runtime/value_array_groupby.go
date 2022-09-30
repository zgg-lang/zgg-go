package runtime

type (
	arrayGroupBy struct {
		arr     []Value
		mappers []*arrayMapper
	}
	arrayGroupByItem struct {
		groupFields []Value
		record      Value
	}
	arrayGroupByList  []arrayGroupByItem
	arrayGroupPredict func(Value, Value, *Context) bool
	arrayGroupRule    []arrayGroupPredict
)

func newGroupBy(arr []Value, mappers []*arrayMapper) *arrayGroupBy {
	return &arrayGroupBy{
		arr:     arr,
		mappers: mappers,
	}
}

func (arrayGroupBy) isSameGroup(g1, g2 []Value, c *Context) bool {
	if len(g1) != len(g2) {
		return false
	}
	for i, v1 := range g1 {
		v2 := g2[i]
		if !c.ValuesEqual(v1, v2) {
			return false
		}
	}
	return true
}

func (gb *arrayGroupBy) Execute(c *Context) [][]Value {
	var (
		arr     = gb.arr
		mappers = gb.mappers
	)
	n := len(arr)
	if n == 0 {
		return [][]Value{}
	}
	raw := make([]arrayGroupByItem, n)
	for i, v := range arr {
		raw[i].record = v
		for _, m := range mappers {
			raw[i].groupFields = append(raw[i].groupFields, m.Map(v, i, c))
		}
	}
	res := make([][]arrayGroupByItem, 0)
	last := raw[0]
	res = append(res, []arrayGroupByItem{last})
	for i := 1; i < len(raw); i++ {
		this := raw[i]
		groupExists := false
		for i, group := range res {
			a := this
			b := group[0]
			if gb.isSameGroup(a.groupFields, b.groupFields, c) {
				res[i] = append(res[i], this)
				groupExists = true
				break
			}
		}
		if !groupExists {
			res = append(res, []arrayGroupByItem{this})
		}
	}
	groups := make([][]Value, len(res))
	for i, g := range res {
		groups[i] = make([]Value, len(g))
		for j, item := range g {
			groups[i][j] = item.record
		}
	}
	return groups
}
