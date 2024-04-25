package repl

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"

	"github.com/chzyer/readline"
)

type ConsoleReplContext struct {
	c             *runtime.Context
	readline      *readline.Instance
	completer     *readline.PrefixCompleter
	shouldRecover bool
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

func (c *ConsoleReplContext) ReadAction(shouldRecover bool) ReplAction {
	var (
		indent = ""
		code   = ""
	)
	for {
		if code == "" {
			c.readline.SetPrompt("zgg> ")
		} else {
			c.readline.SetPrompt(".... ")
		}
		line, err := c.readline.ReadlineWithDefault(indent)
		if err != nil {
			if err == io.EOF {
				return ReplExit{}
			} else {
				return ReplRunCode{Err: fmt.Errorf("readline error: %w", err)}
			}
		}
		if code == "" {
			code = line
		} else {
			code += "\n" + line
		}
		if code == "exit" {
			return ReplExit{}
		}
		if strings.TrimSpace(code) == "" {
			return ReplNoop{}
		}
		compiled, err := ParseInputCode(code, shouldRecover)
		if err != nil {
			return ReplRunCode{Err: err}
		}
		if compiled != nil {
			return ReplRunCode{Compiled: compiled}
		}
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "(") || strings.HasSuffix(line, "[") {
			s := 0
			inputRunes := []rune(line)
			for ; s < len(line) && unicode.IsSpace(inputRunes[s]); s++ {
			}
			indent = string(inputRunes[:s]) + "  "
		}
	}
}

func (c ConsoleReplContext) WriteResult(result interface{}) {
	if result == nil {
		return
	}
	switch v := result.(type) {
	case runtime.Value:
		c.write(v.ToString(c.Context()))
	case string:
		c.write(v)
	default:
		c.write(fmt.Sprint(v))
	}
}

func (c ConsoleReplContext) WriteException(e runtime.Exception) {
	c.write(e.MessageWithStack())
}

func (c ConsoleReplContext) OnEnter() {
	c.write("Welcome to ZGG REPL!")
}

func (c ConsoleReplContext) OnExit() {
	c.write("\nBye!")
}
