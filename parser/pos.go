package parser

import (
	"github.com/zgg-lang/zgg-go/ast"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func getPos(v *ParseVisitor, c antlr.ParserRuleContext) ast.Pos {
	return ast.Pos{
		Line:     c.GetStart().GetLine(),
		FileName: v.FileName,
	}
}
