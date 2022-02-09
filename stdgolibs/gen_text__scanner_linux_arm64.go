package stdgolibs

import (
	pkg "text/scanner"

	"reflect"
)

func init() {
	registerValues("text/scanner", map[string]reflect.Value{
		// Functions
		"TokenString": reflect.ValueOf(pkg.TokenString),

		// Consts

		"ScanIdents":     reflect.ValueOf(pkg.ScanIdents),
		"ScanInts":       reflect.ValueOf(pkg.ScanInts),
		"ScanFloats":     reflect.ValueOf(pkg.ScanFloats),
		"ScanChars":      reflect.ValueOf(pkg.ScanChars),
		"ScanStrings":    reflect.ValueOf(pkg.ScanStrings),
		"ScanRawStrings": reflect.ValueOf(pkg.ScanRawStrings),
		"ScanComments":   reflect.ValueOf(pkg.ScanComments),
		"SkipComments":   reflect.ValueOf(pkg.SkipComments),
		"GoTokens":       reflect.ValueOf(pkg.GoTokens),
		"EOF":            reflect.ValueOf(pkg.EOF),
		"Ident":          reflect.ValueOf(pkg.Ident),
		"Int":            reflect.ValueOf(pkg.Int),
		"Float":          reflect.ValueOf(pkg.Float),
		"Char":           reflect.ValueOf(pkg.Char),
		"String":         reflect.ValueOf(pkg.String),
		"RawString":      reflect.ValueOf(pkg.RawString),
		"Comment":        reflect.ValueOf(pkg.Comment),
		"GoWhitespace":   reflect.ValueOf(pkg.GoWhitespace),

		// Variables

	})
	registerTypes("text/scanner", map[string]reflect.Type{
		// Non interfaces

		"Position": reflect.TypeOf((*pkg.Position)(nil)).Elem(),
		"Scanner":  reflect.TypeOf((*pkg.Scanner)(nil)).Elem(),
	})
}
