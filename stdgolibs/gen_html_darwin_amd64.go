package stdgolibs

import (
	pkg "html"

	"reflect"
)

func init() {
	registerValues("html", map[string]reflect.Value{
		// Functions
		"EscapeString":   reflect.ValueOf(pkg.EscapeString),
		"UnescapeString": reflect.ValueOf(pkg.UnescapeString),

		// Consts

		// Variables

	})
	registerTypes("html", map[string]reflect.Type{
		// Non interfaces

	})
}
