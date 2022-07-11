package runtime

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type GoValue struct {
	*ValueBase
	v reflect.Value
}

type GoFunc struct {
	GoValue
	args []string
	name string
}

func NewGoValue(v interface{}) Value {
	if reflected, ok := v.(reflect.Value); ok {
		return NewReflectedGoValue(reflected)
	}
	return NewReflectedGoValue(reflect.ValueOf(v))
}

func (v GoValue) ReflectedValue() reflect.Value {
	return v.v
}

func NewReflectedGoValue(v reflect.Value) Value {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if v.IsNil() {
			return constNil
		}
	}
	if v.Kind() == reflect.Func {
		vt := v.Type()
		return GoFunc{
			GoValue: GoValue{ValueBase: new(ValueBase), v: v},
			name:    vt.Name(),
		}
	}
	return GoValue{ValueBase: new(ValueBase), v: v}
}

func (v GoValue) Type() ValueType {
	return TypeGoValue
}

func (v GoValue) CompareTo(other Value, c *Context) CompareResult {
	return CompareResultNotEqual
}

func (v GoValue) ToString(*Context) string {
	if !v.v.IsValid() {
		return "undefined"
	}
	if v.v.Kind() == reflect.Slice {
		if v.v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.v.Interface().([]byte))
		}
	}
	return fmt.Sprint(v.v.Interface())
}

func (v GoValue) GoType() reflect.Type {
	return v.v.Type()
}

func (v GoValue) ToGoValue() interface{} {
	return v.v.Interface()
}

func (v GoValue) IsTrue() bool {
	return !v.v.IsZero()
}

func (v GoValue) GetIndex(index int, c *Context) Value {
	switch v.v.Kind() {
	case reflect.Slice:
		if index < 0 || index >= v.v.Len() {
			return constUndefined
		}
		return NewReflectedGoValue(v.v.Index(index))
	case reflect.Array:
		if index < 0 || index >= v.v.Len() {
			return constUndefined
		}
		return NewReflectedGoValue(v.v.Index(index))
	case reflect.String:
		if index < 0 || index >= v.v.Len() {
			return constUndefined
		}
		return NewReflectedGoValue(v.v.Index(index))
	}
	return constUndefined
}

func (v GoValue) SetIndex(index int, val Value, c *Context) {
	switch v.v.Kind() {
	case reflect.Slice, reflect.Array:
		gv := MakeGoValueByArg(c, val, v.v.Type().Elem())
		v.v.Index(index).Set(gv)
	}
}

func (v GoValue) canBeNil() bool {
	return CanBeNil(v.v)
}

func (v GoValue) SetMember(key string, val Value, c *Context) {
	if v.v.Kind() != reflect.Struct {
		c.RaiseRuntimeError("cannot set member on non-struct go value. value kind is %s", v.v.Kind())
	}
	_, fieldFound := v.v.Type().FieldByName(key)
	if !fieldFound {
		c.RaiseRuntimeError("cannot set member %s: not exists", key)
	}
	fieldVal := v.v.FieldByName(key)
	if gv, ok := val.(GoValue); ok {
		fieldVal.Set(gv.v)
	} else {
		toGoValue(c, val, fieldVal)
	}
}

func (v GoValue) Len() int {
	switch v.v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.v.Len()
	}
	return -1
}

func MakeGoValueByArg(c *Context, val Value, typ reflect.Type) reflect.Value {
	if goval, ok := val.(GoValue); ok {
		return goval.ReflectedValue()
	}
	newVal := reflect.New(typ).Elem()
	toGoValue(c, val, newVal)
	return newVal
}

func gochanGetTimeout(val Value) (float64, bool) {
	var timeout float64
	switch timeoutArg := val.(type) {
	case GoValue:
		switch timeoutArg.v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			timeout = float64(timeoutArg.v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			timeout = float64(timeoutArg.v.Uint())
		case reflect.Float32, reflect.Float64:
			timeout = timeoutArg.v.Float()
		default:
			return 0, false
		}
	case ValueInt:
		timeout = float64(timeoutArg.Value())
	case ValueFloat:
		timeout = timeoutArg.Value()
	default:
		return 0, false
	}
	return timeout, true
}

