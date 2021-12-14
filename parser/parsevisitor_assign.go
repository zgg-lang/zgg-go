package parser

import (
	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/runtime"
)

func (v *ParseVisitor) VisitAssignExists(ctx *AssignExistsContext) interface{} {
	op := ctx.GetOp().GetTokenType()
	lval := ctx.Lval().Accept(v).(ast.Lval)
	valExpr := ctx.Expr().Accept(v).(ast.Expr)
	switch op {
	case ZggLexerPLUS_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprPlus{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerMINUS_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprMinus{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerTIMES_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprTimes{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerDIV_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprDiv{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerMOD_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprMod{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerBIT_AND_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprBitAnd{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerBIT_OR_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprBitOr{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerBIT_XOR_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprBitXor{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerBIT_SHL_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprBitShl{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	case ZggLexerBIT_SHR_ASSIGN:
		return &ast.ExprAssign{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Expr: &ast.ExprBitShr{BinOp: ast.BinOp{
				Left:  lval,
				Right: valExpr,
			}},
		}
	}
	return &ast.ExprAssign{
		Pos:  getPos(v, ctx),
		Lval: lval,
		Expr: valExpr,
	}
}

func (v *ParseVisitor) VisitAssignNew(ctx *AssignNewContext) interface{} {
	id := ctx.IDENTIFIER().GetText()
	return &ast.ExprLocalAssign{
		Pos:   getPos(v, ctx),
		Names: []string{id},
		Type:  ast.AssignTypeSingle,
		Expr:  ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitAssignNewDeArray(ctx *AssignNewDeArrayContext) interface{} {
	allIds := ctx.AllIDENTIFIER()
	names := make([]string, len(allIds))
	for i, idToken := range allIds {
		id := idToken.GetText()
		names[i] = id
	}
	return &ast.ExprLocalAssign{
		Pos:        getPos(v, ctx),
		Names:      names,
		Type:       ast.AssignTypeDeArray,
		ExpandLast: ctx.MORE_ARGS() != nil,
		Expr:       ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitAssignNewDeObject(ctx *AssignNewDeObjectContext) interface{} {
	allIds := ctx.AllIDENTIFIER()
	names := make([]string, len(allIds))
	for i, idToken := range allIds {
		id := idToken.GetText()
		names[i] = id
	}
	return &ast.ExprLocalAssign{
		Pos:   getPos(v, ctx),
		Names: names,
		Type:  ast.AssignTypeDeObject,
		Expr:  ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitAssignNewLocal(ctx *AssignNewLocalContext) interface{} {
	return &ast.ExprLocalNewAssign{
		Pos:  getPos(v, ctx),
		Expr: ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitPreIncDec(ctx *PreIncDecContext) interface{} {
	lval := ctx.Lval().Accept(v).(ast.Lval)
	switch ctx.GetOp().GetTokenType() {
	case ZggLexerPLUS_PLUS:
		return &ast.ExprIncDec{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Pre:  true,
			Expr: &ast.ExprPlus{BinOp: ast.BinOp{
				Left:  lval,
				Right: &ast.ExprInt{Value: runtime.NewInt(1)},
			}},
		}
	case ZggLexerMINUS_MINUS:
		return &ast.ExprIncDec{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Pre:  true,
			Expr: &ast.ExprMinus{BinOp: ast.BinOp{
				Left:  lval,
				Right: &ast.ExprInt{Value: runtime.NewInt(1)},
			}},
		}
	}
	panic("should not reach here!")
}

func (v *ParseVisitor) VisitPostIncDec(ctx *PostIncDecContext) interface{} {
	lval := ctx.Lval().Accept(v).(ast.Lval)
	switch ctx.GetOp().GetTokenType() {
	case ZggLexerPLUS_PLUS:
		return &ast.ExprIncDec{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Pre:  false,
			Expr: &ast.ExprPlus{BinOp: ast.BinOp{
				Left:  lval,
				Right: &ast.ExprInt{Value: runtime.NewInt(1)},
			}},
		}
	case ZggLexerMINUS_MINUS:
		return &ast.ExprIncDec{
			Pos:  getPos(v, ctx),
			Lval: lval,
			Pre:  false,
			Expr: &ast.ExprMinus{BinOp: ast.BinOp{
				Left:  lval,
				Right: &ast.ExprInt{Value: runtime.NewInt(1)},
			}},
		}
	}
	panic("should not reach here!")
}
