// Code generated from parser/ZggParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // ZggParser

import "github.com/antlr4-go/antlr/v4"

type BaseZggParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseZggParserVisitor) VisitReplExpr(ctx *ReplExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitReplBlock(ctx *ReplBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitModule(ctx *ModuleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitCodeBlock(ctx *CodeBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtBlock(ctx *StmtBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtPreIncDec(ctx *StmtPreIncDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtPostIncDec(ctx *StmtPostIncDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtAssign(ctx *StmtAssignContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtFuncCall(ctx *StmtFuncCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtFuncDefine(ctx *StmtFuncDefineContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtClassDefine(ctx *StmtClassDefineContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtFor(ctx *StmtForContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtForEach(ctx *StmtForEachContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtDoWhile(ctx *StmtDoWhileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtWhile(ctx *StmtWhileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtContinue(ctx *StmtContinueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtBreak(ctx *StmtBreakContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtIf(ctx *StmtIfContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtSwitch(ctx *StmtSwitchContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtReturnNone(ctx *StmtReturnNoneContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtReturn(ctx *StmtReturnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtExportIdentifier(ctx *StmtExportIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtExportExpr(ctx *StmtExportExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtExportFuncDefine(ctx *StmtExportFuncDefineContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtDefer(ctx *StmtDeferContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtDeferBlock(ctx *StmtDeferBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtTry(ctx *StmtTryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtAssert(ctx *StmtAssertContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStmtExtend(ctx *StmtExtendContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitIfCondition(ctx *IfConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitMemberDef(ctx *MemberDefContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitCallStmt(ctx *CallStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitSwitchDefault(ctx *SwitchDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitComparator(ctx *ComparatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprPlusMinus(ctx *ExprPlusMinusContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprUseCloser(ctx *ExprUseCloserContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprInRange(ctx *ExprInRangeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprAssign(ctx *ExprAssignContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprInContainer(ctx *ExprInContainerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprUseMethod(ctx *ExprUseMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprWhenValue(ctx *ExprWhenValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprLiteral(ctx *ExprLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprCompare(ctx *ExprCompareContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprByField(ctx *ExprByFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprLogicOr(ctx *ExprLogicOrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprBitXor(ctx *ExprBitXorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprPreIncDec(ctx *ExprPreIncDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprUseBlock(ctx *ExprUseBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprPow(ctx *ExprPowContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprWhen(ctx *ExprWhenContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprBitShift(ctx *ExprBitShiftContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprLogicNot(ctx *ExprLogicNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprItByField(ctx *ExprItByFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprLogicAnd(ctx *ExprLogicAndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprIdentifier(ctx *ExprIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprFallback(ctx *ExprFallbackContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprBitAnd(ctx *ExprBitAndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprNegative(ctx *ExprNegativeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprTimesDivMod(ctx *ExprTimesDivModContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprByIndex(ctx *ExprByIndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprBitNot(ctx *ExprBitNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprShortImport(ctx *ExprShortImportContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprSub(ctx *ExprSubContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprCall(ctx *ExprCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprBitOr(ctx *ExprBitOrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprAssertError(ctx *ExprAssertErrorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprQuestion(ctx *ExprQuestionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprPostIncDec(ctx *ExprPostIncDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprBySlice(ctx *ExprBySliceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitExprIsType(ctx *ExprIsTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitWhenConditionInList(ctx *WhenConditionInListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitWhenConditionInRange(ctx *WhenConditionInRangeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitWhenConditionIsType(ctx *WhenConditionIsTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitFuncArgument(ctx *FuncArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitAssignExists(ctx *AssignExistsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitAssignNew(ctx *AssignNewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitAssignNewDeArray(ctx *AssignNewDeArrayContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitAssignNewDeObject(ctx *AssignNewDeObjectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitAssignNewLocal(ctx *AssignNewLocalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitPreIncDec(ctx *PreIncDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitPostIncDec(ctx *PostIncDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLvalById(ctx *LvalByIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLvalByIndex(ctx *LvalByIndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLvalItByField(ctx *LvalItByFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLvalByField(ctx *LvalByFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitIntegerZero(ctx *IntegerZeroContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitIntegerDec(ctx *IntegerDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitIntegerHex(ctx *IntegerHexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitIntegerOct(ctx *IntegerOctContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitIntegerBin(ctx *IntegerBinContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralInteger(ctx *LiteralIntegerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralFloat(ctx *LiteralFloatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralENum(ctx *LiteralENumContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralBigNum(ctx *LiteralBigNumContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralBool(ctx *LiteralBoolContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralString(ctx *LiteralStringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralNil(ctx *LiteralNilContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralUndefined(ctx *LiteralUndefinedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralFunc(ctx *LiteralFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralLambdaExpr(ctx *LiteralLambdaExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralLambdaBlock(ctx *LiteralLambdaBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralObject(ctx *LiteralObjectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitObjectComprehension(ctx *ObjectComprehensionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitLiteralArray(ctx *LiteralArrayContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitArrayComprehension(ctx *ArrayComprehensionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitArrayItem(ctx *ArrayItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitObjItemKV(ctx *ObjItemKVContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitObjItemExpanded(ctx *ObjItemExpandedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitKVIdKey(ctx *KVIdKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitKVStrKey(ctx *KVStrKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitKVExprKey(ctx *KVExprKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitKVKeyFunc(ctx *KVKeyFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitKVIdOnly(ctx *KVIdOnlyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitKVExprOnly(ctx *KVExprOnlyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitTemplateString(ctx *TemplateStringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitTsRaw(ctx *TsRawContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitTsIdentifier(ctx *TsIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseZggParserVisitor) VisitTsExpr(ctx *TsExprContext) interface{} {
	return v.VisitChildren(ctx)
}
