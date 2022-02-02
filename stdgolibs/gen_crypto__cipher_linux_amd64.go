package stdgolibs

import (
	pkg "crypto/cipher"

	"reflect"
)

func init() {
	registerValues("crypto/cipher", map[string]reflect.Value{
		// Functions
		"NewCBCEncrypter":     reflect.ValueOf(pkg.NewCBCEncrypter),
		"NewCBCDecrypter":     reflect.ValueOf(pkg.NewCBCDecrypter),
		"NewCFBEncrypter":     reflect.ValueOf(pkg.NewCFBEncrypter),
		"NewCFBDecrypter":     reflect.ValueOf(pkg.NewCFBDecrypter),
		"NewCTR":              reflect.ValueOf(pkg.NewCTR),
		"NewGCM":              reflect.ValueOf(pkg.NewGCM),
		"NewGCMWithNonceSize": reflect.ValueOf(pkg.NewGCMWithNonceSize),
		"NewGCMWithTagSize":   reflect.ValueOf(pkg.NewGCMWithTagSize),
		"NewOFB":              reflect.ValueOf(pkg.NewOFB),

		// Consts

		// Variables

	})
	registerTypes("crypto/cipher", map[string]reflect.Type{
		// Non interfaces

		"StreamReader": reflect.TypeOf((*pkg.StreamReader)(nil)).Elem(),
		"StreamWriter": reflect.TypeOf((*pkg.StreamWriter)(nil)).Elem(),
	})
}