var (
	gomapEach = NewNativeFunction("gomap.each", func(c *Context, this Value, args []Value) Value {
		m := this.(GoValue).ReflectedValue()
		if m.Kind() != reflect.Map {
			c.RaiseRuntimeError("cannot invoke each from a non-map go value. value kind is %s", m.Kind())
		}
		callback := c.MustCallable(args[0], "callback")
		keys := m.MapKeys()
		for _, k := range keys {
			v := m.MapIndex(k)
			c.Invoke(callback, nil, Args(NewGoValue(k), NewGoValue(v)))
		}
		return Undefined()
	})
	gomapGet = NewNativeFunction("gomap.get", func(c *Context, this Value, args []Value) Value {
		m := this.(GoValue).ReflectedValue()
		if m.Kind() != reflect.Map {
			c.RaiseRuntimeError("cannot get from a non-map go value")
		}
		key := MakeGoValueByArg(c, args[0], m.Type().Key())
		val := m.MapIndex(key)
		if val.IsZero() {
			return constUndefined
		}
		return NewReflectedGoValue(val)
	})
	gomapSet = NewNativeFunction("gomap.set", func(c *Context, this Value, args []Value) Value {
		m := this.(GoValue).ReflectedValue()
		if m.Kind() != reflect.Map {
			c.RaiseRuntimeError("cannot set a non-map go value")
		}
		key := MakeGoValueByArg(c, args[0], m.Type().Key())
		val := MakeGoValueByArg(c, args[1], m.Type().Elem())
		m.SetMapIndex(key, val)
		return this
	})
	gomapDelete = NewNativeFunction("gomap.delete", func(c *Context, this Value, args []Value) Value {
		m := this.(GoValue).ReflectedValue()
		if m.Kind() != reflect.Map {
			c.RaiseRuntimeError("cannot delete from a non-map go value")
		}
		key := MakeGoValueByArg(c, args[0], m.Type().Key())
		val := reflect.Zero(m.Type().Elem())
		m.SetMapIndex(key, val)
		return this
	})
	gochanRecv = NewNativeFunction("gochan.recv", func(c *Context, this Value, args []Value) Value {
		m := this.(GoValue).ReflectedValue()
		if m.Kind() != reflect.Chan {
			c.RaiseRuntimeError("cannot await from a non-chan go value")
		}
		if len(args) > 0 {
			timeout, ok := gochanGetTimeout(args[0])
			if !ok {
				c.RaiseRuntimeError("invalid timeout argument %s", args[0].ToString(c))
			}
			if timeout <= 0 {
				v, ok := m.TryRecv()
				return NewArrayByValues(NewReflectedGoValue(v), NewBool(ok))
			}
			timerCh := time.After(time.Duration(timeout * float64(time.Second)))
			chosen, recv, ok := reflect.Select([]reflect.SelectCase{
				{Dir: reflect.SelectRecv, Chan: m},
				{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timerCh)},
			})
			if chosen == 0 {
				return NewArrayByValues(NewReflectedGoValue(recv), NewBool(ok))
			} else {
				return NewArrayByValues(constUndefined, NewBool(false))
			}
		}
		v, ok := m.Recv()
		return NewArrayByValues(NewReflectedGoValue(v), NewBool(ok))
	})
	gochanSend = NewNativeFunction("gochan.send", func(c *Context, this Value, args []Value) Value {
		m := this.(GoValue).ReflectedValue()
		if m.Kind() != reflect.Chan {
			c.RaiseRuntimeError("cannot send to a non-chan go value")
		}
		v := MakeGoValueByArg(c, args[0], m.Type().Elem())
		if len(args) > 1 {
			timeout, ok := gochanGetTimeout(args[1])
			if !ok {
				c.RaiseRuntimeError("invalid timeout argument %s", args[1].ToString(c))
			}
			if timeout <= 0 {
				return NewBool(m.TrySend(v))
			}
			timerCh := time.After(time.Duration(timeout * float64(time.Second)))
			chosen, _, _ := reflect.Select([]reflect.SelectCase{
				{Dir: reflect.SelectSend, Chan: m, Send: v},
				{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timerCh)},
			})
			return NewBool(chosen == 0)
		}
		m.Send(v)
		return boolTrue
	})
	goMethods = NewNativeFunction("go.methods", func(c *Context, this Value, args []Value) Value {
		var t reflect.Type
		switch thisValue := Unbound(this).(type) {
		case GoValue:
			t = thisValue.v.Type()
		case GoType:
			t = thisValue.typ
		default:
			return NewArray()
		}
		n := t.NumMethod()
		rv := NewArray(n)
		for i := 0; i < n; i++ {
			method := t.Method(i)
			mt := method.Type
			item := NewObject()
			// name
			item.SetMember("name", NewStr(method.Name), c)
			// in arg types
			ins := make([]string, 0, mt.NumIn()-1)
			for j := 1; j < mt.NumIn(); j++ { // skip this
				ins = append(ins, mt.In(j).Name())
			}
			item.SetMember("in_", FromGoValue(reflect.ValueOf(ins), c), c)
			// out types
			outs := make([]string, 0, mt.NumOut())
			for j := 0; j < mt.NumOut(); j++ {
				outs = append(outs, mt.Out(j).Name())
			}
			item.SetMember("out", FromGoValue(reflect.ValueOf(outs), c), c)
			sign := method.Name + "(" + strings.Join(ins, ", ") + ")"
			switch len(outs) {
			case 0:
				// nothing to add
			case 1:
				sign += " " + outs[0]
			default:
				sign += " (" + strings.Join(outs, ", ") + ")"
			}
			item.SetMember("sign", NewStr(sign), c)
			rv.PushBack(item)
		}
		return rv
	})
	goSign = NewNativeFunction("go.sign", func(c *Context, this Value, args []Value) Value {
		gv, ok := Unbound(this).(GoFunc)
		if !ok || gv.v.Kind() != reflect.Func {
			if ok {
				fmt.Println(gv.v.Interface(), gv.v.Kind())
			} else {
				println("get govalue failed", reflect.TypeOf(Unbound(this)).Name())
			}
			return constNil
		}
		mt := gv.v.Type()
		item := NewObject()
		// name
		item.SetMember("name", NewStr(gv.name), c)
		// in arg types
		ins := make([]string, 0, mt.NumIn()-1)
		for j := 0; j < mt.NumIn(); j++ { // skip this
			ins = append(ins, mt.In(j).Name())
		}
		item.SetMember("in_", FromGoValue(reflect.ValueOf(ins), c), c)
		// out types
		outs := make([]string, 0, mt.NumOut())
		for j := 0; j < mt.NumOut(); j++ {
			outs = append(outs, mt.Out(j).Name())
		}
		item.SetMember("out", FromGoValue(reflect.ValueOf(outs), c), c)
		sign := gv.name + "(" + strings.Join(ins, ", ") + ")"
		switch len(outs) {
		case 0:
			// nothing to add
		case 1:
			sign += " " + outs[0]
		default:
			sign += " (" + strings.Join(outs, ", ") + ")"
		}
		item.SetMember("sign", NewStr(sign), c)
		return item
	})
	goFields = NewNativeFunction("go.fields", func(c *Context, this Value, args []Value) Value {
		gv, ok := Unbound(this).(GoValue)
		if !ok || gv.v.Kind() != reflect.Struct {
			return constNil
		}
		t := gv.v.Type()
		n := t.NumField()
		rv := NewArray()
		for i := 0; i < n; i++ {
			item := NewObject()
			f := t.Field(i)
			item.SetMember("name", NewStr(f.Name), c)
			item.SetMember("type", NewStr(f.Type.Name()), c)
			rv.PushBack(item)
		}
		return rv
	})
	goAs = NewNativeFunction("go.as", func(c *Context, this Value, args []Value) Value {
		gv, ok := Unbound(this).(GoValue)
		if !ok {
			return constNil
		}
		var targetType GoType
		EnsureFuncParams(c, "as", args, ArgRuleRequired{"targetType", TypeGoType, &targetType})
		return NewGoValue(gv.ReflectedValue().Convert(targetType.GoType()))
	})
	goIs = NewNativeFunction("go.is", func(c *Context, this Value, args []Value) Value {
		gv, ok := Unbound(this).(GoValue)
		if !ok {
			return constNil
		}
		var targetType GoType
		EnsureFuncParams(c, "as", args, ArgRuleRequired{"targetType", TypeGoType, &targetType})
		return NewBool(gv.ReflectedValue().Type().ConvertibleTo(targetType.GoType()))
	})
)

