package stdgolibs

import (
	pkg "compress/flate"

	"reflect"
)

func init() {
	registerValues("compress/flate", map[string]reflect.Value{
		// Functions
		"NewWriter":     reflect.ValueOf(pkg.NewWriter),
		"NewWriterDict": reflect.ValueOf(pkg.NewWriterDict),
		"NewReader":     reflect.ValueOf(pkg.NewReader),
		"NewReaderDict": reflect.ValueOf(pkg.NewReaderDict),

		// Consts

		"NoCompression":      reflect.ValueOf(pkg.NoCompression),
		"BestSpeed":          reflect.ValueOf(pkg.BestSpeed),
		"BestCompression":    reflect.ValueOf(pkg.BestCompression),
		"DefaultCompression": reflect.ValueOf(pkg.DefaultCompression),
		"HuffmanOnly":        reflect.ValueOf(pkg.HuffmanOnly),

		// Variables

	})
	registerTypes("compress/flate", map[string]reflect.Type{
		// Non interfaces

		"Writer":            reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"CorruptInputError": reflect.TypeOf((*pkg.CorruptInputError)(nil)).Elem(),
		"InternalError":     reflect.TypeOf((*pkg.InternalError)(nil)).Elem(),
		"ReadError":         reflect.TypeOf((*pkg.ReadError)(nil)).Elem(),
		"WriteError":        reflect.TypeOf((*pkg.WriteError)(nil)).Elem(),
	})
}
