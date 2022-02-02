package stdgolibs

import (
	pkg "encoding/asn1"

	"reflect"
)

func init() {
	registerValues("encoding/asn1", map[string]reflect.Value{
		// Functions
		"Unmarshal":           reflect.ValueOf(pkg.Unmarshal),
		"UnmarshalWithParams": reflect.ValueOf(pkg.UnmarshalWithParams),
		"Marshal":             reflect.ValueOf(pkg.Marshal),
		"MarshalWithParams":   reflect.ValueOf(pkg.MarshalWithParams),

		// Consts

		"TagBoolean":           reflect.ValueOf(pkg.TagBoolean),
		"TagInteger":           reflect.ValueOf(pkg.TagInteger),
		"TagBitString":         reflect.ValueOf(pkg.TagBitString),
		"TagOctetString":       reflect.ValueOf(pkg.TagOctetString),
		"TagNull":              reflect.ValueOf(pkg.TagNull),
		"TagOID":               reflect.ValueOf(pkg.TagOID),
		"TagEnum":              reflect.ValueOf(pkg.TagEnum),
		"TagUTF8String":        reflect.ValueOf(pkg.TagUTF8String),
		"TagSequence":          reflect.ValueOf(pkg.TagSequence),
		"TagSet":               reflect.ValueOf(pkg.TagSet),
		"TagNumericString":     reflect.ValueOf(pkg.TagNumericString),
		"TagPrintableString":   reflect.ValueOf(pkg.TagPrintableString),
		"TagT61String":         reflect.ValueOf(pkg.TagT61String),
		"TagIA5String":         reflect.ValueOf(pkg.TagIA5String),
		"TagUTCTime":           reflect.ValueOf(pkg.TagUTCTime),
		"TagGeneralizedTime":   reflect.ValueOf(pkg.TagGeneralizedTime),
		"TagGeneralString":     reflect.ValueOf(pkg.TagGeneralString),
		"TagBMPString":         reflect.ValueOf(pkg.TagBMPString),
		"ClassUniversal":       reflect.ValueOf(pkg.ClassUniversal),
		"ClassApplication":     reflect.ValueOf(pkg.ClassApplication),
		"ClassContextSpecific": reflect.ValueOf(pkg.ClassContextSpecific),
		"ClassPrivate":         reflect.ValueOf(pkg.ClassPrivate),

		// Variables

		"NullRawValue": reflect.ValueOf(&pkg.NullRawValue),
		"NullBytes":    reflect.ValueOf(&pkg.NullBytes),
	})
	registerTypes("encoding/asn1", map[string]reflect.Type{
		// Non interfaces

		"StructuralError":  reflect.TypeOf((*pkg.StructuralError)(nil)).Elem(),
		"SyntaxError":      reflect.TypeOf((*pkg.SyntaxError)(nil)).Elem(),
		"BitString":        reflect.TypeOf((*pkg.BitString)(nil)).Elem(),
		"ObjectIdentifier": reflect.TypeOf((*pkg.ObjectIdentifier)(nil)).Elem(),
		"Enumerated":       reflect.TypeOf((*pkg.Enumerated)(nil)).Elem(),
		"Flag":             reflect.TypeOf((*pkg.Flag)(nil)).Elem(),
		"RawValue":         reflect.TypeOf((*pkg.RawValue)(nil)).Elem(),
		"RawContent":       reflect.TypeOf((*pkg.RawContent)(nil)).Elem(),
	})
}
