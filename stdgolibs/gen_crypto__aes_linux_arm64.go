package stdgolibs

import (
	pkg "crypto/aes"

	"reflect"
)

func init() {
	registerValues("crypto/aes", map[string]reflect.Value{
		// Functions
		"NewCipher": reflect.ValueOf(pkg.NewCipher),

		// Consts

		"BlockSize": reflect.ValueOf(pkg.BlockSize),

		// Variables

	})
	registerTypes("crypto/aes", map[string]reflect.Type{
		// Non interfaces

		"KeySizeError": reflect.TypeOf((*pkg.KeySizeError)(nil)).Elem(),
	})
}
