package stdgolibs

import (
	pkg "unicode/utf8"

	"reflect"
)

func init() {
	registerValues("unicode/utf8", map[string]reflect.Value{
		// Functions
		"FullRune":               reflect.ValueOf(pkg.FullRune),
		"FullRuneInString":       reflect.ValueOf(pkg.FullRuneInString),
		"DecodeRune":             reflect.ValueOf(pkg.DecodeRune),
		"DecodeRuneInString":     reflect.ValueOf(pkg.DecodeRuneInString),
		"DecodeLastRune":         reflect.ValueOf(pkg.DecodeLastRune),
		"DecodeLastRuneInString": reflect.ValueOf(pkg.DecodeLastRuneInString),
		"RuneLen":                reflect.ValueOf(pkg.RuneLen),
		"EncodeRune":             reflect.ValueOf(pkg.EncodeRune),
		"RuneCount":              reflect.ValueOf(pkg.RuneCount),
		"RuneCountInString":      reflect.ValueOf(pkg.RuneCountInString),
		"RuneStart":              reflect.ValueOf(pkg.RuneStart),
		"Valid":                  reflect.ValueOf(pkg.Valid),
		"ValidString":            reflect.ValueOf(pkg.ValidString),
		"ValidRune":              reflect.ValueOf(pkg.ValidRune),

		// Consts

		"RuneError": reflect.ValueOf(pkg.RuneError),
		"RuneSelf":  reflect.ValueOf(pkg.RuneSelf),
		"MaxRune":   reflect.ValueOf(pkg.MaxRune),
		"UTFMax":    reflect.ValueOf(pkg.UTFMax),

		// Variables

	})
	registerTypes("unicode/utf8", map[string]reflect.Type{
		// Non interfaces

	})
}
