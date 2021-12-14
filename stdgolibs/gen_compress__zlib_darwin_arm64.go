package stdgolibs

import (
	pkg "compress/zlib"

	"reflect"
)

func init() {
	registerValues("compress/zlib", map[string]reflect.Value{
		// Functions
		"NewReader":          reflect.ValueOf(pkg.NewReader),
		"NewReaderDict":      reflect.ValueOf(pkg.NewReaderDict),
		"NewWriter":          reflect.ValueOf(pkg.NewWriter),
		"NewWriterLevel":     reflect.ValueOf(pkg.NewWriterLevel),
		"NewWriterLevelDict": reflect.ValueOf(pkg.NewWriterLevelDict),

		// Consts

		"NoCompression":      reflect.ValueOf(pkg.NoCompression),
		"BestSpeed":          reflect.ValueOf(pkg.BestSpeed),
		"BestCompression":    reflect.ValueOf(pkg.BestCompression),
		"DefaultCompression": reflect.ValueOf(pkg.DefaultCompression),
		"HuffmanOnly":        reflect.ValueOf(pkg.HuffmanOnly),

		// Variables

		"ErrChecksum":   reflect.ValueOf(&pkg.ErrChecksum),
		"ErrDictionary": reflect.ValueOf(&pkg.ErrDictionary),
		"ErrHeader":     reflect.ValueOf(&pkg.ErrHeader),
	})
	registerTypes("compress/zlib", map[string]reflect.Type{
		// Non interfaces

		"Writer": reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
	})
}
