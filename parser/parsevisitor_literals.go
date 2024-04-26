package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/runtime"
)

func (v *ParseVisitor) VisitIntegerZero(ctx *IntegerZeroContext) interface{} {
	return &ast.ExprInt{Value: runtime.NewInt(0)}
}

func (v *ParseVisitor) VisitIntegerDec(ctx *IntegerDecContext) interface{} {
	val, _ := strconv.ParseInt(ctx.GetText(), 10, 64)
	return &ast.ExprInt{Value: runtime.NewInt(val)}
}

func (v *ParseVisitor) VisitIntegerHex(ctx *IntegerHexContext) interface{} {
	val, _ := strconv.ParseInt(ctx.GetText()[2:], 16, 64)
	return &ast.ExprInt{Value: runtime.NewInt(val)}
}

func (v *ParseVisitor) VisitIntegerOct(ctx *IntegerOctContext) interface{} {
	val, _ := strconv.ParseInt(ctx.GetText()[1:], 8, 64)
	return &ast.ExprInt{Value: runtime.NewInt(val)}
}

func (v *ParseVisitor) VisitIntegerBin(ctx *IntegerBinContext) interface{} {
	val, _ := strconv.ParseInt(ctx.GetText()[2:], 2, 64)
	return &ast.ExprInt{Value: runtime.NewInt(val)}
}

func (v *ParseVisitor) VisitLiteralInteger(ctx *LiteralIntegerContext) interface{} {
	switch c := ctx.Integer().(type) {
	case *IntegerZeroContext:
		return v.VisitIntegerZero(c)
	case *IntegerDecContext:
		return v.VisitIntegerDec(c)
	case *IntegerHexContext:
		return v.VisitIntegerHex(c)
	case *IntegerOctContext:
		return v.VisitIntegerOct(c)
	case *IntegerBinContext:
		return v.VisitIntegerBin(c)
	}
	panic("should not reach this line")
}

func (v *ParseVisitor) VisitLiteralFloat(ctx *LiteralFloatContext) interface{} {
	val, _ := strconv.ParseFloat(ctx.GetText(), 64)
	return &ast.ExprFloat{Value: runtime.NewFloat(val)}
}

func (v *ParseVisitor) VisitLiteralENum(ctx *LiteralENumContext) interface{} {
	parts := strings.SplitN(ctx.GetText(), "e", 2)
	if len(parts) != 2 {
		panic("parse enum failed")
	}
	base, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		panic("parse enum base failed. " + err.Error())
	}
	exponent, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("parse enum exponent failed. " + err.Error())
	}
	return &ast.ExprFloat{Value: runtime.NewFloat(base * math.Pow10(exponent))}
}

func (v *ParseVisitor) VisitLiteralBool(ctx *LiteralBoolContext) interface{} {
	return &ast.ExprBool{Value: runtime.NewBool(ctx.GetText() == "true")}
}

func (v *ParseVisitor) VisitLiteralString(ctx *LiteralStringContext) interface{} {
	return ctx.StringLiteral().Accept(v)
}

func (v *ParseVisitor) VisitLiteralBigNum(ctx *LiteralBigNumContext) interface{} {
	lit := ctx.BIGNUM().GetText()
	val := big.NewFloat(0).SetPrec(1024)
	val.Parse(lit[:len(lit)-1], 10)
	return &ast.ExprBigNum{Value: runtime.NewBigNum(val)}
}

func (v *ParseVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	if ps := ctx.RSTRING(); ps != nil {
		s := ps.GetText()
		l := len(s)
		return &ast.ExprStr{Value: runtime.NewStr(s[3 : l-3])}
	}
	if ps := ctx.STRING(); ps != nil {
		s := ps.GetText()
		l := len(s)
		return &ast.ExprStr{Value: runtime.NewStr(strings.ReplaceAll(s[2:l-1], "\\'", "'"))}
	}
	return ctx.TemplateString().Accept(v)
}

func (v *ParseVisitor) VisitTemplateString(ctx *TemplateStringContext) interface{} {
	items := ctx.AllTsItem()
	if len(items) == 0 {
		return &ast.ExprStr{Value: runtime.NewStr("")}
	}
	var rv ast.Expr = nil
	for _, item := range items {
		if rv == nil {
			rv = item.Accept(v).(ast.Expr)
		} else {
			rv = &ast.ExprPlus{BinOp: ast.BinOp{
				Left:  rv,
				Right: item.Accept(v).(ast.Expr),
			}}
		}
	}
	return rv
}

