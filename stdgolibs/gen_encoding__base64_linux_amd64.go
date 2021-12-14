package stdgolibs

import (
	pkg "encoding/base64"

	"reflect"
)

func init() {
	registerValues("encoding/base64", map[string]reflect.Value{
		// Functions
		"NewEncoding": reflect.ValueOf(pkg.NewEncoding),
		"NewEncoder":  reflect.ValueOf(pkg.NewEncoder),
		"NewDecoder":  reflect.ValueOf(pkg.NewDecoder),

		// Consts

		"StdPadding": reflect.ValueOf(pkg.StdPadding),
		"NoPadding":  reflect.ValueOf(pkg.NoPadding),

		// Variables

		"StdEncoding":    reflect.ValueOf(&pkg.StdEncoding),
		"URLEncoding":    reflect.ValueOf(&pkg.URLEncoding),
		"RawStdEncoding": reflect.ValueOf(&pkg.RawStdEncoding),
		"RawURLEncoding": reflect.ValueOf(&pkg.RawURLEncoding),
	})
	registerTypes("encoding/base64", map[string]reflect.Type{
		// Non interfaces

		"Encoding":          reflect.TypeOf((*pkg.Encoding)(nil)).Elem(),
		"CorruptInputError": reflect.TypeOf((*pkg.CorruptInputError)(nil)).Elem(),
	})
}
