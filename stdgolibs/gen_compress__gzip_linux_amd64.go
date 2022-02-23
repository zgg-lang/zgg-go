package stdgolibs

import (
	pkg "compress/gzip"

	"reflect"
)

func init() {
	registerValues("compress/gzip", map[string]reflect.Value{
		// Functions
		"NewWriter":      reflect.ValueOf(pkg.NewWriter),
		"NewWriterLevel": reflect.ValueOf(pkg.NewWriterLevel),
		"NewReader":      reflect.ValueOf(pkg.NewReader),

		// Consts

		"NoCompression":      reflect.ValueOf(pkg.NoCompression),
		"BestSpeed":          reflect.ValueOf(pkg.BestSpeed),
		"BestCompression":    reflect.ValueOf(pkg.BestCompression),
		"DefaultCompression": reflect.ValueOf(pkg.DefaultCompression),
		"HuffmanOnly":        reflect.ValueOf(pkg.HuffmanOnly),

		// Variables

		"ErrChecksum": reflect.ValueOf(&pkg.ErrChecksum),
		"ErrHeader":   reflect.ValueOf(&pkg.ErrHeader),
	})
	registerTypes("compress/gzip", map[string]reflect.Type{
		// Non interfaces

		"Writer": reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"Header": reflect.TypeOf((*pkg.Header)(nil)).Elem(),
		"Reader": reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
	})
}
