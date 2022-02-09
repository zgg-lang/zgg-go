package stdgolibs

import (
	pkg "net/http/fcgi"

	"reflect"
)

func init() {
	registerValues("net/http/fcgi", map[string]reflect.Value{
		// Functions
		"Serve":      reflect.ValueOf(pkg.Serve),
		"ProcessEnv": reflect.ValueOf(pkg.ProcessEnv),

		// Consts

		// Variables

		"ErrRequestAborted": reflect.ValueOf(&pkg.ErrRequestAborted),
		"ErrConnClosed":     reflect.ValueOf(&pkg.ErrConnClosed),
	})
	registerTypes("net/http/fcgi", map[string]reflect.Type{
		// Non interfaces

	})
}
