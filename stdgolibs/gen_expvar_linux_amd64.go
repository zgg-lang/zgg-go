package stdgolibs

import (
	pkg "expvar"

	"reflect"
)

func init() {
	registerValues("expvar", map[string]reflect.Value{
		// Functions
		"Publish":   reflect.ValueOf(pkg.Publish),
		"Get":       reflect.ValueOf(pkg.Get),
		"NewInt":    reflect.ValueOf(pkg.NewInt),
		"NewFloat":  reflect.ValueOf(pkg.NewFloat),
		"NewMap":    reflect.ValueOf(pkg.NewMap),
		"NewString": reflect.ValueOf(pkg.NewString),
		"Do":        reflect.ValueOf(pkg.Do),
		"Handler":   reflect.ValueOf(pkg.Handler),

		// Consts

		// Variables

	})
	registerTypes("expvar", map[string]reflect.Type{
		// Non interfaces

		"Int":      reflect.TypeOf((*pkg.Int)(nil)).Elem(),
		"Float":    reflect.TypeOf((*pkg.Float)(nil)).Elem(),
		"Map":      reflect.TypeOf((*pkg.Map)(nil)).Elem(),
		"KeyValue": reflect.TypeOf((*pkg.KeyValue)(nil)).Elem(),
		"String":   reflect.TypeOf((*pkg.String)(nil)).Elem(),
		"Func":     reflect.TypeOf((*pkg.Func)(nil)).Elem(),
	})
}
