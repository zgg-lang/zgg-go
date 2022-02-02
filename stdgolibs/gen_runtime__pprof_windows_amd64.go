package stdgolibs

import (
	pkg "runtime/pprof"

	"reflect"
)

func init() {
	registerValues("runtime/pprof", map[string]reflect.Value{
		// Functions
		"NewProfile":         reflect.ValueOf(pkg.NewProfile),
		"Lookup":             reflect.ValueOf(pkg.Lookup),
		"Profiles":           reflect.ValueOf(pkg.Profiles),
		"WriteHeapProfile":   reflect.ValueOf(pkg.WriteHeapProfile),
		"StartCPUProfile":    reflect.ValueOf(pkg.StartCPUProfile),
		"StopCPUProfile":     reflect.ValueOf(pkg.StopCPUProfile),
		"WithLabels":         reflect.ValueOf(pkg.WithLabels),
		"Labels":             reflect.ValueOf(pkg.Labels),
		"Label":              reflect.ValueOf(pkg.Label),
		"ForLabels":          reflect.ValueOf(pkg.ForLabels),
		"SetGoroutineLabels": reflect.ValueOf(pkg.SetGoroutineLabels),
		"Do":                 reflect.ValueOf(pkg.Do),

		// Consts

		// Variables

	})
	registerTypes("runtime/pprof", map[string]reflect.Type{
		// Non interfaces

		"Profile":  reflect.TypeOf((*pkg.Profile)(nil)).Elem(),
		"LabelSet": reflect.TypeOf((*pkg.LabelSet)(nil)).Elem(),
	})
}
