package stdgolibs

import (
	pkg "compress/lzw"

	"reflect"
)

func init() {
	registerValues("compress/lzw", map[string]reflect.Value{
		// Functions
		"NewWriter": reflect.ValueOf(pkg.NewWriter),
		"NewReader": reflect.ValueOf(pkg.NewReader),

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
