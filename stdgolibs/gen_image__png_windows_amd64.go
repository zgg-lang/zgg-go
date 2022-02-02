package stdgolibs

import (
	pkg "image/png"

	"reflect"
)

func init() {
	registerValues("image/png", map[string]reflect.Value{
		// Functions
		"Decode":       reflect.ValueOf(pkg.Decode),
		"DecodeConfig": reflect.ValueOf(pkg.DecodeConfig),
		"Encode":       reflect.ValueOf(pkg.Encode),

		// Consts

		"DefaultCompression": reflect.ValueOf(pkg.DefaultCompression),
		"NoCompression":      reflect.ValueOf(pkg.NoCompression),
		"BestSpeed":          reflect.ValueOf(pkg.BestSpeed),
		"BestCompression":    reflect.ValueOf(pkg.BestCompression),

		// Variables

	})
	registerTypes("image/png", map[string]reflect.Type{
		// Non interfaces

		"FormatError":      reflect.TypeOf((*pkg.FormatError)(nil)).Elem(),
		"UnsupportedError": reflect.TypeOf((*pkg.UnsupportedError)(nil)).Elem(),
		"Encoder":          reflect.TypeOf((*pkg.Encoder)(nil)).Elem(),
		"EncoderBuffer":    reflect.TypeOf((*pkg.EncoderBuffer)(nil)).Elem(),
		"CompressionLevel": reflect.TypeOf((*pkg.CompressionLevel)(nil)).Elem(),
	})
}
