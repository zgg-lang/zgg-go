package stdgolibs

import (
	pkg "bytes"

	"reflect"
)

func init() {
	registerValues("bytes", map[string]reflect.Value{
		// Functions
		"NewBuffer":       reflect.ValueOf(pkg.NewBuffer),
		"NewBufferString": reflect.ValueOf(pkg.NewBufferString),
		"Equal":           reflect.ValueOf(pkg.Equal),
		"Compare":         reflect.ValueOf(pkg.Compare),
		"Count":           reflect.ValueOf(pkg.Count),
		"Contains":        reflect.ValueOf(pkg.Contains),
		"ContainsAny":     reflect.ValueOf(pkg.ContainsAny),
		"ContainsRune":    reflect.ValueOf(pkg.ContainsRune),
		"IndexByte":       reflect.ValueOf(pkg.IndexByte),
		"LastIndex":       reflect.ValueOf(pkg.LastIndex),
		"LastIndexByte":   reflect.ValueOf(pkg.LastIndexByte),
		"IndexRune":       reflect.ValueOf(pkg.IndexRune),
		"IndexAny":        reflect.ValueOf(pkg.IndexAny),
		"LastIndexAny":    reflect.ValueOf(pkg.LastIndexAny),
		"SplitN":          reflect.ValueOf(pkg.SplitN),
		"SplitAfterN":     reflect.ValueOf(pkg.SplitAfterN),
		"Split":           reflect.ValueOf(pkg.Split),
		"SplitAfter":      reflect.ValueOf(pkg.SplitAfter),
		"Fields":          reflect.ValueOf(pkg.Fields),
		"FieldsFunc":      reflect.ValueOf(pkg.FieldsFunc),
		"Join":            reflect.ValueOf(pkg.Join),
		"HasPrefix":       reflect.ValueOf(pkg.HasPrefix),
		"HasSuffix":       reflect.ValueOf(pkg.HasSuffix),
		"Map":             reflect.ValueOf(pkg.Map),
		"Repeat":          reflect.ValueOf(pkg.Repeat),
		"ToUpper":         reflect.ValueOf(pkg.ToUpper),
		"ToLower":         reflect.ValueOf(pkg.ToLower),
		"ToTitle":         reflect.ValueOf(pkg.ToTitle),
		"ToUpperSpecial":  reflect.ValueOf(pkg.ToUpperSpecial),
		"ToLowerSpecial":  reflect.ValueOf(pkg.ToLowerSpecial),
		"ToTitleSpecial":  reflect.ValueOf(pkg.ToTitleSpecial),
		"ToValidUTF8":     reflect.ValueOf(pkg.ToValidUTF8),
		"Title":           reflect.ValueOf(pkg.Title),
		"TrimLeftFunc":    reflect.ValueOf(pkg.TrimLeftFunc),
		"TrimRightFunc":   reflect.ValueOf(pkg.TrimRightFunc),
		"TrimFunc":        reflect.ValueOf(pkg.TrimFunc),
		"TrimPrefix":      reflect.ValueOf(pkg.TrimPrefix),
		"TrimSuffix":      reflect.ValueOf(pkg.TrimSuffix),
		"IndexFunc":       reflect.ValueOf(pkg.IndexFunc),
		"LastIndexFunc":   reflect.ValueOf(pkg.LastIndexFunc),
		"Trim":            reflect.ValueOf(pkg.Trim),
		"TrimLeft":        reflect.ValueOf(pkg.TrimLeft),
		"TrimRight":       reflect.ValueOf(pkg.TrimRight),
		"TrimSpace":       reflect.ValueOf(pkg.TrimSpace),
		"Runes":           reflect.ValueOf(pkg.Runes),
		"Replace":         reflect.ValueOf(pkg.Replace),
		"ReplaceAll":      reflect.ValueOf(pkg.ReplaceAll),
		"EqualFold":       reflect.ValueOf(pkg.EqualFold),
		"Index":           reflect.ValueOf(pkg.Index),
		"NewReader":       reflect.ValueOf(pkg.NewReader),

		// Consts

		"MinRead": reflect.ValueOf(pkg.MinRead),

		// Variables

		"ErrTooLarge": reflect.ValueOf(&pkg.ErrTooLarge),
	})
	registerTypes("bytes", map[string]reflect.Type{
		// Non interfaces

		"Buffer": reflect.TypeOf((*pkg.Buffer)(nil)).Elem(),
		"Reader": reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
	})
}
