package runtime

import (
	"fmt"
	"reflect"
)

type CompareResult int

const (
	CompareResultNotEqual CompareResult = 0
	CompareResultEqual    CompareResult = 1
	CompareResultLess     CompareResult = 2
	CompareResultGreater  CompareResult = 4
)

type Value interface {
	GetOwner() Value
	SetOwner(Value)
	Type() ValueType
	CompareTo(Value, *Context) CompareResult
	ToString(*Context) string
	ToGoValue() interface{}
	GoType() reflect.Type
	IsTrue() bool
	GetIndex(int, *Context) Value
	GetMember(string, *Context) Value
}

type ValueBase struct {
	Owner Value
}

func (v *ValueBase) GetOwner() Value {
	if v.Owner == nil {
		return constUndefined
	}
	return v.Owner
}

func (v *ValueBase) SetOwner(owner Value) {
	v.Owner = owner
}

func makeMember(owner, member Value, c *Context) Value {
	if callable, ok := member.(ValueCallable); ok {
		if _, isObj := member.(ValueObject); !isObj {
			return NewBoundMethod(owner, callable)
		}
	}
	return member
}

var MakeMember = makeMember

func getExtMember(owner Value, name string, c *Context) Value {
	tid := owner.Type().TypeId
	extName := fmt.Sprintf("%d#%s", tid, name)
	if ext, found := c.FindValue(extName); found {
		return makeMember(owner, ext, c)
	}
	switch name {
	case "$get":
		return makeMember(owner, commonValueGet, c)
	}
	return constUndefined
}

type CanLen interface {
	Len() int
}

type CanSetMember interface {
	SetMember(string, Value, *Context)
}

type CanSetIndex interface {
	SetIndex(int, Value, *Context)
}

type CanHash interface {
	Hash() int64
}

var (
	commonValueGet = NewNativeFunction("$get", func(c *Context, this Value, args []Value) Value {
		var path ValueStr
		EnsureFuncParams(c, "$get", args, ArgRuleRequired("path", TypeStr, &path))
		return GetValueByPath(c, this, path.Value())
	}, "path")
)
