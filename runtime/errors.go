package runtime

import (
	"fmt"
	"strings"
)

type Stack struct {
	FileName string
	Line     int
	Function string
}

type Exception interface {
	GetStack() []Stack
	GetMessage() string
	MessageWithStack() string
}

type RuntimeError struct {
	Stack   []Stack
	Message string
}

func (e *RuntimeError) Error() string {
	return e.GetMessage()
}

func (e *RuntimeError) GetStack() []Stack {
	return e.Stack
}

func (e *RuntimeError) GetMessage() string {
	return e.Message
}

func (e *RuntimeError) MessageWithStack() string {
	var builder strings.Builder
	builder.WriteString("Exception! " + e.Message + "\n")
	for _, s := range e.Stack {
		builder.WriteString(fmt.Sprintf("%s:%d (%s)\n", s.FileName, s.Line, s.Function))
	}
	return builder.String()
}

func ExceptionToValue(e Exception, c *Context) Value {
	if e == nil {
		return constNil
	}
	v := NewObject()
	v.SetMember("message", NewStr(e.GetMessage()), c)
	if re, ok := e.(*RuntimeError); ok {
		stack := NewArray(len(re.Stack))
		for _, s := range re.Stack {
			stack.PushBack(NewArrayByValues(
				NewStr(s.FileName),
				NewInt(int64(s.Line)),
				NewStr(s.Function),
			))
		}
		v.SetMember("stack", stack, c)
	}
	return v
}
