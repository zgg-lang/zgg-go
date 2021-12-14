package stdgolibs

import (
	pkg "image"

	"reflect"
)

func init() {
	registerValues("image", map[string]reflect.Value{
		// Functions
		"NewRGBA":        reflect.ValueOf(pkg.NewRGBA),
		"NewRGBA64":      reflect.ValueOf(pkg.NewRGBA64),
		"NewNRGBA":       reflect.ValueOf(pkg.NewNRGBA),
		"NewNRGBA64":     reflect.ValueOf(pkg.NewNRGBA64),
		"NewAlpha":       reflect.ValueOf(pkg.NewAlpha),
		"NewAlpha16":     reflect.ValueOf(pkg.NewAlpha16),
		"NewGray":        reflect.ValueOf(pkg.NewGray),
		"NewGray16":      reflect.ValueOf(pkg.NewGray16),
		"NewCMYK":        reflect.ValueOf(pkg.NewCMYK),
		"NewPaletted":    reflect.ValueOf(pkg.NewPaletted),
		"NewUniform":     reflect.ValueOf(pkg.NewUniform),
		"NewYCbCr":       reflect.ValueOf(pkg.NewYCbCr),
		"NewNYCbCrA":     reflect.ValueOf(pkg.NewNYCbCrA),
		"RegisterFormat": reflect.ValueOf(pkg.RegisterFormat),
		"Decode":         reflect.ValueOf(pkg.Decode),
		"DecodeConfig":   reflect.ValueOf(pkg.DecodeConfig),
		"Pt":             reflect.ValueOf(pkg.Pt),
		"Rect":           reflect.ValueOf(pkg.Rect),

		// Consts

		"YCbCrSubsampleRatio444": reflect.ValueOf(pkg.YCbCrSubsampleRatio444),
		"YCbCrSubsampleRatio422": reflect.ValueOf(pkg.YCbCrSubsampleRatio422),
		"YCbCrSubsampleRatio420": reflect.ValueOf(pkg.YCbCrSubsampleRatio420),
		"YCbCrSubsampleRatio440": reflect.ValueOf(pkg.YCbCrSubsampleRatio440),
		"YCbCrSubsampleRatio411": reflect.ValueOf(pkg.YCbCrSubsampleRatio411),
		"YCbCrSubsampleRatio410": reflect.ValueOf(pkg.YCbCrSubsampleRatio410),

		// Variables

		"Black":       reflect.ValueOf(&pkg.Black),
		"White":       reflect.ValueOf(&pkg.White),
		"Transparent": reflect.ValueOf(&pkg.Transparent),
		"Opaque":      reflect.ValueOf(&pkg.Opaque),
		"ErrFormat":   reflect.ValueOf(&pkg.ErrFormat),
		"ZP":          reflect.ValueOf(&pkg.ZP),
		"ZR":          reflect.ValueOf(&pkg.ZR),
	})
	registerTypes("image", map[string]reflect.Type{
		// Non interfaces

		"Config":              reflect.TypeOf((*pkg.Config)(nil)).Elem(),
		"RGBA":                reflect.TypeOf((*pkg.RGBA)(nil)).Elem(),
		"RGBA64":              reflect.TypeOf((*pkg.RGBA64)(nil)).Elem(),
		"NRGBA":               reflect.TypeOf((*pkg.NRGBA)(nil)).Elem(),
		"NRGBA64":             reflect.TypeOf((*pkg.NRGBA64)(nil)).Elem(),
		"Alpha":               reflect.TypeOf((*pkg.Alpha)(nil)).Elem(),
		"Alpha16":             reflect.TypeOf((*pkg.Alpha16)(nil)).Elem(),
		"Gray":                reflect.TypeOf((*pkg.Gray)(nil)).Elem(),
		"Gray16":              reflect.TypeOf((*pkg.Gray16)(nil)).Elem(),
		"CMYK":                reflect.TypeOf((*pkg.CMYK)(nil)).Elem(),
		"Paletted":            reflect.TypeOf((*pkg.Paletted)(nil)).Elem(),
		"Uniform":             reflect.TypeOf((*pkg.Uniform)(nil)).Elem(),
		"YCbCrSubsampleRatio": reflect.TypeOf((*pkg.YCbCrSubsampleRatio)(nil)).Elem(),
		"YCbCr":               reflect.TypeOf((*pkg.YCbCr)(nil)).Elem(),
		"NYCbCrA":             reflect.TypeOf((*pkg.NYCbCrA)(nil)).Elem(),
		"Point":               reflect.TypeOf((*pkg.Point)(nil)).Elem(),
		"Rectangle":           reflect.TypeOf((*pkg.Rectangle)(nil)).Elem(),
	})
}
