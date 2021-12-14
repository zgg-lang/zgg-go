package stdgolibs

import (
	pkg "encoding/gob"

	"reflect"
)

func init() {
	registerValues("encoding/gob", map[string]reflect.Value{
		// Functions
		"NewEncoder":   reflect.ValueOf(pkg.NewEncoder),
		"RegisterName": reflect.ValueOf(pkg.RegisterName),
		"Register":     reflect.ValueOf(pkg.Register),
		"NewDecoder":   reflect.ValueOf(pkg.NewDecoder),

		// Consts

		// Variables

	})
	registerTypes("encoding/gob", map[string]reflect.Type{
		// Non interfaces

		"Encoder":    reflect.TypeOf((*pkg.Encoder)(nil)).Elem(),
		"CommonType": reflect.TypeOf((*pkg.CommonType)(nil)).Elem(),
		"Decoder":    reflect.TypeOf((*pkg.Decoder)(nil)).Elem(),
	})
}
