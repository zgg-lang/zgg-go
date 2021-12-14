package stdgolibs

import (
	pkg "os/signal"

	"reflect"
)

func init() {
	registerValues("os/signal", map[string]reflect.Value{
		// Functions
		"Ignore":        reflect.ValueOf(pkg.Ignore),
		"Ignored":       reflect.ValueOf(pkg.Ignored),
		"Notify":        reflect.ValueOf(pkg.Notify),
		"Reset":         reflect.ValueOf(pkg.Reset),
		"Stop":          reflect.ValueOf(pkg.Stop),
		"NotifyContext": reflect.ValueOf(pkg.NotifyContext),

		// Consts

		// Variables

	})
	registerTypes("os/signal", map[string]reflect.Type{
		// Non interfaces

	})
}
