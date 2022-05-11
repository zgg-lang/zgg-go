package builtin_libs

import (
	"reflect"

	. "github.com/zgg-lang/zgg-go/runtime"
)

type DemoStruct struct {
	Name string
	Val  int
}

func libGo(c *Context) ValueObject {
	lib := NewObject()
	// lib.SetMember("", NewNativeFunction("go.ptrOf", func(c *Context, this Value, args []Value) Value {
	// 	return args[0].GetMember("ptr", c)
	// }), c)
	lib.SetMember("DemoStruct", NewGoType(reflect.TypeOf(DemoStruct{})), c)
	lib.SetMember("makeSlice", NewNativeFunction("go.makeSlice", func(c *Context, this Value, args []Value) Value {
		var (
			typ      GoType
			initLen  ValueInt
			capacity ValueInt
		)
		switch len(args) {
		case 1:
			initLen = NewInt(0)
			capacity = NewInt(0)
			EnsureFuncParams(c, "go.makeSlice", args,
				ArgRuleRequired{"elemType", TypeGoType, &typ},
			)
		case 2:
			EnsureFuncParams(c, "go.makeSlice", args,
				ArgRuleRequired{"elemType", TypeGoType, &typ},
				ArgRuleRequired{"initLen", TypeInt, &initLen},
			)
			capacity = initLen
		case 3:
			EnsureFuncParams(c, "go.makeSlice", args,
				ArgRuleRequired{"elemType", TypeGoType, &typ},
				ArgRuleRequired{"initLen", TypeInt, &initLen},
				ArgRuleRequired{"capacity", TypeInt, &capacity},
			)
		}
		return NewReflectedGoValue(reflect.MakeSlice(reflect.SliceOf(typ.GoType()), initLen.AsInt(), capacity.AsInt()))
	}), c)
	lib.SetMember("array", NewNativeFunction("go.array", func(c *Context, this Value, args []Value) Value {
		var (
			typ  GoType
			size ValueInt
		)
		EnsureFuncParams(c, "go.array", args,
			ArgRuleRequired{"elemType", TypeGoType, &typ},
			ArgRuleRequired{"size", TypeInt, &size},
		)
		return NewGoType(reflect.ArrayOf(size.AsInt(), typ.GoType()))
	}), c)
	lib.SetMember("append", NewNativeFunction("go.append", func(c *Context, this Value, args []Value) Value {
		if len(args) < 2 {
			c.RaiseRuntimeError("go.append requires at least 2 arguments")
		}
		slice := args[0].(GoValue).ReflectedValue()
		toAppend := make([]reflect.Value, len(args)-1)
		for i := range toAppend {
			toAppend[i] = MakeGoValueByArg(c, args[i+1], slice.Type().Elem())
		}
		return NewReflectedGoValue(reflect.Append(slice, toAppend...))
	}), c)
	lib.SetMember("int8", NewGoType(reflect.TypeOf(int8(0))), c)
	lib.SetMember("int16", NewGoType(reflect.TypeOf(int16(0))), c)
	lib.SetMember("int32", NewGoType(reflect.TypeOf(int32(0))), c)
	lib.SetMember("int64", NewGoType(reflect.TypeOf(int64(0))), c)
	lib.SetMember("int", NewGoType(reflect.TypeOf(int(0))), c)
	lib.SetMember("uint8", NewGoType(reflect.TypeOf(uint8(0))), c)
	lib.SetMember("uint16", NewGoType(reflect.TypeOf(uint16(0))), c)
	lib.SetMember("uint32", NewGoType(reflect.TypeOf(uint32(0))), c)
	lib.SetMember("uint64", NewGoType(reflect.TypeOf(uint64(0))), c)
	lib.SetMember("uint", NewGoType(reflect.TypeOf(uint(0))), c)
	lib.SetMember("float32", NewGoType(reflect.TypeOf(float32(0))), c)
	lib.SetMember("float64", NewGoType(reflect.TypeOf(float64(0))), c)
	lib.SetMember("string", NewGoType(reflect.TypeOf("")), c)
	lib.SetMember("byte", NewGoType(reflect.TypeOf(byte(0))), c)
	lib.SetMember("rune", NewGoType(reflect.TypeOf(rune(0))), c)
	lib.SetMember("bool", NewGoType(reflect.TypeOf(false)), c)
	lib.SetMember("map", NewNativeFunction("go.map", func(c *Context, this Value, args []Value) Value {
		var (
			keyType GoType
			valType GoType
		)
		EnsureFuncParams(c, "go.map", args,
			ArgRuleRequired{"keyType", TypeGoType, &keyType},
			ArgRuleRequired{"valueType", TypeGoType, &valType},
		)
		return NewGoType(reflect.MapOf(keyType.GoType(), valType.GoType()))
	}), c)
	lib.SetMember("new", NewNativeFunction("go.new", func(c *Context, this Value, args []Value) Value {
		var typ GoType
		EnsureFuncParams(c, "go.new", args,
			ArgRuleRequired{"type", TypeGoType, &typ},
		)
		return NewReflectedGoValue(reflect.New(typ.GoType()))
	}), c)
	lib.SetMember("makeMap", NewNativeFunction("go.makeMap", func(c *Context, this Value, args []Value) Value {
		var mapType reflect.Type
		if len(args) == 1 {
			var mapGoType GoType
			EnsureFuncParams(c, "go.makeMap", args,
				ArgRuleRequired{"mapType", TypeGoType, &mapGoType},
			)
			mapType = mapGoType.GoType()
			if mapType.Kind() != reflect.Map {
				c.RaiseRuntimeError("makeMap: arg mapType is not map")
			}
		} else {
			var (
				keyType GoType
				valType GoType
			)
			EnsureFuncParams(c, "go.makeMap", args,
				ArgRuleRequired{"keyType", TypeGoType, &keyType},
				ArgRuleRequired{"valueType", TypeGoType, &valType},
			)
			mapType = reflect.MapOf(keyType.GoType(), valType.GoType())
		}
		return NewReflectedGoValue(reflect.MakeMap(mapType))
	}), c)
	makeChan := func(name string, dir reflect.ChanDir) *ValueBuiltinFunction {
		return NewNativeFunction(name, func(c *Context, this Value, args []Value) Value {
			var (
				elemType GoType
				size     ValueInt
			)
			switch len(args) {
			case 1:
				EnsureFuncParams(c, name, args, ArgRuleRequired{"elemType", TypeGoType, &elemType})
				size = NewInt(0)
			default:
				EnsureFuncParams(c, name, args, ArgRuleRequired{"elemType", TypeGoType, &elemType}, ArgRuleRequired{"size", TypeInt, &size})
			}
			chanType := reflect.ChanOf(dir, elemType.GoType())
			return NewReflectedGoValue(reflect.MakeChan(chanType, size.AsInt()))
		})
	}
	lib.SetMember("makeChan", makeChan("go.makeChan", reflect.BothDir), c)
	lib.SetMember("makeSendChan", makeChan("go.makeSendChan", reflect.SendDir), c)
	lib.SetMember("makeRecvChan", makeChan("go.makeRecvChan", reflect.RecvDir), c)
	lib.SetMember("type", NewNativeFunction("go.type", func(c *Context, this Value, args []Value) Value {
		var (
			t  GoValue
			tt reflect.Type
			ok bool
		)
		EnsureFuncParams(c, "go.type", args,
			ArgRuleRequired{"t", TypeGoValue, &t},
		)
		if tt, ok = t.ToGoValue().(reflect.Type); !ok {
			c.RaiseRuntimeError("go.type require reflect.Type")
		}
		return NewGoType(tt)
	}), c)
	lib.SetMember("convert", NewNativeFunction("go.convert", func(c *Context, this Value, args []Value) Value {
		var (
			src        Value
			targetType GoType
		)
		EnsureFuncParams(c, "go.convert", args,
			ArgRuleRequired{"src", TypeAny, &src},
			ArgRuleRequired{"targetType", TypeGoType, &targetType},
		)
		var srcGo GoValue
		switch srcVal := src.(type) {
		case GoValue:
			srcGo = srcVal
		case Value:
			srcGo = NewGoValue(srcVal.ToGoValue()).(GoValue)
		}
		return NewGoValue(srcGo.ReflectedValue().Convert(targetType.GoType()))
	}, "src", "targetType"), c)
	return lib
}
