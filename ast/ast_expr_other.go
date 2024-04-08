package ast

import (
	"math"

	"github.com/zgg-lang/zgg-go/runtime"
)

type ExprSlice struct {
	Container, Begin, End Expr
}

func (e *ExprSlice) getBound(c *runtime.Context, expr Expr, defval int64) int64 {
	if expr == nil {
		return defval
	}
	return c.MustInt(evalAndReturn(c, expr))
}

func (e *ExprSlice) Eval(c *runtime.Context) {
	container := evalAndReturn(c, e.Container)
	if cs, is := container.(runtime.CanSlice); is {
		c.RetVal = cs.Slice(c, e.getBound(c, e.Begin, 0), e.getBound(c, e.End, math.MaxInt64))
	} else {
		c.RaiseRuntimeError("Cannot slice")
	}
}
