package stdgolibs

import (
	pkg "net/rpc"

	"reflect"
)

func init() {
	registerValues("net/rpc", map[string]reflect.Value{
		// Functions
		"NewServer":          reflect.ValueOf(pkg.NewServer),
		"Register":           reflect.ValueOf(pkg.Register),
		"RegisterName":       reflect.ValueOf(pkg.RegisterName),
		"ServeConn":          reflect.ValueOf(pkg.ServeConn),
		"ServeCodec":         reflect.ValueOf(pkg.ServeCodec),
		"ServeRequest":       reflect.ValueOf(pkg.ServeRequest),
		"Accept":             reflect.ValueOf(pkg.Accept),
		"HandleHTTP":         reflect.ValueOf(pkg.HandleHTTP),
		"NewClient":          reflect.ValueOf(pkg.NewClient),
		"NewClientWithCodec": reflect.ValueOf(pkg.NewClientWithCodec),
		"DialHTTP":           reflect.ValueOf(pkg.DialHTTP),
		"DialHTTPPath":       reflect.ValueOf(pkg.DialHTTPPath),
		"Dial":               reflect.ValueOf(pkg.Dial),

		// Consts

		"DefaultRPCPath":   reflect.ValueOf(pkg.DefaultRPCPath),
		"DefaultDebugPath": reflect.ValueOf(pkg.DefaultDebugPath),

		// Variables

		"DefaultServer": reflect.ValueOf(&pkg.DefaultServer),
		"ErrShutdown":   reflect.ValueOf(&pkg.ErrShutdown),
	})
	registerTypes("net/rpc", map[string]reflect.Type{
		// Non interfaces

		"Request":     reflect.TypeOf((*pkg.Request)(nil)).Elem(),
		"Response":    reflect.TypeOf((*pkg.Response)(nil)).Elem(),
		"Server":      reflect.TypeOf((*pkg.Server)(nil)).Elem(),
		"ServerError": reflect.TypeOf((*pkg.ServerError)(nil)).Elem(),
		"Call":        reflect.TypeOf((*pkg.Call)(nil)).Elem(),
		"Client":      reflect.TypeOf((*pkg.Client)(nil)).Elem(),
	})
}
