package runtime

import (
	"fmt"
	"math"
	"reflect"
	"time"
	"unsafe"
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

func (v ValueFloat) ToGoValue(*Context) interface{} {
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

func (v ValueFloat) Hash() int64 {
	return *(*int64)(unsafe.Pointer(&v.v))
}

func floatToDuration(c *Context, v Value, unit time.Duration) Value {
	timeMod := c.ImportModule("time", false, ImportTypeScript).(ValueObject)
	du := time.Duration(c.MustFloat(v) * float64(unit))
	duClass := Unbound(timeMod.GetMember("Duration", c)).(ValueType)
	return NewObjectAndInit(duClass, c, NewGoValue(du))
}

var builtinFloatMethods = map[string]ValueCallable{
	// math methods
	"floor": NewNativeFunction("floor", func(c *Context, this Value, args []Value) Value {
		return NewInt(int64(math.Floor(c.MustFloat(this))))
	}),
	"ceil": NewNativeFunction("ceil", func(c *Context, this Value, args []Value) Value {
		return NewInt(int64(math.Ceil(c.MustFloat(this))))
	}),
	"round": NewNativeFunction("round", func(c *Context, this Value, args []Value) Value {
		return NewInt(int64(math.Round(c.MustFloat(this))))
	}),
	// duration methods
	"ms": NewNativeFunction("ms", func(c *Context, this Value, args []Value) Value {
		return floatToDuration(c, this, time.Millisecond)
	}),
	"seconds": NewNativeFunction("seconds", func(c *Context, this Value, args []Value) Value {
		return floatToDuration(c, this, time.Second)
	}),
	"minutes": NewNativeFunction("seconds", func(c *Context, this Value, args []Value) Value {
		return floatToDuration(c, this, time.Minute)
	}),
	"hours": NewNativeFunction("seconds", func(c *Context, this Value, args []Value) Value {
		return floatToDuration(c, this, time.Hour)
	}),
	"days": NewNativeFunction("days", func(c *Context, this Value, args []Value) Value {
		return floatToDuration(c, this, 24*time.Hour)
	}),
	"weeks": NewNativeFunction("weeks", func(c *Context, this Value, args []Value) Value {
		return floatToDuration(c, this, 7*24*time.Hour)
	}),
}

func init() {
	addMembersAndStatics(TypeFloat, builtinFloatMethods)
}
