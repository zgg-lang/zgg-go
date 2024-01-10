package runtime

import "reflect"

type ValueNil struct {
	*ValueBase
}

func Nil() ValueNil {
	return constNil
}

func (ValueNil) GetIndex(int, *Context) Value {
	return constUndefined
}

func (v ValueNil) GetMember(name string, c *Context) Value {
	return getMemberByType(c, v, name)
}

func (ValueNil) IsTrue() bool {
	return false
}

func (ValueNil) Type() ValueType {
	return TypeNil
}

func (ValueNil) ToGoValue(*Context) interface{} {
	return nil
}

func (v ValueNil) GoType() reflect.Type {
	return reflect.TypeOf(nil)
}

func (ValueNil) ToString(*Context) string {
	return "nil"
}

func (ValueNil) CompareTo(other Value, c *Context) CompareResult {
	if _, ok := other.(ValueNil); ok {
		return CompareResultEqual
	}
	return CompareResultNotEqual
}

var constNil = ValueNil{ValueBase: &ValueBase{}}
