package ast

import (
	"testing"

	"github.com/zgg-lang/zgg-go/runtime"
)

func assertResult(t *testing.T, expected runtime.Value, expr Expr) {
	c := runtime.NewContext(true, false, false)
	expr.Eval(c)
	if c.RetVal != expected {
		t.Errorf("计算错误 %s != %s", c.RetVal, expected)
	}
}

func TestBinOp(t *testing.T) {
	assertResult(t, runtime.NewInt(100), &ExprTimes{
		BinOp{
			Left:  &ExprInt{Value: runtime.NewInt(10)},
			Right: &ExprInt{Value: runtime.NewInt(10)},
		},
	})
	assertResult(t, runtime.NewFloat(4000), &ExprTimes{
		BinOp{
			Left: &ExprInt{Value: runtime.NewInt(200)},
			Right: &ExprPlus{
				BinOp{
					Left:  &ExprInt{Value: runtime.NewInt(10)},
					Right: &ExprFloat{Value: runtime.NewFloat(10)},
				},
			},
		},
	})
	assertResult(t, runtime.NewStr("aaabbb"), &ExprPlus{
		BinOp{
			Left:  &ExprStr{Value: runtime.NewStr("aaa")},
			Right: &ExprStr{Value: runtime.NewStr("bbb")},
		},
	})
}
