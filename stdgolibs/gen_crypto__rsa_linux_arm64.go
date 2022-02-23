package stdgolibs

import (
	pkg "crypto/rsa"

	"reflect"
)

func init() {
	registerValues("crypto/rsa", map[string]reflect.Value{
		// Functions
		"SignPSS":                   reflect.ValueOf(pkg.SignPSS),
		"VerifyPSS":                 reflect.ValueOf(pkg.VerifyPSS),
		"GenerateKey":               reflect.ValueOf(pkg.GenerateKey),
		"GenerateMultiPrimeKey":     reflect.ValueOf(pkg.GenerateMultiPrimeKey),
		"EncryptOAEP":               reflect.ValueOf(pkg.EncryptOAEP),
		"DecryptOAEP":               reflect.ValueOf(pkg.DecryptOAEP),
		"EncryptPKCS1v15":           reflect.ValueOf(pkg.EncryptPKCS1v15),
		"DecryptPKCS1v15":           reflect.ValueOf(pkg.DecryptPKCS1v15),
		"DecryptPKCS1v15SessionKey": reflect.ValueOf(pkg.DecryptPKCS1v15SessionKey),
		"SignPKCS1v15":              reflect.ValueOf(pkg.SignPKCS1v15),
		"VerifyPKCS1v15":            reflect.ValueOf(pkg.VerifyPKCS1v15),

		// Consts

		"PSSSaltLengthAuto":       reflect.ValueOf(pkg.PSSSaltLengthAuto),
		"PSSSaltLengthEqualsHash": reflect.ValueOf(pkg.PSSSaltLengthEqualsHash),

		// Variables

		"ErrMessageTooLong": reflect.ValueOf(&pkg.ErrMessageTooLong),
		"ErrDecryption":     reflect.ValueOf(&pkg.ErrDecryption),
		"ErrVerification":   reflect.ValueOf(&pkg.ErrVerification),
	})
	registerTypes("crypto/rsa", map[string]reflect.Type{
		// Non interfaces

		"PSSOptions":             reflect.TypeOf((*pkg.PSSOptions)(nil)).Elem(),
		"PublicKey":              reflect.TypeOf((*pkg.PublicKey)(nil)).Elem(),
		"OAEPOptions":            reflect.TypeOf((*pkg.OAEPOptions)(nil)).Elem(),
		"PrivateKey":             reflect.TypeOf((*pkg.PrivateKey)(nil)).Elem(),
		"PrecomputedValues":      reflect.TypeOf((*pkg.PrecomputedValues)(nil)).Elem(),
		"CRTValue":               reflect.TypeOf((*pkg.CRTValue)(nil)).Elem(),
		"PKCS1v15DecryptOptions": reflect.TypeOf((*pkg.PKCS1v15DecryptOptions)(nil)).Elem(),
	})
}
