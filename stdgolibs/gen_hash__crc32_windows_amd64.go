package stdgolibs

import (
	pkg "hash/crc32"

	"reflect"
)

func init() {
	registerValues("hash/crc32", map[string]reflect.Value{
		// Functions
		"MakeTable":    reflect.ValueOf(pkg.MakeTable),
		"New":          reflect.ValueOf(pkg.New),
		"NewIEEE":      reflect.ValueOf(pkg.NewIEEE),
		"Update":       reflect.ValueOf(pkg.Update),
		"Checksum":     reflect.ValueOf(pkg.Checksum),
		"ChecksumIEEE": reflect.ValueOf(pkg.ChecksumIEEE),

		// Consts

		"Size":       reflect.ValueOf(pkg.Size),
		"IEEE":       reflect.ValueOf(pkg.IEEE),
		"Castagnoli": reflect.ValueOf(pkg.Castagnoli),
		"Koopman":    reflect.ValueOf(pkg.Koopman),

		// Variables

		"IEEETable": reflect.ValueOf(&pkg.IEEETable),
	})
	registerTypes("hash/crc32", map[string]reflect.Type{
		// Non interfaces

		"Table": reflect.TypeOf((*pkg.Table)(nil)).Elem(),
	})
}
