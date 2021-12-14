package stdgolibs

import (
	pkg "archive/zip"

	"reflect"
)

func init() {
	registerValues("archive/zip", map[string]reflect.Value{
		// Functions
		"RegisterDecompressor": reflect.ValueOf(pkg.RegisterDecompressor),
		"RegisterCompressor":   reflect.ValueOf(pkg.RegisterCompressor),
		"FileInfoHeader":       reflect.ValueOf(pkg.FileInfoHeader),
		"NewWriter":            reflect.ValueOf(pkg.NewWriter),
		"OpenReader":           reflect.ValueOf(pkg.OpenReader),
		"NewReader":            reflect.ValueOf(pkg.NewReader),

		// Consts

		"Store":   reflect.ValueOf(pkg.Store),
		"Deflate": reflect.ValueOf(pkg.Deflate),

		// Variables

		"ErrFormat":    reflect.ValueOf(&pkg.ErrFormat),
		"ErrAlgorithm": reflect.ValueOf(&pkg.ErrAlgorithm),
		"ErrChecksum":  reflect.ValueOf(&pkg.ErrChecksum),
	})
	registerTypes("archive/zip", map[string]reflect.Type{
		// Non interfaces

		"Compressor":   reflect.TypeOf((*pkg.Compressor)(nil)).Elem(),
		"Decompressor": reflect.TypeOf((*pkg.Decompressor)(nil)).Elem(),
		"FileHeader":   reflect.TypeOf((*pkg.FileHeader)(nil)).Elem(),
		"Writer":       reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"Reader":       reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"ReadCloser":   reflect.TypeOf((*pkg.ReadCloser)(nil)).Elem(),
		"File":         reflect.TypeOf((*pkg.File)(nil)).Elem(),
	})
}