func (v GoValue) GetMember(key string, c *Context) Value {
	if v.canBeNil() && v.v.IsNil() {
		return constNil
	}
	// process builtin methods
	switch key {
	case "int":
		switch v.v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewInt(v.v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return NewInt(int64(v.v.Uint()))
		case reflect.Float32, reflect.Float64:
			return NewInt(int64(v.v.Float()))
		default:
			c.RaiseRuntimeError("cannot convert type %s to int", v.v.Kind())
		}
	case "float":
		switch v.v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewFloat(float64(v.v.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return NewFloat(float64(v.v.Uint()))
		case reflect.Float32, reflect.Float64:
			return NewFloat(float64(v.v.Float()))
		default:
			c.RaiseRuntimeError("cannot convert type %s to float", v.v.Kind())
		}
	case "str":
		return NewStr(fmt.Sprint(v.v.Interface()))
	case "bool":
		switch v.v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewBool(v.v.Int() != 0)
		case reflect.Bool:
			return NewBool(v.v.Bool())
		default:
			c.RaiseRuntimeError("cannot convert type %s to bool", v.v.Kind())
		}
	case "bytes":
		if bs, ok := v.v.Interface().([]byte); ok {
			return NewBytes(bs)
		}
		return NewBytes([]byte(fmt.Sprint(v.v.Interface())))
	case "zgg":
		return FromGoValue(v.v, c)
	case "ptr":
		return NewReflectedGoValue(v.v.Addr())
	case "el":
		return NewReflectedGoValue(v.v.Elem())
	case "kind":
		return NewStr(v.v.Kind().String())
	// for map
	case "each":
		return makeMember(v, gomapEach, c)
	case "get":
		return makeMember(v, gomapGet, c)
	case "set":
		return makeMember(v, gomapSet, c)
	case "delete":
		return makeMember(v, gomapDelete, c)
	// for chan
	case "send":
		return makeMember(v, gochanSend, c)
	case "recv":
		return makeMember(v, gochanRecv, c)
	case "methods":
		return makeMember(v, goMethods, c)
	case "fields":
		return makeMember(v, goFields, c)
	case "sign":
		return makeMember(v, goSign, c)
	case "as":
		return makeMember(v, goAs, c)
	case "is":
		return makeMember(v, goIs, c)
	}
	// isPtr := false
	goval := v.v
	gotype := goval.Type()
	// if gotype.Kind() == reflect.Ptr {
	// 	// isPtr = true
	// 	gotype = gotype.Elem()
	// 	goval = goval.Elem()
	// }
	switch gotype.Kind() {
	case reflect.Map:
		if gotype.Key().Kind() == reflect.String {
			rv := goval.MapIndex(reflect.ValueOf(key))
			if rv.IsValid() {
				return makeMember(v, NewReflectedGoValue(rv), c)
			}
		}
	case reflect.Struct:
		if _, found := gotype.FieldByName(key); found {
			return makeMember(v, NewReflectedGoValue(goval.FieldByName(key)), c)
		}
	case reflect.Ptr:
		elType := gotype.Elem()
		switch elType.Kind() {
		case reflect.Struct:
			if _, found := elType.FieldByName(key); found {
				return makeMember(v, NewReflectedGoValue(goval.Elem().FieldByName(key)), c)
			}
		}
	}
	if _, found := gotype.MethodByName(key); found {
		return makeMember(v, NewReflectedGoValue(goval.MethodByName(key)), c)
	}
	if goval.CanAddr() {
		pv := goval.Addr()
		if method := pv.MethodByName(key); method.IsValid() {
			return makeMember(v, NewReflectedGoValue(method), c)
		}
	} else if goval.Kind() == reflect.Ptr {
		ev := goval.Elem()
		if method := ev.MethodByName(key); method.IsValid() {
			return makeMember(v, NewReflectedGoValue(method), c)
		}
	}
	return getExtMember(v, key, c)
}