func (v *ParseVisitor) VisitTsRaw(ctx *TsRawContext) interface{} {
	return &ast.ExprStr{Value: runtime.NewStr(parseStrContent(ctx.TS_RAW().GetText(), 0, 0))}
}

func (v *ParseVisitor) VisitTsIdentifier(ctx *TsIdentifierContext) interface{} {
	return &ast.ExprToStr{Expr: &ast.LvalById{Name: ctx.TS_IDENTIFIER().GetText()[1:]}}
}

func (v *ParseVisitor) VisitTsExpr(ctx *TsExprContext) interface{} {
	return &ast.ExprToStr{Expr: ctx.Expr().Accept(v).(ast.Expr)}
}

func (v *ParseVisitor) VisitLiteralNil(ctx *LiteralNilContext) interface{} {
	return &ast.ExprNil{}
}

func (v *ParseVisitor) VisitLiteralUndefined(ctx *LiteralUndefinedContext) interface{} {
	return &ast.ExprUndefined{}
}

func (v *ParseVisitor) VisitLiteralFunc(ctx *LiteralFuncContext) interface{} {
	allArgs := ctx.AllIDENTIFIER()
	args := make([]string, len(allArgs))
	for i, a := range allArgs {
		args[i] = a.GetText()
	}
	body := ctx.CodeBlock().Accept(v).(*ast.Block)
	f := runtime.NewFunc("", args, ctx.MORE_ARGS() != nil, body)
	return &ast.ExprFunc{Value: f}
}

func (v *ParseVisitor) VisitLiteralLambdaExpr(ctx *LiteralLambdaExprContext) interface{} {
	allArgs := ctx.AllIDENTIFIER()
	args := make([]string, len(allArgs))
	for i, a := range allArgs {
		args[i] = a.GetText()
	}
	block := &ast.Block{
		Pos: getPos(v, ctx),
		Stmts: []ast.Stmt{
			&ast.StmtReturn{Pos: getPos(v, ctx), Value: ctx.Expr().Accept(v).(ast.Expr)},
		},
	}
	f := runtime.NewFunc("", args, ctx.MORE_ARGS() != nil, block)
	return &ast.ExprFunc{Value: f}
}

func (v *ParseVisitor) VisitLiteralLambdaBlock(ctx *LiteralLambdaBlockContext) interface{} {
	allArgs := ctx.AllIDENTIFIER()
	args := make([]string, len(allArgs))
	for i, a := range allArgs {
		args[i] = a.GetText()
	}
	body := ctx.CodeBlock().Accept(v).(*ast.Block)
	f := runtime.NewFunc("", args, ctx.MORE_ARGS() != nil, body)
	return &ast.ExprFunc{Value: f}
}

type memberDef struct {
	isStatic bool
	kvPair   kvPair
}

type kvPair struct {
	key, val ast.Expr
}

func (v *ParseVisitor) VisitLiteralObject(ctx *LiteralObjectContext) interface{} {
	items := ctx.AllObjItem()
	rv := &ast.ExprObject{}
	for _, item := range items {
		itemNodeInterface := item.Accept(v)
		switch itemNode := itemNodeInterface.(type) {
		case kvPair:
			rv.Items = append(rv.Items, ast.ExprObjectItemKV{Key: itemNode.key, Value: itemNode.val})
		case ast.Expr:
			rv.Items = append(rv.Items, ast.ExprObjectItemExpandObj{Obj: itemNode})
		}
	}
	return rv
}

func (v *ParseVisitor) VisitObjItemExpanded(ctx *ObjItemExpandedContext) interface{} {
	return ctx.Expr().Accept(v)
}

func (v *ParseVisitor) VisitObjItemKV(ctx *ObjItemKVContext) interface{} {
	return ctx.KeyValue().Accept(v)
}

