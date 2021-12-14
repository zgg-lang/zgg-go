package stdgolibs

import (
	pkg "container/ring"

	"reflect"
)

func init() {
	registerValues("container/ring", map[string]reflect.Value{
		// Functions
		"New": reflect.ValueOf(pkg.New),

		// Consts

		// Variables

	})
	registerTypes("container/ring", map[string]reflect.Type{
		// Non interfaces

		"Ring": reflect.TypeOf((*pkg.Ring)(nil)).Elem(),
	})
}
