package stdgolibs

import (
	pkg "hash/fnv"

	"reflect"
)

func init() {
	registerValues("hash/fnv", map[string]reflect.Value{
		// Functions
		"New32":   reflect.ValueOf(pkg.New32),
		"New32a":  reflect.ValueOf(pkg.New32a),
		"New64":   reflect.ValueOf(pkg.New64),
		"New64a":  reflect.ValueOf(pkg.New64a),
		"New128":  reflect.ValueOf(pkg.New128),
		"New128a": reflect.ValueOf(pkg.New128a),

		// Consts

		// Variables

	})
	registerTypes("hash/fnv", map[string]reflect.Type{
		// Non interfaces

	})
}
