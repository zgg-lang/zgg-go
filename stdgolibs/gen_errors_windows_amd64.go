package stdgolibs

import (
	pkg "errors"

	"reflect"
)

func init() {
	registerValues("errors", map[string]reflect.Value{
		// Functions
		"New":    reflect.ValueOf(pkg.New),
		"Unwrap": reflect.ValueOf(pkg.Unwrap),
		"Is":     reflect.ValueOf(pkg.Is),
		"As":     reflect.ValueOf(pkg.As),

		// Consts

		// Variables

	})
	registerTypes("errors", map[string]reflect.Type{
		// Non interfaces

	})
}
