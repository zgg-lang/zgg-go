package stdgolibs

import (
	pkg "sort"

	"reflect"
)

func init() {
	registerValues("sort", map[string]reflect.Value{
		// Functions
		"Slice":             reflect.ValueOf(pkg.Slice),
		"SliceStable":       reflect.ValueOf(pkg.SliceStable),
		"SliceIsSorted":     reflect.ValueOf(pkg.SliceIsSorted),
		"Sort":              reflect.ValueOf(pkg.Sort),
		"Reverse":           reflect.ValueOf(pkg.Reverse),
		"IsSorted":          reflect.ValueOf(pkg.IsSorted),
		"Ints":              reflect.ValueOf(pkg.Ints),
		"Float64s":          reflect.ValueOf(pkg.Float64s),
		"Strings":           reflect.ValueOf(pkg.Strings),
		"IntsAreSorted":     reflect.ValueOf(pkg.IntsAreSorted),
		"Float64sAreSorted": reflect.ValueOf(pkg.Float64sAreSorted),
		"StringsAreSorted":  reflect.ValueOf(pkg.StringsAreSorted),
		"Stable":            reflect.ValueOf(pkg.Stable),
		"Search":            reflect.ValueOf(pkg.Search),
		"SearchInts":        reflect.ValueOf(pkg.SearchInts),
		"SearchFloat64s":    reflect.ValueOf(pkg.SearchFloat64s),
		"SearchStrings":     reflect.ValueOf(pkg.SearchStrings),

		// Consts

		// Variables

	})
	registerTypes("sort", map[string]reflect.Type{
		// Non interfaces

		"IntSlice":     reflect.TypeOf((*pkg.IntSlice)(nil)).Elem(),
		"Float64Slice": reflect.TypeOf((*pkg.Float64Slice)(nil)).Elem(),
		"StringSlice":  reflect.TypeOf((*pkg.StringSlice)(nil)).Elem(),
	})
}
