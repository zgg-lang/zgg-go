package stdgolibs

import (
	pkg "path/filepath"

	"reflect"
)

func init() {
	registerValues("path/filepath", map[string]reflect.Value{
		// Functions
		"Match":        reflect.ValueOf(pkg.Match),
		"Glob":         reflect.ValueOf(pkg.Glob),
		"Clean":        reflect.ValueOf(pkg.Clean),
		"ToSlash":      reflect.ValueOf(pkg.ToSlash),
		"FromSlash":    reflect.ValueOf(pkg.FromSlash),
		"SplitList":    reflect.ValueOf(pkg.SplitList),
		"Split":        reflect.ValueOf(pkg.Split),
		"Join":         reflect.ValueOf(pkg.Join),
		"Ext":          reflect.ValueOf(pkg.Ext),
		"EvalSymlinks": reflect.ValueOf(pkg.EvalSymlinks),
		"Abs":          reflect.ValueOf(pkg.Abs),
		"Rel":          reflect.ValueOf(pkg.Rel),
		"WalkDir":      reflect.ValueOf(pkg.WalkDir),
		"Walk":         reflect.ValueOf(pkg.Walk),
		"Base":         reflect.ValueOf(pkg.Base),
		"Dir":          reflect.ValueOf(pkg.Dir),
		"VolumeName":   reflect.ValueOf(pkg.VolumeName),
		"IsAbs":        reflect.ValueOf(pkg.IsAbs),
		"HasPrefix":    reflect.ValueOf(pkg.HasPrefix),

		// Consts

		"Separator":     reflect.ValueOf(pkg.Separator),
		"ListSeparator": reflect.ValueOf(pkg.ListSeparator),

		// Variables

		"ErrBadPattern": reflect.ValueOf(&pkg.ErrBadPattern),
		"SkipDir":       reflect.ValueOf(&pkg.SkipDir),
	})
	registerTypes("path/filepath", map[string]reflect.Type{
		// Non interfaces

		"WalkFunc": reflect.TypeOf((*pkg.WalkFunc)(nil)).Elem(),
	})
}
