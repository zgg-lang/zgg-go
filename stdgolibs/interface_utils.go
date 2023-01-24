package stdgolibs

import (
	"reflect"

	. "github.com/zgg-lang/zgg-go/runtime"
)

type interfaceDelegate struct {
	c   *Context
	obj ValueObject
}

func interfaceMethodBridge(c *Context, obj ValueObject, receiver interface{}, method string, ins []interface{}, outs []interface{}) {
	var (
		v    = reflect.ValueOf(receiver)
		vt   = v.Type()
		m, _ = vt.MethodByName(method)
		args = make([]Value, len(ins))
	)
	for i, a := range ins {
		args[i] = NewGoValue(a)
	}
	c.InvokeMethod(obj, method, Args(args...))
	ret := c.RetVal
	var rets []Value
	switch m.Type.NumOut() {
	case 0:
		rets = []Value{}
	case 1:
		rets = []Value{ret}
	default:
		rets = *ret.(ValueArray).Values
	}
	if len(rets) != m.Type.NumOut() {
		c.RaiseRuntimeError("call method %s should return %d value(s), but got %d", m.Type.NumOut(), len(rets))
	}
	for i, ret := range rets {
		rv := reflect.New(m.Type.Out(i)).Elem()
		ToGoValue(c, ret, rv)
		reflect.ValueOf(outs[i]).Elem().Set(rv)
	}
}
