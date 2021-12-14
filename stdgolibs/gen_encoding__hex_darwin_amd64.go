package stdgolibs

import (
	pkg "encoding/hex"

	"reflect"
)

func init() {
	registerValues("encoding/hex", map[string]reflect.Value{
		// Functions
		"EncodedLen":     reflect.ValueOf(pkg.EncodedLen),
		"Encode":         reflect.ValueOf(pkg.Encode),
		"DecodedLen":     reflect.ValueOf(pkg.DecodedLen),
		"Decode":         reflect.ValueOf(pkg.Decode),
		"EncodeToString": reflect.ValueOf(pkg.EncodeToString),
		"DecodeString":   reflect.ValueOf(pkg.DecodeString),
		"Dump":           reflect.ValueOf(pkg.Dump),
		"NewEncoder":     reflect.ValueOf(pkg.NewEncoder),
		"NewDecoder":     reflect.ValueOf(pkg.NewDecoder),
		"Dumper":         reflect.ValueOf(pkg.Dumper),

		// Consts

		// Variables

		"ErrLength": reflect.ValueOf(&pkg.ErrLength),
	})
	registerTypes("encoding/hex", map[string]reflect.Type{
		// Non interfaces

		"InvalidByteError": reflect.TypeOf((*pkg.InvalidByteError)(nil)).Elem(),
	})
}
