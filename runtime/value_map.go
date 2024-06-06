package runtime

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

type mapElem struct {
	key, value Value
}

type valueMap struct {
	*ValueBase
	rw sync.RWMutex
	ma map[int64]*list.List
}

type ValueMap = *valueMap

func NewMap() ValueMap {
	return &valueMap{
		ValueBase: &ValueBase{},
		ma:        make(map[int64]*list.List),
	}
}

func NewMapWithPairs(c *Context, kvs []Value) Value {
	n := len(kvs)
	if n%2 != 0 {
		c.RaiseRuntimeError("NewMap arguments number must be an even number")
	}
	rv := NewMap()
	for i := 0; i < n; i += 2 {
		rv.set(c, kvs[i], kvs[i+1])
	}
	return rv
}

func (m ValueMap) getHash(c *Context, key Value) (int64, bool) {
	if h, ok := key.(CanHash); ok {
		return h.Hash(), true
	}
	if h, ok := c.GetCallable(key.GetMember("hash", c)); ok {
		c.Invoke(h, nil, NoArgs)
		if r, ok := c.RetVal.(ValueInt); !ok {
			c.RaiseRuntimeError("hash method returns an non-integer value!")
		} else {
			return r.Value(), true
		}
	}
	return 0, false
}

func (m ValueMap) get(c *Context, key Value) (Value, bool) {
	hash, ok := m.getHash(c, key)
	if !ok {
		c.RaiseRuntimeError("key %s is not hashable!", key.ToString(c))
	}
	m.rw.RLock()
	defer m.rw.RUnlock()
	if l, found := m.ma[hash]; found {
		for p := l.Front(); p != nil; p = p.Next() {
			pe := p.Value.(mapElem)
			if c.ValuesEqual(key, pe.key) {
				return pe.value, true
			}
		}
	}
	return constUndefined, false
}

func (m ValueMap) set(c *Context, key Value, value Value) {
	hash, ok := m.getHash(c, key)
	if !ok {
		c.RaiseRuntimeError("key %s is not hashable!", key.ToString(c))
	}
	isDelete := IsUndefined(value)
	m.rw.Lock()
	defer m.rw.Unlock()
	if l, found := m.ma[hash]; found {
		for p := l.Front(); p != nil; p = p.Next() {
			pe := p.Value.(mapElem)
			if c.ValuesEqual(key, pe.key) {
				if isDelete {
					l.Remove(p)
				} else {
					pe.value = value
				}
				return
			}
		}
		if !isDelete {
			l.PushFront(mapElem{key: key, value: value})
		}
	} else if !isDelete {
		l = list.New()
		l.PushBack(mapElem{key: key, value: value})
		m.ma[hash] = l
	}
	return
}

func (m ValueMap) Each(handle func(key, value Value) bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for _, v := range m.ma {
		for p := v.Front(); p != nil; p = p.Next() {
			e := p.Value.(mapElem)
			if !handle(e.key, e.value) {
				return
			}
		}
	}
}

// Implements Value

func (m ValueMap) CompareTo(other Value, c *Context) CompareResult {
	return CompareResultNotEqual
}

func (m ValueMap) GetMember(name string, c *Context) Value {
	return getMemberByType(c, m, name)
}

func (m ValueMap) GetIndex(int, *Context) Value {
	return constUndefined
}

func (m ValueMap) Type() ValueType {
	return TypeMap
}

func (m ValueMap) GoType() reflect.Type {
	var vv map[interface{}]interface{}
	return reflect.TypeOf(vv)
}

func (m ValueMap) ToGoValue(c *Context) interface{} {
	vv := make(map[interface{}]interface{}, len(m.ma))
	m.rw.RLock()
	defer m.rw.RUnlock()
	for _, v := range m.ma {
		for p := v.Front(); p != nil; p = p.Next() {
			e := p.Value.(mapElem)
			vv[e.key.ToGoValue(c)] = e.value.ToGoValue(c)
		}
	}
	return vv
}

func (m ValueMap) ToString(c *Context) string {
	return fmt.Sprint(m.ToGoValue(c))
}

func (m ValueMap) IsTrue() bool {
	return len(m.ma) > 0
}

// Implements CanLen
func (m ValueMap) Len() int { return len(m.ma) }

func (m ValueMap) Contains(c *Context, v Value) bool {
	_, found := m.get(c, v)
	return found
}

var builtinMapMethods = map[string]ValueCallable{
	"get": NewNativeFunction("get", func(c *Context, this Value, args []Value) Value {
		var key Value
		EnsureFuncParams(c, "Map.get", args, ArgRuleRequired("key", TypeAny, &key))
		m := this.(ValueMap)
		if rv, found := m.get(c, key); found {
			return rv
		} else {
			return constUndefined
		}
	}, "key"),
	"put": NewNativeFunction("put", func(c *Context, this Value, args []Value) Value {
		var key, value Value
		EnsureFuncParams(c, "Map.put", args,
			ArgRuleRequired("key", TypeAny, &key),
			ArgRuleRequired("value", TypeAny, &value),
		)
		m := this.(ValueMap)
		m.set(c, key, value)
		return this
	}, "key", "value"),
	"keys": NewNativeFunction("keys", func(c *Context, this Value, args []Value) Value {
		m := this.(ValueMap)
		rv := NewArray(m.Len())
		m.Each(func(key, _ Value) bool {
			rv.PushBack(key)
			return true
		})
		return rv
	}),
	"values": NewNativeFunction("values", func(c *Context, this Value, args []Value) Value {
		m := this.(ValueMap)
		rv := NewArray(m.Len())
		m.Each(func(_, value Value) bool {
			rv.PushBack(value)
			return true
		})
		return rv
	}),
	"pairs": NewNativeFunction("pairs", func(c *Context, this Value, args []Value) Value {
		m := this.(ValueMap)
		rv := NewArray(m.Len())
		m.Each(func(key, value Value) bool {
			p := NewObject()
			p.SetMember("key", key, c)
			p.SetMember("value", value, c)
			rv.PushBack(p)
			return true
		})
		return rv
	}),
}

func init() {
	addMembersAndStatics(TypeMap, builtinMapMethods)
}

// check builtin hashable types

var (
	_ CanHash = NewInt(0)
	_ CanHash = NewFloat(0)
	_ CanHash = NewStr("")
	_ CanHash = NewBool(false)
	_ CanHash = TypeStr
)
