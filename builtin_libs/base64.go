package builtin_libs

import (
	"encoding/base64"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libBase64(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("encode", NewNativeFunction("base64.encode", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.RaiseRuntimeError("base64.encode: invalid arugments num")
			return nil
		}
		var bs []byte
		switch v := args[0].(type) {
		case ValueBytes:
			bs = v.Value()
		case ValueStr:
			bs = []byte(v.Value())
		default:
			c.RaiseRuntimeError("base64.encode: argument requires str or bytes, but got %s", args[0].Type().Name)
			return nil
		}
		return NewStr(base64.StdEncoding.EncodeToString(bs))
	}), c)
	lib.SetMember("decode", NewNativeFunction("base64.decode", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.RaiseRuntimeError("base64.decode: invalid arugments num")
			return nil
		}
		bs, err := base64.StdEncoding.DecodeString(args[0].ToString(c))
		if err != nil {
			c.RaiseRuntimeError("base64.decode: decode failed %s", err)
			return nil
		}
		return NewBytes(bs)
	}), c)
	lib.SetMember("encodeN", NewNativeFunction("base64.encodeN", func(c *Context, this Value, args []Value) Value {
		if len(args) != 2 {
			c.RaiseRuntimeError("base64.encodeN: invalid arugments num")
			return nil
		}
		var bs []byte
		switch v := args[0].(type) {
		case ValueBytes:
			bs = v.Value()
		case ValueStr:
			bs = []byte(v.Value())
		default:
			c.RaiseRuntimeError("base64.encode: argument requires str or bytes, but got %s", args[0].Type().Name)
			return nil
		}
		encoding := base64.NewEncoding(args[1].ToString(c))
		return NewStr(encoding.EncodeToString(bs))
	}), c)
	lib.SetMember("decodeN", NewNativeFunction("base64.decodeN", func(c *Context, this Value, args []Value) Value {
		if len(args) != 2 {
			c.RaiseRuntimeError("base64.decodeN: invalid arugments num")
			return nil
		}
		encoding := base64.NewEncoding(args[1].ToString(c))
		bs, err := encoding.DecodeString(args[0].ToString(c))
		if err != nil {
			c.RaiseRuntimeError("base64.decode: decode failed %s", err)
			return nil
		}
		return NewBytes(bs)
	}), c)
	return lib
}
