package stdgolibs

import (
	pkg "sync/atomic"

	"reflect"
)

func init() {
	registerValues("sync/atomic", map[string]reflect.Value{
		// Functions
		"SwapInt32":             reflect.ValueOf(pkg.SwapInt32),
		"SwapInt64":             reflect.ValueOf(pkg.SwapInt64),
		"SwapUint32":            reflect.ValueOf(pkg.SwapUint32),
		"SwapUint64":            reflect.ValueOf(pkg.SwapUint64),
		"SwapUintptr":           reflect.ValueOf(pkg.SwapUintptr),
		"SwapPointer":           reflect.ValueOf(pkg.SwapPointer),
		"CompareAndSwapInt32":   reflect.ValueOf(pkg.CompareAndSwapInt32),
		"CompareAndSwapInt64":   reflect.ValueOf(pkg.CompareAndSwapInt64),
		"CompareAndSwapUint32":  reflect.ValueOf(pkg.CompareAndSwapUint32),
		"CompareAndSwapUint64":  reflect.ValueOf(pkg.CompareAndSwapUint64),
		"CompareAndSwapUintptr": reflect.ValueOf(pkg.CompareAndSwapUintptr),
		"CompareAndSwapPointer": reflect.ValueOf(pkg.CompareAndSwapPointer),
		"AddInt32":              reflect.ValueOf(pkg.AddInt32),
		"AddUint32":             reflect.ValueOf(pkg.AddUint32),
		"AddInt64":              reflect.ValueOf(pkg.AddInt64),
		"AddUint64":             reflect.ValueOf(pkg.AddUint64),
		"AddUintptr":            reflect.ValueOf(pkg.AddUintptr),
		"LoadInt32":             reflect.ValueOf(pkg.LoadInt32),
		"LoadInt64":             reflect.ValueOf(pkg.LoadInt64),
		"LoadUint32":            reflect.ValueOf(pkg.LoadUint32),
		"LoadUint64":            reflect.ValueOf(pkg.LoadUint64),
		"LoadUintptr":           reflect.ValueOf(pkg.LoadUintptr),
		"LoadPointer":           reflect.ValueOf(pkg.LoadPointer),
		"StoreInt32":            reflect.ValueOf(pkg.StoreInt32),
		"StoreInt64":            reflect.ValueOf(pkg.StoreInt64),
		"StoreUint32":           reflect.ValueOf(pkg.StoreUint32),
		"StoreUint64":           reflect.ValueOf(pkg.StoreUint64),
		"StoreUintptr":          reflect.ValueOf(pkg.StoreUintptr),
		"StorePointer":          reflect.ValueOf(pkg.StorePointer),

		// Consts

		// Variables

	})
	registerTypes("sync/atomic", map[string]reflect.Type{
		// Non interfaces

		"Value": reflect.TypeOf((*pkg.Value)(nil)).Elem(),
	})
}
