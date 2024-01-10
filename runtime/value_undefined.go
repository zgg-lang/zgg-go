package runtime

import "reflect"

type ValueUndefined struct {
	*ValueBase
}

func Undefined() ValueUndefined {
	return constUndefined
}

func (ValueUndefined) IsTrue() bool {
	return false
}

func (ValueUndefined) ToGoValue(*Context) interface{} {
	return nil
}

func (ValueUndefined) GoType() reflect.Type {
	return reflect.TypeOf(nil)
}

func (ValueUndefined) GetIndex(int, *Context) Value {
	return constUndefined
}

func (v ValueUndefined) GetMember(name string, c *Context) Value {
	return getExtMember(v, name, c)
}

func (ValueUndefined) Type() ValueType {
	return TypeUndefined
}

func (ValueUndefined) ToString(*Context) string {
	return "undefined"
}

func (ValueUndefined) CompareTo(other Value, c *Context) CompareResult {
	if _, ok := other.(ValueUndefined); ok {
		return CompareResultEqual
	}
	return CompareResultNotEqual
}

var constUndefined = ValueUndefined{ValueBase: &ValueBase{}}

func IsUndefined(v Value) bool {
	_, r := v.(ValueUndefined)
	return r
}
