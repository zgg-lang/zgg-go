package stdgolibs

import (
	pkg "go/format"

	"reflect"
)

func init() {
	registerValues("go/format", map[string]reflect.Value{
		// Functions
		"Node":   reflect.ValueOf(pkg.Node),
		"Source": reflect.ValueOf(pkg.Source),

		// Consts

		// Variables

	})
	registerTypes("go/format", map[string]reflect.Type{
		// Non interfaces

	})
}
