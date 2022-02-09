package stdgolibs

import (
	pkg "io/ioutil"

	"reflect"
)

func init() {
	registerValues("io/ioutil", map[string]reflect.Value{
		// Functions
		"ReadAll":   reflect.ValueOf(pkg.ReadAll),
		"ReadFile":  reflect.ValueOf(pkg.ReadFile),
		"WriteFile": reflect.ValueOf(pkg.WriteFile),
		"ReadDir":   reflect.ValueOf(pkg.ReadDir),
		"NopCloser": reflect.ValueOf(pkg.NopCloser),
		"TempFile":  reflect.ValueOf(pkg.TempFile),
		"TempDir":   reflect.ValueOf(pkg.TempDir),

		// Consts

		// Variables

		"Discard": reflect.ValueOf(&pkg.Discard),
	})
	registerTypes("io/ioutil", map[string]reflect.Type{
		// Non interfaces

	})
}
