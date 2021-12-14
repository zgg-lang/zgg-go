package stdgolibs

import (
	pkg "net/rpc/jsonrpc"

	"reflect"
)

func init() {
	registerValues("net/rpc/jsonrpc", map[string]reflect.Value{
		// Functions
		"NewClientCodec": reflect.ValueOf(pkg.NewClientCodec),
		"NewClient":      reflect.ValueOf(pkg.NewClient),
		"Dial":           reflect.ValueOf(pkg.Dial),
		"NewServerCodec": reflect.ValueOf(pkg.NewServerCodec),
		"ServeConn":      reflect.ValueOf(pkg.ServeConn),

		// Consts

		// Variables

	})
	registerTypes("net/rpc/jsonrpc", map[string]reflect.Type{
		// Non interfaces

	})
}
