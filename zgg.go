package zgg

import (
	"errors"
	"fmt"

	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"
)

type Runner struct {
	context  *runtime.Context
	filename string
}

func NewRunner() *Runner {
	context := runtime.NewContext(true, false, false)
	context.ImportFunc = parser.SimpleImport
	return &Runner{
		context: context,
	}
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
	r.context.SetLocalValue(name, runtime.NewGoValue(value))
	return r
}

func (r *Runner) Run(code string) (err error) {
	astNode, syntaxErrors := parser.ParseFromString(r.filename, code, r.context.IsDebug)
	if len(syntaxErrors) > 0 {
		err = errors.New(syntaxErrors[0].String())
		return
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
	astNode.Eval(r.context)
	return
}

func (r *Runner) Eval(expr string) (val interface{}, err error) {
	astNode, syntaxErrors := parser.ParseReplFromString(expr, r.context.IsDebug)
	if len(syntaxErrors) > 0 {
		err = errors.New(syntaxErrors[0].String())
		return
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
	astNode.Eval(r.context)
	val = r.context.RetVal.ToGoValue()
	return
}

func RunCode(code string) error {
	return NewRunner().Run(code)
}

func Eval(expr string) (interface{}, error) {
	return NewRunner().Eval(expr)
}
