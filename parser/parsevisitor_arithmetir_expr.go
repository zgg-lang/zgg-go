package parser

import (
	"math"
	"math/big"
	"strings"

	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/runtime"
)

var CanCalcInCompileTime = true

func precalc(left, right ast.Expr, op int) ast.Expr {
	if !CanCalcInCompileTime {
		return nil
	}
	l, is := left.(ast.IsLiteral)
	if !is {
		return nil
	}
	r, is := right.(ast.IsLiteral)
	if !is {
		return nil
	}
	switch op {
	case ZggLexerPLUS:
		switch lv := l.ConstValue().(type) {
		case int64:
			switch rv := r.ConstValue().(type) {
			case int64:
				return &ast.ExprInt{Value: runtime.NewInt(lv + rv)}
			case float64:
				return &ast.ExprFloat{Value: runtime.NewFloat(float64(lv) + rv)}
			case string:
				return &ast.ExprStr{Value: runtime.NewStr("%v%s", lv, rv)}
			}
		case float64:
			switch rv := r.ConstValue().(type) {
			case int64:
				return &ast.ExprFloat{Value: runtime.NewFloat(lv + float64(rv))}
			case float64:
				return &ast.ExprFloat{Value: runtime.NewFloat(lv + rv)}
			case string:
				return &ast.ExprStr{Value: runtime.NewStr("%v%s", lv, rv)}
			}
		case string:
			return &ast.ExprStr{Value: runtime.NewStr("%s%v", lv, r.ConstValue())}
		}
	case ZggLexerMINUS:
		switch lv := l.ConstValue().(type) {
		case int64:
			switch rv := r.ConstValue().(type) {
			case int64:
				return &ast.ExprInt{Value: runtime.NewInt(lv - rv)}
			case float64:
				return &ast.ExprFloat{Value: runtime.NewFloat(float64(lv) - rv)}
			}
		case float64:
			switch rv := r.ConstValue().(type) {
			case int64:
				return &ast.ExprFloat{Value: runtime.NewFloat(lv - float64(rv))}
			case float64:
				return &ast.ExprFloat{Value: runtime.NewFloat(lv - rv)}
			}
		}
	case ZggLexerTIMES:
		switch lv := l.ConstValue().(type) {
		case int64:
			switch rv := r.ConstValue().(type) {
			case int64:
				return &ast.ExprInt{Value: runtime.NewInt(lv * rv)}
			case float64:
				return &ast.ExprFloat{Value: runtime.NewFloat(float64(lv) * rv)}
			}
		case float64:
			switch rv := r.ConstValue().(type) {
			case int64:
				return &ast.ExprFloat{Value: runtime.NewFloat(lv * float64(rv))}
			case float64:
				return &ast.ExprFloat{Value: runtime.NewFloat(lv * rv)}
			}
		case string:
			switch rv := r.ConstValue().(type) {
			case int64:
				return &ast.ExprStr{Value: runtime.NewStr(strings.Repeat(lv, int(rv)))}
			}
		}
	case ZggLexerDIV:
		switch lv := l.ConstValue().(type) {
		case int64:
			switch rv := r.ConstValue().(type) {
			case int64:
				if rv != 0 {
					return &ast.ExprInt{Value: runtime.NewInt(lv / rv)}
				}
			case float64:
				if rv != 0 {
					return &ast.ExprFloat{Value: runtime.NewFloat(float64(lv) / rv)}
				}
			}
		case float64:
			switch rv := r.ConstValue().(type) {
			case int64:
				if rv != 0 {
					return &ast.ExprFloat{Value: runtime.NewFloat(lv / float64(rv))}
				}
			case float64:
				if rv != 0 {
					return &ast.ExprFloat{Value: runtime.NewFloat(lv / rv)}
				}
			}
		}
	case ZggLexerMOD:
		switch lv := l.ConstValue().(type) {
		case int64:
			switch rv := r.ConstValue().(type) {
			case int64:
				if rv != 0 {
					return &ast.ExprInt{Value: runtime.NewInt(lv % rv)}
				}
			}
		}
	case ZggLexerPOW:
		switch lv := l.ConstValue().(type) {
		case int64:
			switch rv := r.ConstValue().(type) {
			case int64:
				result := math.Pow(float64(lv), float64(rv))
				if rv > 0 {
					return &ast.ExprInt{Value: runtime.NewInt(int64(result))}
				}
				return &ast.ExprFloat{Value: runtime.NewFloat(result)}
			case float64:
				result := math.Pow(float64(lv), float64(rv))
				return &ast.ExprFloat{Value: runtime.NewFloat(result)}
			}
		case float64:
			switch rv := r.ConstValue().(type) {
			case int64:
				result := math.Pow(float64(lv), float64(rv))
				return &ast.ExprFloat{Value: runtime.NewFloat(result)}
			case float64:
				result := math.Pow(float64(lv), float64(rv))
				return &ast.ExprFloat{Value: runtime.NewFloat(result)}
			}
		}
	}
	return nil
}

