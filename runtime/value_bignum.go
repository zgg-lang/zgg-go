package runtime

import (
	"math/big"
	"reflect"
)

type ValueBigNum struct {
	*ValueBase
	v *big.Float
}

func NewBigNum(v *big.Float) ValueBigNum {
	return ValueBigNum{ValueBase: &ValueBase{}, v: v}
}

func (v ValueBigNum) GoType() reflect.Type {
	return reflect.TypeOf(v.v)
}

func (v ValueBigNum) ToGoValue() interface{} {
	return v.Value()
}

func (v ValueBigNum) Value() *big.Float {
	return v.v
}

func (v ValueBigNum) IsTrue() bool {
	return v.v.Cmp(big.NewFloat(0)) != 0
}

func (v ValueBigNum) CompareTo(other Value, c *Context) CompareResult {
	var otherFloat *big.Float
	switch v2 := other.(type) {
	case ValueInt:
		otherFloat = big.NewFloat(float64(v2.Value()))
	case ValueFloat:
		otherFloat = big.NewFloat(v2.Value())
	case ValueBigNum:
		otherFloat = v2.v
	default:
		return CompareResultNotEqual
	}
	switch v.v.Cmp(otherFloat) {
	case 0:
		return CompareResultEqual
	case 1:
		return CompareResultGreater
	}
	return CompareResultLess
}

func (v ValueBigNum) GetMember(name string, c *Context) Value {
	return getMemberByType(c, v, name)
}

func (ValueBigNum) GetIndex(int, *Context) Value {
	return constUndefined
}

func (ValueBigNum) Type() ValueType {
	return TypeBigNum
}

func (v ValueBigNum) ToString(*Context) string {
	return v.v.Text('f', -1)
}

var builtinBigNumMethods = map[string]ValueCallable{}
