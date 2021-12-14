package stdgolibs

import (
	pkg "crypto/x509/pkix"

	"reflect"
)

func init() {
	registerValues("crypto/x509/pkix", map[string]reflect.Value{
		// Functions

		// Consts

		// Variables

	})
	registerTypes("crypto/x509/pkix", map[string]reflect.Type{
		// Non interfaces

		"AlgorithmIdentifier":          reflect.TypeOf((*pkg.AlgorithmIdentifier)(nil)).Elem(),
		"RDNSequence":                  reflect.TypeOf((*pkg.RDNSequence)(nil)).Elem(),
		"RelativeDistinguishedNameSET": reflect.TypeOf((*pkg.RelativeDistinguishedNameSET)(nil)).Elem(),
		"AttributeTypeAndValue":        reflect.TypeOf((*pkg.AttributeTypeAndValue)(nil)).Elem(),
		"AttributeTypeAndValueSET":     reflect.TypeOf((*pkg.AttributeTypeAndValueSET)(nil)).Elem(),
		"Extension":                    reflect.TypeOf((*pkg.Extension)(nil)).Elem(),
		"Name":                         reflect.TypeOf((*pkg.Name)(nil)).Elem(),
		"CertificateList":              reflect.TypeOf((*pkg.CertificateList)(nil)).Elem(),
		"TBSCertificateList":           reflect.TypeOf((*pkg.TBSCertificateList)(nil)).Elem(),
		"RevokedCertificate":           reflect.TypeOf((*pkg.RevokedCertificate)(nil)).Elem(),
	})
}
