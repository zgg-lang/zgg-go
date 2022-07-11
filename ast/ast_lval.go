package ast

import (
	"fmt"

	"github.com/zgg-lang/zgg-go/runtime"
)

type Lval interface {
	Expr
	GetOwner(*runtime.Context) runtime.Value
	SetValue(*runtime.Context, runtime.Value)
}

type LvalById struct {
	Name string
}

func (v *LvalById) Eval(c *runtime.Context) {
	val, found := c.FindValue(v.Name)
	if found {
		if val == nil {
			c.RaiseRuntimeError(fmt.Sprintf("use variable %s before initialized", v.Name))
			return
		} else {
			c.RetVal = val
		}
	} else {
		c.RetVal = runtime.Undefined()
	}
}

func (v *LvalById) GetOwner(*runtime.Context) runtime.Value {
	return runtime.Undefined()
}

func (v *LvalById) SetValue(c *runtime.Context, val runtime.Value) {
	c.ModifyValue(v.Name, val)
}

type LvalByField struct {
	Owner Expr
	Field Expr
}

func (v *LvalByField) Eval(c *runtime.Context) {
	v.Field.Eval(c)
	fieldVal := c.RetVal
	switch field := fieldVal.(type) {
	case runtime.ValueStr:
		v.Owner.Eval(c)
		owner := c.RetVal
		c.RetVal = owner.GetMember(field.Value(), c)
		return
	case runtime.ValueInt:
		{
			index := int(field.Value())
			v.Owner.Eval(c)
			owner := c.RetVal
			c.RetVal = owner.GetIndex(index, c)
			return
		}
	}
	c.RetVal = runtime.Undefined()
}

func (v *LvalByField) GetOwner(c *runtime.Context) runtime.Value {
	v.Owner.Eval(c)
	return c.RetVal
}

func (v *LvalByField) SetValue(c *runtime.Context, val runtime.Value) {
	v.Field.Eval(c)
	fieldVal := c.RetVal
	v.Owner.Eval(c)
	ownerVal := c.RetVal
	switch field := fieldVal.(type) {
	case runtime.ValueStr:
		if owner, ok := ownerVal.(runtime.CanSetMember); ok {
			owner.SetMember(field.Value(), runtime.MakeMember(ownerVal, val, c), c)
		}
	case runtime.ValueInt:
		if owner, ok := ownerVal.(runtime.CanSetIndex); ok {
			owner.SetIndex(int(field.Value()), runtime.MakeMember(ownerVal, val, c), c)
		}
	}
}
