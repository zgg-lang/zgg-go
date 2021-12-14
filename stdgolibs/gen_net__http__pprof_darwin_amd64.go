package stdgolibs

import (
	pkg "net/http/pprof"

	"reflect"
)

func init() {
	registerValues("net/http/pprof", map[string]reflect.Value{
		// Functions
		"Cmdline": reflect.ValueOf(pkg.Cmdline),
		"Profile": reflect.ValueOf(pkg.Profile),
		"Trace":   reflect.ValueOf(pkg.Trace),
		"Symbol":  reflect.ValueOf(pkg.Symbol),
		"Handler": reflect.ValueOf(pkg.Handler),
		"Index":   reflect.ValueOf(pkg.Index),

		// Consts

		// Variables

	})
	registerTypes("net/http/pprof", map[string]reflect.Type{
		// Non interfaces

	})
}
