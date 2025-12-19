package runtime

import (
	"reflect"
	"strconv"
	"sync"
	"time"
)

const (
	poolMin = -1000
	poolMax = 1000
)

type ValueInt struct {
	*ValueBase
	v int64
}

var intPoolInitOnce sync.Once

func NewInt(v int64) ValueInt {
	intPoolInitOnce.Do(initIntPool)
	var r ValueInt
	if v >= poolMin && v <= poolMax {
		r = intPool[v-poolMin]
	} else {
		r = ValueInt{ValueBase: &ValueBase{}, v: v}
	}
	return r
}

func (v ValueInt) GoType() reflect.Type {
	return reflect.TypeOf(v.v)
}

func (v ValueInt) ToGoValue(*Context) interface{} {
	return v.Value()
}

func (v ValueInt) Value() int64 {
	return v.v
}

func (v ValueInt) AsInt() int {
	return int(v.v)
}

func (v ValueInt) IsTrue() bool {
	return v.v != 0
}

func (v ValueInt) CompareTo(other Value, c *Context) CompareResult {
	if v2, ok := other.(ValueFloat); ok {
		vf := float64(v.Value())
		if vf == v2.Value() {
			return CompareResultEqual
		} else if vf < v2.Value() {
			return CompareResultLess
		} else {
			return CompareResultGreater
		}
	}
	if v2, ok := other.(ValueInt); ok {
		if v.Value() == v2.Value() {
			return CompareResultEqual
		} else if v.Value() < v2.Value() {
			return CompareResultLess
		} else {
			return CompareResultGreater
		}
	}
	return CompareResultNotEqual
}

func (v ValueInt) GetMember(name string, c *Context) Value {
	return getMemberByType(c, v, name)
}

func (ValueInt) GetIndex(int, *Context) Value {
	return constUndefined
}

func (ValueInt) Type() ValueType {
	return TypeInt
}

func (v ValueInt) ToString(*Context) string {
	return strconv.FormatInt(v.Value(), 10)
}

func (v ValueInt) Hash() int64 {
	return v.v
}

var intPool [poolMax - poolMin + 1]ValueInt

func initIntPool() {
	for i := range intPool {
		intPool[i] = ValueInt{ValueBase: &ValueBase{}, v: int64(i + poolMin)}
	}
}

func init() {
	intPoolInitOnce.Do(initIntPool)
	addMembersAndStatics(TypeInt, builtinIntMethods)
}

func intToDuration(c *Context, v Value, unit time.Duration) Value {
	timeMod := c.ImportModule("time", false, ImportTypeScript).(ValueObject)
	du := time.Duration(c.MustInt(v)) * unit
	duClass := Unbound(timeMod.GetMember("Duration", c)).(ValueType)
	return NewObjectAndInit(duClass, c, NewGoValue(du))
}

var builtinIntMethods = map[string]ValueCallable{
	"times": NewNativeFunction("times", func(c *Context, this Value, args []Value) Value {
		times := c.MustInt(this)
		if times < 0 {
			c.RaiseRuntimeError("int.times times must not less than 0, but got %d", times)
		}
		if len(args) != 1 {
			c.RaiseRuntimeError("int.times requires 1 argument")
		}
		callback := c.MustCallable(args[0])
		for i := int64(0); i < times; i++ {
			c.Invoke(callback, nil, Args(NewInt(i)))
		}
		return constUndefined
	}),
	"str": NewNativeFunction("str", func(c *Context, this Value, args []Value) Value {
		var base ValueInt
		EnsureFuncParams(c, "str", args, ArgRuleOptional("base", TypeInt, &base, NewInt(10)))
		b := base.AsInt()
		if b < 2 || b > 36 {
			c.RaiseRuntimeError("str: 2 <= base <= 36")
		}
		val := c.MustInt(this)
		return NewStr(strconv.FormatInt(val, b))
	}, "base"),
	"parse": NewNativeFunction("parse", func(c *Context, this Value, args []Value) Value {
		var (
			value ValueStr
			base  ValueInt
		)
		EnsureFuncParams(c, "parse", args,
			ArgRuleRequired("value", TypeStr, &value),
			ArgRuleOptional("base", TypeInt, &base, NewInt(10)),
		)
		b := base.AsInt()
		if b < 2 || b > 36 {
			c.RaiseRuntimeError("parse: 2 <= base <= 36")
		}
		v, err := strconv.ParseInt(value.Value(), b, 64)
		if err != nil {
			c.RaiseRuntimeError("parse: %s", err)
		}
		return NewInt(v)
	}, "value", "base"),
	"__next__": NewNativeFunction("__next__", func(c *Context, this Value, args []Value) Value {
		return NewInt(c.MustInt(this) + 1)
	}),
	// duration methods
	"ms": NewNativeFunction("ms", func(c *Context, this Value, args []Value) Value {
		return intToDuration(c, this, time.Millisecond)
	}),
	"seconds": NewNativeFunction("seconds", func(c *Context, this Value, args []Value) Value {
		return intToDuration(c, this, time.Second)
	}),
	"minutes": NewNativeFunction("minutes", func(c *Context, this Value, args []Value) Value {
		return intToDuration(c, this, time.Minute)
	}),
	"hours": NewNativeFunction("hours", func(c *Context, this Value, args []Value) Value {
		return intToDuration(c, this, time.Hour)
	}),
	"days": NewNativeFunction("days", func(c *Context, this Value, args []Value) Value {
		return intToDuration(c, this, 24*time.Hour)
	}),
	"weeks": NewNativeFunction("weeks", func(c *Context, this Value, args []Value) Value {
		return intToDuration(c, this, 7*24*time.Hour)
	}),
}
