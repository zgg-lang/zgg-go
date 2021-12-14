package repl

import (
	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"
)

type ChannelReplContext struct {
	c      *runtime.Context
	Input  chan string
	Output chan string
}

func NewChannelReplContext(isDebug, canEval bool) *ChannelReplContext {
	c := runtime.NewContext(true, isDebug, canEval)
	c.ImportFunc = parser.SimpleImport
	return &ChannelReplContext{
		c:      c,
		Input:  make(chan string),
		Output: make(chan string),
	}
}

func (c *ChannelReplContext) Context() *runtime.Context {
	return c.c
}

func (c *ChannelReplContext) ReadCode(newCode bool, indent string) (string, bool) {
	code, ok := <-c.Input
	if !ok {
		return "", false
	}
	return indent + code, true
}

func (c *ChannelReplContext) WriteResult(result string) {
	c.Output <- result
}

func (ChannelReplContext) OnEnter() {
}

func (ChannelReplContext) OnExit() {
}
