package stdgolibs

import (
	pkg "runtime/debug"

	"reflect"
)

func init() {
	registerValues("runtime/debug", map[string]reflect.Value{
		// Functions
		"ReadBuildInfo":   reflect.ValueOf(pkg.ReadBuildInfo),
		"PrintStack":      reflect.ValueOf(pkg.PrintStack),
		"Stack":           reflect.ValueOf(pkg.Stack),
		"ReadGCStats":     reflect.ValueOf(pkg.ReadGCStats),
		"SetGCPercent":    reflect.ValueOf(pkg.SetGCPercent),
		"FreeOSMemory":    reflect.ValueOf(pkg.FreeOSMemory),
		"SetMaxStack":     reflect.ValueOf(pkg.SetMaxStack),
		"SetMaxThreads":   reflect.ValueOf(pkg.SetMaxThreads),
		"SetPanicOnFault": reflect.ValueOf(pkg.SetPanicOnFault),
		"WriteHeapDump":   reflect.ValueOf(pkg.WriteHeapDump),
		"SetTraceback":    reflect.ValueOf(pkg.SetTraceback),

		// Consts

		// Variables

	})
	registerTypes("runtime/debug", map[string]reflect.Type{
		// Non interfaces

		"BuildInfo": reflect.TypeOf((*pkg.BuildInfo)(nil)).Elem(),
		"Module":    reflect.TypeOf((*pkg.Module)(nil)).Elem(),
		"GCStats":   reflect.TypeOf((*pkg.GCStats)(nil)).Elem(),
	})
}
