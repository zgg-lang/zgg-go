package parser

import (
	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/runtime"
)

func (v *ParseVisitor) VisitStmtBlock(ctx *StmtBlockContext) interface{} {
	return ctx.CodeBlock().Accept(v)
}

func (v *ParseVisitor) VisitStmtAssign(ctx *StmtAssignContext) interface{} {
	return ctx.AssignExpr().Accept(v)
}

func (v *ParseVisitor) VisitStmtFor(ctx *StmtForContext) interface{} {
	r := &ast.StmtFor{
		Pos:   getPos(v, ctx),
		Init:  ctx.GetInitExpr().Accept(v).(ast.Expr),
		Check: ctx.GetCheckExpr().Accept(v).(ast.Expr),
		Next:  ctx.GetNextExpr().Accept(v).(ast.Expr),
		Exec:  ctx.GetExecBlock().Accept(v).(*ast.Block),
	}
	if l := ctx.GetLabel(); l != nil {
		r.Label = l.GetText()
	}
	r.Exec.Type = ast.BlockTypeLoopTop
	return r
}

func (v *ParseVisitor) VisitStmtForEach(ctx *StmtForEachContext) interface{} {
	if !tokenExpected(ctx, ctx.GetInword(), "in") {
		return nil
	}
	r := &ast.StmtForEach{
		Pos:     getPos(v, ctx),
		IdValue: ctx.GetIdValue().GetText(),
		Exec:    ctx.GetExecBlock().Accept(v).(*ast.Block),
	}
	r.Exec.Type = ast.BlockTypeLoopTop
	if id := ctx.GetIdIndex(); id != nil {
		r.IdIndex = id.GetText()
	}
	if endExpr := ctx.GetEnd(); endExpr == nil {
		r.Iteratable = ctx.GetBegin().Accept(v).(ast.Expr)
	} else {
		r.RangeBegin = ctx.GetBegin().Accept(v).(ast.Expr)
		r.RangeEnd = endExpr.Accept(v).(ast.Expr)
		r.RangeIncludingEnd = ctx.RANGE_WITH_END() != nil
	}
	if l := ctx.GetLabel(); l != nil {
		r.Label = l.GetText()
	}
	return r
}

func (v *ParseVisitor) VisitStmtDoWhile(ctx *StmtDoWhileContext) interface{} {
	r := &ast.StmtDoWhile{
		Pos:   getPos(v, ctx),
		Check: ctx.GetCheckExpr().Accept(v).(ast.Expr),
		Exec:  ctx.GetExecBlock().Accept(v).(*ast.Block),
	}
	r.Exec.Type = ast.BlockTypeLoopTop
	if l := ctx.GetLabel(); l != nil {
		r.Label = l.GetText()
	}
	return r
}

func (v *ParseVisitor) VisitStmtWhile(ctx *StmtWhileContext) interface{} {
	r := &ast.StmtWhile{
		Pos:   getPos(v, ctx),
		Check: ctx.GetCheckExpr().Accept(v).(ast.Expr),
		Exec:  ctx.GetExecBlock().Accept(v).(*ast.Block),
	}
	if l := ctx.GetLabel(); l != nil {
		r.Label = l.GetText()
	}
	r.Exec.Type = ast.BlockTypeLoopTop
	return r
}

func (v *ParseVisitor) VisitStmtBreak(ctx *StmtBreakContext) interface{} {
	r := &ast.StmtBreak{Pos: getPos(v, ctx)}
	if l := ctx.GetLabel(); l != nil {
		r.ToLabel = l.GetText()
	}
	return r
}

func (v *ParseVisitor) VisitStmtContinue(ctx *StmtContinueContext) interface{} {
	r := &ast.StmtContinue{Pos: getPos(v, ctx)}
	if l := ctx.GetLabel(); l != nil {
		r.ToLabel = l.GetText()
	}
	return r
}

func (v *ParseVisitor) VisitStmtIf(ctx *StmtIfContext) interface{} {
	r := &ast.StmtIf{Pos: getPos(v, ctx)}
	conds := ctx.AllIfCondition()
	blocks := ctx.AllCodeBlock()
	for i, cond := range conds {
		ifc := cond.Accept(v).(*ast.IfCase)
		ifc.Do = blocks[i].Accept(v).(*ast.Block)
		r.Cases = append(r.Cases, ifc)
	}
	if len(blocks) > len(conds) {
		r.ElseDo = blocks[len(blocks)-1].Accept(v).(*ast.Block)
	}
	return r
}

func (v *ParseVisitor) VisitIfCondition(ctx *IfConditionContext) interface{} {
	r := &ast.IfCase{
		Check: ctx.Expr().Accept(v).(ast.Expr),
	}
	if ass := ctx.AssignExpr(); ass != nil {
		r.Assignment = ass.Accept(v).(ast.Expr)
	}
	return r
}

