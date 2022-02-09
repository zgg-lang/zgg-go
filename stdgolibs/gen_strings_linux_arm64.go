package stdgolibs

import (
	pkg "strings"

	"reflect"
)

func init() {
	registerValues("strings", map[string]reflect.Value{
		// Functions
		"NewReplacer":    reflect.ValueOf(pkg.NewReplacer),
		"Count":          reflect.ValueOf(pkg.Count),
		"Contains":       reflect.ValueOf(pkg.Contains),
		"ContainsAny":    reflect.ValueOf(pkg.ContainsAny),
		"ContainsRune":   reflect.ValueOf(pkg.ContainsRune),
		"LastIndex":      reflect.ValueOf(pkg.LastIndex),
		"IndexByte":      reflect.ValueOf(pkg.IndexByte),
		"IndexRune":      reflect.ValueOf(pkg.IndexRune),
		"IndexAny":       reflect.ValueOf(pkg.IndexAny),
		"LastIndexAny":   reflect.ValueOf(pkg.LastIndexAny),
		"LastIndexByte":  reflect.ValueOf(pkg.LastIndexByte),
		"SplitN":         reflect.ValueOf(pkg.SplitN),
		"SplitAfterN":    reflect.ValueOf(pkg.SplitAfterN),
		"Split":          reflect.ValueOf(pkg.Split),
		"SplitAfter":     reflect.ValueOf(pkg.SplitAfter),
		"Fields":         reflect.ValueOf(pkg.Fields),
		"FieldsFunc":     reflect.ValueOf(pkg.FieldsFunc),
		"Join":           reflect.ValueOf(pkg.Join),
		"HasPrefix":      reflect.ValueOf(pkg.HasPrefix),
		"HasSuffix":      reflect.ValueOf(pkg.HasSuffix),
		"Map":            reflect.ValueOf(pkg.Map),
		"Repeat":         reflect.ValueOf(pkg.Repeat),
		"ToUpper":        reflect.ValueOf(pkg.ToUpper),
		"ToLower":        reflect.ValueOf(pkg.ToLower),
		"ToTitle":        reflect.ValueOf(pkg.ToTitle),
		"ToUpperSpecial": reflect.ValueOf(pkg.ToUpperSpecial),
		"ToLowerSpecial": reflect.ValueOf(pkg.ToLowerSpecial),
		"ToTitleSpecial": reflect.ValueOf(pkg.ToTitleSpecial),
		"ToValidUTF8":    reflect.ValueOf(pkg.ToValidUTF8),
		"Title":          reflect.ValueOf(pkg.Title),
		"TrimLeftFunc":   reflect.ValueOf(pkg.TrimLeftFunc),
		"TrimRightFunc":  reflect.ValueOf(pkg.TrimRightFunc),
		"TrimFunc":       reflect.ValueOf(pkg.TrimFunc),
		"IndexFunc":      reflect.ValueOf(pkg.IndexFunc),
		"LastIndexFunc":  reflect.ValueOf(pkg.LastIndexFunc),
		"Trim":           reflect.ValueOf(pkg.Trim),
		"TrimLeft":       reflect.ValueOf(pkg.TrimLeft),
		"TrimRight":      reflect.ValueOf(pkg.TrimRight),
		"TrimSpace":      reflect.ValueOf(pkg.TrimSpace),
		"TrimPrefix":     reflect.ValueOf(pkg.TrimPrefix),
		"TrimSuffix":     reflect.ValueOf(pkg.TrimSuffix),
		"Replace":        reflect.ValueOf(pkg.Replace),
		"ReplaceAll":     reflect.ValueOf(pkg.ReplaceAll),
		"EqualFold":      reflect.ValueOf(pkg.EqualFold),
		"Index":          reflect.ValueOf(pkg.Index),
		"Compare":        reflect.ValueOf(pkg.Compare),
		"NewReader":      reflect.ValueOf(pkg.NewReader),

		// Consts

		// Variables

	})
	registerTypes("strings", map[string]reflect.Type{
		// Non interfaces

		"Replacer": reflect.TypeOf((*pkg.Replacer)(nil)).Elem(),
		"Builder":  reflect.TypeOf((*pkg.Builder)(nil)).Elem(),
		"Reader":   reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
	})
}
