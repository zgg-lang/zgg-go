package stdgolibs

import (
	pkg "net/http/httptest"

	"reflect"
)

func init() {
	registerValues("net/http/httptest", map[string]reflect.Value{
		// Functions
		"NewRequest":         reflect.ValueOf(pkg.NewRequest),
		"NewRecorder":        reflect.ValueOf(pkg.NewRecorder),
		"NewServer":          reflect.ValueOf(pkg.NewServer),
		"NewUnstartedServer": reflect.ValueOf(pkg.NewUnstartedServer),
		"NewTLSServer":       reflect.ValueOf(pkg.NewTLSServer),

		// Consts

		"DefaultRemoteAddr": reflect.ValueOf(pkg.DefaultRemoteAddr),

		// Variables

	})
	registerTypes("net/http/httptest", map[string]reflect.Type{
		// Non interfaces

		"ResponseRecorder": reflect.TypeOf((*pkg.ResponseRecorder)(nil)).Elem(),
		"Server":           reflect.TypeOf((*pkg.Server)(nil)).Elem(),
	})
}
