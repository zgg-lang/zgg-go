package builtin_libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strings"

	. "github.com/zgg-lang/zgg-go/runtime"
	"gopkg.in/yaml.v2"

	"github.com/oliveagle/jsonpath"
)

func jsonMarshal(v interface{}) ([]byte, error) {
	var ob bytes.Buffer
	enc := json.NewEncoder(&ob)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return ob.Bytes(), nil
}

func jsonMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	var ob bytes.Buffer
	enc := json.NewEncoder(&ob)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, indent)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return ob.Bytes(), nil
}

func libJson(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("encode", NewNativeFunction("json.encode", func(c *Context, this Value, args []Value) Value {
		var bs []byte
		var err error
		switch len(args) {
		case 1:
			bs, err = jsonMarshal(args[0].ToGoValue(c))
		case 2:
			var indent string
			switch iv := args[1].(type) {
			case ValueInt:
				indent = strings.Repeat(" ", iv.AsInt())
			case ValueStr:
				indent = iv.Value()
			default:
				c.RaiseRuntimeError("json.encode(value, indent): indent must be a string or an integer, not %s", args[1].Type().Name)
				return nil
			}
			bs, err = jsonMarshalIndent(args[0].ToGoValue(c), "", indent)
		case 3:
			bs, err = jsonMarshalIndent(args[0].ToGoValue(c), args[1].ToString(c), args[2].ToString(c))
		default:
			c.RaiseRuntimeError("json.encode: requires 1 to 3 argument(s)")
			return nil
		}
		if err != nil {
			c.RaiseRuntimeError("json.encode: " + err.Error())
			return nil
		}
		return NewStr(string(bs))
	}), nil)
	lib.SetMember("decode", NewNativeFunction("json.decode", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.RaiseRuntimeError("json.decode: requires 1 argument")
			return nil
		}
		var bs []byte
		switch arg := args[0].(type) {
		case ValueStr:
			bs = []byte(arg.Value())
		case ValueBytes:
			bs = arg.Value()
		default:
			c.RaiseRuntimeError("json.decode: argument must be a string or a bytes")
			return nil
		}
		var j interface{}
		if err := json.Unmarshal(bs, &j); err != nil {
			c.RaiseRuntimeError("json.decode: " + err.Error())
			return nil
		}
		return jsonToValue(j, c)
	}), nil)
	lib.SetMember("format", NewNativeFunction("json.format", func(c *Context, this Value, args []Value) Value {
		var (
			jsonStr    Value
			indentSize ValueInt
			indentStr  ValueStr
			indentType int
			bs         []byte
			isCmd      bool
		)
		if len(args) == 0 {
			bs, _ = io.ReadAll(os.Stdin)
			indentSize = NewInt(4)
			indentType = 0
			isCmd = true
		} else {
			EnsureFuncParams(c, "json.format", args,
				ArgRuleRequired("value", TypeAny, &jsonStr),
				ArgRuleOneOf(
					"indent",
					[]ValueType{TypeInt, TypeStr},
					[]interface{}{&indentSize, &indentStr},
					&indentType, &indentSize, NewInt(0)),
			)
			switch s := jsonStr.(type) {
			case ValueBytes:
				bs = s.Value()
			default:
				bs = []byte(s.ToString(c))
			}
		}
		var j interface{}
		if err := json.Unmarshal(bs, &j); err != nil {
			c.RaiseRuntimeError("json.format: decode failed %s", err)
		}
		var (
			outs []byte
			err  error
		)
		switch indentType {
		case 0:
			outs, err = jsonMarshalIndent(j, "", strings.Repeat(" ", indentSize.AsInt()))
		case 1:
			outs, err = jsonMarshalIndent(j, "", indentStr.Value())
		default:
			outs, err = jsonMarshal(j)
		}
		if err != nil {
			c.RaiseRuntimeError("json.format: marshal failed %s", err)
		}
		if isCmd {
			c.Stdout.Write(outs)
		}
		return NewStr(string(outs))
	}), nil)
	lib.SetMember("find", NewNativeFunction("json.find", func(c *Context, this Value, args []Value) Value {
		var (
			path  ValueStr
			value Value
			data  interface{}
			err   error
		)
		EnsureFuncParams(c, "json.find", args,
			ArgRuleRequired("jsonpath", TypeStr, &path),
			ArgRuleRequired("object", TypeAny, &value),
		)
		switch v := value.(type) {
		case ValueStr:
			err = json.Unmarshal([]byte(v.Value()), &data)
		case ValueBytes:
			err = json.Unmarshal(v.Value(), &data)
		default:
			data = v.ToGoValue(c)
		}
		if err != nil {
			c.RaiseRuntimeError("json.find: parse value error %s", err)
			return nil
		}
		res, err := jsonpath.JsonPathLookup(data, path.Value())
		if err != nil {
			c.RaiseRuntimeError("json.find: jsonpath %s lookup error: %s", path.Value(), err)
			return nil
		}
		return FromGoValue(reflect.ValueOf(res), c)
	}), nil)
	lib.SetMember("e2", func() Value {
		r := NewObject()
		r.SetMember("__rbitOr__", NewNativeFunction("e2.__rbitOr__", func(c *Context, this Value, args []Value) Value {
			var left Value
			EnsureFuncParams(c, "e2.__rbitOr__", args, ArgRuleRequired("left", TypeAny, &left))
			return c.InvokeMethod(lib, "encode", Args(left, NewInt(2)))
		}), nil)
		return r
	}(), nil)
	lib.SetMember("e4", func() Value {
		r := NewObject()
		r.SetMember("__rbitOr__", NewNativeFunction("e4.__rbitOr__", func(c *Context, this Value, args []Value) Value {
			var left Value
			EnsureFuncParams(c, "e4.__rbitOr__", args, ArgRuleRequired("left", TypeAny, &left))
			return c.InvokeMethod(lib, "encode", Args(left, NewInt(4)))
		}), nil)
		return r
	}(), nil)
	lib.SetMember("d", func() Value {
		r := NewObject()
		r.SetMember("__rbitOr__", NewNativeFunction("d.__rbitOr__", func(c *Context, this Value, args []Value) Value {
			return c.InvokeMethod(lib, "decode", Args(args...))
		}), nil)
		return r
	}(), nil)
	lib.SetMember("f2", func() Value {
		r := NewObject()
		r.SetMember("__rbitOr__", NewNativeFunction("f2.__rbitOr__", func(c *Context, this Value, args []Value) Value {
			var left Value
			EnsureFuncParams(c, "f2.__rbitOr__", args, ArgRuleRequired("left", TypeAny, &left))
			return c.InvokeMethod(lib, "format", Args(left, NewInt(2)))
		}), nil)
		return r
	}(), nil)
	lib.SetMember("f4", func() Value {
		r := NewObject()
		r.SetMember("__rbitOr__", NewNativeFunction("f4.__rbitOr__", func(c *Context, this Value, args []Value) Value {
			var left Value
			EnsureFuncParams(c, "f4.__rbitOr__", args, ArgRuleRequired("left", TypeAny, &left))
			return c.InvokeMethod(lib, "format", Args(left, NewInt(4)))
		}), nil)
		return r
	}(), nil)
	return lib
}

