package stdgolibs

import (
	pkg "mime/multipart"

	"reflect"
)

func init() {
	registerValues("mime/multipart", map[string]reflect.Value{
		// Functions
		"NewReader": reflect.ValueOf(pkg.NewReader),
		"NewWriter": reflect.ValueOf(pkg.NewWriter),

		// Consts

		// Variables

		"ErrMessageTooLarge": reflect.ValueOf(&pkg.ErrMessageTooLarge),
	})
	registerTypes("mime/multipart", map[string]reflect.Type{
		// Non interfaces

		"Part":       reflect.TypeOf((*pkg.Part)(nil)).Elem(),
		"Reader":     reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"Writer":     reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"Form":       reflect.TypeOf((*pkg.Form)(nil)).Elem(),
		"FileHeader": reflect.TypeOf((*pkg.FileHeader)(nil)).Elem(),
	})
}
