package repl

import (
	"fmt"

	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"

	"github.com/chzyer/readline"
)

type ConsoleReplContext struct {
	c         *runtime.Context
	readline  *readline.Instance
	completer *readline.PrefixCompleter
}

func NewConsoleReplContext(isDebug, canEval bool) *ConsoleReplContext {
	ac := readline.NewPrefixCompleter(
		readline.PcItem("func"),
		readline.PcItem("class"),
		readline.PcItem("if"),
		readline.PcItem("for"),
		readline.PcItem("... := import('"),
		readline.PcItem("export"),
	)
	r, _ := readline.NewEx(&readline.Config{
		Prompt:       "zgg> ",
		HistoryFile:  "/tmp/zgg_history",
		AutoComplete: ac,
	})
	c := runtime.NewContext(true, isDebug, canEval)
	c.ImportFunc = parser.SimpleImport
	c.AutoImport()
	return &ConsoleReplContext{
		c:         c,
		readline:  r,
		completer: ac,
	}
}

func (c *ConsoleReplContext) Context() *runtime.Context {
	return c.c
}

func (c *ConsoleReplContext) ReadCode(newCode bool, initCode string) (string, bool) {
	if newCode {
		c.readline.SetPrompt("zgg> ")
	} else {
		c.readline.SetPrompt(".... ")
	}
	if line, err := c.readline.ReadlineWithDefault(initCode); err == nil {
		return line, true
	} else {
		return "", false
	}
}

func (ConsoleReplContext) WriteResult(result string) {
	fmt.Println(result)
}

func (ConsoleReplContext) OnEnter() {
	fmt.Println("Welcome to Zgg REPL!")
}

func (ConsoleReplContext) OnExit() {
	fmt.Println("\nBye!")
}