func jsonToValue(src interface{}, c *Context) (_r Value) {
	if src == nil {
		return Nil()
	}
	switch srcVal := src.(type) {
	case float64:
		if math.Floor(srcVal) == srcVal {
			return NewInt(int64(srcVal))
		}
		return NewFloat(srcVal)
	case int64:
		return NewInt(srcVal)
	case int:
		return NewInt(int64(srcVal))
	case int8:
		return NewInt(int64(srcVal))
	case int16:
		return NewInt(int64(srcVal))
	case int32:
		return NewInt(int64(srcVal))
	case uint:
		return NewInt(int64(srcVal))
	case uint8:
		return NewInt(int64(srcVal))
	case uint16:
		return NewInt(int64(srcVal))
	case uint32:
		return NewInt(int64(srcVal))
	case uint64:
		return NewInt(int64(srcVal))
	case bool:
		return NewBool(srcVal)
	case string:
		return NewStr(srcVal)
	case []interface{}:
		{
			rv := NewArray(len(srcVal))
			for _, elemVal := range srcVal {
				rv.PushBack(jsonToValue(elemVal, c))
			}
			return rv
		}
	case map[string]interface{}:
		{
			rv := NewObject()
			for k, elemVal := range srcVal {
				rv.SetMember(k, jsonToValue(elemVal, c), c)
			}
			return rv
		}
	case map[interface{}]interface{}:
		{
			rv := NewObject()
			for k, elemVal := range srcVal {
				c.DebugLog("-- key %v value %v", fmt.Sprint(k), elemVal)
				rv.SetMember(fmt.Sprint(k), jsonToValue(elemVal, c), c)
			}
			return rv
		}
	}
	return Undefined()
}

func libYaml(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("encode", NewNativeFunction("yaml.encode", func(c *Context, this Value, args []Value) Value {
		var val Value
		EnsureFuncParams(c, "yaml.encode", args, ArgRuleRequired("value", TypeAny, &val))
		bs, err := yaml.Marshal(val.ToGoValue(c))
		if err != nil {
			c.RaiseRuntimeError("yaml.encode error %v", err)
		}
		return NewStr(string(bs))
	}, "value"), nil)
	lib.SetMember("decode", NewNativeFunction("yaml.decode", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.RaiseRuntimeError("yaml.decode: requires 1 argument")
			return nil
		}
		var bs []byte
		switch arg := args[0].(type) {
		case ValueStr:
			bs = []byte(arg.Value())
		case ValueBytes:
			bs = arg.Value()
		default:
			c.RaiseRuntimeError("yaml.decode: argument must be a string or a bytes")
			return nil
		}
		var j interface{}
		if err := yaml.Unmarshal(bs, &j); err != nil {
			c.RaiseRuntimeError("yaml.decode error %v", err)
			return nil
		}
		rv := jsonToValue(j, c)
		return rv
	}, "yaml"), nil)
	return lib
}
