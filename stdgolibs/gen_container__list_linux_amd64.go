package stdgolibs

import (
	pkg "container/list"

	"reflect"
)

func init() {
	registerValues("container/list", map[string]reflect.Value{
		// Functions
		"New": reflect.ValueOf(pkg.New),

		// Consts

		// Variables

	})
	registerTypes("container/list", map[string]reflect.Type{
		// Non interfaces

		"Element": reflect.TypeOf((*pkg.Element)(nil)).Elem(),
		"List":    reflect.TypeOf((*pkg.List)(nil)).Elem(),
	})
}
