package stdgolibs

import (
	pkg "net/rpc/jsonrpc"

	"reflect"
)

func init() {
	registerValues("net/rpc/jsonrpc", map[string]reflect.Value{
		// Functions
		"NewServerCodec": reflect.ValueOf(pkg.NewServerCodec),
		"ServeConn":      reflect.ValueOf(pkg.ServeConn),
		"NewClientCodec": reflect.ValueOf(pkg.NewClientCodec),
		"NewClient":      reflect.ValueOf(pkg.NewClient),
		"Dial":           reflect.ValueOf(pkg.Dial),

		// Consts

		// Variables

	})
	registerTypes("net/rpc/jsonrpc", map[string]reflect.Type{
		// Non interfaces

	})
}
