package stdgolibs

import (
	pkg "go/importer"

	"reflect"
)

func init() {
	registerValues("go/importer", map[string]reflect.Value{
		// Functions
		"ForCompiler": reflect.ValueOf(pkg.ForCompiler),
		"For":         reflect.ValueOf(pkg.For),
		"Default":     reflect.ValueOf(pkg.Default),

		// Consts

		// Variables

	})
	registerTypes("go/importer", map[string]reflect.Type{
		// Non interfaces

		"Lookup": reflect.TypeOf((*pkg.Lookup)(nil)).Elem(),
	})
}
