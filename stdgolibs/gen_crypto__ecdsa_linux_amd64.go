package stdgolibs

import (
	pkg "crypto/ecdsa"

	"reflect"
)

func init() {
	registerValues("crypto/ecdsa", map[string]reflect.Value{
		// Functions
		"GenerateKey": reflect.ValueOf(pkg.GenerateKey),
		"Sign":        reflect.ValueOf(pkg.Sign),
		"SignASN1":    reflect.ValueOf(pkg.SignASN1),
		"Verify":      reflect.ValueOf(pkg.Verify),
		"VerifyASN1":  reflect.ValueOf(pkg.VerifyASN1),

		// Consts

		// Variables

	})
	registerTypes("crypto/ecdsa", map[string]reflect.Type{
		// Non interfaces

		"PublicKey":  reflect.TypeOf((*pkg.PublicKey)(nil)).Elem(),
		"PrivateKey": reflect.TypeOf((*pkg.PrivateKey)(nil)).Elem(),
	})
}
