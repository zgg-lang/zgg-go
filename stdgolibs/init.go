package stdgolibs

import (
	"reflect"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var values map[string]map[string]reflect.Value
var types map[string]map[string]reflect.Type

func registerValues(importPath string, content map[string]reflect.Value) {
	if values == nil {
		values = make(map[string]map[string]reflect.Value)
	}
	values[importPath] = content
}

func registerTypes(importPath string, content map[string]reflect.Type) {
	if types == nil {
		types = make(map[string]map[string]reflect.Type)
	}
	types[importPath] = content
}

func FindLib(c *Context, importPath string) (Value, bool) {
	rv := NewObject()
	libFound := false
	if content, found := values[importPath]; found {
		for name, val := range content {
			rv.SetMember(name, NewReflectedGoValue(val), c)
		}
		libFound = true
	}
	if content, found := types[importPath]; found {
		for name, typ := range content {
			rv.SetMember(name, NewGoType(typ), c)
		}
		libFound = true
	}
	return rv, libFound
}
