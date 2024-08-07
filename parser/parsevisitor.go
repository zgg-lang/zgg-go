package parser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"time"

	"github.com/zgg-lang/zgg-go/stdgolibs"

	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/builtin_libs"

	"github.com/zgg-lang/zgg-go/runtime"

	"github.com/antlr4-go/antlr/v4"
)

type varRec struct {
	ref   map[string]bool
	local []map[string]bool
}

type ParseVisitor struct {
	BaseZggParserVisitor
	FileName string
}

func (v *ParseVisitor) Init() {
}

// func (v *ParseVisitor) VisitReplStmt(ctx *ReplStmtContext) interface{} {
// 	return ctx.Stmt().Accept(v)
// }

func (v *ParseVisitor) VisitReplBlock(ctx *ReplBlockContext) interface{} {
	return ctx.Block().Accept(v)
}

func (v *ParseVisitor) VisitReplExpr(ctx *ReplExprContext) interface{} {
	return ctx.Expr().Accept(v)
}

func (v *ParseVisitor) VisitModule(ctx *ModuleContext) interface{} {
	r := &ast.Module{Block: ctx.Block().Accept(v).(*ast.Block)}
	r.Block.Type = ast.BlockTypeModuleTop
	return r
}

func (v *ParseVisitor) VisitCodeBlock(ctx *CodeBlockContext) interface{} {
	return ctx.Block().Accept(v)
}

func (v *ParseVisitor) VisitBlock(ctx *BlockContext) interface{} {
	r := &ast.Block{Pos: getPos(v, ctx)}
	for _, s := range ctx.AllStmt() {
		stmt := s.Accept(v).(ast.Stmt)
		r.Stmts = append(r.Stmts, stmt)
	}
	return r
}

func (v *ParseVisitor) VisitExprCall(ctx *ExprCallContext) interface{} {
	args := ctx.Arguments().Accept(v).([]ast.CallArgument)
	return &ast.ExprCall{
		Pos:       getPos(v, ctx),
		Optional:  ctx.OPTIONAL_CALL() != nil,
		Callee:    ctx.Expr().Accept(v).(ast.Expr),
		Arguments: args,
		IsBind:    ast.IsBindList(args),
	}
}

func (v *ParseVisitor) VisitExprIdentifier(ctx *ExprIdentifierContext) interface{} {
	id := ctx.IDENTIFIER().GetText()
	switch id {
	case "__file__":
		return &ast.ExprStr{Value: runtime.NewStr(v.FileName)}
	case "__dir__":
		return &ast.ExprStr{Value: runtime.NewStr(filepath.Dir(v.FileName))}
	case "__line__":
		return &ast.ExprInt{Value: runtime.NewInt(int64(ctx.GetStart().GetLine()))}
	}
	return &ast.LvalById{Name: id}
}

func (v *ParseVisitor) VisitExprByField(ctx *ExprByFieldContext) interface{} {
	return &ast.LvalByField{
		Owner: ctx.Expr().Accept(v).(ast.Expr),
		Field: &ast.ExprStr{Value: runtime.NewStr(ctx.IDENTIFIER().GetText())},
	}
}

func (v *ParseVisitor) VisitExprItByField(ctx *ExprItByFieldContext) interface{} {
	return &ast.LvalByField{
		Owner: &ast.ExprIdentifier{Name: "it"},
		Field: &ast.ExprStr{Value: runtime.NewStr(ctx.IDENTIFIER().GetText())},
	}
}

func (v *ParseVisitor) VisitExprByIndex(ctx *ExprByIndexContext) interface{} {
	return &ast.LvalByField{
		Owner: ctx.Expr(0).Accept(v).(ast.Expr),
		Field: ctx.Expr(1).Accept(v).(ast.Expr),
	}
}

func (v *ParseVisitor) VisitExprBySlice(ctx *ExprBySliceContext) interface{} {
	r := &ast.ExprSlice{
		Container: ctx.GetContainer().Accept(v).(ast.Expr),
	}
	if b := ctx.GetBegin(); b != nil {
		r.Begin = b.Accept(v).(ast.Expr)
	}
	if e := ctx.GetEnd(); e != nil {
		r.End = e.Accept(v).(ast.Expr)
	}
	return r
}

