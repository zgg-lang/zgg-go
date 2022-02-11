package parser

import (
	"path/filepath"

	"github.com/zgg-lang/zgg-go/ast"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func getPos(v *ParseVisitor, c antlr.ParserRuleContext) ast.Pos {
	filename := v.FileName
	if n, err := filepath.Abs(filename); err == nil {
		filename = n
	}
	return ast.Pos{
		Line:     c.GetStart().GetLine(),
		FileName: filename,
	}
}
