package stdgolibs

import (
	pkg "crypto/rand"

	"reflect"
)

func init() {
	registerValues("crypto/rand", map[string]reflect.Value{
		// Functions
		"Read":  reflect.ValueOf(pkg.Read),
		"Prime": reflect.ValueOf(pkg.Prime),
		"Int":   reflect.ValueOf(pkg.Int),

		// Consts

		// Variables

		"Reader": reflect.ValueOf(&pkg.Reader),
	})
	registerTypes("crypto/rand", map[string]reflect.Type{
		// Non interfaces

	})
}
