package runtime

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

func (c *Context) ValuesPlus(left, right Value) Value {
	if _, isStr := right.(ValueStr); isStr {
		c.RetVal = NewStr(left.ToString(c) + right.ToString(c))
		return c.RetVal
	}
	switch val1 := left.(type) {
	case ValueInt:
		switch val2 := right.(type) {
		case ValueInt:
			c.RetVal = NewInt(val1.Value() + val2.Value())
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(float64(val1.Value()) + val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Add(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueFloat:
		switch val2 := right.(type) {
		case ValueInt:
			c.RetVal = NewFloat(val1.Value() + float64(val2.Value()))
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(val1.Value() + val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Add(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueBigNum:
		var big2, bigR big.Float
		switch val2 := right.(type) {
		case ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Add(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Add(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueBigNum:
			var bigR big.Float
			bigR.Add(val1.Value(), val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueStr:
		c.RetVal = NewStr(val1.Value() + right.ToString(c))
		return c.RetVal
	case ValueArray:
		switch val2 := right.(type) {
		case ValueArray:
			{
				rv := NewArray(val1.Len() + val2.Len())
				for i := 0; i < val1.Len(); i++ {
					rv.PushBack(val1.GetIndex(i, c))
				}
				for i := 0; i < val2.Len(); i++ {
					rv.PushBack(val2.GetIndex(i, c))
				}
				c.RetVal = rv
				return c.RetVal
			}
		}
	default:
		{
			if opFn, ok := val1.GetMember("__add__", c).(ValueCallable); ok {
				c.Invoke(opFn, val1, func() []Value { return []Value{right} })
				return c.RetVal
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot plus between %s and %s", left.Type().Name, right.Type().Name))
	return nil
}

func (c *Context) ValuesMinus(left, right Value) Value {
	switch val1 := left.(type) {
	case ValueInt:
		switch val2 := right.(type) {
		case ValueInt:
			c.RetVal = NewInt(val1.Value() - val2.Value())
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(float64(val1.Value()) - val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Sub(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueFloat:
		switch val2 := right.(type) {
		case ValueInt:
			c.RetVal = NewFloat(val1.Value() - float64(val2.Value()))
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(val1.Value() - val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Sub(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueBigNum:
		var big2, bigR big.Float
		switch val2 := right.(type) {
		case ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Sub(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Sub(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueBigNum:
			var bigR big.Float
			bigR.Sub(val1.Value(), val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	default:
		{
			if opFn, ok := val1.GetMember("__sub__", c).(ValueCallable); ok {
				c.Invoke(opFn, val1, func() []Value { return []Value{right} })
				return c.RetVal
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot minus between %s and %s", left.Type().Name, right.Type().Name))
	return nil
}

func (c *Context) ValuesTimes(left, right Value) (ret Value) {
	switch val1 := left.(type) {
	case ValueInt:
		switch val2 := right.(type) {
		case ValueInt:
			c.RetVal = NewInt(val1.Value() * val2.Value())
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(float64(val1.Value()) * val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Mul(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueFloat:
		switch val2 := right.(type) {
		case ValueInt:
			c.RetVal = NewFloat(val1.Value() * float64(val2.Value()))
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(val1.Value() * val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Mul(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueBigNum:
		var big2, bigR big.Float
		switch val2 := right.(type) {
		case ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Mul(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Mul(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueBigNum:
			var bigR big.Float
			bigR.Mul(val1.Value(), val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueStr:
		switch val2 := right.(type) {
		case ValueInt:
			{
				var sb strings.Builder
				item := val1.Value()
				for times := int(val2.Value()); times > 0; times-- {
					sb.WriteString(item)
				}
				c.RetVal = NewStr(sb.String())
				return c.RetVal
			}
		}
	case ValueArray:
		switch val2 := right.(type) {
		case ValueInt:
			{
				rv := NewArray(val1.Len() * val2.AsInt())
				for times := val2.AsInt(); times > 0; times-- {
					for i := 0; i < val1.Len(); i++ {
						rv.PushBack(val1.GetIndex(i, c))
					}
				}
				c.RetVal = rv
				return c.RetVal
			}
		}
	default:
		{
			if opFn, ok := val1.GetMember("__mul__", c).(ValueCallable); ok {
				c.Invoke(opFn, val1, func() []Value { return []Value{right} })
				return c.RetVal
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot times between %s and %s", left.Type().Name, right.Type().Name))
	return nil
}

func (c *Context) ValuesDiv(left, right Value) Value {
	switch val1 := left.(type) {
	case ValueInt:
		switch val2 := right.(type) {
		case ValueInt:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("division by zero")
				return c.RetVal
			}
			c.RetVal = NewInt(val1.Value() / val2.Value())
			return c.RetVal
		case ValueFloat:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return c.RetVal
			}
			c.RetVal = NewFloat(float64(val1.Value()) / val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Quo(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueFloat:
		switch val2 := right.(type) {
		case ValueInt:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return c.RetVal
			}
			c.RetVal = NewFloat(val1.Value() / float64(val2.Value()))
			return c.RetVal
		case ValueFloat:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return c.RetVal
			}
			c.RetVal = NewFloat(val1.Value() / val2.Value())
			return c.RetVal
		case ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Quo(big1, val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	case ValueBigNum:
		var big2, bigR big.Float
		big2.SetPrec(1024)
		switch val2 := right.(type) {
		case ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Quo(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Quo(val1.Value(), &big2)
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		case ValueBigNum:
			var bigR big.Float
			bigR.Quo(val1.Value(), val2.Value())
			c.RetVal = NewBigNum(&bigR)
			return c.RetVal
		}
	default:
		{
			if opFn, ok := val1.GetMember("__div__", c).(ValueCallable); ok {
				c.Invoke(opFn, val1, func() []Value { return []Value{right} })
				return c.RetVal
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot div between %s and %s", left.Type().Name, right.Type().Name))
	return nil
}

func (c *Context) ValuesMod(left, right Value) Value {
	switch val1 := left.(type) {
	case ValueInt:
		switch val2 := right.(type) {
		case ValueInt:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return c.RetVal
			}
			c.RetVal = NewInt(val1.Value() % val2.Value())
			return c.RetVal
		}
	case ValueStr:
		switch val2 := right.(type) {
		case ValueArray:
			fargs := make([]interface{}, val2.Len())
			for i := range fargs {
				fargs[i] = val2.GetIndex(i, c).ToGoValue()
			}
			c.RetVal = NewStr(fmt.Sprintf(val1.Value(), fargs...))
		default:
			c.RetVal = NewStr(fmt.Sprintf(val1.Value(), val2.ToGoValue()))
		}
		return c.RetVal
	default:
		{
			if opFn, ok := val1.GetMember("__mod__", c).(ValueCallable); ok {
				c.Invoke(opFn, val1, func() []Value { return []Value{right} })
				return c.RetVal
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot mod between %s and %s", left.Type().Name, right.Type().Name))
	return nil
}

func (c *Context) ValuesPow(left, right Value) Value {
	switch val1 := left.(type) {
	case ValueInt:
		switch val2 := right.(type) {
		case ValueInt:
			if val2.Value() >= 0 {
				c.RetVal = NewInt(int64(math.Pow(float64(val1.Value()), float64(val2.Value()))))
			} else {
				c.RetVal = NewFloat(math.Pow(float64(val1.Value()), float64(val2.Value())))
			}
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(math.Pow(float64(val1.Value()), val2.Value()))
			return c.RetVal
		}
	case ValueFloat:
		switch val2 := right.(type) {
		case ValueInt:
			c.RetVal = NewFloat(math.Pow(val1.Value(), float64(val2.Value())))
			return c.RetVal
		case ValueFloat:
			c.RetVal = NewFloat(float64(val1.Value()) - val2.Value())
			return c.RetVal
		}
	default:
		{
			if opFn, ok := val1.GetMember("__pow__", c).(ValueCallable); ok {
				c.Invoke(opFn, val1, func() []Value { return []Value{right} })
				return c.RetVal
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot pow between %s and %s", left.Type().Name, right.Type().Name))
	return nil
}
