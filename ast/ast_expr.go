package ast

import "github.com/zgg-lang/zgg-go/runtime"

type Expr interface {
	Node
}

func ensureZgg(v runtime.Value, c *runtime.Context) runtime.Value {
	if gv, ok := v.(runtime.GoValue); ok {
		return runtime.FromGoValue(gv.ReflectedValue(), c)
	}
	return v
}