func (v *ParseVisitor) VisitStmtSwitch(ctx *StmtSwitchContext) interface{} {
	rv := &ast.StmtSwitch{
		Pos: getPos(v, ctx),
		Val: ctx.GetTestValue().Accept(v).(ast.Expr),
	}
	cases := ctx.AllSwitchCase()
	for _, ic := range cases {
		c := ic.(*SwitchCaseContext)
		rv.Cases = append(rv.Cases, ast.SwitchCase{
			Condition:   c.WhenCondition().Accept(v).(ast.ValueCondition),
			Code:        c.Block().Accept(v).(*ast.Block),
			Fallthrough: c.FALLTHROUGH() != nil,
		})
	}
	if d := ctx.SwitchDefault(); d != nil {
		rv.Default = d.(*SwitchDefaultContext).Block().Accept(v).(*ast.Block)
	}
	return rv
}

func (v *ParseVisitor) VisitStmtReturn(ctx *StmtReturnContext) interface{} {
	expr := ctx.Expr()
	rv := &ast.StmtReturn{Pos: getPos(v, ctx), Value: &ast.ExprUndefined{}}
	if expr != nil {
		rv.Value = ctx.Expr().Accept(v).(ast.Expr)
	}
	return rv
}

func (v *ParseVisitor) VisitStmtReturnNone(ctx *StmtReturnNoneContext) interface{} {
	return &ast.StmtReturn{Pos: getPos(v, ctx), Value: &ast.ExprUndefined{}}
}

func (v *ParseVisitor) VisitStmtFuncCall(ctx *StmtFuncCallContext) interface{} {
	return ctx.CallStmt().Accept(v)
}

func (v *ParseVisitor) VisitCallStmt(ctx *CallStmtContext) interface{} {
	args := ctx.Arguments().Accept(v).([]ast.CallArgument)
	var rv ast.Stmt
	rv = &ast.ExprCall{
		Pos:       getPos(v, ctx),
		Optional:  ctx.OPTIONAL_CALL() != nil,
		Callee:    ctx.Expr().Accept(v).(ast.Expr),
		Arguments: args,
		IsBind:    ast.IsBindList(args),
	}
	if ctx.OPTIONAL_ELSE() != nil {
		rv = &ast.StmtFallback{
			Pos:      getPos(v, ctx),
			Stmt:     rv,
			Fallback: ctx.CodeBlock().Accept(v).(*ast.Block),
		}
	}
	return rv
}

func (v *ParseVisitor) VisitStmtExportIdentifier(ctx *StmtExportIdentifierContext) interface{} {
	return &ast.StmtExport{
		Pos:  getPos(v, ctx),
		Name: ctx.IDENTIFIER().GetText(),
		Expr: &ast.ExprIdentifier{
			Name: ctx.IDENTIFIER().GetText(),
		},
	}
}

func (v *ParseVisitor) VisitStmtExportExpr(ctx *StmtExportExprContext) interface{} {
	name := ctx.IDENTIFIER().GetText()
	return &ast.StmtExport{
		Pos:  getPos(v, ctx),
		Name: name,
		Expr: &ast.ExprLocalAssign{
			Names: []string{name},
			Type:  ast.AssignTypeSingle,
			Expr:  ctx.Expr().Accept(v).(ast.Expr),
		},
	}
}

func (v *ParseVisitor) VisitStmtExportFuncDefine(ctx *StmtExportFuncDefineContext) interface{} {
	ids := ctx.AllIDENTIFIER()
	name := ids[0].GetText()
	argIds := ids[1:]
	args := make([]string, len(argIds))
	for i, a := range argIds {
		args[i] = a.GetText()
	}
	body := ctx.CodeBlock().Accept(v).(ast.Node)
	funcExpr := &ast.ExprFunc{
		Value: runtime.NewFunc(name, args, ctx.MORE_ARGS() != nil, body),
	}
	return &ast.StmtExport{
		Pos:  getPos(v, ctx),
		Name: name,
		Expr: &ast.ExprLocalAssign{
			Names: []string{name},
			Type:  ast.AssignTypeSingle,
			Expr:  funcExpr,
		},
	}
}

func (v *ParseVisitor) VisitStmtPreIncDec(ctx *StmtPreIncDecContext) interface{} {
	return ctx.PreIncDec().Accept(v)
}

func (v *ParseVisitor) VisitStmtPostIncDec(ctx *StmtPostIncDecContext) interface{} {
	return ctx.PostIncDec().Accept(v)
}

func (v *ParseVisitor) VisitStmtFuncDefine(ctx *StmtFuncDefineContext) interface{} {
	ids := ctx.AllIDENTIFIER()
	name := ids[0].GetText()
	argIds := ids[1:]
	args := make([]string, len(argIds))
	for i, a := range argIds {
		args[i] = a.GetText()
	}
	body := ctx.CodeBlock().Accept(v).(ast.Node)
	funcExpr := &ast.ExprFunc{
		Value: runtime.NewFunc(name, args, ctx.MORE_ARGS() != nil, body),
	}
	return &ast.ExprLocalAssign{
		Names: []string{name},
		Type:  ast.AssignTypeSingle,
		Expr:  funcExpr,
	}
}

