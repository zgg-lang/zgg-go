package stdgolibs

import (
	pkg "crypto/rc4"

	"reflect"
)

func init() {
	registerValues("crypto/rc4", map[string]reflect.Value{
		// Functions
		"NewCipher": reflect.ValueOf(pkg.NewCipher),

		// Consts

		// Variables

	})
	registerTypes("crypto/rc4", map[string]reflect.Type{
		// Non interfaces

		"Cipher":       reflect.TypeOf((*pkg.Cipher)(nil)).Elem(),
		"KeySizeError": reflect.TypeOf((*pkg.KeySizeError)(nil)).Elem(),
	})
}
