package runtime

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func CanBeNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	}
	return false
}

func FromGoValue(v reflect.Value, c *Context) Value {
	if !v.IsValid() || (CanBeNil(v) && v.IsNil()) {
		return constNil
	}
	vi := v.Interface()
	vt := v.Type()
	if vt.Kind() == reflect.Interface {
		v = v.Elem()
		vi = v.Interface()
		vt = v.Type()
	}
	switch vt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewInt(int64(v.Uint()))
	case reflect.Float32, reflect.Float64:
		return NewFloat(v.Float())
	case reflect.Bool:
		return NewBool(v.Bool())
	case reflect.String:
		return NewStr(v.String())
	case reflect.Slice, reflect.Array:
		if vt.Elem().Kind() == reflect.Uint8 {
			var bs []byte
			if vt.Kind() == reflect.Slice {
				bs = vi.([]byte)
			} else {
				bs = v.Slice(0, v.Len()).Interface().([]byte)
			}
			return NewBytes(bs)
		} else {
			rv := NewArray(v.Len())
			for i := 0; i < v.Len(); i++ {
				rv.PushBack(FromGoValue(v.Index(i), c))
			}
			return rv
		}
	case reflect.Struct:
		{
			rv := NewObject()
			for i := 0; i < vt.NumField(); i++ {
				f := vt.Field(i)
				if name := f.Name; name != "" && name[0] >= 'A' && name[0] <= 'Z' {
					rv.SetMember(name, FromGoValue(v.Field(i), c), c)
				}
			}
			return rv
		}
	case reflect.Map:
		{
			rv := NewObject()
			keys := v.MapKeys()
			for _, k := range keys {
				value := v.MapIndex(k)
				rv.SetMember(fmt.Sprint(k.Interface()), FromGoValue(value, c), c)
			}
			return rv
		}
	case reflect.Ptr:
		return FromGoValue(v.Elem(), c)
	}
	return NewReflectedGoValue(v)
}

func toGoValue(c *Context, v Value, goVal reflect.Value) {
	if Nullish(v) {
		switch goVal.Kind() {
		case reflect.Ptr:
			return
		case reflect.Slice:
			return
		case reflect.Map:
			return
		case reflect.Interface:
			return
		}
	}
	goTyp := goVal.Type()
	if v.GoType() == goTyp {
		goVal.Set(reflect.ValueOf(v.ToGoValue()))
		return
	}
	switch goTyp.Kind() {
	case reflect.Int:
		goVal.Set(reflect.ValueOf(int(c.MustInt(v))))
	case reflect.Int8:
		goVal.Set(reflect.ValueOf(int8(c.MustInt(v))))
	case reflect.Int16:
		goVal.Set(reflect.ValueOf(int16(c.MustInt(v))))
	case reflect.Int32:
		goVal.Set(reflect.ValueOf(int32(c.MustInt(v))))
	case reflect.Int64:
		goVal.Set(reflect.ValueOf(int64(c.MustInt(v))))
	case reflect.Uint:
		goVal.Set(reflect.ValueOf(uint(c.MustInt(v))))
	case reflect.Uint8:
		goVal.Set(reflect.ValueOf(uint8(c.MustInt(v))))
	case reflect.Uint16:
		goVal.Set(reflect.ValueOf(uint16(c.MustInt(v))))
	case reflect.Uint32:
		goVal.Set(reflect.ValueOf(uint32(c.MustInt(v))))
	case reflect.Uint64:
		goVal.Set(reflect.ValueOf(uint64(c.MustInt(v))))
	case reflect.Float32:
		goVal.Set(reflect.ValueOf(float32(c.MustFloat(v))))
	case reflect.Float64:
		goVal.Set(reflect.ValueOf(c.MustFloat(v)))
	case reflect.Bool:
		goVal.Set(reflect.ValueOf(c.MustBool(v)))
	case reflect.String:
		goVal.Set(reflect.ValueOf(c.MustStr(v)))
	case reflect.Slice:
		elType := goTyp.Elem()
		arr := c.MustArray(v)
		outVal := reflect.MakeSlice(reflect.SliceOf(elType), 0, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			el := reflect.New(elType)
			toGoValue(c, arr.GetIndex(i, c), el.Elem())
			outVal = reflect.Append(outVal, el.Elem())
		}
		goVal.Set(outVal)
	case reflect.Struct:
		nf := goTyp.NumField()
		for i := 0; i < nf; i++ {
			f := goTyp.Field(i)
			fieldVal := v.GetMember(f.Name, c)
			if IsUndefined(fieldVal) {
				continue
			}
			toGoValue(c, fieldVal, goVal.Field(i))
		}
	case reflect.Func:
		callable := c.MustCallable(v)
		goVal.Set(reflect.MakeFunc(goTyp, func(args []reflect.Value) []reflect.Value {
			zggArgs := make([]Value, len(args))
			for i, a := range args {
				zggArgs[i] = NewGoValue(a)
			}
			c.Invoke(callable, nil, Args(zggArgs...))
			outN := goTyp.NumOut()
			rv := make([]reflect.Value, outN)
			switch outN {
			case 0:
			case 1:
				rv[0] = reflect.Zero(goTyp.Out(0))
				toGoValue(c, c.RetVal, rv[0])
			default:
				if retArr, ok := c.RetVal.(ValueArray); ok {
					var i int
					for i = 0; i < retArr.Len(); i++ {
						rv[i] = reflect.Zero(goTyp.Out(i))
						toGoValue(c, retArr.GetIndex(i, c), rv[i])
					}
					for ; i < outN; i++ {
						rv[i] = reflect.Zero(goTyp.Out(i))
					}
				}
			}
			return rv
		}))
	default:
		{
			goIntr := goVal.Addr().Interface()
			jsonBs, err := json.Marshal(buildJson(v, c))
			if err != nil {
				c.RaiseRuntimeError("value %s to go value error: " + err.Error())
				return
			}
			if err := json.Unmarshal(jsonBs, goIntr); err != nil {
				c.RaiseRuntimeError("value %s, to go value %s error: %s", v.ToString(c), goVal.Type().String(), err.Error())
				return
			}
		}
	}
}

func WrapGoFunction(f interface{}) *ValueBuiltinFunction {
	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		panic("WrapGoFunction argument must be a function")
	}
	fTyp := fVal.Type()
	rv := &ValueBuiltinFunction{
		name: fTyp.Name(),
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != fTyp.NumIn() {
				c.RaiseRuntimeError(fmt.Sprintf("%s arguments requires %d get %d", fTyp.Name(), fTyp.NumIn(), len(args)))
				return nil
			}
			in := make([]reflect.Value, len(args))
			for i, arg := range args {
				in[i] = reflect.New(fTyp.In(i)).Elem()
				toGoValue(c, arg, in[i])
			}
			out := fVal.Call(in)
			if len(out) == 0 {
				return constUndefined
			} else if len(out) == 1 {
				if outVal, isVal := out[0].Interface().(Value); isVal {
					return outVal
				}
				return FromGoValue(out[0], c)
			} else {
				rv := NewArray(len(out))
				for _, o := range out {
					rv.PushBack(FromGoValue(o, c))
				}
				return rv
			}
		},
	}
	return rv
}
