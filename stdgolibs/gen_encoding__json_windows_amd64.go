package stdgolibs

import (
	pkg "encoding/json"

	"reflect"
)

func init() {
	registerValues("encoding/json", map[string]reflect.Value{
		// Functions
		"Marshal":       reflect.ValueOf(pkg.Marshal),
		"MarshalIndent": reflect.ValueOf(pkg.MarshalIndent),
		"HTMLEscape":    reflect.ValueOf(pkg.HTMLEscape),
		"Compact":       reflect.ValueOf(pkg.Compact),
		"Indent":        reflect.ValueOf(pkg.Indent),
		"Valid":         reflect.ValueOf(pkg.Valid),
		"NewDecoder":    reflect.ValueOf(pkg.NewDecoder),
		"NewEncoder":    reflect.ValueOf(pkg.NewEncoder),
		"Unmarshal":     reflect.ValueOf(pkg.Unmarshal),

		// Consts

		// Variables

	})
	registerTypes("encoding/json", map[string]reflect.Type{
		// Non interfaces

		"UnsupportedTypeError":  reflect.TypeOf((*pkg.UnsupportedTypeError)(nil)).Elem(),
		"UnsupportedValueError": reflect.TypeOf((*pkg.UnsupportedValueError)(nil)).Elem(),
		"InvalidUTF8Error":      reflect.TypeOf((*pkg.InvalidUTF8Error)(nil)).Elem(),
		"MarshalerError":        reflect.TypeOf((*pkg.MarshalerError)(nil)).Elem(),
		"SyntaxError":           reflect.TypeOf((*pkg.SyntaxError)(nil)).Elem(),
		"Decoder":               reflect.TypeOf((*pkg.Decoder)(nil)).Elem(),
		"Encoder":               reflect.TypeOf((*pkg.Encoder)(nil)).Elem(),
		"RawMessage":            reflect.TypeOf((*pkg.RawMessage)(nil)).Elem(),
		"Delim":                 reflect.TypeOf((*pkg.Delim)(nil)).Elem(),
		"UnmarshalTypeError":    reflect.TypeOf((*pkg.UnmarshalTypeError)(nil)).Elem(),
		"UnmarshalFieldError":   reflect.TypeOf((*pkg.UnmarshalFieldError)(nil)).Elem(),
		"InvalidUnmarshalError": reflect.TypeOf((*pkg.InvalidUnmarshalError)(nil)).Elem(),
		"Number":                reflect.TypeOf((*pkg.Number)(nil)).Elem(),
	})
}
