package stdgolibs

import (
	pkg "regexp"

	"reflect"
)

func init() {
	registerValues("regexp", map[string]reflect.Value{
		// Functions
		"Compile":          reflect.ValueOf(pkg.Compile),
		"CompilePOSIX":     reflect.ValueOf(pkg.CompilePOSIX),
		"MustCompile":      reflect.ValueOf(pkg.MustCompile),
		"MustCompilePOSIX": reflect.ValueOf(pkg.MustCompilePOSIX),
		"MatchReader":      reflect.ValueOf(pkg.MatchReader),
		"MatchString":      reflect.ValueOf(pkg.MatchString),
		"Match":            reflect.ValueOf(pkg.Match),
		"QuoteMeta":        reflect.ValueOf(pkg.QuoteMeta),

		// Consts

		// Variables

	})
	registerTypes("regexp", map[string]reflect.Type{
		// Non interfaces

		"Regexp": reflect.TypeOf((*pkg.Regexp)(nil)).Elem(),
	})
}
