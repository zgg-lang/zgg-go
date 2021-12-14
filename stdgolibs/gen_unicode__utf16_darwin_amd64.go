package stdgolibs

import (
	pkg "unicode/utf16"

	"reflect"
)

func init() {
	registerValues("unicode/utf16", map[string]reflect.Value{
		// Functions
		"IsSurrogate": reflect.ValueOf(pkg.IsSurrogate),
		"DecodeRune":  reflect.ValueOf(pkg.DecodeRune),
		"EncodeRune":  reflect.ValueOf(pkg.EncodeRune),
		"Encode":      reflect.ValueOf(pkg.Encode),
		"Decode":      reflect.ValueOf(pkg.Decode),

		// Consts

		// Variables

	})
	registerTypes("unicode/utf16", map[string]reflect.Type{
		// Non interfaces

	})
}
