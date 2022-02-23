package stdgolibs

import (
	pkg "encoding/csv"

	"reflect"
)

func init() {
	registerValues("encoding/csv", map[string]reflect.Value{
		// Functions
		"NewWriter": reflect.ValueOf(pkg.NewWriter),
		"NewReader": reflect.ValueOf(pkg.NewReader),

		// Consts

		// Variables

		"ErrTrailingComma": reflect.ValueOf(&pkg.ErrTrailingComma),
		"ErrBareQuote":     reflect.ValueOf(&pkg.ErrBareQuote),
		"ErrQuote":         reflect.ValueOf(&pkg.ErrQuote),
		"ErrFieldCount":    reflect.ValueOf(&pkg.ErrFieldCount),
	})
	registerTypes("encoding/csv", map[string]reflect.Type{
		// Non interfaces

		"Writer":     reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"ParseError": reflect.TypeOf((*pkg.ParseError)(nil)).Elem(),
		"Reader":     reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
	})
}
