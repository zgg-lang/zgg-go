package parser

import (
	"strings"

	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/runtime"
)

func (v *ParseVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	allArgs := ctx.AllFuncArgument()
	rv := make([]ast.CallArgument, 0, len(allArgs))
	for _, arg := range allArgs {
		rv = append(rv, arg.Accept(v).(ast.CallArgument))
	}
	return rv
}

func (v *ParseVisitor) VisitFuncArgument(ctx *FuncArgumentContext) interface{} {
	var rv ast.CallArgument
	if e := ctx.Expr(); e != nil {
		rv.Arg = e.Accept(v).(ast.Expr)
		rv.ShouldExpand = ctx.MORE_ARGS() != nil
	} else if c := ctx.CodeBlock(); c != nil {
		body := c.Accept(v).(*ast.Block)
		f := runtime.NewFunc("", []string{"it"}, false, body)
		rv.Arg = &ast.ExprFunc{Value: f}
	}
	if id := ctx.IDENTIFIER(); id != nil {
		rv.Keyword = id.GetText()
	}
	return rv
}

func (v *ParseVisitor) VisitExprCompare(ctx *ExprCompareContext) interface{} {
	rv := &ast.ExprCompare{
		First: ctx.Expr(0).Accept(v).(ast.Expr),
	}
	// for i, opToken := range ctx.AllComparator() {
	opToken := ctx.Comparator()
	var op int
	switch opToken.GetText() {
	case "==":
		op = ast.CompareOpEQ
	case "!=":
		op = ast.CompareOpNE
	case ">":
		op = ast.CompareOpGT
	case "<":
		op = ast.CompareOpLT
	case ">=":
		op = ast.CompareOpGE
	case "<=":
		op = ast.CompareOpLE
	default:
		panic("invalid op " + opToken.GetText())
	}
	rv.Ops = append(rv.Ops, op)
	target := ctx.Expr(1).Accept(v).(ast.Expr)
	if next, ok := target.(*ast.ExprCompare); ok {
		rv.Ops = append(rv.Ops, next.Ops...)
		rv.Targets = append(rv.Targets, next.First)
		rv.Targets = append(rv.Targets, next.Targets...)
	} else {
		rv.Targets = append(rv.Targets, target)
	}
	// }
	return rv
}

func (v *ParseVisitor) VisitExprIsType(ctx *ExprIsTypeContext) interface{} {
	return &ast.ExprIsType{BinOp: ast.BinOp{
		Left:  ctx.Expr(0).Accept(v).(ast.Expr),
		Right: ctx.Expr(1).Accept(v).(ast.Expr),
	}}
}

func (v *ParseVisitor) VisitExprAssign(ctx *ExprAssignContext) interface{} {
	return ctx.AssignExpr().Accept(v)
}

func (v *ParseVisitor) VisitExprSub(ctx *ExprSubContext) interface{} {
	return ctx.Expr().Accept(v)
}

func (v *ParseVisitor) VisitExprPreIncDec(ctx *ExprPreIncDecContext) interface{} {
	return ctx.PreIncDec().Accept(v)
}

func (v *ParseVisitor) VisitExprPostIncDec(ctx *ExprPostIncDecContext) interface{} {
	return ctx.PostIncDec().Accept(v)
}

func (v *ParseVisitor) VisitExprFallback(ctx *ExprFallbackContext) interface{} {
	return &ast.ExprFallback{
		BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		},
	}
}

func (v *ParseVisitor) VisitExprBitNot(ctx *ExprBitNotContext) interface{} {
	return &ast.ExprBitNot{
		Expr: ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitExprBitShift(ctx *ExprBitShiftContext) interface{} {
	switch ctx.GetOp().GetTokenType() {
	case ZggLexerBIT_SHL:
		return &ast.ExprBitShl{
			BinOp: ast.BinOp{
				Left:  ctx.Expr(0).Accept(v).(ast.Expr),
				Right: ctx.Expr(1).Accept(v).(ast.Expr),
			},
		}
	case ZggLexerBIT_SHR:
		return &ast.ExprBitShr{
			BinOp: ast.BinOp{
				Left:  ctx.Expr(0).Accept(v).(ast.Expr),
				Right: ctx.Expr(1).Accept(v).(ast.Expr),
			},
		}
	}
	panic("!")
}

func (v *ParseVisitor) VisitExprBitAnd(ctx *ExprBitAndContext) interface{} {
	return &ast.ExprBitAnd{
		BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		},
	}
}

func (v *ParseVisitor) VisitExprBitOr(ctx *ExprBitOrContext) interface{} {
	return &ast.ExprBitOr{
		BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		},
	}
}

