package stdgolibs

import (
	pkg "io"

	"reflect"
)

func init() {
	registerValues("io", map[string]reflect.Value{
		// Functions
		"MultiReader":      reflect.ValueOf(pkg.MultiReader),
		"MultiWriter":      reflect.ValueOf(pkg.MultiWriter),
		"Pipe":             reflect.ValueOf(pkg.Pipe),
		"WriteString":      reflect.ValueOf(pkg.WriteString),
		"ReadAtLeast":      reflect.ValueOf(pkg.ReadAtLeast),
		"ReadFull":         reflect.ValueOf(pkg.ReadFull),
		"CopyN":            reflect.ValueOf(pkg.CopyN),
		"Copy":             reflect.ValueOf(pkg.Copy),
		"CopyBuffer":       reflect.ValueOf(pkg.CopyBuffer),
		"LimitReader":      reflect.ValueOf(pkg.LimitReader),
		"NewSectionReader": reflect.ValueOf(pkg.NewSectionReader),
		"TeeReader":        reflect.ValueOf(pkg.TeeReader),
		"NopCloser":        reflect.ValueOf(pkg.NopCloser),
		"ReadAll":          reflect.ValueOf(pkg.ReadAll),

		// Consts

		"SeekStart":   reflect.ValueOf(pkg.SeekStart),
		"SeekCurrent": reflect.ValueOf(pkg.SeekCurrent),
		"SeekEnd":     reflect.ValueOf(pkg.SeekEnd),

		// Variables

		"ErrClosedPipe":    reflect.ValueOf(&pkg.ErrClosedPipe),
		"ErrShortWrite":    reflect.ValueOf(&pkg.ErrShortWrite),
		"ErrShortBuffer":   reflect.ValueOf(&pkg.ErrShortBuffer),
		"EOF":              reflect.ValueOf(&pkg.EOF),
		"ErrUnexpectedEOF": reflect.ValueOf(&pkg.ErrUnexpectedEOF),
		"ErrNoProgress":    reflect.ValueOf(&pkg.ErrNoProgress),
		"Discard":          reflect.ValueOf(&pkg.Discard),
	})
	registerTypes("io", map[string]reflect.Type{
		// Non interfaces

		"PipeReader":    reflect.TypeOf((*pkg.PipeReader)(nil)).Elem(),
		"PipeWriter":    reflect.TypeOf((*pkg.PipeWriter)(nil)).Elem(),
		"LimitedReader": reflect.TypeOf((*pkg.LimitedReader)(nil)).Elem(),
		"SectionReader": reflect.TypeOf((*pkg.SectionReader)(nil)).Elem(),
	})
}