func (v *ParseVisitor) VisitLvalById(ctx *LvalByIdContext) interface{} {
	id := ctx.IDENTIFIER().GetText()
	return &ast.LvalById{Name: id}
}

func (v *ParseVisitor) VisitLvalByField(ctx *LvalByFieldContext) interface{} {
	return &ast.LvalByField{
		Owner: ctx.Lval().Accept(v).(ast.Expr),
		Field: &ast.ExprStr{Value: runtime.NewStr(ctx.IDENTIFIER().GetText())},
	}
}

func (v *ParseVisitor) VisitLvalItByField(ctx *LvalItByFieldContext) interface{} {
	return &ast.LvalByField{
		Owner: &ast.LvalById{Name: "it"},
		Field: &ast.ExprStr{Value: runtime.NewStr(ctx.IDENTIFIER().GetText())},
	}
}

func (v *ParseVisitor) VisitLvalByIndex(ctx *LvalByIndexContext) interface{} {
	return &ast.LvalByField{
		Owner: ctx.Lval().Accept(v).(ast.Expr),
		Field: ctx.Expr().Accept(v).(ast.Expr),
	}
}

// func (v *ParseVisitor) VisitExprOnceCall(ctx *ExprOnceCallContext) interface{} {
// 	funcBody := ctx.Block().Accept(v).(*ast.Block)
// 	funcBody.Type = ast.BlockTypeFuncTop
// 	return &ast.ExprCall{
// 		Pos:      getPos(v, ctx),
// 		Optional: false,
// 		Callee: &ast.ExprFunc{
// 			Value: runtime.NewFunc("", []string{}, false, funcBody),
// 		},
// 		Arguments: []ast.CallArgument{},
// 	}
// }

func (v *ParseVisitor) VisitExprLiteral(ctx *ExprLiteralContext) interface{} {
	switch c := ctx.Literal().(type) {
	case *LiteralIntegerContext:
		return v.VisitLiteralInteger(c)
	case *LiteralFloatContext:
		return v.VisitLiteralFloat(c)
	case *LiteralENumContext:
		return v.VisitLiteralENum(c)
	case *LiteralBoolContext:
		return v.VisitLiteralBool(c)
	case *LiteralNilContext:
		return v.VisitLiteralNil(c)
	case *LiteralUndefinedContext:
		return v.VisitLiteralUndefined(c)
	case *LiteralStringContext:
		return v.VisitLiteralString(c)
	case *LiteralFuncContext:
		return v.VisitLiteralFunc(c)
	case *LiteralLambdaBlockContext:
		return v.VisitLiteralLambdaBlock(c)
	case *LiteralLambdaExprContext:
		return v.VisitLiteralLambdaExpr(c)
	case *LiteralObjectContext:
		return v.VisitLiteralObject(c)
	case *LiteralArrayContext:
		return v.VisitLiteralArray(c)
	case *LiteralBigNumContext:
		return v.VisitLiteralBigNum(c)
	case *ArrayComprehensionContext:
		return v.VisitArrayComprehension(c)
	case *ObjectComprehensionContext:
		return v.VisitObjectComprehension(c)
	}
	panic("should not reach this line")
}

func (v *ParseVisitor) VisitExprQuestion(ctx *ExprQuestionContext) interface{} {
	return &ast.ExprWhen{
		Cases: []ast.Case{ast.Case{
			Condition: ctx.GetCondition().Accept(v).(ast.Expr),
			Action:    ctx.GetTrueExpr().Accept(v).(ast.Expr),
		}},
		Else: ctx.GetFalseExpr().Accept(v).(ast.Expr),
	}
}

func safeParse(tree antlr.ParseTree, v *ParseVisitor, shouldRecover bool) (node ast.Node) {
	defer func() {
		if shouldRecover {
			if err := recover(); err != nil {
				node = nil
			}
		}
	}()
	node = tree.Accept(v).(ast.Node)
	return
}

