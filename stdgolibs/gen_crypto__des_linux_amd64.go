package stdgolibs

import (
	pkg "crypto/des"

	"reflect"
)

func init() {
	registerValues("crypto/des", map[string]reflect.Value{
		// Functions
		"NewCipher":          reflect.ValueOf(pkg.NewCipher),
		"NewTripleDESCipher": reflect.ValueOf(pkg.NewTripleDESCipher),

		// Consts

		"BlockSize": reflect.ValueOf(pkg.BlockSize),

		// Variables

	})
	registerTypes("crypto/des", map[string]reflect.Type{
		// Non interfaces

		"KeySizeError": reflect.TypeOf((*pkg.KeySizeError)(nil)).Elem(),
	})
}
