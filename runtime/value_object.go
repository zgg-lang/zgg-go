package runtime

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type valueObject struct {
	*ValueBase
	Reserved any
	t        *valueType
	this     *valueObject
	m        sync.Map
	size     int
	sizeLock sync.Mutex
}

type ValueObject = *valueObject

func NewObject(objType ...ValueType) ValueObject {
	t := TypeObject
	if len(objType) > 0 {
		t = objType[0]
	}
	return &valueObject{
		ValueBase: &ValueBase{},
		t:         t,
	}
}

func NewObjectAndInit(objType ValueType, c *Context, initArgs ...Value) ValueObject {
	if objType == nil {
		objType = TypeObject
	}
	rv := NewObject(objType)
	rv.Init(c, initArgs)
	return rv
}

func (v *valueObject) Init(c *Context, args []Value) {
	if initFn := v.t.getInitFunc(c); initFn != nil {
		c.Invoke(initFn, v, func() []Value { return args })
	}
}

func (v *valueObject) Super(fromType ValueType) ValueObject {
	this := v
	if v.this != nil {
		this = v.this
	}
	return &valueObject{
		ValueBase: &ValueBase{},
		t:         fromType.Super(),
		this:      this,
		m:         v.m,
	}
}

func (v *valueObject) Each(f func(string, Value) bool) {
	v.m.Range(func(k interface{}, v interface{}) bool {
		return f(k.(string), v.(Value))
	})
}

func (v *valueObject) GoType() reflect.Type {
	var vv map[string]interface{}
	return reflect.TypeOf(vv)
}

func (v *valueObject) ToGoValue() interface{} {
	rv := map[string]interface{}{}
	v.Iterate(func(key string, value Value) {
		rv[key] = value.ToGoValue()
	})
	return rv
}

func (v *valueObject) GetIndex(index int, c *Context) Value {
	getItem := v.GetMember("__getItem__", c)
	if getItemFunc, callable := c.GetCallable(getItem); callable {
		getItemFunc.Invoke(c, v, []Value{NewInt(int64(index))})
		return c.RetVal
	}
	return constUndefined
}

func (v *valueObject) SetMember(name string, value Value, c *Context) {
	_, isUndefined := value.(ValueUndefined)
	v.sizeLock.Lock()
	defer v.sizeLock.Unlock()
	if _, found := v.m.Load(name); found {
		if isUndefined {
			v.size--
			v.m.Delete(name)
		} else {
			v.m.Store(name, value)
		}
	} else {
		if !isUndefined {
			v.size++
			v.m.Store(name, value)
		}
	}
}

func (v *valueObject) GetMember(name string, c *Context) Value {
	if val, found := v.m.Load(name); found {
		return makeMember(v, val.(Value), c)
	}
	return getMemberByType(c, v, name)
}

func (v *valueObject) Len() int {
	return v.size
}

func (v *valueObject) IsTrue() (isTrue bool) {
	return v.size != 0
}

func (v *valueObject) Type() ValueType {
	return v.t
}

func (v *valueObject) CompareTo(other Value, c *Context) CompareResult {
	v2, isObj := other.(ValueObject)
	if !isObj {
		return CompareResultNotEqual
	}
	if v.Len() != v.Len() {
		return CompareResultNotEqual
	}
	rv := CompareResultEqual
	v.Each(func(k1 string, elem1 Value) bool {
		elem2, found := v2.m.Load(k1)
		if !found {
			rv = CompareResultNotEqual
			return false
		}
		if !c.ValuesEqual(elem1, elem2.(Value)) {
			rv = CompareResultNotEqual
			return false
		}
		return true
	})
	return rv
}

func (v *valueObject) ToString(c *Context) string {
	if strFn, ok := c.GetCallable(v.t.findMember("__str__")); ok {
		c.Invoke(strFn, v, func() []Value { return []Value{} })
		return c.RetVal.ToString(c)
	}
	var sb strings.Builder
	sb.WriteRune('{')
	i := 0
	v.Iterate(func(k string, v Value) {
		i++
		if i > 1 {
			sb.WriteString(", ")
		}
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v.ToString(c))
	})
	sb.WriteRune('}')
	return sb.String()
}

func (v *valueObject) Iterate(f func(key string, value Value)) {
	v.Each(func(k string, v Value) bool {
		f(k, v)
		return true
	})
}

func (v *valueObject) GetName() string {
	return ""
}

func (ValueObject) GetArgNames(*Context) []string {
	return []string{}
}

func (ValueObject) GetRefs() []string {
	return []string{}
}

func (v *valueObject) Invoke(c *Context, this Value, args []Value) {
	callMethod, ok := c.GetCallable(v.GetMember("__call__", c))
	if !ok {
		c.RaiseRuntimeError("invoked object is not callable")
		return
	}
	c.Invoke(callMethod, v, Args(args...))
}

var builtinObjMethods = map[string]ValueCallable{
	"keys": NewNativeFunction("object.keys", func(c *Context, thisArg Value, args []Value) Value {
		thisObj := thisArg.(ValueObject)
		rv := NewArray()
		thisObj.Iterate(func(k string, v Value) {
			rv.PushBack(NewStr(k))
		})
		return rv
	}),
	"values": &ValueBuiltinFunction{
		name: "object.values",
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisObj := thisArg.(ValueObject)
			rv := NewArray()
			thisObj.Iterate(func(k string, v Value) {
				rv.PushBack(v)
			})
			return rv
		},
	},
	"pairs": NewNativeFunction("object.pairs", func(c *Context, thisArg Value, args []Value) Value {
		thisObj := thisArg.(ValueObject)
		rv := NewArray()
		thisObj.Iterate(func(k string, v Value) {
			pair := NewObject()
			pair.SetMember("key", NewStr(k), c)
			pair.SetMember("value", v, c)
			rv.PushBack(pair)
		})
		return rv
	}),
	"each": NewNativeFunction("object.each", func(c *Context, this Value, args []Value) Value {
		o := c.MustObject(this)
		if len(args) != 1 {
			c.RaiseRuntimeError("object.each requires 1 argument")
		}
		handleFunc := c.MustCallable(args[0], "object.each handleFunc")
		o.Each(func(key string, value Value) bool {
			c.Invoke(handleFunc, o, func() []Value {
				return []Value{NewStr(key), value}
			})
			return true
		})
		return constUndefined
	}),
	"printReserved": NewNativeFunction("object.printReserved", func(c *Context, this Value, args []Value) Value {
		o := c.MustObject(this)
		fmt.Printf("%#v\n", o.Reserved)
		return constUndefined
	}),
}

func init() {
	addMembersAndStatics(TypeObject, builtinObjMethods)
}
