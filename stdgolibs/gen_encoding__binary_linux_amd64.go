package stdgolibs

import (
	pkg "encoding/binary"

	"reflect"
)

func init() {
	registerValues("encoding/binary", map[string]reflect.Value{
		// Functions
		"Read":        reflect.ValueOf(pkg.Read),
		"Write":       reflect.ValueOf(pkg.Write),
		"Size":        reflect.ValueOf(pkg.Size),
		"PutUvarint":  reflect.ValueOf(pkg.PutUvarint),
		"Uvarint":     reflect.ValueOf(pkg.Uvarint),
		"PutVarint":   reflect.ValueOf(pkg.PutVarint),
		"Varint":      reflect.ValueOf(pkg.Varint),
		"ReadUvarint": reflect.ValueOf(pkg.ReadUvarint),
		"ReadVarint":  reflect.ValueOf(pkg.ReadVarint),

		// Consts

		"MaxVarintLen16": reflect.ValueOf(pkg.MaxVarintLen16),
		"MaxVarintLen32": reflect.ValueOf(pkg.MaxVarintLen32),
		"MaxVarintLen64": reflect.ValueOf(pkg.MaxVarintLen64),

		// Variables

		"LittleEndian": reflect.ValueOf(&pkg.LittleEndian),
		"BigEndian":    reflect.ValueOf(&pkg.BigEndian),
	})
	registerTypes("encoding/binary", map[string]reflect.Type{
		// Non interfaces

	})
}
