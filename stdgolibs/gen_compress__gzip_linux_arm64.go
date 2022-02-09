package stdgolibs

import (
	pkg "compress/gzip"

	"reflect"
)

func init() {
	registerValues("compress/gzip", map[string]reflect.Value{
		// Functions
		"NewReader":      reflect.ValueOf(pkg.NewReader),
		"NewWriter":      reflect.ValueOf(pkg.NewWriter),
		"NewWriterLevel": reflect.ValueOf(pkg.NewWriterLevel),

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

		"Header": reflect.TypeOf((*pkg.Header)(nil)).Elem(),
		"Reader": reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"Writer": reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
	})
}
