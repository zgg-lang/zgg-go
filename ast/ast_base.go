package ast

import (
	"fmt"

	"github.com/zgg-lang/zgg-go/runtime"
)

type Node interface {
	runtime.IEval
}

type Position interface {
	Position() (fileName string, lineNum int)
}

type Pos struct {
	FileName string
	Line     int
}

func (p *Pos) Position() (string, int) {
	return p.FileName, p.Line
}

func (p *Pos) PositionStr() string {
	return fmt.Sprintf("%s:%d", p.FileName, p.Line)
}

func evalAndReturn(c *runtime.Context, n Node) runtime.Value {
	n.Eval(c)
	return c.RetVal
}
