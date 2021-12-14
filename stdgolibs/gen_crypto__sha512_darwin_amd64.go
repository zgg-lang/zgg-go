package stdgolibs

import (
	pkg "crypto/sha512"

	"reflect"
)

func init() {
	registerValues("crypto/sha512", map[string]reflect.Value{
		// Functions
		"New":        reflect.ValueOf(pkg.New),
		"New512_224": reflect.ValueOf(pkg.New512_224),
		"New512_256": reflect.ValueOf(pkg.New512_256),
		"New384":     reflect.ValueOf(pkg.New384),
		"Sum512":     reflect.ValueOf(pkg.Sum512),
		"Sum384":     reflect.ValueOf(pkg.Sum384),
		"Sum512_224": reflect.ValueOf(pkg.Sum512_224),
		"Sum512_256": reflect.ValueOf(pkg.Sum512_256),

		// Consts

		"Size":      reflect.ValueOf(pkg.Size),
		"Size224":   reflect.ValueOf(pkg.Size224),
		"Size256":   reflect.ValueOf(pkg.Size256),
		"Size384":   reflect.ValueOf(pkg.Size384),
		"BlockSize": reflect.ValueOf(pkg.BlockSize),

		// Variables

	})
	registerTypes("crypto/sha512", map[string]reflect.Type{
		// Non interfaces

	})
}
