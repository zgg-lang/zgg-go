package stdgolibs

import (
	pkg "go/build"

	"reflect"
)

func init() {
	registerValues("go/build", map[string]reflect.Value{
		// Functions
		"Import":        reflect.ValueOf(pkg.Import),
		"ImportDir":     reflect.ValueOf(pkg.ImportDir),
		"IsLocalImport": reflect.ValueOf(pkg.IsLocalImport),
		"ArchChar":      reflect.ValueOf(pkg.ArchChar),

		// Consts

		"FindOnly":      reflect.ValueOf(pkg.FindOnly),
		"AllowBinary":   reflect.ValueOf(pkg.AllowBinary),
		"ImportComment": reflect.ValueOf(pkg.ImportComment),
		"IgnoreVendor":  reflect.ValueOf(pkg.IgnoreVendor),

		// Variables

		"Default": reflect.ValueOf(&pkg.Default),
		"ToolDir": reflect.ValueOf(&pkg.ToolDir),
	})
	registerTypes("go/build", map[string]reflect.Type{
		// Non interfaces

		"Context":              reflect.TypeOf((*pkg.Context)(nil)).Elem(),
		"ImportMode":           reflect.TypeOf((*pkg.ImportMode)(nil)).Elem(),
		"Package":              reflect.TypeOf((*pkg.Package)(nil)).Elem(),
		"NoGoError":            reflect.TypeOf((*pkg.NoGoError)(nil)).Elem(),
		"MultiplePackageError": reflect.TypeOf((*pkg.MultiplePackageError)(nil)).Elem(),
	})
}
