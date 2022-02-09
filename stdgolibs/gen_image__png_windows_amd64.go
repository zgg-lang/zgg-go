package stdgolibs

import (
	pkg "image/png"

	"reflect"
)

func init() {
	registerValues("image/png", map[string]reflect.Value{
		// Functions
		"Encode":       reflect.ValueOf(pkg.Encode),
		"Decode":       reflect.ValueOf(pkg.Decode),
		"DecodeConfig": reflect.ValueOf(pkg.DecodeConfig),

		// Consts

		"DefaultCompression": reflect.ValueOf(pkg.DefaultCompression),
		"NoCompression":      reflect.ValueOf(pkg.NoCompression),
		"BestSpeed":          reflect.ValueOf(pkg.BestSpeed),
		"BestCompression":    reflect.ValueOf(pkg.BestCompression),

		// Variables

	})
	registerTypes("image/png", map[string]reflect.Type{
		// Non interfaces

		"Encoder":          reflect.TypeOf((*pkg.Encoder)(nil)).Elem(),
		"EncoderBuffer":    reflect.TypeOf((*pkg.EncoderBuffer)(nil)).Elem(),
		"CompressionLevel": reflect.TypeOf((*pkg.CompressionLevel)(nil)).Elem(),
		"FormatError":      reflect.TypeOf((*pkg.FormatError)(nil)).Elem(),
		"UnsupportedError": reflect.TypeOf((*pkg.UnsupportedError)(nil)).Elem(),
	})
}
