package stdgolibs

import (
	pkg "net/http/cookiejar"

	"reflect"
)

func init() {
	registerValues("net/http/cookiejar", map[string]reflect.Value{
		// Functions
		"New": reflect.ValueOf(pkg.New),

		// Consts

		// Variables

	})
	registerTypes("net/http/cookiejar", map[string]reflect.Type{
		// Non interfaces

		"Options": reflect.TypeOf((*pkg.Options)(nil)).Elem(),
		"Jar":     reflect.TypeOf((*pkg.Jar)(nil)).Elem(),
	})
}