func (v GoFunc) GetName() string {
	return v.v.Type().Name()
}

func (v GoFunc) GetArgNames(*Context) []string {
	return []string{}
}

func (v GoFunc) GetRefs() []string {
	return []string{}
}

func (v GoFunc) Invoke(c *Context, this Value, args []Value) {
	c.PushFuncStack(v.GetName())
	defer c.PopStack()
	method := v.GoValue.v.Type()
	numArgs := method.NumIn()
	if len(args) != numArgs {
		if numArgs == 0 {
			c.RaiseRuntimeError("invoke %s fail: argument num not match. required %d got %d", v.GetName(), numArgs, len(args))
			return
		} else if len(args) < numArgs-1 || !method.IsVariadic() {
			c.RaiseRuntimeError("invoke %s fail: argument num not match. required %d got %d", v.GetName(), numArgs, len(args))
			return
		}
	}
	goArgs := make([]reflect.Value, len(args))
	for i, a := range args {
		var inType reflect.Type
		if i < numArgs-1 || !method.IsVariadic() {
			inType = method.In(i)
		} else {
			inType = method.In(numArgs - 1).Elem()
		}
		if goVal, isGoVal := a.(GoValue); isGoVal {
			goArgs[i] = goVal.v
		} else if goVal, isGoVal := a.(GoFunc); isGoVal {
			goArgs[i] = goVal.GoValue.v
		} else if !Nullish(v) {
			goVal := reflect.New(inType)
			toGoValue(c, a, goVal.Elem())
			goArgs[i] = goVal.Elem()
		} else {
			goArgs[i] = reflect.Zero(inType)
		}
	}
	rv := v.GoValue.v.Call(goArgs)
	switch len(rv) {
	case 0:
		c.RetVal = constUndefined
	case 1:
		c.RetVal = NewReflectedGoValue(rv[0])
	default:
		{
			retArr := NewArray(len(rv))
			for _, rvItem := range rv {
				retArr.PushBack(NewReflectedGoValue(rvItem))
			}
			c.RetVal = retArr
		}
	}
}

