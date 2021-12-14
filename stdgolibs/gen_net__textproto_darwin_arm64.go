package stdgolibs

import (
	pkg "net/textproto"

	"reflect"
)

func init() {
	registerValues("net/textproto", map[string]reflect.Value{
		// Functions
		"NewWriter":              reflect.ValueOf(pkg.NewWriter),
		"NewReader":              reflect.ValueOf(pkg.NewReader),
		"CanonicalMIMEHeaderKey": reflect.ValueOf(pkg.CanonicalMIMEHeaderKey),
		"NewConn":                reflect.ValueOf(pkg.NewConn),
		"Dial":                   reflect.ValueOf(pkg.Dial),
		"TrimString":             reflect.ValueOf(pkg.TrimString),
		"TrimBytes":              reflect.ValueOf(pkg.TrimBytes),

		// Consts

		// Variables

	})
	registerTypes("net/textproto", map[string]reflect.Type{
		// Non interfaces

		"Writer":        reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"MIMEHeader":    reflect.TypeOf((*pkg.MIMEHeader)(nil)).Elem(),
		"Pipeline":      reflect.TypeOf((*pkg.Pipeline)(nil)).Elem(),
		"Reader":        reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"Error":         reflect.TypeOf((*pkg.Error)(nil)).Elem(),
		"ProtocolError": reflect.TypeOf((*pkg.ProtocolError)(nil)).Elem(),
		"Conn":          reflect.TypeOf((*pkg.Conn)(nil)).Elem(),
	})
}
