package stdgolibs

import (
	pkg "sync"

	"reflect"
)

func init() {
	registerValues("sync", map[string]reflect.Value{
		// Functions
		"NewCond": reflect.ValueOf(pkg.NewCond),

		// Consts

		// Variables

	})
	registerTypes("sync", map[string]reflect.Type{
		// Non interfaces

		"Pool":      reflect.TypeOf((*pkg.Pool)(nil)).Elem(),
		"Map":       reflect.TypeOf((*pkg.Map)(nil)).Elem(),
		"Mutex":     reflect.TypeOf((*pkg.Mutex)(nil)).Elem(),
		"RWMutex":   reflect.TypeOf((*pkg.RWMutex)(nil)).Elem(),
		"WaitGroup": reflect.TypeOf((*pkg.WaitGroup)(nil)).Elem(),
		"Cond":      reflect.TypeOf((*pkg.Cond)(nil)).Elem(),
		"Once":      reflect.TypeOf((*pkg.Once)(nil)).Elem(),
	})
}
