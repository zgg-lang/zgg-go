package stdgolibs

import (
	pkg "image/jpeg"

	"reflect"
)

func init() {
	registerValues("image/jpeg", map[string]reflect.Value{
		// Functions
		"Decode":       reflect.ValueOf(pkg.Decode),
		"DecodeConfig": reflect.ValueOf(pkg.DecodeConfig),
		"Encode":       reflect.ValueOf(pkg.Encode),

		// Consts

		"DefaultQuality": reflect.ValueOf(pkg.DefaultQuality),

		// Variables

	})
	registerTypes("image/jpeg", map[string]reflect.Type{
		// Non interfaces

		"FormatError":      reflect.TypeOf((*pkg.FormatError)(nil)).Elem(),
		"UnsupportedError": reflect.TypeOf((*pkg.UnsupportedError)(nil)).Elem(),
		"Options":          reflect.TypeOf((*pkg.Options)(nil)).Elem(),
	})
}