type goType struct {
	*ValueBase
	typ reflect.Type
}

type GoType = *goType

func NewGoType(typ reflect.Type) GoType {
	return &goType{
		ValueBase: new(ValueBase),
		typ:       typ,
	}
}

func (GoType) Type() ValueType {
	return TypeGoType
}

func (t GoType) GoType() reflect.Type {
	return t.typ
}

func (t GoType) CompareTo(other Value, c *Context) CompareResult {
	if otherType, ok := other.(GoType); ok {
		if otherType.typ == t.typ {
			return CompareResultEqual
		}
	}
	return CompareResultNotEqual
}

func (t GoType) ToString(*Context) string {
	return t.typ.String()
}

func (t GoType) ToGoValue() interface{} {
	return t.typ
}

func (GoType) IsTrue() bool {
	return true
}

func (t GoType) GetMember(name string, c *Context) Value {
	switch name {
	case "slice":
		return NewGoType(reflect.SliceOf(t.typ))
	case "ptr":
		return NewGoType(reflect.PtrTo(t.typ))
	}
	return getExtMember(t, name, c)
}

func (GoType) GetIndex(int, *Context) Value {
	return constUndefined
}

func (t GoType) GetName() string {
	return t.typ.Name()
}

func (t GoType) GetArgNames(*Context) []string {
	return []string{}
}

func (t GoType) GetRefs() []string {
	return []string{}
}

func (t GoType) Invoke(c *Context, this Value, args []Value) {
	val := reflect.New(t.typ).Elem()
	if len(args) > 0 {
		if gv, ok := args[0].(GoValue); ok {
			val.Set(gv.v)
		} else {
			toGoValue(c, args[0], val)
		}
	}
	c.RetVal = NewReflectedGoValue(val)
}
