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
		New       func(c *Context, args []Value) Value
		Members   *sync.Map
		Statics   *sync.Map
	}
	ValueType = *valueType
)

func NewTypeWithCreator(id int, name string, creator func(*Context, []Value) Value) ValueType {
	rv := &valueType{ValueBase: &ValueBase{}, TypeId: id, Name: name, Members: new(sync.Map), Statics: new(sync.Map)}
	rv.Statics.Store("__name__", NewStr(name))
	rv.New = creator
	return rv
}

func NewType(id int, name string) ValueType {
	return NewTypeWithCreator(id, name, nil)
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
					if initFn := base.getInitFunc(c); initFn != nil {
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
		if getattr, ok := c.GetCallable(m); ok {
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

func (t *valueType) getInitFunc(c *Context) ValueCallable {
	if initer, found := t.Members.Load("__init__"); found {
		if initFn, callable := c.GetCallable(initer); callable {
			return initFn
		}
	}
	return nil
}

func (t *valueType) GetArgNames(c *Context) []string {
	if initFn := t.getInitFunc(c); initFn != nil {
		return initFn.GetArgNames(c)
	}
	return []string{}
}

func (t *valueType) Invoke(c *Context, this Value, args []Value) {
	if t.New != nil {
		c.RetVal = t.New(c, args)
	} else {
		rv := NewObject(t)
		if initFn := t.getInitFunc(c); initFn != nil {
			c.Invoke(initFn, rv, func() []Value { return args })
		}
		c.RetVal = rv
	}
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

func (t *valueType) Hash() int64 {
	return int64(t.TypeId)
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

func (h *ClassBuilder) Methods(names []string, f func(*Context, ValueObject, []Value) Value, args ...string) *ClassBuilder {
	for _, name := range names {
		h.Method(name, f, args...)
	}
	return h
}

func (h *ClassBuilder) StaticMethod(name string, f func(*Context, Value, []Value) Value, args ...string) *ClassBuilder {
	h.t.Statics.Store(name, NewNativeFunction(h.t.Name+"."+name, func(c *Context, this Value, args []Value) Value {
		return f(c, this, args)
	}, args...))
	return h
}

func (h *ClassBuilder) Build() ValueType {
	return h.t
}

func addMembersAndStatics(vt ValueType, m map[string]ValueCallable) {
	for name, memberFunc := range m {
		vt.Members.Store(name, memberFunc)
		nativeFunc := memberFunc.Invoke
		staticFunc := NewNativeFunction(memberFunc.GetName(), func(c *Context, this Value, args []Value) Value {
			target := args[0]
			args = args[1:]
			nativeFunc(c, target, args)
			return c.RetVal
		})
		vt.Statics.Store(name, staticFunc)
	}
}
