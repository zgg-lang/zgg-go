package runtime

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

type valueSet struct {
	*ValueBase
	rw sync.RWMutex
	ma map[int64]*list.List
}

type ValueSet = *valueSet

func NewSet() ValueMap {
	return &valueMap{
		ValueBase: &ValueBase{},
		ma:        make(map[int64]*list.List),
	}
}

func (m ValueSet) getHash(c *Context, key Value) (int64, bool) {
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

func (m ValueSet) get(c *Context, key Value) (Value, bool) {
	hash, ok := m.getHash(c, key)
	if !ok {
		c.RaiseRuntimeError("key %s is not hashable!", key.ToString(c))
	}
	m.rw.RLock()
	defer m.rw.RUnlock()
	if l, found := m.ma[hash]; found {
		for p := l.Front(); p != nil; p = p.Next() {
			pe := p.Value.(Value)
			if c.ValuesEqual(key, pe) {
				return key, true
			}
		}
	}
	return constUndefined, false
}

func (m ValueSet) set(c *Context, key Value, isDelete bool) {
	hash, ok := m.getHash(c, key)
	if !ok {
		c.RaiseRuntimeError("key %s is not hashable!", key.ToString(c))
	}
	m.rw.Lock()
	defer m.rw.Unlock()
	if l, found := m.ma[hash]; found {
		for p := l.Front(); p != nil; p = p.Next() {
			pe := p.Value.(Value)
			if c.ValuesEqual(key, pe) {
				if isDelete {
					l.Remove(p)
				}
				return
			}
		}
		if !isDelete {
			l.PushFront(key)
		}
	} else if !isDelete {
		l = list.New()
		l.PushBack(key)
		m.ma[hash] = l
	}
	return
}

func (m ValueSet) Each(handle func(key, value Value) bool) {
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

func (m ValueSet) CompareTo(other Value, c *Context) CompareResult {
	return CompareResultNotEqual
}

func (m ValueSet) GetMember(name string, c *Context) Value {
	return getMemberByType(c, m, name)
}

func (m ValueSet) GetIndex(int, *Context) Value {
	return constUndefined
}

func (m ValueSet) Type() ValueType {
	return TypeMap
}

func (m ValueSet) GoType() reflect.Type {
	var vv map[interface{}]interface{}
	return reflect.TypeOf(vv)
}

func (m ValueSet) ToGoValue() interface{} {
	vv := make(map[interface{}]interface{}, len(m.ma))
	m.rw.RLock()
	defer m.rw.RUnlock()
	for _, v := range m.ma {
		for p := v.Front(); p != nil; p = p.Next() {
			e := p.Value.(mapElem)
			vv[e.key.ToGoValue()] = e.key.ToGoValue()
		}
	}
	return vv
}

func (m ValueSet) ToString(c *Context) string {
	return fmt.Sprint(m.ToGoValue())
}

func (m ValueSet) IsTrue() bool {
	return len(m.ma) > 0
}

// Implements CanLen
func (m ValueSet) Len() int { return len(m.ma) }

var builtinSetMethods = map[string]ValueCallable{
	"load": NewNativeFunction("load", func(c *Context, this Value, args []Value) Value {
		var key Value
		EnsureFuncParams(c, "Map.load", args, ArgRuleRequired("key", TypeAny, &key))
		m := this.(ValueMap)
		if rv, found := m.get(c, key); found {
			return rv
		} else {
			return constUndefined
		}
	}, "key"),
	"store": NewNativeFunction("store", func(c *Context, this Value, args []Value) Value {
		var key, value Value
		EnsureFuncParams(c, "Map.store", args,
			ArgRuleRequired("key", TypeAny, &key),
			ArgRuleRequired("value", TypeAny, &value),
		)
		m := this.(ValueMap)
		m.set(c, key, value)
		return constUndefined
	}, "key", "value"),
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
