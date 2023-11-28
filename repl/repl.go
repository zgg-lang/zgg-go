package repl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zgg-lang/zgg-go/ast"
	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"
)

type (
	ReplAction interface {
		Handle(context ReplContext, shouldRecover bool) (continueRunning bool)
	}

	ReplRunCode struct {
		Err      error
		Compiled runtime.IEval
	}

	ReplHintCode string

	ReplExit struct{}
	ReplNoop struct{}

	ReplContext interface {
		Context() *runtime.Context
		ReadAction(bool) ReplAction
		WriteResult(interface{})
		WriteException(runtime.Exception)
		OnEnter()
		OnExit()
	}

	ReplContextWithShouldWriteResult interface {
		ShouldWriteResult(codeAst ast.Node) bool
	}
)

func shouldWriteResult(c ReplContext, codeAst ast.Node) bool {
	if swr, is := c.(ReplContextWithShouldWriteResult); is {
		return swr.ShouldWriteResult(codeAst)
	}
	if _, is := codeAst.(ast.Expr); !is {
		return false
	}
	if _, is := codeAst.(ast.IsAssign); is {
		return false
	}
	return true
}

func ParseInputCode(code string, shouldRecover bool) (compiled runtime.IEval, err error) {
	codeAst, errs := parser.ParseReplFromString(code, shouldRecover)
	if len(errs) > 0 {
		e := errs[0]
		lines := strings.Split(code, "\n")
		if e.Line != len(lines) || e.Column != len(lines[len(lines)-1]) {
			err = &e
		}
	} else if codeAst == nil {
		err = errors.New("parse code fail")
	} else {
		compiled = codeAst
	}
	return
}

func ReplLoop(context ReplContext, shouldRecover bool) {
	context.OnEnter()
	for {
		action := context.ReadAction(shouldRecover)
		if action == nil || !action.Handle(context, shouldRecover) {
			break
		}
	}
	context.OnExit()
}

func (rrc ReplRunCode) Handle(context ReplContext, shouldRecover bool) (shouldContinue bool) {
	shouldContinue = true
	if rrc.Err != nil {
		context.WriteResult(rrc.Err)
		return
	}
	codeAst := rrc.Compiled
	if codeAst == nil {
		return
	}
	c := context.Context()
	defer func() {
		if shouldRecover {
			if err := recover(); err != nil {
				if exc, ok := err.(runtime.Exception); ok {
					context.WriteException(exc)
				} else {
					context.WriteResult(fmt.Sprintf("ERR! %s", err))
				}
			}
		}
	}()
	codeAst.Eval(c)
	retVal := c.RetVal
	if shouldWriteResult(context, codeAst) {
		context.WriteResult(retVal)
	} else {
		context.WriteResult(nil)
	}
	c.ForceSetLocalValue("__last__", retVal)
	return
}

func (rrc ReplHintCode) Handle(context ReplContext, shouldRecover bool) bool {
	c := context.Context()
	code := string(rrc)
	defer func() {
		if shouldRecover {
			if err := recover(); err != nil {
				if exc, ok := err.(runtime.Exception); ok {
					context.WriteResult(exc.MessageWithStack())
				} else {
					context.WriteResult(fmt.Sprintf("ERR! %s", err))
				}
			}
		}
	}()
	codeAst, errs := parser.ParseReplFromString(code, shouldRecover)
	if len(errs) > 0 {
		context.WriteResult(errs[0].String())
	} else if codeAst == nil {
		context.WriteResult("parse code fail")
	} else {
		c.EvalConst(codeAst)
		context.WriteResult(c.RetVal)
	}
	return true
}

func (ReplExit) Handle(context ReplContext, shouldRecover bool) bool {
	return false
}

func (ReplNoop) Handle(context ReplContext, shouldRecover bool) bool {
	return true
}
