package stdgolibs

import (
	pkg "crypto/tls"

	"reflect"
)

func init() {
	registerValues("crypto/tls", map[string]reflect.Value{
		// Functions
		"Server":                   reflect.ValueOf(pkg.Server),
		"Client":                   reflect.ValueOf(pkg.Client),
		"NewListener":              reflect.ValueOf(pkg.NewListener),
		"Listen":                   reflect.ValueOf(pkg.Listen),
		"DialWithDialer":           reflect.ValueOf(pkg.DialWithDialer),
		"Dial":                     reflect.ValueOf(pkg.Dial),
		"LoadX509KeyPair":          reflect.ValueOf(pkg.LoadX509KeyPair),
		"X509KeyPair":              reflect.ValueOf(pkg.X509KeyPair),
		"CipherSuites":             reflect.ValueOf(pkg.CipherSuites),
		"InsecureCipherSuites":     reflect.ValueOf(pkg.InsecureCipherSuites),
		"CipherSuiteName":          reflect.ValueOf(pkg.CipherSuiteName),
		"NewLRUClientSessionCache": reflect.ValueOf(pkg.NewLRUClientSessionCache),

		// Consts

		"TLS_RSA_WITH_RC4_128_SHA":                      reflect.ValueOf(pkg.TLS_RSA_WITH_RC4_128_SHA),
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":                 reflect.ValueOf(pkg.TLS_RSA_WITH_3DES_EDE_CBC_SHA),
		"TLS_RSA_WITH_AES_128_CBC_SHA":                  reflect.ValueOf(pkg.TLS_RSA_WITH_AES_128_CBC_SHA),
		"TLS_RSA_WITH_AES_256_CBC_SHA":                  reflect.ValueOf(pkg.TLS_RSA_WITH_AES_256_CBC_SHA),
		"TLS_RSA_WITH_AES_128_CBC_SHA256":               reflect.ValueOf(pkg.TLS_RSA_WITH_AES_128_CBC_SHA256),
		"TLS_RSA_WITH_AES_128_GCM_SHA256":               reflect.ValueOf(pkg.TLS_RSA_WITH_AES_128_GCM_SHA256),
		"TLS_RSA_WITH_AES_256_GCM_SHA384":               reflect.ValueOf(pkg.TLS_RSA_WITH_AES_256_GCM_SHA384),
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":              reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA),
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA),
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA),
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":                reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_RC4_128_SHA),
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":           reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA),
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":            reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA),
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":            reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA),
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":       reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256),
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":         reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256),
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256),
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256),
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":         reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384),
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384),
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256),
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256),
		"TLS_AES_128_GCM_SHA256":                        reflect.ValueOf(pkg.TLS_AES_128_GCM_SHA256),
		"TLS_AES_256_GCM_SHA384":                        reflect.ValueOf(pkg.TLS_AES_256_GCM_SHA384),
		"TLS_CHACHA20_POLY1305_SHA256":                  reflect.ValueOf(pkg.TLS_CHACHA20_POLY1305_SHA256),
		"TLS_FALLBACK_SCSV":                             reflect.ValueOf(pkg.TLS_FALLBACK_SCSV),
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":          reflect.ValueOf(pkg.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305),
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":        reflect.ValueOf(pkg.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305),
		"VersionTLS10":                                  reflect.ValueOf(pkg.VersionTLS10),
		"VersionTLS11":                                  reflect.ValueOf(pkg.VersionTLS11),
		"VersionTLS12":                                  reflect.ValueOf(pkg.VersionTLS12),
		"VersionTLS13":                                  reflect.ValueOf(pkg.VersionTLS13),
		"VersionSSL30":                                  reflect.ValueOf(pkg.VersionSSL30),
		"CurveP256":                                     reflect.ValueOf(pkg.CurveP256),
		"CurveP384":                                     reflect.ValueOf(pkg.CurveP384),
		"CurveP521":                                     reflect.ValueOf(pkg.CurveP521),
		"X25519":                                        reflect.ValueOf(pkg.X25519),
		"NoClientCert":                                  reflect.ValueOf(pkg.NoClientCert),
		"RequestClientCert":                             reflect.ValueOf(pkg.RequestClientCert),
		"RequireAnyClientCert":                          reflect.ValueOf(pkg.RequireAnyClientCert),
		"VerifyClientCertIfGiven":                       reflect.ValueOf(pkg.VerifyClientCertIfGiven),
		"RequireAndVerifyClientCert":                    reflect.ValueOf(pkg.RequireAndVerifyClientCert),
		"PKCS1WithSHA256":                               reflect.ValueOf(pkg.PKCS1WithSHA256),
		"PKCS1WithSHA384":                               reflect.ValueOf(pkg.PKCS1WithSHA384),
		"PKCS1WithSHA512":                               reflect.ValueOf(pkg.PKCS1WithSHA512),
		"PSSWithSHA256":                                 reflect.ValueOf(pkg.PSSWithSHA256),
		"PSSWithSHA384":                                 reflect.ValueOf(pkg.PSSWithSHA384),
		"PSSWithSHA512":                                 reflect.ValueOf(pkg.PSSWithSHA512),
		"ECDSAWithP256AndSHA256":                        reflect.ValueOf(pkg.ECDSAWithP256AndSHA256),
		"ECDSAWithP384AndSHA384":                        reflect.ValueOf(pkg.ECDSAWithP384AndSHA384),
		"ECDSAWithP521AndSHA512":                        reflect.ValueOf(pkg.ECDSAWithP521AndSHA512),
		"Ed25519":                                       reflect.ValueOf(pkg.Ed25519),
		"PKCS1WithSHA1":                                 reflect.ValueOf(pkg.PKCS1WithSHA1),
		"ECDSAWithSHA1":                                 reflect.ValueOf(pkg.ECDSAWithSHA1),
		"RenegotiateNever":                              reflect.ValueOf(pkg.RenegotiateNever),
		"RenegotiateOnceAsClient":                       reflect.ValueOf(pkg.RenegotiateOnceAsClient),
		"RenegotiateFreelyAsClient":                     reflect.ValueOf(pkg.RenegotiateFreelyAsClient),

		// Variables

	})
	registerTypes("crypto/tls", map[string]reflect.Type{
		// Non interfaces

		"Dialer":                 reflect.TypeOf((*pkg.Dialer)(nil)).Elem(),
		"Conn":                   reflect.TypeOf((*pkg.Conn)(nil)).Elem(),
		"RecordHeaderError":      reflect.TypeOf((*pkg.RecordHeaderError)(nil)).Elem(),
		"CipherSuite":            reflect.TypeOf((*pkg.CipherSuite)(nil)).Elem(),
		"CurveID":                reflect.TypeOf((*pkg.CurveID)(nil)).Elem(),
		"ConnectionState":        reflect.TypeOf((*pkg.ConnectionState)(nil)).Elem(),
		"ClientAuthType":         reflect.TypeOf((*pkg.ClientAuthType)(nil)).Elem(),
		"ClientSessionState":     reflect.TypeOf((*pkg.ClientSessionState)(nil)).Elem(),
		"SignatureScheme":        reflect.TypeOf((*pkg.SignatureScheme)(nil)).Elem(),
		"ClientHelloInfo":        reflect.TypeOf((*pkg.ClientHelloInfo)(nil)).Elem(),
		"CertificateRequestInfo": reflect.TypeOf((*pkg.CertificateRequestInfo)(nil)).Elem(),
		"RenegotiationSupport":   reflect.TypeOf((*pkg.RenegotiationSupport)(nil)).Elem(),
		"Config":                 reflect.TypeOf((*pkg.Config)(nil)).Elem(),
		"Certificate":            reflect.TypeOf((*pkg.Certificate)(nil)).Elem(),
	})
}
