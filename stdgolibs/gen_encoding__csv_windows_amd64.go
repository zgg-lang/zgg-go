package stdgolibs

import (
	pkg "encoding/csv"

	"reflect"
)

func init() {
	registerValues("encoding/csv", map[string]reflect.Value{
		// Functions
		"NewReader": reflect.ValueOf(pkg.NewReader),
		"NewWriter": reflect.ValueOf(pkg.NewWriter),

		// Consts

		// Variables

		"ErrTrailingComma": reflect.ValueOf(&pkg.ErrTrailingComma),
		"ErrBareQuote":     reflect.ValueOf(&pkg.ErrBareQuote),
		"ErrQuote":         reflect.ValueOf(&pkg.ErrQuote),
		"ErrFieldCount":    reflect.ValueOf(&pkg.ErrFieldCount),
	})
	registerTypes("encoding/csv", map[string]reflect.Type{
		// Non interfaces

		"ParseError": reflect.TypeOf((*pkg.ParseError)(nil)).Elem(),
		"Reader":     reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"Writer":     reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
	})
}
