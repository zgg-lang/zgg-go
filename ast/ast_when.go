package ast

import (
	"github.com/zgg-lang/zgg-go/runtime"
)

type Case struct {
	Condition Expr
	Action    Expr
}

type ExprWhen struct {
	Cases []Case
	Else  Expr
}

func (expr *ExprWhen) Eval(c *runtime.Context) {
	for _, whenCase := range expr.Cases {
		whenCase.Condition.Eval(c)
		if c.ReturnTrue() {
			whenCase.Action.Eval(c)
			return
		}
	}
	expr.Else.Eval(c)
}

type ValueCondition interface {
	IsMatch(*runtime.Context, runtime.Value) bool
	Return(*runtime.Context)
}

type ValueConditionInList struct {
	ValueList []Expr
	Ret       Expr
}

func (vc *ValueConditionInList) IsMatch(c *runtime.Context, v runtime.Value) bool {
	for _, expected := range vc.ValueList {
		expected.Eval(c)
		if c.ValuesEqual(v, c.RetVal) {
			return true
		}
	}
	return false
}

func (vc *ValueConditionInList) Return(c *runtime.Context) {
	vc.Ret.Eval(c)
}

type ValueConditionInRange struct {
	Min, Max               Expr
	IncludeMin, IncludeMax bool
	Ret                    Expr
}

func (vc *ValueConditionInRange) IsMatch(c *runtime.Context, v runtime.Value) bool {
	var min, max runtime.Value
	if v := vc.Min; v != nil {
		v.Eval(c)
		min = c.RetVal
	} else {
		min = runtime.Undefined()
	}
	if _, isUndefined := min.(runtime.ValueUndefined); !isUndefined {
		if c.ValuesLess(v, min) {
			return false
		}
		if !vc.IncludeMin && c.ValuesEqual(v, min) {
			return false
		}
	}
	if v := vc.Max; v != nil {
		v.Eval(c)
		max = c.RetVal
	} else {
		max = runtime.Undefined()
	}
	if _, isUndefined := max.(runtime.ValueUndefined); !isUndefined {
		if c.ValuesGreater(v, max) {
			return false
		}
		if !vc.IncludeMax && c.ValuesEqual(v, max) {
			return false
		}
	}
	return true
}

func (vc *ValueConditionInRange) Return(c *runtime.Context) {
	vc.Ret.Eval(c)
}

type ExprWhenValue struct {
	Input Expr
	Cases []ValueCondition

	Else Expr
}

func (expr *ExprWhenValue) Eval(c *runtime.Context) {
	expr.Input.Eval(c)
	v := c.RetVal
	for _, whenCase := range expr.Cases {
		if whenCase.IsMatch(c, v) {
			whenCase.Return(c)
			return
		}
	}
	expr.Else.Eval(c)
}
