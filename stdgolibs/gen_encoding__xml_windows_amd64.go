package stdgolibs

import (
	pkg "encoding/xml"

	"reflect"
)

func init() {
	registerValues("encoding/xml", map[string]reflect.Value{
		// Functions
		"CopyToken":       reflect.ValueOf(pkg.CopyToken),
		"NewDecoder":      reflect.ValueOf(pkg.NewDecoder),
		"NewTokenDecoder": reflect.ValueOf(pkg.NewTokenDecoder),
		"EscapeText":      reflect.ValueOf(pkg.EscapeText),
		"Escape":          reflect.ValueOf(pkg.Escape),
		"Marshal":         reflect.ValueOf(pkg.Marshal),
		"MarshalIndent":   reflect.ValueOf(pkg.MarshalIndent),
		"NewEncoder":      reflect.ValueOf(pkg.NewEncoder),
		"Unmarshal":       reflect.ValueOf(pkg.Unmarshal),

		// Consts

		"Header": reflect.ValueOf(pkg.Header),

		// Variables

		"HTMLEntity":    reflect.ValueOf(&pkg.HTMLEntity),
		"HTMLAutoClose": reflect.ValueOf(&pkg.HTMLAutoClose),
	})
	registerTypes("encoding/xml", map[string]reflect.Type{
		// Non interfaces

		"TagPathError":         reflect.TypeOf((*pkg.TagPathError)(nil)).Elem(),
		"SyntaxError":          reflect.TypeOf((*pkg.SyntaxError)(nil)).Elem(),
		"Name":                 reflect.TypeOf((*pkg.Name)(nil)).Elem(),
		"Attr":                 reflect.TypeOf((*pkg.Attr)(nil)).Elem(),
		"StartElement":         reflect.TypeOf((*pkg.StartElement)(nil)).Elem(),
		"EndElement":           reflect.TypeOf((*pkg.EndElement)(nil)).Elem(),
		"CharData":             reflect.TypeOf((*pkg.CharData)(nil)).Elem(),
		"Comment":              reflect.TypeOf((*pkg.Comment)(nil)).Elem(),
		"ProcInst":             reflect.TypeOf((*pkg.ProcInst)(nil)).Elem(),
		"Directive":            reflect.TypeOf((*pkg.Directive)(nil)).Elem(),
		"Decoder":              reflect.TypeOf((*pkg.Decoder)(nil)).Elem(),
		"Encoder":              reflect.TypeOf((*pkg.Encoder)(nil)).Elem(),
		"UnsupportedTypeError": reflect.TypeOf((*pkg.UnsupportedTypeError)(nil)).Elem(),
		"UnmarshalError":       reflect.TypeOf((*pkg.UnmarshalError)(nil)).Elem(),
	})
}
