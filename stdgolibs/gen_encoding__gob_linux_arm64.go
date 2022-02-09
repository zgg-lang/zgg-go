package stdgolibs

import (
	pkg "encoding/gob"

	"reflect"
)

func init() {
	registerValues("encoding/gob", map[string]reflect.Value{
		// Functions
		"NewDecoder":   reflect.ValueOf(pkg.NewDecoder),
		"RegisterName": reflect.ValueOf(pkg.RegisterName),
		"Register":     reflect.ValueOf(pkg.Register),
		"NewEncoder":   reflect.ValueOf(pkg.NewEncoder),

		// Consts

		// Variables

	})
	registerTypes("encoding/gob", map[string]reflect.Type{
		// Non interfaces

		"Decoder":    reflect.TypeOf((*pkg.Decoder)(nil)).Elem(),
		"CommonType": reflect.TypeOf((*pkg.CommonType)(nil)).Elem(),
		"Encoder":    reflect.TypeOf((*pkg.Encoder)(nil)).Elem(),
	})
}
