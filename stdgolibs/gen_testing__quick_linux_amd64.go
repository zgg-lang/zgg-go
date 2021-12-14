package stdgolibs

import (
	pkg "testing/quick"

	"reflect"
)

func init() {
	registerValues("testing/quick", map[string]reflect.Value{
		// Functions
		"Value":      reflect.ValueOf(pkg.Value),
		"Check":      reflect.ValueOf(pkg.Check),
		"CheckEqual": reflect.ValueOf(pkg.CheckEqual),

		// Consts

		// Variables

	})
	registerTypes("testing/quick", map[string]reflect.Type{
		// Non interfaces

		"Config":          reflect.TypeOf((*pkg.Config)(nil)).Elem(),
		"SetupError":      reflect.TypeOf((*pkg.SetupError)(nil)).Elem(),
		"CheckError":      reflect.TypeOf((*pkg.CheckError)(nil)).Elem(),
		"CheckEqualError": reflect.TypeOf((*pkg.CheckEqualError)(nil)).Elem(),
	})
}
