package stdgolibs

import (
	pkg "os/exec"

	"reflect"
)

func init() {
	registerValues("os/exec", map[string]reflect.Value{
		// Functions
		"LookPath":       reflect.ValueOf(pkg.LookPath),
		"Command":        reflect.ValueOf(pkg.Command),
		"CommandContext": reflect.ValueOf(pkg.CommandContext),

		// Consts

		// Variables

		"ErrNotFound": reflect.ValueOf(&pkg.ErrNotFound),
	})
	registerTypes("os/exec", map[string]reflect.Type{
		// Non interfaces

		"Error":     reflect.TypeOf((*pkg.Error)(nil)).Elem(),
		"Cmd":       reflect.TypeOf((*pkg.Cmd)(nil)).Elem(),
		"ExitError": reflect.TypeOf((*pkg.ExitError)(nil)).Elem(),
	})
}
