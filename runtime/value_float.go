package runtime

import (
	"fmt"
	"reflect"
)

type ValueFloat struct {
	*ValueBase
	v float64
}

func NewFloat(v float64) ValueFloat {
	return ValueFloat{ValueBase: &ValueBase{}, v: v}
}

func (v ValueFloat) GoType() reflect.Type {
	return reflect.TypeOf(v.v)
}

func (v ValueFloat) ToGoValue() interface{} {
	return v.Value()
}

func (v ValueFloat) Value() float64 {
	return v.v
}

func (v ValueFloat) IsTrue() bool {
	return v.v != 0
}

func (v ValueFloat) CompareTo(other Value, c *Context) CompareResult {
	if v2, ok := other.(ValueFloat); ok {
		if v.Value() == v2.Value() {
			return CompareResultEqual
		} else if v.Value() < v2.Value() {
			return CompareResultLess
		} else {
			return CompareResultGreater
		}
	}
	if v2, ok := other.(ValueInt); ok {
		v2f := float64(v2.Value())
		if v.Value() == v2f {
			return CompareResultEqual
		} else if v.Value() < v2f {
			return CompareResultLess
		} else {
			return CompareResultGreater
		}
	}
	return CompareResultNotEqual
}

func (v ValueFloat) GetMember(name string, c *Context) Value {
	return getMemberByType(c, v, name)
}

func (ValueFloat) GetIndex(int, *Context) Value {
	return constUndefined
}

func (ValueFloat) Type() ValueType {
	return TypeFloat
}

func (v ValueFloat) ToString(*Context) string {
	return fmt.Sprint(v.Value())
}
