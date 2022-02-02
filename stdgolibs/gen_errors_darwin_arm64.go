package stdgolibs

import (
	pkg "errors"

	"reflect"
)

func init() {
	registerValues("errors", map[string]reflect.Value{
		// Functions
		"Unwrap": reflect.ValueOf(pkg.Unwrap),
		"Is":     reflect.ValueOf(pkg.Is),
		"As":     reflect.ValueOf(pkg.As),
		"New":    reflect.ValueOf(pkg.New),

		// Consts

		// Variables

	})
	registerTypes("errors", map[string]reflect.Type{
		// Non interfaces

	})
}
