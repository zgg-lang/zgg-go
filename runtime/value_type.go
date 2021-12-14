package runtime

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	valueType struct {
		*ValueBase
		Bases     []*valueType
		superInit sync.Once
		superType *valueType
		TypeId    int
		Name      string
		Members   *sync.Map
		Statics   *sync.Map
	}
	ValueType = *valueType
)

func NewType(id int, name string) ValueType {
	rv := &valueType{ValueBase: &ValueBase{}, TypeId: id, Name: name, Members: new(sync.Map), Statics: new(sync.Map)}
	rv.Statics.Store("__name__", NewStr(name))
	return rv
}

func (*valueType) Type() ValueType {
	return TypeType
}

func (t *valueType) Super() ValueType {
	if t.TypeId < 0 { // 已经是一个super类型
		return t
	}
	t.superInit.Do(func() {
		t.superType = &valueType{ValueBase: new(ValueBase), TypeId: -t.TypeId, Name: t.Name + ".super", Members: new(sync.Map), Statics: new(sync.Map), Bases: t.Bases}
		t.superType.Members.Store("__init__", NewNativeFunction("__init__", func(c *Context, this Value, args []Value) Value {
			getArgs := func() []Value { return args }
			thisObj, ok := this.(ValueObject)
			if ok {
				for _, base := range t.Bases {
					if initFn := base.getInitFunc(); initFn != nil {
						thisObj.t = base
						c.Invoke(initFn, thisObj, getArgs)
					}
				}
			}
			return constUndefined
		}))
	})
	return t.superType
}

func (t *valueType) CompareTo(other Value, c *Context) CompareResult {
	otherType, isType := other.(ValueType)
	if !isType || t.TypeId != otherType.TypeId {
		return CompareResultNotEqual
	}
	return CompareResultEqual
}

func (t *valueType) IsSubOf(t2 *valueType) bool {
	if t.TypeId == t2.TypeId {
		return true
	}
	for _, b := range t.Bases {
		if b.IsSubOf(t2) {
			return true
		}
	}
	return false
}

func (t *valueType) ToString(*Context) string {
	return fmt.Sprintf("<type %s>", t.Name)
}

func (*valueType) ToGoValue() interface{} {
	return nil
}

func (*valueType) GoType() reflect.Type {
	return reflect.TypeOf(nil)
}

func (*valueType) IsTrue() bool {
	return true
}

func (*valueType) GetIndex(int, *Context) Value {
	return constUndefined
}

func (t *valueType) GetMember(name string, c *Context) Value {
	if m, e := t.findStatic(t, name, c); e {
		return m
	}
	return getExtMember(t, name, c)
}

func (t *valueType) SetMember(name string, val Value, c *Context) {
	t.Statics.Store(name, val)
}

func (t *valueType) findStatic(vt *valueType, name string, c *Context) (Value, bool) {
	if m, ok := t.Statics.Load(name); ok {
		return m.(Value), true
	}
	if m, ok := t.Statics.Load("__getAttr__"); ok {
		if getattr, ok := m.(ValueCallable); ok {
			c.Invoke(getattr, vt, Args(NewStr(name)))
			return c.RetVal, true
		}
	}
	for _, b := range t.Bases {
		r, e := b.findStatic(vt, name, c)
		if e {
			return r, e
		}
	}
	return constUndefined, false
}

func (t *valueType) GetName() string {
	return t.Name
}

func (t *valueType) getInitFunc() ValueCallable {
	if initer, found := t.Members.Load("__init__"); found {
		if initFn, callable := initer.(ValueCallable); callable {
			return initFn
		}
	}
	return nil
}

func (t *valueType) GetArgNames() []string {
	if initFn := t.getInitFunc(); initFn != nil {
		return initFn.GetArgNames()
	}
	return []string{}
}

func (t *valueType) Invoke(c *Context, this Value, args []Value) {
	rv := NewObject(t)
	if initFn := t.getInitFunc(); initFn != nil {
		c.Invoke(initFn, rv, func() []Value { return args })
	}
	c.RetVal = rv
}

func (t *valueType) findMember(name string) Value {
	if member, found := t.Members.Load(name); found {
		return member.(Value)
	}
	for _, baseCls := range t.Bases {
		if member := baseCls.findMember(name); member != nil {
			return member.(Value)
		}
	}
	return nil
}

var nextTypeId int32 = 100000

func NextTypeId() int {
	return int(atomic.AddInt32(&nextTypeId, 1))
}

type ClassBuilder struct {
	t ValueType
}

func NewClassBuilder(name string, bases ...ValueType) *ClassBuilder {
	t := NewType(NextTypeId(), name)
	if len(bases) > 0 {
		t.Bases = bases
	} else {
		t.Bases = []ValueType{TypeObject}
	}
	return &ClassBuilder{t}
}

func (h *ClassBuilder) Constructor(f func(*Context, ValueObject, []Value)) *ClassBuilder {
	h.t.Members.Store("__init__", NewNativeFunction(h.t.Name+".__init__", func(c *Context, this Value, args []Value) Value {
		thisObj := c.MustObject(this)
		f(c, thisObj, args)
		return constUndefined
	}))
	return h
}

func (h *ClassBuilder) Method(name string, f func(*Context, ValueObject, []Value) Value, args ...string) *ClassBuilder {
	h.t.Members.Store(name, NewNativeFunction(h.t.Name+"."+name, func(c *Context, this Value, args []Value) Value {
		return f(c, c.MustObject(this), args)
	}, args...))
	return h
}

func (h *ClassBuilder) Build() ValueType {
	return h.t
}
