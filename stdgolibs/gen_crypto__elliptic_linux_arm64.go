package stdgolibs

import (
	pkg "crypto/elliptic"

	"reflect"
)

func init() {
	registerValues("crypto/elliptic", map[string]reflect.Value{
		// Functions
		"GenerateKey":         reflect.ValueOf(pkg.GenerateKey),
		"Marshal":             reflect.ValueOf(pkg.Marshal),
		"MarshalCompressed":   reflect.ValueOf(pkg.MarshalCompressed),
		"Unmarshal":           reflect.ValueOf(pkg.Unmarshal),
		"UnmarshalCompressed": reflect.ValueOf(pkg.UnmarshalCompressed),
		"P256":                reflect.ValueOf(pkg.P256),
		"P384":                reflect.ValueOf(pkg.P384),
		"P521":                reflect.ValueOf(pkg.P521),
		"P224":                reflect.ValueOf(pkg.P224),

		// Consts

		// Variables

	})
	registerTypes("crypto/elliptic", map[string]reflect.Type{
		// Non interfaces

		"CurveParams": reflect.TypeOf((*pkg.CurveParams)(nil)).Elem(),
	})
}
