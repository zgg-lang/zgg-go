// Code generated from parser/ZggParser.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // ZggParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ZggParser.
type ZggParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ZggParser#replExpr.
	VisitReplExpr(ctx *ReplExprContext) interface{}

	// Visit a parse tree produced by ZggParser#replBlock.
	VisitReplBlock(ctx *ReplBlockContext) interface{}

	// Visit a parse tree produced by ZggParser#module.
	VisitModule(ctx *ModuleContext) interface{}

	// Visit a parse tree produced by ZggParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by ZggParser#codeBlock.
	VisitCodeBlock(ctx *CodeBlockContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtBlock.
	VisitStmtBlock(ctx *StmtBlockContext) interface{}

	// Visit a parse tree produced by ZggParser#StmtPreIncDec.
	VisitStmtPreIncDec(ctx *StmtPreIncDecContext) interface{}

	// Visit a parse tree produced by ZggParser#StmtPostIncDec.
	VisitStmtPostIncDec(ctx *StmtPostIncDecContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtAssign.
	VisitStmtAssign(ctx *StmtAssignContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtFuncCall.
	VisitStmtFuncCall(ctx *StmtFuncCallContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtFuncDefine.
	VisitStmtFuncDefine(ctx *StmtFuncDefineContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtClassDefine.
	VisitStmtClassDefine(ctx *StmtClassDefineContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtFor.
	VisitStmtFor(ctx *StmtForContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtForEach.
	VisitStmtForEach(ctx *StmtForEachContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtDoWhile.
	VisitStmtDoWhile(ctx *StmtDoWhileContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtWhile.
	VisitStmtWhile(ctx *StmtWhileContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtContinue.
	VisitStmtContinue(ctx *StmtContinueContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtBreak.
	VisitStmtBreak(ctx *StmtBreakContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtIf.
	VisitStmtIf(ctx *StmtIfContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtSwitch.
	VisitStmtSwitch(ctx *StmtSwitchContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtReturnNone.
	VisitStmtReturnNone(ctx *StmtReturnNoneContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtReturn.
	VisitStmtReturn(ctx *StmtReturnContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtExportIdentifier.
	VisitStmtExportIdentifier(ctx *StmtExportIdentifierContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtExportExpr.
	VisitStmtExportExpr(ctx *StmtExportExprContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtExportFuncDefine.
	VisitStmtExportFuncDefine(ctx *StmtExportFuncDefineContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtDefer.
	VisitStmtDefer(ctx *StmtDeferContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtDeferBlock.
	VisitStmtDeferBlock(ctx *StmtDeferBlockContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtTry.
	VisitStmtTry(ctx *StmtTryContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtAssert.
	VisitStmtAssert(ctx *StmtAssertContext) interface{}

	// Visit a parse tree produced by ZggParser#stmtExtend.
	VisitStmtExtend(ctx *StmtExtendContext) interface{}

	// Visit a parse tree produced by ZggParser#ifCondition.
	VisitIfCondition(ctx *IfConditionContext) interface{}

	// Visit a parse tree produced by ZggParser#memberDef.
	VisitMemberDef(ctx *MemberDefContext) interface{}

	// Visit a parse tree produced by ZggParser#callStmt.
	VisitCallStmt(ctx *CallStmtContext) interface{}

	// Visit a parse tree produced by ZggParser#switchCase.
	VisitSwitchCase(ctx *SwitchCaseContext) interface{}

	// Visit a parse tree produced by ZggParser#switchDefault.
	VisitSwitchDefault(ctx *SwitchDefaultContext) interface{}

	// Visit a parse tree produced by ZggParser#comparator.
	VisitComparator(ctx *ComparatorContext) interface{}

	// Visit a parse tree produced by ZggParser#exprPlusMinus.
	VisitExprPlusMinus(ctx *ExprPlusMinusContext) interface{}

	// Visit a parse tree produced by ZggParser#exprUseCloser.
	VisitExprUseCloser(ctx *ExprUseCloserContext) interface{}

	// Visit a parse tree produced by ZggParser#exprInRange.
	VisitExprInRange(ctx *ExprInRangeContext) interface{}

	// Visit a parse tree produced by ZggParser#exprAssign.
	VisitExprAssign(ctx *ExprAssignContext) interface{}

	// Visit a parse tree produced by ZggParser#exprInContainer.
	VisitExprInContainer(ctx *ExprInContainerContext) interface{}

	// Visit a parse tree produced by ZggParser#exprUseMethod.
	VisitExprUseMethod(ctx *ExprUseMethodContext) interface{}

	// Visit a parse tree produced by ZggParser#exprWhenValue.
	VisitExprWhenValue(ctx *ExprWhenValueContext) interface{}

	// Visit a parse tree produced by ZggParser#exprLiteral.
	VisitExprLiteral(ctx *ExprLiteralContext) interface{}

	// Visit a parse tree produced by ZggParser#exprCompare.
	VisitExprCompare(ctx *ExprCompareContext) interface{}

	// Visit a parse tree produced by ZggParser#exprByField.
	VisitExprByField(ctx *ExprByFieldContext) interface{}

	// Visit a parse tree produced by ZggParser#exprLogicOr.
	VisitExprLogicOr(ctx *ExprLogicOrContext) interface{}

	// Visit a parse tree produced by ZggParser#exprBitXor.
	VisitExprBitXor(ctx *ExprBitXorContext) interface{}

	// Visit a parse tree produced by ZggParser#exprPreIncDec.
	VisitExprPreIncDec(ctx *ExprPreIncDecContext) interface{}

	// Visit a parse tree produced by ZggParser#exprUseBlock.
	VisitExprUseBlock(ctx *ExprUseBlockContext) interface{}

	// Visit a parse tree produced by ZggParser#exprPow.
	VisitExprPow(ctx *ExprPowContext) interface{}

	// Visit a parse tree produced by ZggParser#exprWhen.
	VisitExprWhen(ctx *ExprWhenContext) interface{}

	// Visit a parse tree produced by ZggParser#exprBitShift.
	VisitExprBitShift(ctx *ExprBitShiftContext) interface{}

	// Visit a parse tree produced by ZggParser#exprLogicNot.
	VisitExprLogicNot(ctx *ExprLogicNotContext) interface{}

	// Visit a parse tree produced by ZggParser#exprLogicAnd.
	VisitExprLogicAnd(ctx *ExprLogicAndContext) interface{}

	// Visit a parse tree produced by ZggParser#exprIdentifier.
	VisitExprIdentifier(ctx *ExprIdentifierContext) interface{}

	// Visit a parse tree produced by ZggParser#exprFallback.
	VisitExprFallback(ctx *ExprFallbackContext) interface{}

	// Visit a parse tree produced by ZggParser#exprBitAnd.
	VisitExprBitAnd(ctx *ExprBitAndContext) interface{}

	// Visit a parse tree produced by ZggParser#exprNegative.
	VisitExprNegative(ctx *ExprNegativeContext) interface{}

	// Visit a parse tree produced by ZggParser#exprTimesDivMod.
	VisitExprTimesDivMod(ctx *ExprTimesDivModContext) interface{}

	// Visit a parse tree produced by ZggParser#exprByIndex.
	VisitExprByIndex(ctx *ExprByIndexContext) interface{}

	// Visit a parse tree produced by ZggParser#exprBitNot.
	VisitExprBitNot(ctx *ExprBitNotContext) interface{}

	// Visit a parse tree produced by ZggParser#exprShortImport.
	VisitExprShortImport(ctx *ExprShortImportContext) interface{}

	// Visit a parse tree produced by ZggParser#exprSub.
	VisitExprSub(ctx *ExprSubContext) interface{}

	// Visit a parse tree produced by ZggParser#exprCall.
	VisitExprCall(ctx *ExprCallContext) interface{}

	// Visit a parse tree produced by ZggParser#exprBitOr.
	VisitExprBitOr(ctx *ExprBitOrContext) interface{}

	// Visit a parse tree produced by ZggParser#exprAssertError.
	VisitExprAssertError(ctx *ExprAssertErrorContext) interface{}

	// Visit a parse tree produced by ZggParser#exprQuestion.
	VisitExprQuestion(ctx *ExprQuestionContext) interface{}

	// Visit a parse tree produced by ZggParser#exprPostIncDec.
	VisitExprPostIncDec(ctx *ExprPostIncDecContext) interface{}

	// Visit a parse tree produced by ZggParser#exprBySlice.
	VisitExprBySlice(ctx *ExprBySliceContext) interface{}

	// Visit a parse tree produced by ZggParser#exprIsType.
	VisitExprIsType(ctx *ExprIsTypeContext) interface{}

	// Visit a parse tree produced by ZggParser#whenConditionInList.
	VisitWhenConditionInList(ctx *WhenConditionInListContext) interface{}

	// Visit a parse tree produced by ZggParser#whenConditionInRange.
	VisitWhenConditionInRange(ctx *WhenConditionInRangeContext) interface{}

	// Visit a parse tree produced by ZggParser#whenConditionIsType.
	VisitWhenConditionIsType(ctx *WhenConditionIsTypeContext) interface{}

	// Visit a parse tree produced by ZggParser#arguments.
	VisitArguments(ctx *ArgumentsContext) interface{}

	// Visit a parse tree produced by ZggParser#funcArgument.
	VisitFuncArgument(ctx *FuncArgumentContext) interface{}

	// Visit a parse tree produced by ZggParser#assignExists.
	VisitAssignExists(ctx *AssignExistsContext) interface{}

	// Visit a parse tree produced by ZggParser#assignNew.
	VisitAssignNew(ctx *AssignNewContext) interface{}

	// Visit a parse tree produced by ZggParser#assignNewDeArray.
	VisitAssignNewDeArray(ctx *AssignNewDeArrayContext) interface{}

	// Visit a parse tree produced by ZggParser#assignNewDeObject.
	VisitAssignNewDeObject(ctx *AssignNewDeObjectContext) interface{}

	// Visit a parse tree produced by ZggParser#assignNewLocal.
	VisitAssignNewLocal(ctx *AssignNewLocalContext) interface{}

	// Visit a parse tree produced by ZggParser#preIncDec.
	VisitPreIncDec(ctx *PreIncDecContext) interface{}

	// Visit a parse tree produced by ZggParser#postIncDec.
	VisitPostIncDec(ctx *PostIncDecContext) interface{}

	// Visit a parse tree produced by ZggParser#lvalById.
	VisitLvalById(ctx *LvalByIdContext) interface{}

	// Visit a parse tree produced by ZggParser#lvalByIndex.
	VisitLvalByIndex(ctx *LvalByIndexContext) interface{}

	// Visit a parse tree produced by ZggParser#lvalByField.
	VisitLvalByField(ctx *LvalByFieldContext) interface{}

	// Visit a parse tree produced by ZggParser#IntegerZero.
	VisitIntegerZero(ctx *IntegerZeroContext) interface{}

	// Visit a parse tree produced by ZggParser#IntegerDec.
	VisitIntegerDec(ctx *IntegerDecContext) interface{}

	// Visit a parse tree produced by ZggParser#IntegerHex.
	VisitIntegerHex(ctx *IntegerHexContext) interface{}

	// Visit a parse tree produced by ZggParser#IntegerOct.
	VisitIntegerOct(ctx *IntegerOctContext) interface{}

	// Visit a parse tree produced by ZggParser#IntegerBin.
	VisitIntegerBin(ctx *IntegerBinContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralInteger.
	VisitLiteralInteger(ctx *LiteralIntegerContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralFloat.
	VisitLiteralFloat(ctx *LiteralFloatContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralBigNum.
	VisitLiteralBigNum(ctx *LiteralBigNumContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralBool.
	VisitLiteralBool(ctx *LiteralBoolContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralString.
	VisitLiteralString(ctx *LiteralStringContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralNil.
	VisitLiteralNil(ctx *LiteralNilContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralUndefined.
	VisitLiteralUndefined(ctx *LiteralUndefinedContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralFunc.
	VisitLiteralFunc(ctx *LiteralFuncContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralLambdaExpr.
	VisitLiteralLambdaExpr(ctx *LiteralLambdaExprContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralLambdaBlock.
	VisitLiteralLambdaBlock(ctx *LiteralLambdaBlockContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralObject.
	VisitLiteralObject(ctx *LiteralObjectContext) interface{}

	// Visit a parse tree produced by ZggParser#ObjectComprehension.
	VisitObjectComprehension(ctx *ObjectComprehensionContext) interface{}

	// Visit a parse tree produced by ZggParser#LiteralArray.
	VisitLiteralArray(ctx *LiteralArrayContext) interface{}

	// Visit a parse tree produced by ZggParser#ArrayComprehension.
	VisitArrayComprehension(ctx *ArrayComprehensionContext) interface{}

	// Visit a parse tree produced by ZggParser#arrayItem.
	VisitArrayItem(ctx *ArrayItemContext) interface{}

	// Visit a parse tree produced by ZggParser#ObjItemKV.
	VisitObjItemKV(ctx *ObjItemKVContext) interface{}

	// Visit a parse tree produced by ZggParser#ObjItemExpanded.
	VisitObjItemExpanded(ctx *ObjItemExpandedContext) interface{}

	// Visit a parse tree produced by ZggParser#KVIdKey.
	VisitKVIdKey(ctx *KVIdKeyContext) interface{}

	// Visit a parse tree produced by ZggParser#KVStrKey.
	VisitKVStrKey(ctx *KVStrKeyContext) interface{}

	// Visit a parse tree produced by ZggParser#KVExprKey.
	VisitKVExprKey(ctx *KVExprKeyContext) interface{}

	// Visit a parse tree produced by ZggParser#KVKeyFunc.
	VisitKVKeyFunc(ctx *KVKeyFuncContext) interface{}

	// Visit a parse tree produced by ZggParser#KVIdOnly.
	VisitKVIdOnly(ctx *KVIdOnlyContext) interface{}

	// Visit a parse tree produced by ZggParser#KVExprOnly.
	VisitKVExprOnly(ctx *KVExprOnlyContext) interface{}

	// Visit a parse tree produced by ZggParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by ZggParser#templateString.
	VisitTemplateString(ctx *TemplateStringContext) interface{}

	// Visit a parse tree produced by ZggParser#tsRaw.
	VisitTsRaw(ctx *TsRawContext) interface{}

	// Visit a parse tree produced by ZggParser#tsIdentifier.
	VisitTsIdentifier(ctx *TsIdentifierContext) interface{}

	// Visit a parse tree produced by ZggParser#tsExpr.
	VisitTsExpr(ctx *TsExprContext) interface{}
}
