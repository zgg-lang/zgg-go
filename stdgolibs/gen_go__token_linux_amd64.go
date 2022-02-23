package stdgolibs

import (
	pkg "go/token"

	"reflect"
)

func init() {
	registerValues("go/token", map[string]reflect.Value{
		// Functions
		"Lookup":       reflect.ValueOf(pkg.Lookup),
		"IsExported":   reflect.ValueOf(pkg.IsExported),
		"IsKeyword":    reflect.ValueOf(pkg.IsKeyword),
		"IsIdentifier": reflect.ValueOf(pkg.IsIdentifier),
		"NewFileSet":   reflect.ValueOf(pkg.NewFileSet),

		// Consts

		"ILLEGAL":        reflect.ValueOf(pkg.ILLEGAL),
		"EOF":            reflect.ValueOf(pkg.EOF),
		"COMMENT":        reflect.ValueOf(pkg.COMMENT),
		"IDENT":          reflect.ValueOf(pkg.IDENT),
		"INT":            reflect.ValueOf(pkg.INT),
		"FLOAT":          reflect.ValueOf(pkg.FLOAT),
		"IMAG":           reflect.ValueOf(pkg.IMAG),
		"CHAR":           reflect.ValueOf(pkg.CHAR),
		"STRING":         reflect.ValueOf(pkg.STRING),
		"ADD":            reflect.ValueOf(pkg.ADD),
		"SUB":            reflect.ValueOf(pkg.SUB),
		"MUL":            reflect.ValueOf(pkg.MUL),
		"QUO":            reflect.ValueOf(pkg.QUO),
		"REM":            reflect.ValueOf(pkg.REM),
		"AND":            reflect.ValueOf(pkg.AND),
		"OR":             reflect.ValueOf(pkg.OR),
		"XOR":            reflect.ValueOf(pkg.XOR),
		"SHL":            reflect.ValueOf(pkg.SHL),
		"SHR":            reflect.ValueOf(pkg.SHR),
		"AND_NOT":        reflect.ValueOf(pkg.AND_NOT),
		"ADD_ASSIGN":     reflect.ValueOf(pkg.ADD_ASSIGN),
		"SUB_ASSIGN":     reflect.ValueOf(pkg.SUB_ASSIGN),
		"MUL_ASSIGN":     reflect.ValueOf(pkg.MUL_ASSIGN),
		"QUO_ASSIGN":     reflect.ValueOf(pkg.QUO_ASSIGN),
		"REM_ASSIGN":     reflect.ValueOf(pkg.REM_ASSIGN),
		"AND_ASSIGN":     reflect.ValueOf(pkg.AND_ASSIGN),
		"OR_ASSIGN":      reflect.ValueOf(pkg.OR_ASSIGN),
		"XOR_ASSIGN":     reflect.ValueOf(pkg.XOR_ASSIGN),
		"SHL_ASSIGN":     reflect.ValueOf(pkg.SHL_ASSIGN),
		"SHR_ASSIGN":     reflect.ValueOf(pkg.SHR_ASSIGN),
		"AND_NOT_ASSIGN": reflect.ValueOf(pkg.AND_NOT_ASSIGN),
		"LAND":           reflect.ValueOf(pkg.LAND),
		"LOR":            reflect.ValueOf(pkg.LOR),
		"ARROW":          reflect.ValueOf(pkg.ARROW),
		"INC":            reflect.ValueOf(pkg.INC),
		"DEC":            reflect.ValueOf(pkg.DEC),
		"EQL":            reflect.ValueOf(pkg.EQL),
		"LSS":            reflect.ValueOf(pkg.LSS),
		"GTR":            reflect.ValueOf(pkg.GTR),
		"ASSIGN":         reflect.ValueOf(pkg.ASSIGN),
		"NOT":            reflect.ValueOf(pkg.NOT),
		"NEQ":            reflect.ValueOf(pkg.NEQ),
		"LEQ":            reflect.ValueOf(pkg.LEQ),
		"GEQ":            reflect.ValueOf(pkg.GEQ),
		"DEFINE":         reflect.ValueOf(pkg.DEFINE),
		"ELLIPSIS":       reflect.ValueOf(pkg.ELLIPSIS),
		"LPAREN":         reflect.ValueOf(pkg.LPAREN),
		"LBRACK":         reflect.ValueOf(pkg.LBRACK),
		"LBRACE":         reflect.ValueOf(pkg.LBRACE),
		"COMMA":          reflect.ValueOf(pkg.COMMA),
		"PERIOD":         reflect.ValueOf(pkg.PERIOD),
		"RPAREN":         reflect.ValueOf(pkg.RPAREN),
		"RBRACK":         reflect.ValueOf(pkg.RBRACK),
		"RBRACE":         reflect.ValueOf(pkg.RBRACE),
		"SEMICOLON":      reflect.ValueOf(pkg.SEMICOLON),
		"COLON":          reflect.ValueOf(pkg.COLON),
		"BREAK":          reflect.ValueOf(pkg.BREAK),
		"CASE":           reflect.ValueOf(pkg.CASE),
		"CHAN":           reflect.ValueOf(pkg.CHAN),
		"CONST":          reflect.ValueOf(pkg.CONST),
		"CONTINUE":       reflect.ValueOf(pkg.CONTINUE),
		"DEFAULT":        reflect.ValueOf(pkg.DEFAULT),
		"DEFER":          reflect.ValueOf(pkg.DEFER),
		"ELSE":           reflect.ValueOf(pkg.ELSE),
		"FALLTHROUGH":    reflect.ValueOf(pkg.FALLTHROUGH),
		"FOR":            reflect.ValueOf(pkg.FOR),
		"FUNC":           reflect.ValueOf(pkg.FUNC),
		"GO":             reflect.ValueOf(pkg.GO),
		"GOTO":           reflect.ValueOf(pkg.GOTO),
		"IF":             reflect.ValueOf(pkg.IF),
		"IMPORT":         reflect.ValueOf(pkg.IMPORT),
		"INTERFACE":      reflect.ValueOf(pkg.INTERFACE),
		"MAP":            reflect.ValueOf(pkg.MAP),
		"PACKAGE":        reflect.ValueOf(pkg.PACKAGE),
		"RANGE":          reflect.ValueOf(pkg.RANGE),
		"RETURN":         reflect.ValueOf(pkg.RETURN),
		"SELECT":         reflect.ValueOf(pkg.SELECT),
		"STRUCT":         reflect.ValueOf(pkg.STRUCT),
		"SWITCH":         reflect.ValueOf(pkg.SWITCH),
		"TYPE":           reflect.ValueOf(pkg.TYPE),
		"VAR":            reflect.ValueOf(pkg.VAR),
		"LowestPrec":     reflect.ValueOf(pkg.LowestPrec),
		"UnaryPrec":      reflect.ValueOf(pkg.UnaryPrec),
		"HighestPrec":    reflect.ValueOf(pkg.HighestPrec),
		"NoPos":          reflect.ValueOf(pkg.NoPos),

		// Variables

	})
	registerTypes("go/token", map[string]reflect.Type{
		// Non interfaces

		"Token":    reflect.TypeOf((*pkg.Token)(nil)).Elem(),
		"Position": reflect.TypeOf((*pkg.Position)(nil)).Elem(),
		"Pos":      reflect.TypeOf((*pkg.Pos)(nil)).Elem(),
		"File":     reflect.TypeOf((*pkg.File)(nil)).Elem(),
		"FileSet":  reflect.TypeOf((*pkg.FileSet)(nil)).Elem(),
	})
}
