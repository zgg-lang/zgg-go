package stdgolibs

import (
	pkg "encoding/ascii85"

	"reflect"
)

func init() {
	registerValues("encoding/ascii85", map[string]reflect.Value{
		// Functions
		"Encode":        reflect.ValueOf(pkg.Encode),
		"MaxEncodedLen": reflect.ValueOf(pkg.MaxEncodedLen),
		"NewEncoder":    reflect.ValueOf(pkg.NewEncoder),
		"Decode":        reflect.ValueOf(pkg.Decode),
		"NewDecoder":    reflect.ValueOf(pkg.NewDecoder),

		// Consts

		// Variables

	})
	registerTypes("encoding/ascii85", map[string]reflect.Type{
		// Non interfaces

		"CorruptInputError": reflect.TypeOf((*pkg.CorruptInputError)(nil)).Elem(),
	})
}
