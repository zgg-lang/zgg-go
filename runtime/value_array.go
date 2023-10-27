package runtime

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type ValueArray struct {
	*ValueBase
	Values *[]Value
}

func (v ValueArray) IsTrue() bool {
	return v.Len() > 0
}

func NewArray(n ...int) ValueArray {
	cap := 0
	if len(n) > 0 {
		cap = n[0]
	}
	values := make([]Value, 0, cap)
	return ValueArray{ValueBase: &ValueBase{}, Values: &values}
}

func NewArrayByValues(values ...Value) ValueArray {
	return ValueArray{ValueBase: &ValueBase{}, Values: &values}
}

func (v ValueArray) GetIndex(index int, c *Context) Value {
	if index < 0 {
		index += v.Len()
	}
	if index < 0 || index >= len(*v.Values) {
		return Undefined()
	}
	return (*v.Values)[index]
}

func (v ValueArray) SetIndex(index int, value Value, c *Context) {
	if index < 0 {
		index += v.Len()
	}
	if index < 0 || index >= v.Len() {
		c.RaiseRuntimeError(fmt.Sprintf("set array item error: Out of bound length %d index %d", v.Len(), index))
		return
	}
	(*v.Values)[index] = value
}

func (v ValueArray) GetMember(name string, c *Context) Value {
	return getMemberByType(c, v, name)
}

func (v ValueArray) Type() ValueType {
	return TypeArray
}

func (v ValueArray) ToString(c *Context) string {
	var b strings.Builder
	b.WriteString("[")
	for i, v := range *v.Values {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(v.ToString(c))
	}
	b.WriteString("]")
	return b.String()
}

func (v ValueArray) ToGoValue() interface{} {
	rv := make([]interface{}, v.Len())
	for i := range rv {
		rv[i] = v.GetIndex(i, nil).ToGoValue()
	}
	return rv
}

func (ValueArray) GoType() reflect.Type {
	var v []interface{}
	return reflect.TypeOf(v)
}

func (v ValueArray) CompareTo(other Value, c *Context) CompareResult {
	otherArr, isArray := other.(ValueArray)
	if !isArray {
		return CompareResultNotEqual
	}
	l1, l2 := v.Len(), otherArr.Len()
	minLen := l1
	if l1 > l2 {
		minLen = l2
	}
	for i := 0; i < minLen; i++ {
		v1 := (*v.Values)[i]
		v2 := (*otherArr.Values)[i]
		r := v1.CompareTo(v2, c)
		if r != CompareResultEqual {
			return r
		}
	}
	if l1 > l2 {
		return CompareResultGreater
	} else if l1 < l2 {
		return CompareResultLess
	}
	return CompareResultEqual
}

func (v ValueArray) PushBack(el Value) {
	*v.Values = append(*v.Values, el)
}

func (v ValueArray) Len() int {
	return len(*v.Values)
}

func (v ValueArray) slice(i, j int) ValueArray {
	rslice := (*v.Values)[i:j]
	rv := NewArray(len(rslice))
	for _, v := range rslice {
		rv.PushBack(v)
	}
	return rv
}

func arrayFindByPredictFunc(c *Context, arr ValueArray, predict ValueCallable, start int) (int, Value) {
	items := *(arr.Values)
	n := len(items)
	for i := start; i < n; i++ {
		c.Invoke(predict, nil, func() []Value { return items[i : i+1] })
		if c.RetVal.IsTrue() {
			return i, items[i]
		}
	}
	return -1, nil
}

func arrayFindByValue(c *Context, arr ValueArray, expected Value, start int) (int, Value) {
	items := *(arr.Values)
	n := len(items)
	for i := start; i < n; i++ {
		if c.ValuesEqual(items[i], expected) {
			return i, items[i]
		}
	}
	return -1, nil
}

func arrayFind(c *Context, arr ValueArray, predict Value, start int) (int, Value) {
	if pd, ok := c.GetCallable(predict); ok {
		return arrayFindByPredictFunc(c, arr, pd, start)
	}
	return arrayFindByValue(c, arr, predict, start)
}

type arrayMapper struct {
	f      ValueCallable
	k      ValueStr
	i      ValueInt
	d      ValueInt
	which  int
	mapper func(Value, int, *Context) Value
}

