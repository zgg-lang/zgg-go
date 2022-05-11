package builtin_libs

import (
	"reflect"

	. "github.com/zgg-lang/zgg-go/runtime"

	"github.com/vmihailenco/msgpack"
)

func libMsgpack(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("encode", NewNativeFunction(
		"msgpack.encode",
		func(c *Context, this Value, args []Value) Value {
			if len(args) != 1 {
				c.RaiseRuntimeError("msgpack.encode: requires 1 argument")
				return nil
			}
			bs, err := msgpack.Marshal(args[0].ToGoValue())
			if err != nil {
				c.RaiseRuntimeError("msgpack.encode: " + err.Error())
				return nil
			}
			return NewBytes(bs)
		},
	), nil)
	lib.SetMember("decode", NewNativeFunction(
		"msgpack.decode",
		func(c *Context, this Value, args []Value) Value {
			if len(args) != 1 {
				c.RaiseRuntimeError("msgpack.decode: requires 1 argument")
				return nil
			}
			var bs []byte
			switch arg := args[0].(type) {
			case ValueStr:
				bs = []byte(arg.Value())
			case ValueBytes:
				bs = arg.Value()
			default:
				c.RaiseRuntimeError("msgpack.decode: argument must be a string or a bytes")
				return nil
			}
			var j interface{}
			if err := msgpack.Unmarshal(bs, &j); err != nil {
				c.RaiseRuntimeError("msgpack.decode: " + err.Error())
				return nil
			}
			return FromGoValue(reflect.ValueOf(j), c)
		},
	), nil)
	return lib
}
