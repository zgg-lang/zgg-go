package runtime

import (
	"reflect"
	"strings"
	"sync"
)

type ValueObject struct {
	*ValueBase
	t        ValueType
	this     *ValueObject
	m        *sync.Map
	size     *int
	sizeLock *sync.Mutex
}

func NewObject(objType ...ValueType) ValueObject {
	fields := new(sync.Map)
	t := TypeObject
	if len(objType) > 0 {
		t = objType[0]
	}
	size := new(int64)
	*size = 0
	return ValueObject{
		ValueBase: &ValueBase{},
		t:         t,
		m:         fields,
		size:      new(int),
		sizeLock:  &sync.Mutex{},
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

func (v *ValueObject) Init(c *Context, args []Value) {
	if initFn := v.t.getInitFunc(c); initFn != nil {
		c.Invoke(initFn, *v, func() []Value { return args })
	}
}

func (v *ValueObject) Super(fromType ValueType) ValueObject {
	this := v
	if v.this != nil {
		this = v.this
	}
	return ValueObject{
		ValueBase: &ValueBase{},
		t:         fromType.Super(),
		this:      this,
		m:         v.m,
	}
}

func (v ValueObject) Each(f func(string, Value) bool) {
	v.m.Range(func(k interface{}, v interface{}) bool {
		return f(k.(string), v.(Value))
	})
}

func (v ValueObject) GoType() reflect.Type {
	var vv map[string]interface{}
	return reflect.TypeOf(vv)
}

func (v ValueObject) ToGoValue() interface{} {
	rv := map[string]interface{}{}
	v.Iterate(func(key string, value Value) {
		rv[key] = value.ToGoValue()
	})
	return rv
}

func (v ValueObject) GetIndex(index int, c *Context) Value {
	getItem := v.GetMember("__getItem__", c)
	if getItemFunc, callable := c.GetCallable(getItem); callable {
		getItemFunc.Invoke(c, v, []Value{NewInt(int64(index))})
		return c.RetVal
	}
	return constUndefined
}

func (v ValueObject) SetMember(name string, value Value, c *Context) {
	_, isUndefined := value.(ValueUndefined)
	v.sizeLock.Lock()
	defer v.sizeLock.Unlock()
	if _, found := v.m.Load(name); found {
		if isUndefined {
			(*v.size)--
			v.m.Delete(name)
		} else {
			v.m.Store(name, value)
		}
	} else {
		if !isUndefined {
			(*v.size)++
			v.m.Store(name, value)
		}

	}
}

func (v ValueObject) GetMember(name string, c *Context) Value {
	if val, found := v.m.Load(name); found {
		return makeMember(v, val.(Value), c)
	}
	return getMemberByType(c, v, name)
}

func (v ValueObject) Len() int {
	return *v.size
}

func (v ValueObject) IsTrue() (isTrue bool) {
	return *v.size != 0
}

func (v ValueObject) Type() ValueType {
	// if name, isStr := v["__name__"].(ValueStr); isStr {
	// 	return fmt.Sprintf("<class %s>", name.Value())
	// }
	// if proto, isObj := v["__proto__"].(ValueObject); isObj {
	// 	if name, isStr := proto["__name__"].(ValueStr); isStr {
	// 		return fmt.Sprintf("<object %s>", name.Value())
	// 	}
	// }
	return v.t
}

func (v ValueObject) CompareTo(other Value, c *Context) CompareResult {
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

func (v ValueObject) ToString(c *Context) string {
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

func (v ValueObject) Iterate(f func(key string, value Value)) {
	v.Each(func(k string, v Value) bool {
		f(k, v)
		return true
	})
}

func (v ValueObject) GetName() string {
	return ""
}

func (ValueObject) GetArgNames(*Context) []string {
	return []string{}
}

func (ValueObject) GetRefs() []string {
	return []string{}
}

func (v ValueObject) Invoke(c *Context, this Value, args []Value) {
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
}

func init() {
	addMembersAndStatics(TypeObject, builtinObjMethods)
}