func (v *ParseVisitor) VisitExprBitXor(ctx *ExprBitXorContext) interface{} {
	return &ast.ExprBitXor{
		BinOp: ast.BinOp{
			Left:  ctx.Expr(0).Accept(v).(ast.Expr),
			Right: ctx.Expr(1).Accept(v).(ast.Expr),
		},
	}
}

func (v *ParseVisitor) VisitExprUseMethod(ctx *ExprUseMethodContext) interface{} {
	return &ast.ExprUse{
		Expr:       ctx.Expr().Accept(v).(ast.Expr),
		Identifier: ctx.IDENTIFIER().GetText(),
	}
}

func (v *ParseVisitor) VisitExprUseCloser(ctx *ExprUseCloserContext) interface{} {
	return &ast.ExprUse{
		Expr: ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitExprUseBlock(ctx *ExprUseBlockContext) interface{} {
	return &ast.ExprUse{
		Expr: ctx.Expr().Accept(v).(ast.Expr),
		DeferFunc: ast.ExprFunc{
			Value: runtime.NewFunc("", []string{}, false,
				ctx.CodeBlock().Accept(v).(*ast.Block),
			),
		},
	}
}

func (v *ParseVisitor) VisitExprShortImport(ctx *ExprShortImportContext) interface{} {
	path := strings.ReplaceAll(ctx.IDENTIFIER().GetText(), "__", "/")
	if isGoStd := ctx.DOUBLE_AT() != nil; isGoStd {
		path = "gostd/" + path
	}
	return &ast.ExprCall{
		Pos:      getPos(v, ctx),
		Optional: false,
		Callee:   &ast.ExprIdentifier{Name: "import"},
		Arguments: []ast.CallArgument{
			{Arg: &ast.ExprStr{Value: runtime.NewStr(path)}, ShouldExpand: false},
		},
	}
	// return &ast.ExprShortImport{ImportPath: path}
}

func (v *ParseVisitor) VisitExprWhen(ctx *ExprWhenContext) interface{} {
	all := ctx.AllExpr()
	node := &ast.ExprWhen{}
	i := 1
	n := len(all)
	for ; i < n; i += 2 {
		condition := all[i-1].Accept(v).(ast.Expr)
		action := all[i].Accept(v).(ast.Expr)
		node.Cases = append(node.Cases, ast.Case{
			Condition: condition,
			Action:    action,
		})
	}
	if i <= n {
		node.Else = all[n-1].Accept(v).(ast.Expr)
	} else {
		node.Else = &ast.ExprNil{}
	}
	return node
}

func (v *ParseVisitor) VisitExprWhenValue(ctx *ExprWhenValueContext) interface{} {
	exprs := ctx.AllExpr()
	rv := &ast.ExprWhenValue{
		Input: exprs[0].Accept(v).(ast.Expr),
	}
	conds := ctx.AllWhenCondition()
	for i, cond := range conds {
		wc := cond.Accept(v).(ast.ValueCondition)
		switch wcv := wc.(type) {
		case *ast.ValueConditionInList:
			wcv.Ret = exprs[i+1].Accept(v).(ast.Expr)
		case *ast.ValueConditionInRange:
			wcv.Ret = exprs[i+1].Accept(v).(ast.Expr)
		case *ast.ValueConditionIsType:
			wcv.Ret = exprs[i+1].Accept(v).(ast.Expr)
		}
		rv.Cases = append(rv.Cases, wc)
	}
	if n := len(exprs); len(conds)+1 < n {
		rv.Else = exprs[n-1].Accept(v).(ast.Expr)
	} else {
		rv.Else = &ast.ExprNil{}
	}
	return rv
}

func (v *ParseVisitor) VisitWhenConditionInList(ctx *WhenConditionInListContext) interface{} {
	all := ctx.AllExpr()
	rv := &ast.ValueConditionInList{ValueList: make([]ast.Expr, len(all))}
	for i, e := range all {
		rv.ValueList[i] = e.Accept(v).(ast.Expr)
	}
	return rv
}

func (v *ParseVisitor) VisitWhenConditionInRange(ctx *WhenConditionInRangeContext) interface{} {
	rv := &ast.ValueConditionInRange{
		IncludeMin: true,
		IncludeMax: ctx.RANGE_WITH_END() != nil,
	}
	if lb := ctx.GetLowerBound(); lb != nil {
		rv.Min = lb.Accept(v).(ast.Expr)
	}
	if ub := ctx.GetUpperBound(); ub != nil {
		rv.Max = ub.Accept(v).(ast.Expr)
	}
	return rv
}

func (v *ParseVisitor) VisitWhenConditionIsType(ctx *WhenConditionIsTypeContext) interface{} {
	return &ast.ValueConditionIsType{
		ExpectedType: ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitExprAssertError(ctx *ExprAssertErrorContext) interface{} {
	return &ast.ExprAssertError{
		Expr: ctx.Expr().Accept(v).(ast.Expr),
	}
}
