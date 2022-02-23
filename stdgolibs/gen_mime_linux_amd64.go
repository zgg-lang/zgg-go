package stdgolibs

import (
	pkg "mime"

	"reflect"
)

func init() {
	registerValues("mime", map[string]reflect.Value{
		// Functions
		"FormatMediaType":  reflect.ValueOf(pkg.FormatMediaType),
		"ParseMediaType":   reflect.ValueOf(pkg.ParseMediaType),
		"TypeByExtension":  reflect.ValueOf(pkg.TypeByExtension),
		"ExtensionsByType": reflect.ValueOf(pkg.ExtensionsByType),
		"AddExtensionType": reflect.ValueOf(pkg.AddExtensionType),

		// Consts

		"BEncoding": reflect.ValueOf(pkg.BEncoding),
		"QEncoding": reflect.ValueOf(pkg.QEncoding),

		// Variables

		"ErrInvalidMediaParameter": reflect.ValueOf(&pkg.ErrInvalidMediaParameter),
	})
	registerTypes("mime", map[string]reflect.Type{
		// Non interfaces

		"WordEncoder": reflect.TypeOf((*pkg.WordEncoder)(nil)).Elem(),
		"WordDecoder": reflect.TypeOf((*pkg.WordDecoder)(nil)).Elem(),
	})
}
