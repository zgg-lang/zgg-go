package repl

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"
)

type ReplContext interface {
	Context() *runtime.Context
	ReadCode(bool, string) (string, bool)
	WriteResult(string)
	OnEnter()
	OnExit()
}

func ReplLoop(context ReplContext, shouldRecover bool) {
	c := context.Context()
	context.OnEnter()
	code := ""
	continueInput := false
	indent := ""
	for {
		if !continueInput {
			code = ""
			indent = ""
		}
		inputCode, shouldRun := context.ReadCode(!continueInput, indent)
		if !shouldRun {
			break
		}
		code = code + "\n" + inputCode
		continueInput = false
		if strings.TrimSpace(code) == "" {
			continue
		}
		func() {
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
				e := errs[0]
				lines := strings.Split(code, "\n")
				if e.Line == len(lines) && e.Column == len(lines[len(lines)-1]) {
					s := 0
					inputRunes := []rune(inputCode)
					for ; s < len(inputCode) && unicode.IsSpace(inputRunes[s]); s++ {
					}
					indent = string(inputRunes[:s])
					if strings.HasSuffix(inputCode, "{") || strings.HasSuffix(inputCode, "(") || strings.HasSuffix(inputCode, "[") {
						indent += "  "
					}
					continueInput = true
				} else {
					context.WriteResult(errs[0].String())
				}
			} else if codeAst == nil {
				context.WriteResult("parse code fail")
			} else {
				codeAst.Eval(c)
				context.WriteResult(c.RetVal.ToString(c))
			}
		}()
	}
	context.OnExit()
}
