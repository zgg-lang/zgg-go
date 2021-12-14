package stdgolibs

import (
	pkg "crypto/sha1"

	"reflect"
)

func init() {
	registerValues("crypto/sha1", map[string]reflect.Value{
		// Functions
		"New": reflect.ValueOf(pkg.New),
		"Sum": reflect.ValueOf(pkg.Sum),

		// Consts

		"Size":      reflect.ValueOf(pkg.Size),
		"BlockSize": reflect.ValueOf(pkg.BlockSize),

		// Variables

	})
	registerTypes("crypto/sha1", map[string]reflect.Type{
		// Non interfaces

	})
}
