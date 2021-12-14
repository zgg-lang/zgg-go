package stdgolibs

import (
	pkg "debug/gosym"

	"reflect"
)

func init() {
	registerValues("debug/gosym", map[string]reflect.Value{
		// Functions
		"NewLineTable": reflect.ValueOf(pkg.NewLineTable),
		"NewTable":     reflect.ValueOf(pkg.NewTable),

		// Consts

		// Variables

	})
	registerTypes("debug/gosym", map[string]reflect.Type{
		// Non interfaces

		"LineTable":        reflect.TypeOf((*pkg.LineTable)(nil)).Elem(),
		"Sym":              reflect.TypeOf((*pkg.Sym)(nil)).Elem(),
		"Func":             reflect.TypeOf((*pkg.Func)(nil)).Elem(),
		"Obj":              reflect.TypeOf((*pkg.Obj)(nil)).Elem(),
		"Table":            reflect.TypeOf((*pkg.Table)(nil)).Elem(),
		"UnknownFileError": reflect.TypeOf((*pkg.UnknownFileError)(nil)).Elem(),
		"UnknownLineError": reflect.TypeOf((*pkg.UnknownLineError)(nil)).Elem(),
		"DecodingError":    reflect.TypeOf((*pkg.DecodingError)(nil)).Elem(),
	})
}
