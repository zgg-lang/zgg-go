package parser

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type SyntaxErrorInfo struct {
	FileName string
	Line     int
	Column   int
	Msg      string
}

func (e *SyntaxErrorInfo) String() string {
	return fmt.Sprintf("%s: line %d:%d: %s", e.FileName, e.Line, e.Column, e.Msg)
}

func (e *SyntaxErrorInfo) Error() string {
	return e.String()
}

type zggErrorListener struct {
	antlr.DefaultErrorListener
	FileName string
	Errors   []SyntaxErrorInfo
}

func (l *zggErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.Errors = append(l.Errors, SyntaxErrorInfo{
		FileName: l.FileName,
		Line:     line,
		Column:   column,
		Msg:      msg,
	})
}
