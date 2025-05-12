package ast

import (
	"github.com/zgg-lang/zgg-go/runtime"
)

type IsLiteral interface {
	ConstValue() any
}

type ExprInt struct {
	Value runtime.ValueInt
}

func (n *ExprInt) Eval(c *runtime.Context) {
	c.RetVal = n.Value
}

func (n *ExprInt) ConstValue() any {
	return n.Value.Value()
}

type ExprStr struct {
	Value runtime.ValueStr
}

func (n *ExprStr) Eval(c *runtime.Context) {
	c.RetVal = n.Value
}

func (n *ExprStr) ConstValue() any {
	return n.Value.Value()
}

type ExprToStr struct {
	Expr Expr
}

func (n *ExprToStr) Eval(c *runtime.Context) {
	n.Expr.Eval(c)
	if _, isStr := c.RetVal.(runtime.ValueStr); isStr {
		return
	}
	c.RetVal = runtime.NewStr(c.RetVal.ToString(c))
}

type ExprFloat struct {
	Value runtime.ValueFloat
}

func (n *ExprFloat) Eval(c *runtime.Context) {
	c.RetVal = n.Value
}

func (n *ExprFloat) ConstValue() any {
	return n.Value.Value()
}

type ExprBool struct {
	Value runtime.ValueBool
}

func (n *ExprBool) Eval(c *runtime.Context) {
	c.RetVal = n.Value
}

func (n *ExprBool) ConstValue() any {
	return n.Value.Value()
}

type ExprNil struct {
}

func (n *ExprNil) Eval(c *runtime.Context) {
	c.RetVal = runtime.Nil()
}

type ExprUndefined struct {
}

func (n *ExprUndefined) Eval(c *runtime.Context) {
	c.RetVal = runtime.Undefined()
}

type ExprFunc struct {
	Value *runtime.ValueFunc
	Refs  map[string]runtime.Value
}

func (e *ExprFunc) Eval(c *runtime.Context) {
	c.RetVal = e.Value.CloneWithEnv(c)
}

// type ExprObject struct {
// Keys       []Expr
// Values     []Expr
// ExpandObjs []Expr
// }
//
// func (e *ExprObject) Eval(c *runtime.Context) {
// if len(e.Keys) != len(e.Values) {
// panic("!!")
// }
// rv := runtime.NewObject()
// for i, ke := range e.Keys {
// ke.Eval(c)
// k := c.RetVal
// e.Values[i].Eval(c)
// rv.SetMember(k.ToString(c), c.RetVal, c)
// }
// for _, expExpr := range e.ExpandObjs {
// expExpr.Eval(c)
// expObj, isObj := c.RetVal.(runtime.ValueObject)
// if !isObj {
// c.RaiseRuntimeError("object: expand item must be an object")
// return
// }
// expObj.Iterate(func(k string, v runtime.Value) {
// rv.SetMember(k, v, c)
// })
// }
// c.RetVal = rv
// }

type ExprObjectItemKV struct {
	Key, Value Expr
}

type ExprObjectItemExpandObj struct {
	Obj Expr
}

type ExprObject struct {
	Items []interface{}
}

func (e *ExprObject) Eval(c *runtime.Context) {
	rv := runtime.NewObject()
	for _, item := range e.Items {
		switch it := item.(type) {
		case ExprObjectItemKV:
			it.Key.Eval(c)
			k := c.RetVal
			it.Value.Eval(c)
			v := c.RetVal
			rv.SetMember(k.ToString(c), v, c)
		case ExprObjectItemExpandObj:
			it.Obj.Eval(c)
			o, ok := c.RetVal.(runtime.ValueObject)
			if !ok {
				c.RaiseRuntimeError("object: expand item must be an object")
			}
			o.Iterate(func(k string, v runtime.Value) {
				rv.SetMember(k, v, c)
			})
		default:
			panic("!!")
		}
	}
	c.RetVal = rv
}

type ArrayItem struct {
	Expr         Expr
	Condition    Expr
	ShouldExpand bool
}

type ExprArray struct {
	Items []*ArrayItem
}

func (e *ExprArray) Eval(c *runtime.Context) {
	rv := runtime.NewArray(len(e.Items))
	for _, item := range e.Items {
		if item.Condition != nil {
			item.Condition.Eval(c)
			if !c.RetVal.IsTrue() {
				continue
			}
		}
		item.Expr.Eval(c)
		val := c.RetVal
		if item.ShouldExpand {
			expanded, isArr := val.(runtime.ValueArray)
			if !isArr {
				c.RaiseRuntimeError("array: expanded item must be an array")
				return
			}
			for i := 0; i < expanded.Len(); i++ {
				rv.PushBack(expanded.GetIndex(i, c))
			}
		} else {
			rv.PushBack(val)
		}
	}
	c.RetVal = rv
}

type ExprBigNum struct {
	Value runtime.ValueBigNum
}

func (n *ExprBigNum) Eval(c *runtime.Context) {
	c.RetVal = n.Value
}