func ParseFromString(filename, in string, shouldRecover bool) (ast.Node, []SyntaxErrorInfo) {
	if strings.HasPrefix(in, "#!") {
		pos := strings.Index(in, "\n")
		if pos < 0 {
			in = ""
		} else {
			in = in[pos+1:]
		}
	}
	ins := antlr.NewInputStream(in)
	lexer := NewZggLexer(ins)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := NewZggParser(stream)
	var errListener zggErrorListener
	errListener.FileName = filename
	p.RemoveErrorListeners()
	p.AddErrorListener(&errListener)
	var v ParseVisitor
	v.Init()
	v.FileName = filename
	treeRoot := p.Module()
	astNode := safeParse(treeRoot, &v, shouldRecover)
	return astNode, errListener.Errors
}

func ParseReplFromString(in string, shouldRecover bool) (ast.Node, []SyntaxErrorInfo) {
	ins := antlr.NewInputStream(in)
	lexer := NewZggLexer(ins)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := NewZggParser(stream)
	var errListener zggErrorListener
	errListener.FileName = "input"
	p.RemoveErrorListeners()
	p.AddErrorListener(&errListener)
	var v ParseVisitor
	v.Init()
	v.FileName = "-"
	treeRoot := p.ReplItem()
	astNode := safeParse(treeRoot, &v, shouldRecover)
	return astNode, errListener.Errors
}

func tryModulePath(prefix string) (string, bool) {
	if fi, err := os.Stat(prefix); err == nil {
		if fi.IsDir() {
			return tryModulePath(filepath.Join(prefix, "index"))
		}
		return prefix, true
	}
	if _, err := os.Stat(prefix + ".zgg"); err == nil {
		return prefix + ".zgg", true
	}
	if _, err := os.Stat(prefix + ".so"); err == nil {
		return prefix + ".so", true
	}
	return "", false
}

func GetModulePath(c *runtime.Context, name string) string {
	if strings.HasPrefix(name, ".") && c != nil {
		var root, curFile string
		if c != nil {
			root = c.Path
			curFile, _ = c.GetPosition()
			if curFile != "" {
				root = filepath.Dir(curFile)
			}
		}
		filename, _ := tryModulePath(filepath.Join(root, name))
		return filename
	}
	if filepath.IsAbs(name) {
		f, _ := tryModulePath(name)
		return f
	}
	var roots []string
	if c != nil {
		roots = c.ImportPaths
	} else {
		roots = runtime.GetImportPaths()
	}
	for _, root := range roots {
		testname := filepath.Join(root, name)
		if f, ok := tryModulePath(testname); ok {
			return f
		}
	}
	return ""
}

var debugLogger = log.New(os.Stderr, "[DBG]", log.Ltime|log.Lshortfile)

func debugTrace(beginTime time.Time) {
	gap := time.Since(beginTime)
	debugLogger.Output(2, "TRACE"+gap.String())
}

