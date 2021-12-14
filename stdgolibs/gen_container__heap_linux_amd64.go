package stdgolibs

import (
	pkg "container/heap"

	"reflect"
)

func init() {
	registerValues("container/heap", map[string]reflect.Value{
		// Functions
		"Init":   reflect.ValueOf(pkg.Init),
		"Push":   reflect.ValueOf(pkg.Push),
		"Pop":    reflect.ValueOf(pkg.Pop),
		"Remove": reflect.ValueOf(pkg.Remove),
		"Fix":    reflect.ValueOf(pkg.Fix),

		// Consts

		// Variables

	})
	registerTypes("container/heap", map[string]reflect.Type{
		// Non interfaces

	})
}
