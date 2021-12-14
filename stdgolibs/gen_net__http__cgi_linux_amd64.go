package stdgolibs

import (
	pkg "net/http/cgi"

	"reflect"
)

func init() {
	registerValues("net/http/cgi", map[string]reflect.Value{
		// Functions
		"Request":        reflect.ValueOf(pkg.Request),
		"RequestFromMap": reflect.ValueOf(pkg.RequestFromMap),
		"Serve":          reflect.ValueOf(pkg.Serve),

		// Consts

		// Variables

	})
	registerTypes("net/http/cgi", map[string]reflect.Type{
		// Non interfaces

		"Handler": reflect.TypeOf((*pkg.Handler)(nil)).Elem(),
	})
}
