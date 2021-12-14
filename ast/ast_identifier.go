package ast

import (
	"github.com/zgg-lang/zgg-go/runtime"
)

type ExprIdentifier struct {
	Name string
}

func (expr *ExprIdentifier) Eval(c *runtime.Context) {
	val, found := c.FindValue(expr.Name)
	if found {
		c.RetVal = val
	} else {
		c.RetVal = runtime.Undefined()
	}
}
