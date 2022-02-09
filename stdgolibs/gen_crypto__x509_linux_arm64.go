package stdgolibs

import (
	pkg "crypto/x509"

	"reflect"
)

func init() {
	registerValues("crypto/x509", map[string]reflect.Value{
		// Functions
		"NewCertPool":              reflect.ValueOf(pkg.NewCertPool),
		"SystemCertPool":           reflect.ValueOf(pkg.SystemCertPool),
		"IsEncryptedPEMBlock":      reflect.ValueOf(pkg.IsEncryptedPEMBlock),
		"DecryptPEMBlock":          reflect.ValueOf(pkg.DecryptPEMBlock),
		"EncryptPEMBlock":          reflect.ValueOf(pkg.EncryptPEMBlock),
		"ParsePKCS8PrivateKey":     reflect.ValueOf(pkg.ParsePKCS8PrivateKey),
		"MarshalPKCS8PrivateKey":   reflect.ValueOf(pkg.MarshalPKCS8PrivateKey),
		"ParsePKIXPublicKey":       reflect.ValueOf(pkg.ParsePKIXPublicKey),
		"MarshalPKIXPublicKey":     reflect.ValueOf(pkg.MarshalPKIXPublicKey),
		"ParseCertificate":         reflect.ValueOf(pkg.ParseCertificate),
		"ParseCertificates":        reflect.ValueOf(pkg.ParseCertificates),
		"CreateCertificate":        reflect.ValueOf(pkg.CreateCertificate),
		"ParseCRL":                 reflect.ValueOf(pkg.ParseCRL),
		"ParseDERCRL":              reflect.ValueOf(pkg.ParseDERCRL),
		"CreateCertificateRequest": reflect.ValueOf(pkg.CreateCertificateRequest),
		"ParseCertificateRequest":  reflect.ValueOf(pkg.ParseCertificateRequest),
		"CreateRevocationList":     reflect.ValueOf(pkg.CreateRevocationList),
		"ParsePKCS1PrivateKey":     reflect.ValueOf(pkg.ParsePKCS1PrivateKey),
		"MarshalPKCS1PrivateKey":   reflect.ValueOf(pkg.MarshalPKCS1PrivateKey),
		"ParsePKCS1PublicKey":      reflect.ValueOf(pkg.ParsePKCS1PublicKey),
		"MarshalPKCS1PublicKey":    reflect.ValueOf(pkg.MarshalPKCS1PublicKey),
		"ParseECPrivateKey":        reflect.ValueOf(pkg.ParseECPrivateKey),
		"MarshalECPrivateKey":      reflect.ValueOf(pkg.MarshalECPrivateKey),

		// Consts

		"PEMCipherDES":                              reflect.ValueOf(pkg.PEMCipherDES),
		"PEMCipher3DES":                             reflect.ValueOf(pkg.PEMCipher3DES),
		"PEMCipherAES128":                           reflect.ValueOf(pkg.PEMCipherAES128),
		"PEMCipherAES192":                           reflect.ValueOf(pkg.PEMCipherAES192),
		"PEMCipherAES256":                           reflect.ValueOf(pkg.PEMCipherAES256),
		"UnknownSignatureAlgorithm":                 reflect.ValueOf(pkg.UnknownSignatureAlgorithm),
		"MD2WithRSA":                                reflect.ValueOf(pkg.MD2WithRSA),
		"MD5WithRSA":                                reflect.ValueOf(pkg.MD5WithRSA),
		"SHA1WithRSA":                               reflect.ValueOf(pkg.SHA1WithRSA),
		"SHA256WithRSA":                             reflect.ValueOf(pkg.SHA256WithRSA),
		"SHA384WithRSA":                             reflect.ValueOf(pkg.SHA384WithRSA),
		"SHA512WithRSA":                             reflect.ValueOf(pkg.SHA512WithRSA),
		"DSAWithSHA1":                               reflect.ValueOf(pkg.DSAWithSHA1),
		"DSAWithSHA256":                             reflect.ValueOf(pkg.DSAWithSHA256),
		"ECDSAWithSHA1":                             reflect.ValueOf(pkg.ECDSAWithSHA1),
		"ECDSAWithSHA256":                           reflect.ValueOf(pkg.ECDSAWithSHA256),
		"ECDSAWithSHA384":                           reflect.ValueOf(pkg.ECDSAWithSHA384),
		"ECDSAWithSHA512":                           reflect.ValueOf(pkg.ECDSAWithSHA512),
		"SHA256WithRSAPSS":                          reflect.ValueOf(pkg.SHA256WithRSAPSS),
		"SHA384WithRSAPSS":                          reflect.ValueOf(pkg.SHA384WithRSAPSS),
		"SHA512WithRSAPSS":                          reflect.ValueOf(pkg.SHA512WithRSAPSS),
		"PureEd25519":                               reflect.ValueOf(pkg.PureEd25519),
		"UnknownPublicKeyAlgorithm":                 reflect.ValueOf(pkg.UnknownPublicKeyAlgorithm),
		"RSA":                                       reflect.ValueOf(pkg.RSA),
		"DSA":                                       reflect.ValueOf(pkg.DSA),
		"ECDSA":                                     reflect.ValueOf(pkg.ECDSA),
		"Ed25519":                                   reflect.ValueOf(pkg.Ed25519),
		"KeyUsageDigitalSignature":                  reflect.ValueOf(pkg.KeyUsageDigitalSignature),
		"KeyUsageContentCommitment":                 reflect.ValueOf(pkg.KeyUsageContentCommitment),
		"KeyUsageKeyEncipherment":                   reflect.ValueOf(pkg.KeyUsageKeyEncipherment),
		"KeyUsageDataEncipherment":                  reflect.ValueOf(pkg.KeyUsageDataEncipherment),
		"KeyUsageKeyAgreement":                      reflect.ValueOf(pkg.KeyUsageKeyAgreement),
		"KeyUsageCertSign":                          reflect.ValueOf(pkg.KeyUsageCertSign),
		"KeyUsageCRLSign":                           reflect.ValueOf(pkg.KeyUsageCRLSign),
		"KeyUsageEncipherOnly":                      reflect.ValueOf(pkg.KeyUsageEncipherOnly),
		"KeyUsageDecipherOnly":                      reflect.ValueOf(pkg.KeyUsageDecipherOnly),
		"ExtKeyUsageAny":                            reflect.ValueOf(pkg.ExtKeyUsageAny),
		"ExtKeyUsageServerAuth":                     reflect.ValueOf(pkg.ExtKeyUsageServerAuth),
		"ExtKeyUsageClientAuth":                     reflect.ValueOf(pkg.ExtKeyUsageClientAuth),
		"ExtKeyUsageCodeSigning":                    reflect.ValueOf(pkg.ExtKeyUsageCodeSigning),
		"ExtKeyUsageEmailProtection":                reflect.ValueOf(pkg.ExtKeyUsageEmailProtection),
		"ExtKeyUsageIPSECEndSystem":                 reflect.ValueOf(pkg.ExtKeyUsageIPSECEndSystem),
		"ExtKeyUsageIPSECTunnel":                    reflect.ValueOf(pkg.ExtKeyUsageIPSECTunnel),
		"ExtKeyUsageIPSECUser":                      reflect.ValueOf(pkg.ExtKeyUsageIPSECUser),
		"ExtKeyUsageTimeStamping":                   reflect.ValueOf(pkg.ExtKeyUsageTimeStamping),
		"ExtKeyUsageOCSPSigning":                    reflect.ValueOf(pkg.ExtKeyUsageOCSPSigning),
		"ExtKeyUsageMicrosoftServerGatedCrypto":     reflect.ValueOf(pkg.ExtKeyUsageMicrosoftServerGatedCrypto),
		"ExtKeyUsageNetscapeServerGatedCrypto":      reflect.ValueOf(pkg.ExtKeyUsageNetscapeServerGatedCrypto),
		"ExtKeyUsageMicrosoftCommercialCodeSigning": reflect.ValueOf(pkg.ExtKeyUsageMicrosoftCommercialCodeSigning),
		"ExtKeyUsageMicrosoftKernelCodeSigning":     reflect.ValueOf(pkg.ExtKeyUsageMicrosoftKernelCodeSigning),
		"NotAuthorizedToSign":                       reflect.ValueOf(pkg.NotAuthorizedToSign),
		"Expired":                                   reflect.ValueOf(pkg.Expired),
		"CANotAuthorizedForThisName":                reflect.ValueOf(pkg.CANotAuthorizedForThisName),
		"TooManyIntermediates":                      reflect.ValueOf(pkg.TooManyIntermediates),
		"IncompatibleUsage":                         reflect.ValueOf(pkg.IncompatibleUsage),
		"NameMismatch":                              reflect.ValueOf(pkg.NameMismatch),
		"NameConstraintsWithoutSANs":                reflect.ValueOf(pkg.NameConstraintsWithoutSANs),
		"UnconstrainedName":                         reflect.ValueOf(pkg.UnconstrainedName),
		"TooManyConstraints":                        reflect.ValueOf(pkg.TooManyConstraints),
		"CANotAuthorizedForExtKeyUsage":             reflect.ValueOf(pkg.CANotAuthorizedForExtKeyUsage),

		// Variables

		"IncorrectPasswordError":  reflect.ValueOf(&pkg.IncorrectPasswordError),
		"ErrUnsupportedAlgorithm": reflect.ValueOf(&pkg.ErrUnsupportedAlgorithm),
	})
	registerTypes("crypto/x509", map[string]reflect.Type{
		// Non interfaces

		"CertPool":                   reflect.TypeOf((*pkg.CertPool)(nil)).Elem(),
		"PEMCipher":                  reflect.TypeOf((*pkg.PEMCipher)(nil)).Elem(),
		"SignatureAlgorithm":         reflect.TypeOf((*pkg.SignatureAlgorithm)(nil)).Elem(),
		"PublicKeyAlgorithm":         reflect.TypeOf((*pkg.PublicKeyAlgorithm)(nil)).Elem(),
		"KeyUsage":                   reflect.TypeOf((*pkg.KeyUsage)(nil)).Elem(),
		"ExtKeyUsage":                reflect.TypeOf((*pkg.ExtKeyUsage)(nil)).Elem(),
		"Certificate":                reflect.TypeOf((*pkg.Certificate)(nil)).Elem(),
		"InsecureAlgorithmError":     reflect.TypeOf((*pkg.InsecureAlgorithmError)(nil)).Elem(),
		"ConstraintViolationError":   reflect.TypeOf((*pkg.ConstraintViolationError)(nil)).Elem(),
		"UnhandledCriticalExtension": reflect.TypeOf((*pkg.UnhandledCriticalExtension)(nil)).Elem(),
		"CertificateRequest":         reflect.TypeOf((*pkg.CertificateRequest)(nil)).Elem(),
		"RevocationList":             reflect.TypeOf((*pkg.RevocationList)(nil)).Elem(),
		"InvalidReason":              reflect.TypeOf((*pkg.InvalidReason)(nil)).Elem(),
		"CertificateInvalidError":    reflect.TypeOf((*pkg.CertificateInvalidError)(nil)).Elem(),
		"HostnameError":              reflect.TypeOf((*pkg.HostnameError)(nil)).Elem(),
		"UnknownAuthorityError":      reflect.TypeOf((*pkg.UnknownAuthorityError)(nil)).Elem(),
		"SystemRootsError":           reflect.TypeOf((*pkg.SystemRootsError)(nil)).Elem(),
		"VerifyOptions":              reflect.TypeOf((*pkg.VerifyOptions)(nil)).Elem(),
	})
}
