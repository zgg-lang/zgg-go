package builtin_libs

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"hash/crc32"
	"hash/fnv"
	"io"
	"strings"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libHash(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("md5", hashMakeHashFunc("md5", md5.New), nil)
	lib.SetMember("sha1", hashMakeHashFunc("sha1", sha1.New), nil)
	lib.SetMember("sha256", hashMakeHashFunc("sha256", sha256.New), nil)

	lib.SetMember("crc32", hashMakeHashFunc("crc32", crc32.NewIEEE), nil)

	lib.SetMember("fnv32", hashMakeHashFunc("fnv32", fnv.New32), nil)
	lib.SetMember("fnv32a", hashMakeHashFunc("fnv32a", fnv.New32a), nil)
	lib.SetMember("fnv64", hashMakeHashFunc("fnv64", fnv.New64), nil)
	lib.SetMember("fnv64a", hashMakeHashFunc("fnv64a", fnv.New64a), nil)
	lib.SetMember("fnv128", hashMakeHashFunc("fnv128", fnv.New128), nil)
	lib.SetMember("fnv128a", hashMakeHashFunc("fnv128a", fnv.New128a), nil)
	{
		hmac := NewObject()
		hmac.SetMember("sha1", hashMakeHmacFunc("sha1", sha1.New), nil)
		hmac.SetMember("sha256", hashMakeHmacFunc("sha256", sha256.New), nil)
		hmac.SetMember("md5", hashMakeHmacFunc("md5", md5.New), nil)
		lib.SetMember("hmac", hmac, nil)
	}
	return lib
}

func hashMakeHashFunc[T hash.Hash](name string, getHash func() T) *ValueBuiltinFunction {
	return NewNativeFunction(name, func(c *Context, this Value, args []Value) Value {
		var rd io.Reader
		EnsureFuncParams(c, "hash."+name, args, NewOneOfHelper("value").
			On(TypeStr, func(a Value) {
				rd = strings.NewReader(a.ToString(c))
			}).
			On(TypeBytes, func(a Value) {
				rd = bytes.NewReader(a.(ValueBytes).Value())
			}).
			On(TypeGoValue, func(a Value) {
				if _rd, ok := a.ToGoValue(c).(io.Reader); ok {
					rd = _rd
				} else {
					c.RaiseRuntimeError("hash.%s: value is not a reader", name)
				}
			}))
		h := getHash()
		io.Copy(h, rd)
		return NewBytes(h.Sum(nil))
	}, "value")
}

func hashMakeHmacFunc(name string, getHash func() hash.Hash) *ValueBuiltinFunction {
	return NewNativeFunction(name, func(c *Context, this Value, args []Value) Value {
		var rd io.Reader
		var key []byte
		EnsureFuncParams(c, "hash."+name, args,
			NewOneOfHelper("key").
				On(TypeStr, func(a Value) {
					key = []byte(a.ToString(c))
				}).
				On(TypeBytes, func(a Value) {
					key = a.(ValueBytes).Value()
				}),
			NewOneOfHelper("value").
				On(TypeStr, func(a Value) {
					rd = strings.NewReader(a.ToString(c))
				}).
				On(TypeBytes, func(a Value) {
					rd = bytes.NewReader(a.(ValueBytes).Value())
				}).
				On(TypeGoValue, func(a Value) {
					if _rd, ok := a.ToGoValue(c).(io.Reader); ok {
						rd = _rd
					} else {
						c.RaiseRuntimeError("hash.%s: value is not a reader", name)
					}
				}))
		h := hmac.New(getHash, key)
		io.Copy(h, rd)
		return NewBytes(h.Sum(nil))
	}, "value")
}
