package stdgolibs

import (
	pkg "crypto/dsa"

	"reflect"
)

func init() {
	registerValues("crypto/dsa", map[string]reflect.Value{
		// Functions
		"GenerateParameters": reflect.ValueOf(pkg.GenerateParameters),
		"GenerateKey":        reflect.ValueOf(pkg.GenerateKey),
		"Sign":               reflect.ValueOf(pkg.Sign),
		"Verify":             reflect.ValueOf(pkg.Verify),

		// Consts

		"L1024N160": reflect.ValueOf(pkg.L1024N160),
		"L2048N224": reflect.ValueOf(pkg.L2048N224),
		"L2048N256": reflect.ValueOf(pkg.L2048N256),
		"L3072N256": reflect.ValueOf(pkg.L3072N256),

		// Variables

		"ErrInvalidPublicKey": reflect.ValueOf(&pkg.ErrInvalidPublicKey),
	})
	registerTypes("crypto/dsa", map[string]reflect.Type{
		// Non interfaces

		"Parameters":     reflect.TypeOf((*pkg.Parameters)(nil)).Elem(),
		"PublicKey":      reflect.TypeOf((*pkg.PublicKey)(nil)).Elem(),
		"PrivateKey":     reflect.TypeOf((*pkg.PrivateKey)(nil)).Elem(),
		"ParameterSizes": reflect.TypeOf((*pkg.ParameterSizes)(nil)).Elem(),
	})
}
