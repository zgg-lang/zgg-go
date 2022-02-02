package stdgolibs

import (
	pkg "runtime/trace"

	"reflect"
)

func init() {
	registerValues("runtime/trace", map[string]reflect.Value{
		// Functions
		"Start":       reflect.ValueOf(pkg.Start),
		"Stop":        reflect.ValueOf(pkg.Stop),
		"NewTask":     reflect.ValueOf(pkg.NewTask),
		"Log":         reflect.ValueOf(pkg.Log),
		"Logf":        reflect.ValueOf(pkg.Logf),
		"WithRegion":  reflect.ValueOf(pkg.WithRegion),
		"StartRegion": reflect.ValueOf(pkg.StartRegion),
		"IsEnabled":   reflect.ValueOf(pkg.IsEnabled),

		// Consts

		// Variables

	})
	registerTypes("runtime/trace", map[string]reflect.Type{
		// Non interfaces

		"Task":   reflect.TypeOf((*pkg.Task)(nil)).Elem(),
		"Region": reflect.TypeOf((*pkg.Region)(nil)).Elem(),
	})
}
