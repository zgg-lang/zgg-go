package stdgolibs

import (
	pkg "runtime/debug"

	"reflect"
)

func init() {
	registerValues("runtime/debug", map[string]reflect.Value{
		// Functions
		"ReadGCStats":     reflect.ValueOf(pkg.ReadGCStats),
		"SetGCPercent":    reflect.ValueOf(pkg.SetGCPercent),
		"FreeOSMemory":    reflect.ValueOf(pkg.FreeOSMemory),
		"SetMaxStack":     reflect.ValueOf(pkg.SetMaxStack),
		"SetMaxThreads":   reflect.ValueOf(pkg.SetMaxThreads),
		"SetPanicOnFault": reflect.ValueOf(pkg.SetPanicOnFault),
		"WriteHeapDump":   reflect.ValueOf(pkg.WriteHeapDump),
		"SetTraceback":    reflect.ValueOf(pkg.SetTraceback),
		"ReadBuildInfo":   reflect.ValueOf(pkg.ReadBuildInfo),
		"PrintStack":      reflect.ValueOf(pkg.PrintStack),
		"Stack":           reflect.ValueOf(pkg.Stack),

		// Consts

		// Variables

	})
	registerTypes("runtime/debug", map[string]reflect.Type{
		// Non interfaces

		"GCStats":   reflect.TypeOf((*pkg.GCStats)(nil)).Elem(),
		"BuildInfo": reflect.TypeOf((*pkg.BuildInfo)(nil)).Elem(),
		"Module":    reflect.TypeOf((*pkg.Module)(nil)).Elem(),
	})
}
