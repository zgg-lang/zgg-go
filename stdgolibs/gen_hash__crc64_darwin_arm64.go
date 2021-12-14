package stdgolibs

import (
	pkg "hash/crc64"

	"reflect"
)

func init() {
	registerValues("hash/crc64", map[string]reflect.Value{
		// Functions
		"MakeTable": reflect.ValueOf(pkg.MakeTable),
		"New":       reflect.ValueOf(pkg.New),
		"Update":    reflect.ValueOf(pkg.Update),
		"Checksum":  reflect.ValueOf(pkg.Checksum),

		// Consts

		"Size": reflect.ValueOf(pkg.Size),
		"ISO":  reflect.ValueOf(uint64(pkg.ISO)),
		"ECMA": reflect.ValueOf(uint64(pkg.ECMA)),

		// Variables

	})
	registerTypes("hash/crc64", map[string]reflect.Type{
		// Non interfaces

		"Table": reflect.TypeOf((*pkg.Table)(nil)).Elem(),
	})
}
