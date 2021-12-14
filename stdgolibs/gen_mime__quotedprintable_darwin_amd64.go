package stdgolibs

import (
	pkg "mime/quotedprintable"

	"reflect"
)

func init() {
	registerValues("mime/quotedprintable", map[string]reflect.Value{
		// Functions
		"NewWriter": reflect.ValueOf(pkg.NewWriter),
		"NewReader": reflect.ValueOf(pkg.NewReader),

		// Consts

		// Variables

	})
	registerTypes("mime/quotedprintable", map[string]reflect.Type{
		// Non interfaces

		"Writer": reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"Reader": reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
	})
}