func (m *arrayMapper) ArgRule(argName string, required bool) argRuleOneOf {
	rv := argRuleOneOf{
		ArgName:       argName,
		ExpectedTypes: []ValueType{TypeCallable, TypeStr, TypeInt},
		StoreTos:      []interface{}{&(m.f), &(m.k), &(m.i)},
		Selected:      &(m.which),
	}
	if !required {
		rv.DefaultStore = &(m.d)
		rv.DefaultValue = NewInt(0)
	}
	return rv
}

func (m *arrayMapper) Build() {
	switch m.which {
	case 0:
		m.mapper = func(item Value, index int, c *Context) Value {
			c.Invoke(m.f, nil, Args(item, NewInt(int64(index))))
			return c.RetVal
		}
	case 1:
		key := m.k.Value()
		m.mapper = func(item Value, index int, c *Context) Value {
			return item.GetMember(key, c)
		}
	case 2:
		i := m.i.AsInt()
		m.mapper = func(item Value, index int, c *Context) Value {
			return item.GetIndex(i, c)
		}
	default:
		m.mapper = func(item Value, index int, c *Context) Value {
			return item
		}
	}
}

func (m *arrayMapper) Map(item Value, index int, c *Context) Value {
	return m.mapper(item, index, c)
}

var builtinArrayMethods = map[string]ValueCallable{
	"map": NewNativeFunction("map", func(c *Context, thisArg Value, args []Value) Value {
		var (
			mapper arrayMapper
		)
		EnsureFuncParams(c, "array.map", args,
			mapper.ArgRule("mapper", true),
		)
		mapper.Build()
		thisArr := thisArg.(ValueArray)
		l := thisArr.Len()
		rv := NewArray(l)
		for i := 0; i < l; i++ {
			v := thisArr.GetIndex(i, c)
			rv.PushBack(mapper.Map(v, i, c))
		}
		return rv
	}),
	"filter": NewNativeFunction("filter", func(c *Context, thisArg Value, args []Value) Value {
		if len(args) != 1 {
			c.RaiseRuntimeError("array.filter: arguments length must be 1")
			return nil
		}
		f, isCallable := c.GetCallable(args[0])
		if !isCallable {
			c.RaiseRuntimeError("array.filter: argument 0 must be callable")
			return nil
		}
		thisArr := thisArg.(ValueArray)
		l := thisArr.Len()
		rv := NewArray(l)
		for i := 0; i < l; i++ {
			v := thisArr.GetIndex(i, c)
			f.Invoke(c, constUndefined, []Value{v, NewInt(int64(i))})
			if c.ReturnTrue() {
				rv.PushBack(v)
			}
		}
		return rv
	}),
	"flatMap": NewNativeFunction("flatMap", func(c *Context, thisArg Value, args []Value) Value {
		var (
			mapper arrayMapper
		)
		EnsureFuncParams(c, "array.flatMap", args,
			mapper.ArgRule("mapper", true),
		)
		mapper.Build()
		thisArr := thisArg.(ValueArray)
		l := thisArr.Len()
		rv := NewArray(l)
		for i := 0; i < l; i++ {
			v := thisArr.GetIndex(i, c)
			mapped := mapper.Map(v, i, c)
			if mappedArray, is := mapped.(ValueArray); !is {
				c.RaiseRuntimeError("flatMap's mapper must return an array")
			} else {
				for _, v := range *mappedArray.Values {
					rv.PushBack(v)
				}
			}
		}
		return rv
	}),
	"filterMap": NewNativeFunction("filterMap", func(c *Context, thisArg Value, args []Value) Value {
		if len(args) != 1 {
			c.RaiseRuntimeError("array.filterMap: arguments length must be 1")
			return nil
		}
		f, isCallable := c.GetCallable(args[0])
		if !isCallable {
			c.RaiseRuntimeError("array.filterMap: argument 0 must be callable")
			return nil
		}
		thisArr := thisArg.(ValueArray)
		l := thisArr.Len()
		rv := NewArray(l)
		for i := 0; i < l; i++ {
			v := thisArr.GetIndex(i, c)
			f.Invoke(c, constUndefined, []Value{v, NewInt(int64(i))})
			if retArr, isArr := c.RetVal.(ValueArray); !isArr || retArr.Len() != 2 {
				c.RaiseRuntimeError("filterMap arg 0 must return an array with 2 elements")
			} else {
				targetVal, accepted := retArr.GetIndex(0, c), retArr.GetIndex(1, c)
				if accepted.IsTrue() {
					rv.PushBack(targetVal)
				}
			}
		}
		return rv
	}),
	"reduce": NewNativeFunction("reduce", func(c *Context, thisArg Value, args []Value) Value {
		thisArr := thisArg.(ValueArray)
		l := thisArr.Len()
		if len(args) < 1 {
			args = []Value{c.Eval("(prev, cur) => prev + cur", true)}
		}
		f, isCallable := c.GetCallable(args[0])
		if !isCallable {
			c.RaiseRuntimeError("array.reduce: argument 0 must be callable")
			return nil
		}
		var initVal Value
		start := 0
		if len(args) > 1 {
			initVal = args[1]
		} else {
			if thisArr.Len() < 1 {
				return constUndefined
			}
			initVal = thisArr.GetIndex(start, c)
			start++
		}
		for i := start; i < l; i++ {
			v := thisArr.GetIndex(i, c)
			f.Invoke(c, constUndefined, []Value{initVal, v, NewInt(int64(i)), thisArr})
			initVal = c.RetVal
		}
		return initVal
	}),
	"sum": NewNativeFunction("sum", func(c *Context, thisArg Value, args []Value) Value {
		var (
			mapper arrayMapper
		)
		EnsureFuncParams(c, "array.sum", args,
			mapper.ArgRule("mapper", false),
		)
		thisArr := *thisArg.(ValueArray).Values
		if l := len(thisArr); l > 0 {
			rv := mapper.Map(thisArr[0], 0, c)
			for i := 1; i < l; i++ {
				cur := mapper.Map(thisArr[i], i, c)
				rv = c.ValuesPlus(rv, cur)
			}
			return rv
		}
		return constUndefined
	}),
	"each": NewNativeFunction("each", func(c *Context, thisArg Value, args []Value) Value {
		if len(args) != 1 {
			c.RaiseRuntimeError("array.each: arguments length must be 1")
			return nil
		}
		f, isCallable := c.GetCallable(args[0])
		if !isCallable {
			c.RaiseRuntimeError("array.each: argument 0 must be callable")
			return nil
		}
		rv := NewArray()
		thisArr := thisArg.(ValueArray)
		l := thisArr.Len()
		for i := 0; i < l; i++ {
			v := thisArr.GetIndex(i, c)
			f.Invoke(c, constUndefined, []Value{v, NewInt(int64(i))})
		}
		return rv
	}),
	"push": NewNativeFunction("push", func(c *Context, thisArg Value, args []Value) Value {
		thisArr := thisArg.(ValueArray)
		for _, arg := range args {
			thisArr.PushBack(arg)
		}
		return NewInt(int64(len(args)))
	}),
	"slice": NewNativeFunction("slice", func(c *Context, thisArg Value, args []Value) Value {
		thisArr := thisArg.(ValueArray)
		arrLen := thisArr.Len()
		begin, end := 0, arrLen
		switch len(args) {
		case 2:
			{
				endArg, isInt := args[1].(ValueInt)
				if !isInt {
					c.RaiseRuntimeError("array.slice arg 1 must be int")
					return nil
				}
				end = int(endArg.Value())
				if end < 0 {
					end += arrLen
				}
			}
			fallthrough
		case 1:
			{
				beginArg, isInt := args[0].(ValueInt)
				if !isInt {
					c.RaiseRuntimeError("array.slice arg 0 must be int")
					return nil
				}
				begin = int(beginArg.Value())
				if begin < 0 {
					begin += arrLen
				}
			}
		case 0:
		default:
			c.RaiseRuntimeError("array.slice arguments num error")
			return nil
		}
		if end > arrLen {
			end = arrLen
		}
		if begin < 0 {
			begin = 0
		}
		if begin <= end && begin < arrLen {
			return thisArr.slice(begin, end)
		}
		return NewArray()
	}),
	"sort": NewNativeFunction("array.sort", func(c *Context, thisArg Value, args []Value) Value {
		thisArr, isArr := thisArg.(ValueArray)
		if !isArr {
			c.RaiseRuntimeError("array.sort: not an array")
			return nil
		}
		var (
			lessFn  ValueCallable
			reverse ValueBool
		)
		EnsureFuncParams(c, "array.sort", args,
			ArgRuleOptional("lessFunc", TypeCallable, &lessFn, nil),
			ArgRuleOptional("reverse", TypeBool, &reverse, NewBool(false)),
		)
		var s sort.Interface = &ArraySort{
			Array:   thisArr,
			LessFn:  lessFn,
			Context: c,
		}
		if reverse.Value() {
			s = sort.Reverse(s)
		}
		sort.Sort(s)
		return thisArg
	}, "lessFunc", "reverse"),
	"join": NewNativeFunction("array.join", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		joinStr := " "
		if len(args) >= 1 {
			joinStr = c.MustStr(args[0])
		}
		elems := make([]string, thisArr.Len())
		for i := range elems {
			elems[i] = thisArr.GetIndex(i, c).ToString(c)
		}
		return NewStr(strings.Join(elems, joinStr))
	}),
	"toMap": NewNativeFunction("array.toMap", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		var (
			keyMapper arrayMapper
			valMapper arrayMapper
		)
		EnsureFuncParams(c, "array.toMap", args,
			keyMapper.ArgRule("keyMapper", false),
			valMapper.ArgRule("valMapper", false),
		)
		keyMapper.Build()
		valMapper.Build()
		rv := NewObject()
		for i, item := range *(thisArr.Values) {
			k := keyMapper.Map(item, i, c).ToString(c)
			v := valMapper.Map(item, i, c)
			rv.SetMember(k, v, c)
		}
		return rv
	}),
	"count": NewNativeFunction("array.count", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		var (
			keyMapper arrayMapper
		)
		EnsureFuncParams(c, "array.count", args,
			keyMapper.ArgRule("keyMapper", false),
		)
		keyMapper.Build()
		cm := map[string]int64{}
		for i, item := range *(thisArr.Values) {
			k := keyMapper.Map(item, i, c).ToString(c)
			cm[k]++
		}
		rv := NewObject()
		for k, cnt := range cm {
			rv.SetMember(k, NewInt(cnt), c)
		}
		return rv
	}),
	"countIf": NewNativeFunction("array.countIf", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		var (
			predict ValueCallable
		)
		EnsureFuncParams(c, "array.countIf", args,
			ArgRuleRequired("predict", TypeCallable, &predict),
		)
		cnt := 0
		for i, item := range *(thisArr.Values) {
			c.Invoke(predict, nil, Args(item, NewInt(int64(i))))
			if c.RetVal.IsTrue() {
				cnt++
			}
		}
		return NewInt(int64(cnt))
	}),
	"toGroup": NewNativeFunction("array.toGroup", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		var (
			keyMapper arrayMapper
			valMapper arrayMapper
		)
		EnsureFuncParams(c, "array.toMap", args,
			keyMapper.ArgRule("keyMapper", false),
			valMapper.ArgRule("valMapper", false),
		)
		keyMapper.Build()
		valMapper.Build()
		rv := NewObject()
		for i, item := range *(thisArr.Values) {
			k := keyMapper.Map(item, i, c).ToString(c)
			group := rv.GetMember(k, c)
			v := valMapper.Map(item, i, c)
			if groupArr, ok := group.(ValueArray); ok {
				groupArr.PushBack(v)
			} else {
				rv.SetMember(k, NewArrayByValues(v), c)
			}
		}
		return rv
	}),
	"uniq": NewNativeFunction("array.uniq", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		var (
			valMapper arrayMapper
		)
		EnsureFuncParams(c, "array.uniq", args,
			valMapper.ArgRule("valMapper", false),
		)
		valMapper.Build()
		valMap := NewMap()
		rv := NewArray()
		for i, item := range *(thisArr.Values) {
			v := valMapper.Map(item, i, c)
			if _, found := valMap.get(c, v); found {
				continue
			}
			valMap.set(c, v, constNil)
			rv.PushBack(v)
		}
		return rv
	}),
	"chunk": NewNativeFunction("chunk", func(c *Context, this Value, args []Value) Value {
		var (
			chunkSize ValueInt
		)
		thisArr := c.MustArray(this)
		EnsureFuncParams(c, "array.chunk", args,
			ArgRuleRequired("chunkSize", TypeInt, &chunkSize),
		)
		items := *(thisArr.Values)
		n := len(items)
		cs := chunkSize.AsInt()
		rv := NewArray(n/cs + 1)
		for i := 0; i < n; i += cs {
			begin := i
			end := i + cs
			if end > n {
				end = n
			}
			chunk := NewArrayByValues(items[begin:end]...)
			rv.PushBack(chunk)
		}
		return rv
	}),
	"find": NewNativeFunction("array.find", func(c *Context, this Value, args []Value) Value {
		var (
			predict Value
			start   ValueInt
		)
		EnsureFuncParams(c, "array.find", args,
			ArgRuleRequired("predict", TypeAny, &predict),
			ArgRuleOptional("start", TypeInt, &start, NewInt(0)),
		)
		thisArr := c.MustArray(this)
		index, item := arrayFind(c, thisArr, predict, start.AsInt())
		if index >= 0 {
			return item
		} else {
			return constUndefined
		}
	}, "predict", "start"),
	"findIndex": NewNativeFunction("array.findIndex", func(c *Context, this Value, args []Value) Value {
		var (
			predict Value
			start   ValueInt
		)
		EnsureFuncParams(c, "array.findIndex", args,
			ArgRuleRequired("predict", TypeAny, &predict),
			ArgRuleOptional("start", TypeInt, &start, NewInt(0)),
		)
		thisArr := c.MustArray(this)
		index, _ := arrayFind(c, thisArr, predict, start.AsInt())
		return NewInt(int64(index))
	}),
	"times": NewNativeFunction("array.times", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		if thisArr.Len() < 1 {
			return constUndefined
		}
		if len(args) != 1 {
			c.RaiseRuntimeError("array.times requires 1 argument")
		}
		cb := c.MustCallable(args[0])
		ends := make([]int, thisArr.Len())
		for i := range ends {
			nv, ok := thisArr.GetIndex(i, c).(ValueInt)
			if !ok {
				c.RaiseRuntimeError("array.times: all arguments must be integer")
			}
			n := nv.AsInt()
			if n < 1 {
				c.RaiseRuntimeError("array.times: all arguments must be positive number")
			}
			ends[i] = n
		}
		curs := make([]Value, len(ends))
		for i := range curs {
			curs[i] = NewInt(0)
		}
		for {
			c.Invoke(cb, nil, Args(curs...))
			p := len(curs) - 1
			for ; p >= 0; p-- {
				next := curs[p].(ValueInt).AsInt() + 1
				if next < ends[p] {
					curs[p] = NewInt(int64(next))
					break
				} else {
					curs[p] = NewInt(0)
				}
			}
			if p < 0 {
				break
			}
		}
		return constUndefined
	}),
	"groupBy": NewNativeFunction("array.groupBy", func(c *Context, this Value, args []Value) Value {
		if len(args) == 0 {
			c.RaiseRuntimeError("groupBy requires at least one mapper")
		}
		mappers := make([]*arrayMapper, len(args))
		for i := range mappers {
			mappers[i] = &arrayMapper{}
		}
		argRules := make([]ArgRule, len(mappers))
		for i, m := range mappers {
			argRules[i] = m.ArgRule(fmt.Sprintf("mapper%d", i), true)
		}
		EnsureFuncParams(c, "array.groupBy", args, argRules...)
		for i := range mappers {
			mappers[i].Build()
		}
		groups := newGroupBy(*c.MustArray(this).Values, mappers).Execute(c)
		res := NewArray(len(groups))
		for _, g := range groups {
			item := NewArrayByValues(g...)
			res.PushBack(item)
		}
		return res
	}),
	"reverse": NewNativeFunction("array.reverse", func(c *Context, this Value, args []Value) Value {
		var (
			arr  = *c.MustArray(this).Values
			n    = len(arr)
			last = n - 1
			half = n / 2
		)
		for i := 0; i < half; i++ {
			j := last - i
			t := arr[i]
			arr[i] = arr[j]
			arr[j] = t
		}
		return this
	}),
}

type ArraySort struct {
	Array   ValueArray
	LessFn  ValueCallable
	Context *Context
}

func (s *ArraySort) Len() int {
	return s.Array.Len()
}

func (s *ArraySort) Less(i, j int) bool {
	vi := s.Array.GetIndex(i, s.Context)
	vj := s.Array.GetIndex(j, s.Context)
	if s.LessFn != nil {
		s.Context.Invoke(s.LessFn, nil, Args(vi))
		vi = s.Context.RetVal
		s.Context.Invoke(s.LessFn, nil, Args(vj))
		vj = s.Context.RetVal
	}
	return s.Context.ValuesLess(vi, vj)
}

func (s *ArraySort) Swap(i, j int) {
	vi := s.Array.GetIndex(i, s.Context)
	vj := s.Array.GetIndex(j, s.Context)
	s.Array.SetIndex(i, vj, s.Context)
	s.Array.SetIndex(j, vi, s.Context)
}

func init() {
	addMembersAndStatics(TypeArray, builtinArrayMethods)
}
