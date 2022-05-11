package ast

import (
	"fmt"

	"github.com/zgg-lang/zgg-go/runtime"
)

type CallArgument struct {
	Keyword      string
	Arg          Expr
	ShouldExpand bool
}

type ExprCall struct {
	Pos
	Optional  bool
	Callee    Expr
	Arguments []CallArgument
}

func (expr *ExprCall) GetArgs(c *runtime.Context, callable runtime.ValueCallable) []runtime.Value {
	args := make([]runtime.Value, 0, len(expr.Arguments))
	argNames := callable.GetArgNames()
	argPos := make(map[string]int, len(argNames))
	for i, n := range argNames {
		argPos[n] = i
	}
	for _, arg := range expr.Arguments {
		arg.Arg.Eval(c)
		argVal := c.RetVal
		if arg.ShouldExpand {
			switch moreArgs := argVal.(type) {
			case runtime.ValueArray:
				for i := 0; i < moreArgs.Len(); i++ {
					args = append(args, moreArgs.GetIndex(i, c))
				}
			default:
				c.RaiseRuntimeError("more args must be array")
				return nil
			}
		} else if arg.Keyword != "" {
			pos, found := argPos[arg.Keyword]
			if !found {
				c.RaiseRuntimeError("'%s' is an invalid keyword argument for %s", arg.Keyword, callable.GetName())
			}
			if _, isUndefined := argVal.(runtime.ValueUndefined); !isUndefined {
				if pos < len(args) {
					args[pos] = argVal
				} else {
					for j := len(args); j < pos; j++ {
						args = append(args, runtime.Undefined())
					}
					args = append(args, argVal)
				}
			}
		} else {
			args = append(args, argVal)
		}
	}
	return args
}

func (expr *ExprCall) Eval(c *runtime.Context) {
	expr.Callee.Eval(c)
	calleeVal := c.RetVal
	switch callee := calleeVal.(type) {
	case runtime.ValueCallable:
		c.Invoke(callee, callee.GetOwner(), func() []runtime.Value { return expr.GetArgs(c, callee) })
	default:
		if expr.Optional {
			c.RetVal = runtime.Undefined()
		} else {
			c.RaiseRuntimeError(fmt.Sprintf("%s is not callable", calleeVal.Type().Name))
		}
	}
}

type ExprShortImport struct {
	ImportPath string
}

func (e *ExprShortImport) Eval(c *runtime.Context) {
	c.RetVal = c.ImportModule(e.ImportPath, false, "script")
}
