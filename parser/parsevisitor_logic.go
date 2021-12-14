package parser

import "github.com/zgg-lang/zgg-go/ast"

func (v *ParseVisitor) VisitExprLogicNot(ctx *ExprLogicNotContext) interface{} {
	return &ast.ExprLogicNot{
		Expr: ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitExprLogicAnd(ctx *ExprLogicAndContext) interface{} {
	return &ast.ExprLogicAnd{BinOp: ast.BinOp{
		Left:  ctx.Expr(0).Accept(v).(ast.Expr),
		Right: ctx.Expr(1).Accept(v).(ast.Expr),
	}}
}

func (v *ParseVisitor) VisitExprLogicOr(ctx *ExprLogicOrContext) interface{} {
	return &ast.ExprLogicOr{BinOp: ast.BinOp{
		Left:  ctx.Expr(0).Accept(v).(ast.Expr),
		Right: ctx.Expr(1).Accept(v).(ast.Expr),
	}}
}
