package stdgolibs

import (
	pkg "crypto/sha256"

	"reflect"
)

func init() {
	registerValues("crypto/sha256", map[string]reflect.Value{
		// Functions
		"New":    reflect.ValueOf(pkg.New),
		"New224": reflect.ValueOf(pkg.New224),
		"Sum256": reflect.ValueOf(pkg.Sum256),
		"Sum224": reflect.ValueOf(pkg.Sum224),

		// Consts

		"Size":      reflect.ValueOf(pkg.Size),
		"Size224":   reflect.ValueOf(pkg.Size224),
		"BlockSize": reflect.ValueOf(pkg.BlockSize),

		// Variables

	})
	registerTypes("crypto/sha256", map[string]reflect.Type{
		// Non interfaces

	})
}
