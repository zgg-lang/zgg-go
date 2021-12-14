package stdgolibs

import (
	pkg "debug/plan9obj"

	"reflect"
)

func init() {
	registerValues("debug/plan9obj", map[string]reflect.Value{
		// Functions
		"Open":    reflect.ValueOf(pkg.Open),
		"NewFile": reflect.ValueOf(pkg.NewFile),

		// Consts

		"Magic64":    reflect.ValueOf(pkg.Magic64),
		"Magic386":   reflect.ValueOf(pkg.Magic386),
		"MagicAMD64": reflect.ValueOf(pkg.MagicAMD64),
		"MagicARM":   reflect.ValueOf(pkg.MagicARM),

		// Variables

	})
	registerTypes("debug/plan9obj", map[string]reflect.Type{
		// Non interfaces

		"FileHeader":    reflect.TypeOf((*pkg.FileHeader)(nil)).Elem(),
		"File":          reflect.TypeOf((*pkg.File)(nil)).Elem(),
		"SectionHeader": reflect.TypeOf((*pkg.SectionHeader)(nil)).Elem(),
		"Section":       reflect.TypeOf((*pkg.Section)(nil)).Elem(),
		"Sym":           reflect.TypeOf((*pkg.Sym)(nil)).Elem(),
	})
}
