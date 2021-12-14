package stdgolibs

import (
	pkg "crypto/hmac"

	"reflect"
)

func init() {
	registerValues("crypto/hmac", map[string]reflect.Value{
		// Functions
		"New":   reflect.ValueOf(pkg.New),
		"Equal": reflect.ValueOf(pkg.Equal),

		// Consts

		// Variables

	})
	registerTypes("crypto/hmac", map[string]reflect.Type{
		// Non interfaces

	})
}
