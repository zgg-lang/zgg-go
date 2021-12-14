package stdgolibs

import (
	pkg "mime/quotedprintable"

	"reflect"
)

func init() {
	registerValues("mime/quotedprintable", map[string]reflect.Value{
		// Functions
		"NewReader": reflect.ValueOf(pkg.NewReader),
		"NewWriter": reflect.ValueOf(pkg.NewWriter),

		// Consts

		// Variables

	})
	registerTypes("mime/quotedprintable", map[string]reflect.Type{
		// Non interfaces

		"Reader": reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"Writer": reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
	})
}
