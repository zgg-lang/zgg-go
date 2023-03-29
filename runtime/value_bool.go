package runtime

import (
	"reflect"
	"strconv"
)

type ValueBool struct {
	*ValueBase
	v bool
}

func NewBool(v bool) ValueBool {
	if v {
		return boolTrue
	}
	return boolFalse
}

func (v ValueBool) IsTrue() bool {
	return v.v
}

func (v ValueBool) ToGoValue() interface{} {
	return v.Value()
}

func (v ValueBool) GoType() reflect.Type {
	return reflect.TypeOf(v.v)
}

func (v ValueBool) Value() bool {
	return v.v
}

func (v ValueBool) CompareTo(other Value, c *Context) CompareResult {
	if v2, ok := other.(ValueBool); ok {
		if v.Value() == v2.Value() {
			return CompareResultEqual
		}
	}
	return CompareResultNotEqual
}

func (v ValueBool) GetMember(name string, c *Context) Value {
	return getMemberByType(c, v, name)
}

func (ValueBool) GetIndex(int, *Context) Value {
	return constUndefined
}

func (ValueBool) Type() ValueType {
	return TypeBool
}

func (v ValueBool) ToString(*Context) string {
	return strconv.FormatBool(v.Value())
}

func (v ValueBool) Hash() int64 {
	if v.v {
		return 1
	} else {
		return 0
	}
}

var (
	boolTrue  = ValueBool{ValueBase: &ValueBase{}, v: true}
	boolFalse = ValueBool{ValueBase: &ValueBase{}, v: false}
)
