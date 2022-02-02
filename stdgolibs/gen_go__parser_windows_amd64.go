package stdgolibs

import (
	pkg "go/parser"

	"reflect"
)

func init() {
	registerValues("go/parser", map[string]reflect.Value{
		// Functions
		"ParseFile":     reflect.ValueOf(pkg.ParseFile),
		"ParseDir":      reflect.ValueOf(pkg.ParseDir),
		"ParseExprFrom": reflect.ValueOf(pkg.ParseExprFrom),
		"ParseExpr":     reflect.ValueOf(pkg.ParseExpr),

		// Consts

		"PackageClauseOnly": reflect.ValueOf(pkg.PackageClauseOnly),
		"ImportsOnly":       reflect.ValueOf(pkg.ImportsOnly),
		"ParseComments":     reflect.ValueOf(pkg.ParseComments),
		"Trace":             reflect.ValueOf(pkg.Trace),
		"DeclarationErrors": reflect.ValueOf(pkg.DeclarationErrors),
		"SpuriousErrors":    reflect.ValueOf(pkg.SpuriousErrors),
		"AllErrors":         reflect.ValueOf(pkg.AllErrors),

		// Variables

	})
	registerTypes("go/parser", map[string]reflect.Type{
		// Non interfaces

		"Mode": reflect.TypeOf((*pkg.Mode)(nil)).Elem(),
	})
}
