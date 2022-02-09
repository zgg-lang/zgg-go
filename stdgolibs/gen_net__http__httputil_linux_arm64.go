package stdgolibs

import (
	pkg "net/http/httputil"

	"reflect"
)

func init() {
	registerValues("net/http/httputil", map[string]reflect.Value{
		// Functions
		"DumpRequestOut":            reflect.ValueOf(pkg.DumpRequestOut),
		"DumpRequest":               reflect.ValueOf(pkg.DumpRequest),
		"DumpResponse":              reflect.ValueOf(pkg.DumpResponse),
		"NewChunkedReader":          reflect.ValueOf(pkg.NewChunkedReader),
		"NewChunkedWriter":          reflect.ValueOf(pkg.NewChunkedWriter),
		"NewServerConn":             reflect.ValueOf(pkg.NewServerConn),
		"NewClientConn":             reflect.ValueOf(pkg.NewClientConn),
		"NewProxyClientConn":        reflect.ValueOf(pkg.NewProxyClientConn),
		"NewSingleHostReverseProxy": reflect.ValueOf(pkg.NewSingleHostReverseProxy),

		// Consts

		// Variables

		"ErrLineTooLong": reflect.ValueOf(&pkg.ErrLineTooLong),
		"ErrPersistEOF":  reflect.ValueOf(&pkg.ErrPersistEOF),
		"ErrClosed":      reflect.ValueOf(&pkg.ErrClosed),
		"ErrPipeline":    reflect.ValueOf(&pkg.ErrPipeline),
	})
	registerTypes("net/http/httputil", map[string]reflect.Type{
		// Non interfaces

		"ServerConn":   reflect.TypeOf((*pkg.ServerConn)(nil)).Elem(),
		"ClientConn":   reflect.TypeOf((*pkg.ClientConn)(nil)).Elem(),
		"ReverseProxy": reflect.TypeOf((*pkg.ReverseProxy)(nil)).Elem(),
	})
}
