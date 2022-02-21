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
		c.OnRuntimeError(fmt.Sprintf("set array item error: Out of bound length %d index %d", v.Len(), index))
		return
	}
	(*v.Values)[index] = value
}

func (v ValueArray) GetMember(name string, c *Context) Value {
	if member, found := builtinArrayMethods[name]; found {
		return makeMember(v, member)
	}
	return getExtMember(v, name, c)
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

func arrayFind(c *Context, arr ValueArray, predict ValueCallable, start int) (int, Value) {
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

var builtinArrayMethods = map[string]ValueCallable{
	"map": &ValueBuiltinFunction{
		body: func(c *Context, thisArg Value, args []Value) Value {
			var (
				mapFunc    ValueCallable
				fieldName  ValueStr
				fieldIndex ValueInt
				mapType    int
			)
			EnsureFuncParams(c, "array.map", args,
				ArgRuleOneOf{"mapper",
					[]ValueType{TypeCallable, TypeStr, TypeInt},
					[]interface{}{&mapFunc, &fieldName, &fieldIndex},
					&mapType, nil, nil,
				},
			)
			thisArr := thisArg.(ValueArray)
			l := thisArr.Len()
			rv := NewArray(l)
			switch mapType {
			case 0:
				for i := 0; i < l; i++ {
					v := thisArr.GetIndex(i, c)
					mapFunc.Invoke(c, constUndefined, []Value{v, NewInt(int64(i))})
					rv.PushBack(c.RetVal)
				}
			case 1:
				name := fieldName.Value()
				for i := 0; i < l; i++ {
					v := thisArr.GetIndex(i, c)
					rv.PushBack(v.GetMember(name, c))
				}
			case 2:
				index := fieldIndex.AsInt()
				for i := 0; i < l; i++ {
					v := thisArr.GetIndex(i, c)
					rv.PushBack(v.GetIndex(index, c))
				}
			}
			return rv
		},
	},
	"filter": &ValueBuiltinFunction{
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				c.OnRuntimeError("array.filter: arguments length must be 1")
				return nil
			}
			f, isCallable := args[0].(ValueCallable)
			if !isCallable {
				c.OnRuntimeError("array.filter: argument 0 must be callable")
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
		},
	},
	"reduce": &ValueBuiltinFunction{
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisArr := thisArg.(ValueArray)
			l := thisArr.Len()
			if len(args) < 1 {
				args = []Value{c.Eval("(prev, cur) => prev + cur", true)}
			}
			f, isCallable := args[0].(ValueCallable)
			if !isCallable {
				c.OnRuntimeError("array.reduce: argument 0 must be callable")
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
		},
	},
	"each": &ValueBuiltinFunction{
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				c.OnRuntimeError("array.each: arguments length must be 1")
				return nil
			}
			f, isCallable := args[0].(ValueCallable)
			if !isCallable {
				c.OnRuntimeError("array.each: argument 0 must be callable")
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
		},
	},
	"push": &ValueBuiltinFunction{
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisArr := thisArg.(ValueArray)
			for _, arg := range args {
				thisArr.PushBack(arg)
			}
			return NewInt(int64(len(args)))
		},
	},
	"slice": &ValueBuiltinFunction{
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisArr := thisArg.(ValueArray)
			arrLen := thisArr.Len()
			begin, end := 0, arrLen
			switch len(args) {
			case 2:
				{
					endArg, isInt := args[1].(ValueInt)
					if !isInt {
						c.OnRuntimeError("array.slice arg 1 must be int")
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
						c.OnRuntimeError("array.slice arg 0 must be int")
						return nil
					}
					begin = int(beginArg.Value())
					if begin < 0 {
						begin += arrLen
					}
				}
			case 0:
			default:
				c.OnRuntimeError("array.slice arguments num error")
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
		},
	},
	"sort": NewNativeFunction("array.sort", func(c *Context, thisArg Value, args []Value) Value {
		thisArr, isArr := thisArg.(ValueArray)
		if !isArr {
			c.OnRuntimeError("array.sort: not an array")
			return nil
		}
		var (
			lessFn  ValueCallable
			reverse ValueBool
		)
		EnsureFuncParams(c, "array.sort", args,
			ArgRuleOptional{"lessFunc", TypeCallable, &lessFn, nil},
			ArgRuleOptional{"reverse", TypeBool, &reverse, NewBool(false)},
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
		switch len(args) {
		case 0:
			{
				rv := NewObject()
				for _, item := range *(thisArr.Values) {
					rv.SetMember(item.ToString(c), item, c)
				}
				return rv
			}
		case 1:
			{
				getkey := c.MustCallable(args[0], "getkey")
				rv := NewObject()
				for _, item := range *(thisArr.Values) {
					c.Invoke(getkey, nil, func() []Value { return []Value{item} })
					key := c.RetVal.ToString(c)
					rv.SetMember(key, item, c)
				}
				return rv
			}
		default:
			c.OnRuntimeError("array.toMap requires 0 or 1 argument")
		}
		return nil
	}),
	"toGroup": NewNativeFunction("array.toGroup", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		if len(args) != 1 {
			c.OnRuntimeError("array.toGroup requires 1 argument")
		}
		getkey := c.MustCallable(args[0], "getkey")
		rv := NewObject()
		for _, item := range *(thisArr.Values) {
			c.Invoke(getkey, nil, func() []Value { return []Value{item} })
			key := c.RetVal.ToString(c)
			group := rv.GetMember(key, c)
			if groupArr, ok := group.(ValueArray); ok {
				groupArr.PushBack(item)
			} else {
				rv.SetMember(key, NewArrayByValues(item), c)
			}
		}
		return rv
	}),
	"find": NewNativeFunction("array.find", func(c *Context, this Value, args []Value) Value {
		var (
			predict ValueCallable
			start   ValueInt
		)
		EnsureFuncParams(c, "array.find", args,
			ArgRuleRequired{"predict", TypeCallable, &predict},
			ArgRuleOptional{"start", TypeInt, &start, NewInt(0)},
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
		thisArr := c.MustArray(this)
		start := 0
		var predict ValueCallable
		switch len(args) {
		case 2:
			start = int(c.MustInt(args[1]))
			fallthrough
		case 1:
			predict = c.MustCallable(args[0])
		default:
			c.OnRuntimeError("array.find(predict[, start]) requires 1 or 2 argument")
		}
		index, _ := arrayFind(c, thisArr, predict, start)
		return NewInt(int64(index))
	}),
	"times": NewNativeFunction("array.times", func(c *Context, this Value, args []Value) Value {
		thisArr := c.MustArray(this)
		if thisArr.Len() < 1 {
			return constUndefined
		}
		if len(args) != 1 {
			c.OnRuntimeError("array.times requires 1 argument")
		}
		cb := c.MustCallable(args[0])
		ends := make([]int, thisArr.Len())
		for i := range ends {
			nv, ok := thisArr.GetIndex(i, c).(ValueInt)
			if !ok {
				c.OnRuntimeError("array.times: all arguments must be integer")
			}
			n := nv.AsInt()
			if n < 1 {
				c.OnRuntimeError("array.times: all arguments must be positive number")
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
