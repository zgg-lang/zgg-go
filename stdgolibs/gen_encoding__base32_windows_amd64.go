package stdgolibs

import (
	pkg "encoding/base32"

	"reflect"
)

func init() {
	registerValues("encoding/base32", map[string]reflect.Value{
		// Functions
		"NewEncoding": reflect.ValueOf(pkg.NewEncoding),
		"NewEncoder":  reflect.ValueOf(pkg.NewEncoder),
		"NewDecoder":  reflect.ValueOf(pkg.NewDecoder),

		// Consts

		"StdPadding": reflect.ValueOf(pkg.StdPadding),
		"NoPadding":  reflect.ValueOf(pkg.NoPadding),

		// Variables

		"StdEncoding": reflect.ValueOf(&pkg.StdEncoding),
		"HexEncoding": reflect.ValueOf(&pkg.HexEncoding),
	})
	registerTypes("encoding/base32", map[string]reflect.Type{
		// Non interfaces

		"Encoding":          reflect.TypeOf((*pkg.Encoding)(nil)).Elem(),
		"CorruptInputError": reflect.TypeOf((*pkg.CorruptInputError)(nil)).Elem(),
	})
}