func (v *ParseVisitor) VisitExprPlusMinus(ctx *ExprPlusMinusContext) interface{} {
	var (
		left  = ctx.Expr(0).Accept(v).(ast.Expr)
		right = ctx.Expr(1).Accept(v).(ast.Expr)
	)
	if r := precalc(left, right, ctx.GetOp().GetTokenType()); r != nil {
		return r
	}
	switch ctx.GetOp().GetTokenType() {
	case ZggLexerPLUS:
		return &ast.ExprPlus{BinOp: ast.BinOp{
			Left:  left,
			Right: right,
		}}
	case ZggLexerMINUS:
		return &ast.ExprMinus{BinOp: ast.BinOp{
			Left:  left,
			Right: right,
		}}
	}
	panic("should not reach this line")
}

func (v *ParseVisitor) VisitExprPow(ctx *ExprPowContext) interface{} {
	var (
		left  = ctx.Expr(0).Accept(v).(ast.Expr)
		right = ctx.Expr(1).Accept(v).(ast.Expr)
	)
	if r := precalc(left, right, ZggLexerPOW); r != nil {
		return r
	}
	return &ast.ExprPow{BinOp: ast.BinOp{
		Left:  left,
		Right: right,
	}}
}

func (v *ParseVisitor) VisitExprNegative(ctx *ExprNegativeContext) interface{} {
	sub := ctx.Expr().Accept(v).(ast.Expr)
	switch subExpr := sub.(type) {
	case *ast.ExprInt:
		return &ast.ExprInt{Value: runtime.NewInt(-subExpr.Value.Value())}
	case *ast.ExprFloat:
		return &ast.ExprFloat{Value: runtime.NewFloat(-subExpr.Value.Value())}
	case *ast.ExprBigNum:
		return &ast.ExprBigNum{Value: runtime.NewBigNum(big.NewFloat(0).SetPrec(1024).Neg(subExpr.Value.Value()))}
	}
	return &ast.ExprNegative{Expr: sub}
}

func (v *ParseVisitor) VisitExprTimesDivMod(ctx *ExprTimesDivModContext) interface{} {
	var (
		left  = ctx.Expr(0).Accept(v).(ast.Expr)
		right = ctx.Expr(1).Accept(v).(ast.Expr)
	)
	if r := precalc(left, right, ctx.GetOp().GetTokenType()); r != nil {
		return r
	}
	switch ctx.GetOp().GetTokenType() {
	case ZggLexerTIMES:
		return &ast.ExprTimes{BinOp: ast.BinOp{
			Left:  left,
			Right: right,
		}}
	case ZggLexerDIV:
		return &ast.ExprDiv{BinOp: ast.BinOp{
			Left:  left,
			Right: right,
		}}
	case ZggLexerMOD:
		return &ast.ExprMod{BinOp: ast.BinOp{
			Left:  left,
			Right: right,
		}}
	}
	panic("should not reach this line")
}
