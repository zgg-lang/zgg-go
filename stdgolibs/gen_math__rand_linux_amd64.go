package stdgolibs

import (
	pkg "math/rand"

	"reflect"
)

func init() {
	registerValues("math/rand", map[string]reflect.Value{
		// Functions
		"NewZipf":     reflect.ValueOf(pkg.NewZipf),
		"NewSource":   reflect.ValueOf(pkg.NewSource),
		"New":         reflect.ValueOf(pkg.New),
		"Seed":        reflect.ValueOf(pkg.Seed),
		"Int63":       reflect.ValueOf(pkg.Int63),
		"Uint32":      reflect.ValueOf(pkg.Uint32),
		"Uint64":      reflect.ValueOf(pkg.Uint64),
		"Int31":       reflect.ValueOf(pkg.Int31),
		"Int":         reflect.ValueOf(pkg.Int),
		"Int63n":      reflect.ValueOf(pkg.Int63n),
		"Int31n":      reflect.ValueOf(pkg.Int31n),
		"Intn":        reflect.ValueOf(pkg.Intn),
		"Float64":     reflect.ValueOf(pkg.Float64),
		"Float32":     reflect.ValueOf(pkg.Float32),
		"Perm":        reflect.ValueOf(pkg.Perm),
		"Shuffle":     reflect.ValueOf(pkg.Shuffle),
		"Read":        reflect.ValueOf(pkg.Read),
		"NormFloat64": reflect.ValueOf(pkg.NormFloat64),
		"ExpFloat64":  reflect.ValueOf(pkg.ExpFloat64),

		// Consts

		// Variables

	})
	registerTypes("math/rand", map[string]reflect.Type{
		// Non interfaces

		"Zipf": reflect.TypeOf((*pkg.Zipf)(nil)).Elem(),
		"Rand": reflect.TypeOf((*pkg.Rand)(nil)).Elem(),
	})
}
