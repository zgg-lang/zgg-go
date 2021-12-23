package zgg

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/zgg-lang/zgg-go/ast"
	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"
)

type (
	Runner struct {
		context  *runtime.Context
		filename string
	}
	compileFunc = func(string) (ast.Node, []parser.SyntaxErrorInfo)
	ExecOption  interface {
		Apply(*Runner)
	}
	Var struct {
		Name  string
		Value interface{}
	}
	Val interface{}
)

func NewRunner() *Runner {
	context := runtime.NewContext(true, false, false)
	context.ImportFunc = parser.SimpleImport
	return &Runner{
		context: context,
	}
}

func (r *Runner) Reset() {
	r.context.Reset()
	r.filename = ""
}

func (r *Runner) IsDebug(isDebug bool) *Runner {
	r.context.IsDebug = isDebug
	return r
}

func (r *Runner) CanEval(canEval bool) *Runner {
	r.context.CanEval = canEval
	return r
}

func (r *Runner) Filename(filename string) *Runner {
	r.filename = filename
	return r
}

func (r *Runner) Workdir(wd string) *Runner {
	r.context.Path = wd
	return r
}

func (r *Runner) Args(args ...string) *Runner {
	r.context.Args = args
	return r
}

func (r *Runner) Var(name string, value interface{}) *Runner {
	if zggval, ok := value.(Val); ok {
		r.context.ForceSetLocalValue(name, runtime.FromGoValue(reflect.ValueOf(zggval), r.context))
	} else {
		r.context.ForceSetLocalValue(name, runtime.NewGoValue(value))
	}
	return r
}

func (r *Runner) Stdout(w io.Writer) *Runner {
	r.context.Stdout = w
	return r
}

func (r *Runner) Stderr(w io.Writer) *Runner {
	r.context.Stderr = w
	return r
}

func (r *Runner) compileCode(code string) (ast.Node, []parser.SyntaxErrorInfo) {
	return parser.ParseFromString(r.filename, code, !r.context.IsDebug)
}

func (r *Runner) compileExpr(expr string) (ast.Node, []parser.SyntaxErrorInfo) {
	return parser.ParseReplFromString(expr, !r.context.IsDebug)
}

func (r *Runner) compile(code interface{}, compileFunc compileFunc) (ast.Node, error) {
	var codeNode ast.Node
	switch codeVal := code.(type) {
	case ast.Node:
		codeNode = codeVal
	case string:
		astNode, syntaxErrors := compileFunc(codeVal)
		if len(syntaxErrors) > 0 {
			return nil, errors.New(syntaxErrors[0].String())
		}
		codeNode = astNode
	default:
		return nil, errors.New("invalid code type")
	}
	return codeNode, nil
}

func (r *Runner) CompileCode(code interface{}) (ast.Node, error) {
	return r.compile(code, r.compileCode)
}

func (r *Runner) CompileExpr(code interface{}) (ast.Node, error) {
	return r.compile(code, r.compileExpr)
}

func (r *Runner) execute(code interface{}, compileFunc compileFunc) (rv interface{}, err error) {
	codeNode, err := r.compile(code, compileFunc)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := recover()
		if e != nil {
			switch ee := e.(type) {
			case error:
				err = ee
			default:
				err = errors.New(fmt.Sprint(ee))
			}
		}
	}()
	codeNode.Eval(r.context)
	rv = r.context.RetVal.ToGoValue()
	return
}

func (v Var) Apply(runner *Runner) {
	runner.Var(v.Name, v.Value)
}

func (r *Runner) Run(code interface{}, opts ...ExecOption) (interface{}, error) {
	return r.execute(code, r.compileCode)
}

func (r *Runner) Eval(expr interface{}) (interface{}, error) {
	return r.execute(expr, r.compileExpr)
}

var runnerPool = sync.Pool{
	New: func() interface{} {
		return NewRunner()
	},
}

func RunCode(code interface{}, opts ...ExecOption) (interface{}, error) {
	runner := runnerPool.Get().(*Runner)
	defer runnerPool.Put(runner)
	runner.Reset()
	for _, opt := range opts {
		opt.Apply(runner)
	}
	return runner.Run(code)
}

func Eval(expr interface{}, opts ...ExecOption) (interface{}, error) {
	runner := runnerPool.Get().(*Runner)
	defer runnerPool.Put(runner)
	runner.Reset()
	for _, opt := range opts {
		opt.Apply(runner)
	}
	return runner.Eval(expr)
}
