package stdgolibs

import (
	pkg "compress/lzw"

	"reflect"
)

func init() {
	registerValues("compress/lzw", map[string]reflect.Value{
		// Functions
		"NewReader": reflect.ValueOf(pkg.NewReader),
		"NewWriter": reflect.ValueOf(pkg.NewWriter),

		// Consts

		"LSB": reflect.ValueOf(pkg.LSB),
		"MSB": reflect.ValueOf(pkg.MSB),

		// Variables

	})
	registerTypes("compress/lzw", map[string]reflect.Type{
		// Non interfaces

		"Order": reflect.TypeOf((*pkg.Order)(nil)).Elem(),
	})
}
