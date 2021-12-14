package parser

import (
	"math/big"

	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/runtime"
)

func (v *ParseVisitor) VisitExprPlusMinus(ctx *ExprPlusMinusContext) interface{} {
	switch ctx.GetOp().GetTokenType() {
	case ZggLexerPLUS:
		return &ast.ExprPlus{BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		}}
	case ZggLexerMINUS:
		return &ast.ExprMinus{BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		}}
	}
	panic("should not reach this line")
}

func (v *ParseVisitor) VisitExprPow(ctx *ExprPowContext) interface{} {
	return &ast.ExprPow{BinOp: ast.BinOp{
		Left:  ctx.Expr(0).Accept(v).(ast.Expr),
		Right: ctx.Expr(1).Accept(v).(ast.Expr),
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
	switch ctx.GetOp().GetTokenType() {
	case ZggLexerTIMES:
		return &ast.ExprTimes{BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		}}
	case ZggLexerDIV:
		return &ast.ExprDiv{BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		}}
	case ZggLexerMOD:
		return &ast.ExprMod{BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		}}
	}
	panic("should not reach this line")
}
