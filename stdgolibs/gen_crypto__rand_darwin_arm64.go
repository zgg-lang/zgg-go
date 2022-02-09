package stdgolibs

import (
	pkg "crypto/rand"

	"reflect"
)

func init() {
	registerValues("crypto/rand", map[string]reflect.Value{
		// Functions
		"Prime": reflect.ValueOf(pkg.Prime),
		"Int":   reflect.ValueOf(pkg.Int),
		"Read":  reflect.ValueOf(pkg.Read),

		// Consts

		// Variables

		"Reader": reflect.ValueOf(&pkg.Reader),
	})
	registerTypes("crypto/rand", map[string]reflect.Type{
		// Non interfaces

	})
}
