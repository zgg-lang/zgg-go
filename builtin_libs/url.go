package builtin_libs

import (
	"net/url"
	"reflect"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libUrl(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("encodeURIComponent", NewNativeFunction("encodeURIComponent", func(c *Context, this Value, args []Value) Value {
		var value Value
		EnsureFuncParams(c, "decodeURIComponent", args, ArgRuleRequired("value", TypeAny, &value))
		return NewStr(url.QueryEscape(value.ToString(c)))
	}, "value"), c)
	lib.SetMember("decodeURIComponent", NewNativeFunction("decodeURIComponent", func(c *Context, this Value, args []Value) Value {
		var encoded ValueStr
		EnsureFuncParams(c, "decodeURIComponent", args, ArgRuleRequired("encoded", TypeStr, &encoded))
		ev := encoded.Value()
		r, err := url.QueryUnescape(ev)
		if err != nil {
			c.RaiseRuntimeError("decodeURIComponent: decode error! input \"%s\", error: %s", ev, err)
		}
		return NewStr(r)
	}, "encoded"), c)
	lib.SetMember("encodeForm", NewNativeFunction("encodeForm", func(c *Context, this Value, args []Value) Value {
		var form = url.Values{}
		var formValue ValueObject
		EnsureFuncParams(c, "encodeForm", args, ArgRuleRequired("form", TypeObject, &formValue))
		formValue.Iterate(func(k string, v Value) {
			if varr, ok := v.(ValueArray); ok {
				n := varr.Len()
				for i := 0; i < n; i++ {
					form.Add(k, varr.GetIndex(i, c).ToString(c))
				}
			} else {
				form.Set(k, v.ToString(c))
			}
		})
		return NewStr(form.Encode())
	}), c)
	lib.SetMember("decodeForm", NewNativeFunction("decodeForm", func(c *Context, this Value, args []Value) Value {
		var formStr ValueStr
		EnsureFuncParams(c, "decodeForm", args, ArgRuleRequired("form", TypeStr, &formStr))
		form, err := url.ParseQuery(formStr.Value())
		if err != nil {
			c.RaiseRuntimeError("decodeForm parse error: %s", err)
		}
		return FromGoValue(reflect.ValueOf(form), c)
	}), c)
	return lib
}
