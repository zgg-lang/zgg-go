package stdgolibs

import (
	pkg "runtime"

	"reflect"
)

func init() {
	registerValues("runtime", map[string]reflect.Value{
		// Functions
		"SetCPUProfileRate":       reflect.ValueOf(pkg.SetCPUProfileRate),
		"CPUProfile":              reflect.ValueOf(pkg.CPUProfile),
		"SetBlockProfileRate":     reflect.ValueOf(pkg.SetBlockProfileRate),
		"SetMutexProfileFraction": reflect.ValueOf(pkg.SetMutexProfileFraction),
		"MemProfile":              reflect.ValueOf(pkg.MemProfile),
		"BlockProfile":            reflect.ValueOf(pkg.BlockProfile),
		"MutexProfile":            reflect.ValueOf(pkg.MutexProfile),
		"ThreadCreateProfile":     reflect.ValueOf(pkg.ThreadCreateProfile),
		"GoroutineProfile":        reflect.ValueOf(pkg.GoroutineProfile),
		"Stack":                   reflect.ValueOf(pkg.Stack),
		"StartTrace":              reflect.ValueOf(pkg.StartTrace),
		"StopTrace":               reflect.ValueOf(pkg.StopTrace),
		"ReadTrace":               reflect.ValueOf(pkg.ReadTrace),
		"GOMAXPROCS":              reflect.ValueOf(pkg.GOMAXPROCS),
		"NumCPU":                  reflect.ValueOf(pkg.NumCPU),
		"NumCgoCall":              reflect.ValueOf(pkg.NumCgoCall),
		"NumGoroutine":            reflect.ValueOf(pkg.NumGoroutine),
		"Caller":                  reflect.ValueOf(pkg.Caller),
		"Callers":                 reflect.ValueOf(pkg.Callers),
		"GOROOT":                  reflect.ValueOf(pkg.GOROOT),
		"Version":                 reflect.ValueOf(pkg.Version),
		"SetFinalizer":            reflect.ValueOf(pkg.SetFinalizer),
		"KeepAlive":               reflect.ValueOf(pkg.KeepAlive),
		"ReadMemStats":            reflect.ValueOf(pkg.ReadMemStats),
		"GC":                      reflect.ValueOf(pkg.GC),
		"SetCgoTraceback":         reflect.ValueOf(pkg.SetCgoTraceback),
		"CallersFrames":           reflect.ValueOf(pkg.CallersFrames),
		"FuncForPC":               reflect.ValueOf(pkg.FuncForPC),
		"Goexit":                  reflect.ValueOf(pkg.Goexit),
		"Gosched":                 reflect.ValueOf(pkg.Gosched),
		"Breakpoint":              reflect.ValueOf(pkg.Breakpoint),
		"LockOSThread":            reflect.ValueOf(pkg.LockOSThread),
		"UnlockOSThread":          reflect.ValueOf(pkg.UnlockOSThread),

		// Consts

		"GOOS":     reflect.ValueOf(pkg.GOOS),
		"GOARCH":   reflect.ValueOf(pkg.GOARCH),
		"Compiler": reflect.ValueOf(pkg.Compiler),

		// Variables

		"MemProfileRate": reflect.ValueOf(&pkg.MemProfileRate),
	})
	registerTypes("runtime", map[string]reflect.Type{
		// Non interfaces

		"TypeAssertionError": reflect.TypeOf((*pkg.TypeAssertionError)(nil)).Elem(),
		"StackRecord":        reflect.TypeOf((*pkg.StackRecord)(nil)).Elem(),
		"MemProfileRecord":   reflect.TypeOf((*pkg.MemProfileRecord)(nil)).Elem(),
		"BlockProfileRecord": reflect.TypeOf((*pkg.BlockProfileRecord)(nil)).Elem(),
		"MemStats":           reflect.TypeOf((*pkg.MemStats)(nil)).Elem(),
		"Frames":             reflect.TypeOf((*pkg.Frames)(nil)).Elem(),
		"Frame":              reflect.TypeOf((*pkg.Frame)(nil)).Elem(),
		"Func":               reflect.TypeOf((*pkg.Func)(nil)).Elem(),
	})
}