func SimpleImport(c *runtime.Context, name string, code string, importType string, reloadIfNewer bool) (modVal runtime.Value, thisTime int64, success bool) {
	if name == "" {
		node, errs := ParseReplFromString(code, true)
		if len(errs) > 0 || node == nil {
			c.RaiseRuntimeError("parse code %s fail", code)
			return
		}
		node.Eval(c)
		return c.RetVal, 0, true
	}
	if lib, found := builtin_libs.FindLib(c, name); found {
		return lib, 0, true
	}
	if strings.HasPrefix(name, "gostd/") {
		goName := name[6:]
		if lib, found := stdgolibs.FindLib(c, goName); found {
			return lib, 0, true
		}
		c.RaiseRuntimeError("import: cannot find module file %s", name)
		return
	}
	filename := GetModulePath(c, name)
	if filename == "" {
		c.RaiseRuntimeError("import: cannot find module file %s", name)
		return
	}
	fi, err := os.Stat(filename)
	if err != nil {
		c.RaiseRuntimeError("import: stat file %s err %s", name, err)
		return
	}
	var lastTime int64
	modVal, lastTime = c.GetModule(filename)
	if thisTime = fi.ModTime().UnixNano(); thisTime == lastTime || (lastTime != 0 && !reloadIfNewer) {
		success = true
		return
	}
	modVal = runtime.Undefined()
	success = false
	defer func() {
		if success {
			c.AddModule(filename, modVal, thisTime)
		}
	}()
	if strings.ToLower(filepath.Ext(filename)) == ".so" {
		p, err := plugin.Open(filename)
		if err != nil {
			c.RaiseRuntimeError("import: load %s in %s error %s", name, filename, err)
			return
		}
		s, err := p.Lookup("New")
		if err != nil {
			c.RaiseRuntimeError("import: load %s find entry error %s", name, err)
			return
		}
		switch newFn := s.(type) {
		case func(*runtime.Context) runtime.Value:
			modVal, thisTime, success = newFn(c), 0, true
			return
		case func() map[string]interface{}:
			{
				v := runtime.NewObject()
				for name, val := range newFn() {
					if rt, ok := val.(reflect.Type); ok {
						v.SetMember(name, runtime.NewGoType(rt), c)
					} else if rv, ok := val.(reflect.Value); ok {
						v.SetMember(name, runtime.NewReflectedGoValue(rv), c)
					} else {
						v.SetMember(name, runtime.NewGoValue(val), c)
					}
				}
				if initSymbol, err := p.Lookup("InitScript"); err == nil {
					if initScript, ok := initSymbol.(*string); ok {
						modAst, errs := ParseFromString(filename, *initScript, true)
						if len(errs) > 0 || modAst == nil {
							c.RaiseRuntimeError("parse module %s initScript fail", name)
							return
						}
						modC := c.Clone()
						modC.Path = filepath.Dir(filename)
						modC.SetLocalValue("_native", v)
						modAst.Eval(modC)
						modVal, success = modC.RetVal, true
						return
					}
				}
				modVal, thisTime, success = v, 0, true
				return
			}
		default:
			c.RaiseRuntimeError("import: load %s find entry error", name)
		}
	}
	codeBs, err := os.ReadFile(filename)
	if err != nil {
		c.RaiseRuntimeError("import: read file %s err %s", filename, err)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			if runtimeErr, ok := err.(*runtime.RuntimeError); ok {
				panic(runtimeErr)
			}
			return
		}
	}()
	csvSplitter := ','
	switch importType {
	case runtime.ImportTypeScript:
		modAst, errs := ParseFromString(filename, string(codeBs), true)
		if len(errs) > 0 || modAst == nil {
			c.RaiseRuntimeError("parse module %s fail", name)
			return
		}
		modC := c.Clone()
		modC.Path = filepath.Dir(filename)
		modAst.Eval(modC)
		modVal, success = modC.RetVal, true
		return
	case runtime.ImportTypeText:
		modVal, success = runtime.NewStr(string(codeBs)), true
		return
	case runtime.ImportTypeBytes:
		modVal, success = runtime.NewBytes(codeBs), true
		return
	case runtime.ImportTypeCsvByTab:
		csvSplitter = '\t'
		fallthrough
	case runtime.ImportTypeCsvByComma:
		fallthrough
	case runtime.ImportTypeCsv:
		{
			rd := csv.NewReader(bytes.NewReader(codeBs))
			rd.Comma = csvSplitter
			all, err := rd.ReadAll()
			if err != nil {
				c.RaiseRuntimeError("read csv %s fail: %s", name, err)
			}
			rows := runtime.NewArray(len(all))
			for _, row := range all {
				rowItem := runtime.NewArray(len(row))
				for _, cell := range row {
					rowItem.PushBack(runtime.NewStr(cell))
				}
				rows.PushBack(rowItem)
			}
			modVal, success = rows, true
			return
		}
	case runtime.ImportTypeJson:
		{
			var r interface{}
			if err := json.Unmarshal(codeBs, &r); err != nil {
				c.RaiseRuntimeError("read json %s fail: %s", name, err)
			}
			modVal, success = runtime.FromGoValue(reflect.ValueOf(r), c), true
			return
		}
	}
	return
}
