package stdgolibs

import (
	"reflect"

	"github.com/zgg-lang/zgg-go/runtime"
	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	values map[string]map[string]reflect.Value
	types  map[string]map[string]reflect.Type
	funcs  map[string]map[string]*runtime.ValueBuiltinFunction
)

func registerValues(importPath string, content map[string]reflect.Value) {
	if values == nil {
		values = make(map[string]map[string]reflect.Value)
	}
	if cur, found := values[importPath]; found {
		for k, v := range content {
			cur[k] = v
		}
	} else {
		values[importPath] = content
	}
}

func registerTypes(importPath string, content map[string]reflect.Type) {
	if types == nil {
		types = make(map[string]map[string]reflect.Type)
	}
	if cur, found := types[importPath]; found {
		for k, v := range content {
			cur[k] = v
		}
	} else {
		types[importPath] = content
	}
}

func registerFuncs(importPath string, content map[string]*runtime.ValueBuiltinFunction) {
	if funcs == nil {
		funcs = make(map[string]map[string]*runtime.ValueBuiltinFunction)
	}
	if cur, found := funcs[importPath]; found {
		for k, v := range content {
			cur[k] = v
		}
	} else {
		funcs[importPath] = content
	}
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
	if content, found := funcs[importPath]; found {
		for name, f := range content {
			rv.SetMember(name, f, c)
		}
		libFound = true
	}
	return rv, libFound
}
