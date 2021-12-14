package stdgolibs

import (
	pkg "compress/bzip2"

	"reflect"
)

func init() {
	registerValues("compress/bzip2", map[string]reflect.Value{
		// Functions
		"NewReader": reflect.ValueOf(pkg.NewReader),

		// Consts

		// Variables

	})
	registerTypes("compress/bzip2", map[string]reflect.Type{
		// Non interfaces

		"StructuralError": reflect.TypeOf((*pkg.StructuralError)(nil)).Elem(),
	})
}
