package stdgolibs

import (
	pkg "encoding/pem"

	"reflect"
)

func init() {
	registerValues("encoding/pem", map[string]reflect.Value{
		// Functions
		"Decode":         reflect.ValueOf(pkg.Decode),
		"Encode":         reflect.ValueOf(pkg.Encode),
		"EncodeToMemory": reflect.ValueOf(pkg.EncodeToMemory),

		// Consts

		// Variables

	})
	registerTypes("encoding/pem", map[string]reflect.Type{
		// Non interfaces

		"Block": reflect.TypeOf((*pkg.Block)(nil)).Elem(),
	})
}
