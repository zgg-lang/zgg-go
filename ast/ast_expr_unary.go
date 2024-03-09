package ast

import (
	"math/big"

	"github.com/zgg-lang/zgg-go/runtime"
)

type ExprNegative struct {
	Expr Expr
}

func (e *ExprNegative) Eval(c *runtime.Context) {
	e.Expr.Eval(c)
	switch v := ensureZgg(c.RetVal, c).(type) {
	case runtime.ValueInt:
		c.RetVal = runtime.NewInt(-v.Value())
		return
	case runtime.ValueFloat:
		c.RetVal = runtime.NewFloat(-v.Value())
		return
	case runtime.ValueBigNum:
		r := big.NewFloat(0).SetPrec(1024)
		r.Neg(v.Value())
		c.RetVal = runtime.NewBigNum(r)
		return
	}
	c.RaiseRuntimeError("type error")
}

type ExprIncDec struct {
	Pos
	Lval Lval
	Expr Expr
	Pre  bool
}

func (e *ExprIncDec) Eval(c *runtime.Context) {
	e.Lval.Eval(c)
	oldVal := c.RetVal
	e.Expr.Eval(c)
	newVal := c.RetVal
	e.Lval.SetValue(c, newVal)
	if !e.Pre {
		c.RetVal = oldVal
	}
}

type ExprBitNot struct {
	Expr Expr
}

func (e *ExprBitNot) Eval(c *runtime.Context) {
	e.Expr.Eval(c)
	switch v := ensureZgg(c.RetVal, c).(type) {
	case runtime.ValueInt:
		c.RetVal = runtime.NewInt(^v.Value())
		return
	}
	c.RaiseRuntimeError("type error")
}

type ExprUse struct {
	Expr       Expr
	Identifier string
	DeferFunc  ExprFunc
}

func (e *ExprUse) Eval(c *runtime.Context) {
	e.Expr.Eval(c)
	v := c.RetVal
	if e.DeferFunc.Value != nil {
		closer := e.DeferFunc.Value.CloneWithEnv(c)
		c.AddBlockDefer(closer, []runtime.Value{v}, true)
	} else if e.Identifier != "" {
		if closer := v.GetMember(e.Identifier, c); c.IsCallable(closer) {
			c.AddBlockDefer(closer, []runtime.Value{}, true)
		} else {
			c.RaiseRuntimeError("use value without close/Close method")
		}
	} else {
		if closer := v.GetMember("close", c); c.IsCallable(closer) {
			c.AddBlockDefer(closer, []runtime.Value{}, true)
		} else if closer := v.GetMember("Close", c); c.IsCallable(closer) {
			c.AddBlockDefer(closer, []runtime.Value{}, true)
		} else {
			c.RaiseRuntimeError("use value without close/Close method")
		}
	}
}

type ExprAssertError struct {
	Expr Expr
}

func (e *ExprAssertError) Eval(c *runtime.Context) {
	e.Expr.Eval(c)
	r := c.RetVal
	if rs, ok := r.(runtime.ValueArray); ok {
		n := rs.Len()
		if n > 0 {
			last := rs.GetIndex(n-1, c)
			if err, isErr := last.ToGoValue(c).(error); isErr {
				c.RaiseRuntimeError("assert error fail! %s", err)
			} else {
				n--
			}
			switch n {
			case 0:
				c.RetVal = runtime.Undefined()
			case 1:
				c.RetVal = rs.GetIndex(0, c)
			default:
				newRet := runtime.NewArray(n)
				for i := 0; i < n; i++ {
					newRet.PushBack(rs.GetIndex(i, c))
				}
				c.RetVal = newRet
			}
		}
	} else if err, isErr := r.ToGoValue(c).(error); isErr {
		if err != nil {
			c.RaiseRuntimeError("assert error fail! %s", err)
		} else {
			c.RetVal = runtime.Undefined()
		}
	}
	c.RetVal = r
}