func (v *ParseVisitor) VisitKVIdKey(ctx *KVIdKeyContext) interface{} {
	return kvPair{
		key: &ast.ExprStr{Value: runtime.NewStr(ctx.IDENTIFIER().GetText())},
		val: ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitKVStrKey(ctx *KVStrKeyContext) interface{} {
	return kvPair{
		key: ctx.StringLiteral().Accept(v).(ast.Expr),
		// key: &ast.ExprStr{Value: runtime.NewStr(parseStr(ctx.GetText()))},
		val: ctx.Expr().Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitKVExprKey(ctx *KVExprKeyContext) interface{} {
	return kvPair{
		key: ctx.Expr(0).Accept(v).(ast.Expr),
		val: ctx.Expr(1).Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitKVKeyFunc(ctx *KVKeyFuncContext) interface{} {
	allArgs := ctx.AllIDENTIFIER()
	id := allArgs[0].GetText()
	args := make([]string, len(allArgs))
	for i, a := range allArgs {
		args[i] = a.GetText()
	}
	body := ctx.CodeBlock().Accept(v).(*ast.Block)
	fVal := runtime.NewFunc(id, args[1:], ctx.MORE_ARGS() != nil, body)
	fNode := &ast.ExprFunc{Value: fVal}
	return kvPair{
		key: &ast.ExprStr{Value: runtime.NewStr(id)},
		val: fNode,
	}
}

func (v *ParseVisitor) VisitKVIdOnly(ctx *KVIdOnlyContext) interface{} {
	return kvPair{
		key: &ast.ExprStr{Value: runtime.NewStr(ctx.IDENTIFIER().GetText())},
		val: &ast.LvalById{Name: ctx.IDENTIFIER().GetText()},
	}
}

func (v *ParseVisitor) VisitKVExprOnly(ctx *KVExprOnlyContext) interface{} {
	expr := ctx.Expr()
	return kvPair{
		key: &ast.ExprStr{Value: runtime.NewStr(expr.GetText())},
		val: expr.Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitLiteralArray(ctx *LiteralArrayContext) interface{} {
	allItems := ctx.AllArrayItem()
	itemNum := len(allItems)
	rv := &ast.ExprArray{
		Items: make([]*ast.ArrayItem, itemNum),
	}
	for i, item := range allItems {
		rv.Items[i] = item.Accept(v).(*ast.ArrayItem)
	}
	return rv
}

func (v *ParseVisitor) VisitArrayItem(ctx *ArrayItemContext) interface{} {
	rv := &ast.ArrayItem{
		Expr:         ctx.Expr(0).Accept(v).(ast.Expr),
		ShouldExpand: ctx.MORE_ARGS() != nil,
	}
	if cond := ctx.GetCondition(); cond != nil {
		rv.Condition = cond.Accept(v).(ast.Expr)
	}
	return rv
}

func (v *ParseVisitor) VisitArrayComprehension(ctx *ArrayComprehensionContext) interface{} {
	rv := &ast.ArrayComprehension{
		ItemExpr: ctx.GetItemExpr().Accept(v).(ast.Expr),
	}
	rv.ValueName = ctx.GetValue().GetText()
	if t := ctx.GetIndexer(); t != nil {
		rv.IndexerName = t.GetText()
	}
	begin := ctx.GetBegin().Accept(v).(ast.Expr)
	if t := ctx.GetEnd(); t != nil {
		rv.RangeBegin = begin
		rv.RangeEnd = t.Accept(v).(ast.Expr)
		rv.RangeIncludingEnd = ctx.RANGE_WITH_END() != nil
	} else {
		rv.Iterable = begin
	}
	if t := ctx.GetFilter(); t != nil {
		rv.FilterExpr = t.Accept(v).(ast.Expr)
	}
	return rv
}

func (v *ParseVisitor) VisitObjectComprehension(ctx *ObjectComprehensionContext) interface{} {
	rv := &ast.ObjectComprehension{
		KeyExpr:   ctx.GetKeyExpr().Accept(v).(ast.Expr),
		ValueExpr: ctx.GetValueExpr().Accept(v).(ast.Expr),
	}
	rv.ValueName = ctx.GetValue().GetText()
	if t := ctx.GetIndexer(); t != nil {
		rv.IndexerName = t.GetText()
	}
	begin := ctx.GetBegin().Accept(v).(ast.Expr)
	if t := ctx.GetEnd(); t != nil {
		rv.RangeBegin = begin
		rv.RangeEnd = t.Accept(v).(ast.Expr)
		rv.RangeIncludingEnd = ctx.RANGE_WITH_END() != nil
	} else {
		rv.Iterable = begin
	}
	if t := ctx.GetFilter(); t != nil {
		rv.FilterExpr = t.Accept(v).(ast.Expr)
	}
	return rv
}
