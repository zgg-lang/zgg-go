package stdgolibs

import (
	pkg "go/doc"

	"reflect"
)

func init() {
	registerValues("go/doc", map[string]reflect.Value{
		// Functions
		"Synopsis":      reflect.ValueOf(pkg.Synopsis),
		"ToHTML":        reflect.ValueOf(pkg.ToHTML),
		"ToText":        reflect.ValueOf(pkg.ToText),
		"New":           reflect.ValueOf(pkg.New),
		"NewFromFiles":  reflect.ValueOf(pkg.NewFromFiles),
		"Examples":      reflect.ValueOf(pkg.Examples),
		"IsPredeclared": reflect.ValueOf(pkg.IsPredeclared),

		// Consts

		"AllDecls":    reflect.ValueOf(pkg.AllDecls),
		"AllMethods":  reflect.ValueOf(pkg.AllMethods),
		"PreserveAST": reflect.ValueOf(pkg.PreserveAST),

		// Variables

		"IllegalPrefixes": reflect.ValueOf(&pkg.IllegalPrefixes),
	})
	registerTypes("go/doc", map[string]reflect.Type{
		// Non interfaces

		"Package": reflect.TypeOf((*pkg.Package)(nil)).Elem(),
		"Value":   reflect.TypeOf((*pkg.Value)(nil)).Elem(),
		"Type":    reflect.TypeOf((*pkg.Type)(nil)).Elem(),
		"Func":    reflect.TypeOf((*pkg.Func)(nil)).Elem(),
		"Note":    reflect.TypeOf((*pkg.Note)(nil)).Elem(),
		"Mode":    reflect.TypeOf((*pkg.Mode)(nil)).Elem(),
		"Example": reflect.TypeOf((*pkg.Example)(nil)).Elem(),
		"Filter":  reflect.TypeOf((*pkg.Filter)(nil)).Elem(),
	})
}
