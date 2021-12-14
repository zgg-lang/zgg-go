package stdgolibs

import (
	pkg "crypto"

	"reflect"
)

func init() {
	registerValues("crypto", map[string]reflect.Value{
		// Functions
		"RegisterHash": reflect.ValueOf(pkg.RegisterHash),

		// Consts

		"MD4":         reflect.ValueOf(pkg.MD4),
		"MD5":         reflect.ValueOf(pkg.MD5),
		"SHA1":        reflect.ValueOf(pkg.SHA1),
		"SHA224":      reflect.ValueOf(pkg.SHA224),
		"SHA256":      reflect.ValueOf(pkg.SHA256),
		"SHA384":      reflect.ValueOf(pkg.SHA384),
		"SHA512":      reflect.ValueOf(pkg.SHA512),
		"MD5SHA1":     reflect.ValueOf(pkg.MD5SHA1),
		"RIPEMD160":   reflect.ValueOf(pkg.RIPEMD160),
		"SHA3_224":    reflect.ValueOf(pkg.SHA3_224),
		"SHA3_256":    reflect.ValueOf(pkg.SHA3_256),
		"SHA3_384":    reflect.ValueOf(pkg.SHA3_384),
		"SHA3_512":    reflect.ValueOf(pkg.SHA3_512),
		"SHA512_224":  reflect.ValueOf(pkg.SHA512_224),
		"SHA512_256":  reflect.ValueOf(pkg.SHA512_256),
		"BLAKE2s_256": reflect.ValueOf(pkg.BLAKE2s_256),
		"BLAKE2b_256": reflect.ValueOf(pkg.BLAKE2b_256),
		"BLAKE2b_384": reflect.ValueOf(pkg.BLAKE2b_384),
		"BLAKE2b_512": reflect.ValueOf(pkg.BLAKE2b_512),

		// Variables

	})
	registerTypes("crypto", map[string]reflect.Type{
		// Non interfaces

		"Hash": reflect.TypeOf((*pkg.Hash)(nil)).Elem(),
	})
}