func (v *ParseVisitor) VisitMemberDef(ctx *MemberDefContext) interface{} {
	return memberDef{
		isStatic: ctx.STATIC() != nil,
		kvPair:   ctx.KeyValue().Accept(v).(kvPair),
	}
}

func (v *ParseVisitor) VisitStmtClassDefine(ctx *StmtClassDefineContext) interface{} {
	members := ctx.AllMemberDef()
	body := &ast.ExprObject{}
	staticBody := &ast.ExprObject{}
	for _, member := range members {
		m := member.Accept(v).(memberDef)
		if m.isStatic {
			staticBody.Items = append(staticBody.Items, ast.ExprObjectItemKV{Key: m.kvPair.key, Value: m.kvPair.val})
		} else {
			body.Items = append(body.Items, ast.ExprObjectItemKV{Key: m.kvPair.key, Value: m.kvPair.val})
		}
	}
	rv := &ast.StmtClassDefine{
		Pos:      getPos(v, ctx),
		Exported: ctx.EXPORT() != nil,
		Name:     ctx.GetClassName().GetText(),
		Static:   staticBody,
		Body:     body,
	}
	bases := ctx.GetBaseCls()
	rv.Bases = make([]ast.Expr, 0, len(bases))
	for _, b := range bases {
		// if b := ctx.GetBaseCls(); b != nil {
		rv.Bases = append(rv.Bases, b.Accept(v).(ast.Expr))
	}
	return rv
}

func (v *ParseVisitor) VisitStmtDefer(ctx *StmtDeferContext) interface{} {
	args := ctx.Arguments().Accept(v).([]ast.CallArgument)
	call := &ast.ExprCall{
		Pos:       getPos(v, ctx),
		Optional:  ctx.OPTIONAL_CALL() != nil,
		Callee:    ctx.Expr().Accept(v).(ast.Expr),
		Arguments: args,
		IsBind:    ast.IsBindList(args),
	}
	if ctx.DEFER() != nil {
		return &ast.StmtDefer{
			Call: call,
		}
	}
	if ctx.BLOCK_DEFER() != nil {
		return &ast.StmtBlockDefer{
			Call: call,
		}
	}
	panic("should not reach here")
}

func (v *ParseVisitor) VisitStmtDeferBlock(ctx *StmtDeferBlockContext) interface{} {
	f := &ast.ExprFunc{
		Value: runtime.NewFunc("", []string{}, false, ctx.CodeBlock().Accept(v).(ast.Node)),
	}
	call := &ast.ExprCall{
		Pos:       getPos(v, ctx),
		Optional:  false,
		Callee:    f,
		Arguments: nil,
	}
	if ctx.DEFER() != nil {
		return &ast.StmtDefer{
			Call: call,
		}
	}
	if ctx.BLOCK_DEFER() != nil {
		return &ast.StmtBlockDefer{
			Call: call,
		}
	}
	panic("should not reach here")
}

func (v *ParseVisitor) VisitStmtTry(ctx *StmtTryContext) interface{} {
	rv := &ast.StmtTry{
		Pos: getPos(v, ctx),
		Try: ctx.GetTryBlock().Accept(v).(*ast.Block),
	}
	if b := ctx.GetCatchBlock(); b != nil {
		rv.ExcName = ctx.IDENTIFIER().GetText()
		rv.Catch = b.Accept(v).(*ast.Block)
	}
	if b := ctx.GetFinallyBlock(); b != nil {
		rv.Finally = b.Accept(v).(*ast.Block)
	}
	return rv
}

func (v *ParseVisitor) VisitStmtAssert(ctx *StmtAssertContext) interface{} {
	exprs := ctx.AllExpr()
	s := &ast.StmtAssert{
		Pos:  getPos(v, ctx),
		Expr: exprs[0].Accept(v).(ast.Expr),
	}
	if len(exprs) > 1 {
		s.Message = exprs[1].Accept(v).(ast.Expr)
	} else {
		s.Message = &ast.ExprStr{Value: runtime.NewStr(exprs[0].GetText())}
	}
	return s
}

func (v *ParseVisitor) VisitStmtExtend(ctx *StmtExtendContext) interface{} {
	extType := ctx.Expr().Accept(v).(ast.Expr)
	allExts := ctx.AllKeyValue()
	names := make([]ast.Expr, 0, len(allExts))
	funcs := make([]ast.Expr, 0, len(allExts))
	for _, e := range allExts {
		kv := e.Accept(v).(kvPair)
		names = append(names, kv.key)
		funcs = append(funcs, kv.val)
	}
	return &ast.StmtExtend{
		Pos:      getPos(v, ctx),
		Exported: ctx.EXPORT() != nil,
		Type:     extType,
		Func:     funcs,
		Name:     names,
	}
}
