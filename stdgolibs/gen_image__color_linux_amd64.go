package stdgolibs

import (
	pkg "image/color"

	"reflect"
)

func init() {
	registerValues("image/color", map[string]reflect.Value{
		// Functions
		"RGBToYCbCr": reflect.ValueOf(pkg.RGBToYCbCr),
		"YCbCrToRGB": reflect.ValueOf(pkg.YCbCrToRGB),
		"RGBToCMYK":  reflect.ValueOf(pkg.RGBToCMYK),
		"CMYKToRGB":  reflect.ValueOf(pkg.CMYKToRGB),
		"ModelFunc":  reflect.ValueOf(pkg.ModelFunc),

		// Consts

		// Variables

		"YCbCrModel":   reflect.ValueOf(&pkg.YCbCrModel),
		"NYCbCrAModel": reflect.ValueOf(&pkg.NYCbCrAModel),
		"CMYKModel":    reflect.ValueOf(&pkg.CMYKModel),
		"RGBAModel":    reflect.ValueOf(&pkg.RGBAModel),
		"RGBA64Model":  reflect.ValueOf(&pkg.RGBA64Model),
		"NRGBAModel":   reflect.ValueOf(&pkg.NRGBAModel),
		"NRGBA64Model": reflect.ValueOf(&pkg.NRGBA64Model),
		"AlphaModel":   reflect.ValueOf(&pkg.AlphaModel),
		"Alpha16Model": reflect.ValueOf(&pkg.Alpha16Model),
		"GrayModel":    reflect.ValueOf(&pkg.GrayModel),
		"Gray16Model":  reflect.ValueOf(&pkg.Gray16Model),
		"Black":        reflect.ValueOf(&pkg.Black),
		"White":        reflect.ValueOf(&pkg.White),
		"Transparent":  reflect.ValueOf(&pkg.Transparent),
		"Opaque":       reflect.ValueOf(&pkg.Opaque),
	})
	registerTypes("image/color", map[string]reflect.Type{
		// Non interfaces

		"YCbCr":   reflect.TypeOf((*pkg.YCbCr)(nil)).Elem(),
		"NYCbCrA": reflect.TypeOf((*pkg.NYCbCrA)(nil)).Elem(),
		"CMYK":    reflect.TypeOf((*pkg.CMYK)(nil)).Elem(),
		"RGBA":    reflect.TypeOf((*pkg.RGBA)(nil)).Elem(),
		"RGBA64":  reflect.TypeOf((*pkg.RGBA64)(nil)).Elem(),
		"NRGBA":   reflect.TypeOf((*pkg.NRGBA)(nil)).Elem(),
		"NRGBA64": reflect.TypeOf((*pkg.NRGBA64)(nil)).Elem(),
		"Alpha":   reflect.TypeOf((*pkg.Alpha)(nil)).Elem(),
		"Alpha16": reflect.TypeOf((*pkg.Alpha16)(nil)).Elem(),
		"Gray":    reflect.TypeOf((*pkg.Gray)(nil)).Elem(),
		"Gray16":  reflect.TypeOf((*pkg.Gray16)(nil)).Elem(),
		"Palette": reflect.TypeOf((*pkg.Palette)(nil)).Elem(),
	})
}
