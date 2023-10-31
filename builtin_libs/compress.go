package builtin_libs

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var compressAlgorithmClass = compressBuildAlgorithmClass()

func libCompress(c *Context) ValueObject {
	lib := NewObject()
	// Algorithms
	lib.SetMember("ZLIB", compressNewAlgo(c, "ZLIB", compressZlib{}), c)
	lib.SetMember("GZIP", compressNewAlgo(c, "GZIP", compressGzip{}), c)
	lib.SetMember("FLATE", compressNewAlgo(c, "FLATE", compressFlate{}), c)
	return lib
}

func compressCompress(c *Context, opt ValueObject, algo any, data []byte) ([]byte, error) {
	algoc, is := algo.(compressCompressor)
	if !is {
		return nil, errors.New("not a compress algorithm")
	}
	return algoc.Compress(data, c, opt)
}

func compressDecompress(c *Context, opt ValueObject, algo any, data []byte) ([]byte, error) {
	algoc, is := algo.(compressDecompressor)
	if !is {
		return nil, errors.New("not a decompress algorithm")
	}
	return algoc.Decompress(data, c, opt)
}

type (
	compressCompressor interface {
		Compress([]byte, *Context, ValueObject) ([]byte, error)
	}

	compressDecompressor interface {
		Decompress([]byte, *Context, ValueObject) ([]byte, error)
	}

	compressZlib  struct{}
	compressGzip  struct{}
	compressFlate struct{}
)

func (compressZlib) Compress(data []byte, c *Context, opt ValueObject) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	level := 1
	if optLevel, is := opt.GetMember("level", c).(ValueInt); is {
		level = optLevel.AsInt()
	}
	w, err := zlib.NewWriterLevel(out, level)
	if err != nil {
		return nil, err
	}
	i, l := 0, len(data)
	for i < l {
		n, err := w.Write(data[i:])
		if err != nil {
			return nil, err
		}
		i += n
	}
	w.Close()
	return out.Bytes(), nil
}

func (compressZlib) Decompress(data []byte, c *Context, opt ValueObject) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func (compressGzip) Compress(data []byte, c *Context, opt ValueObject) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	level := 1
	if optLevel, is := opt.GetMember("level", c).(ValueInt); is {
		level = optLevel.AsInt()
	}
	w, err := gzip.NewWriterLevel(out, level)
	if err != nil {
		return nil, err
	}
	i, l := 0, len(data)
	for i < l {
		n, err := w.Write(data[i:])
		if err != nil {
			return nil, err
		}
		i += n
	}
	w.Close()
	return out.Bytes(), nil
}

func (compressGzip) Decompress(data []byte, c *Context, opt ValueObject) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func (compressFlate) Compress(data []byte, c *Context, opt ValueObject) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	level := 1
	if optLevel, is := opt.GetMember("level", c).(ValueInt); is {
		level = optLevel.AsInt()
	}
	w, err := flate.NewWriter(out, level)
	if err != nil {
		return nil, err
	}
	i, l := 0, len(data)
	for i < l {
		n, err := w.Write(data[i:])
		if err != nil {
			return nil, err
		}
		i += n
	}
	w.Close()
	return out.Bytes(), nil
}

func (compressFlate) Decompress(data []byte, c *Context, opt ValueObject) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()
	return io.ReadAll(r)
}

func compressBuildAlgorithmClass() ValueType {
	emptyOpt := NewObject()
	return NewClassBuilder("Algorithm").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				name ValueStr
				algo GoValue
			)
			EnsureFuncParams(c, "Algorithm.__init__", args,
				ArgRuleRequired("name", TypeStr, &name),
				ArgRuleRequired("algo", TypeGoValue, &algo),
			)
			a := algo.ToGoValue()
			_, canCompress := a.(compressCompressor)
			_, canDecompress := a.(compressDecompressor)
			if !canCompress && !canDecompress {
				c.RaiseRuntimeError("algo is neither a Compressor nor a Decompressor")
			}
			this.SetMember("__name", name, c)
			this.SetMember("__algo", algo, c)
		}).
		Method("compress", func(c *Context, this ValueObject, args []Value) Value {
			var (
				dataBytes ValueBytes
				dataStr   ValueStr
				dataWhich int
				data      []byte
				opt       ValueObject
			)
			EnsureFuncParams(c, "compress", args,
				ArgRuleOneOf("data",
					[]ValueType{TypeBytes, TypeStr},
					[]any{&dataBytes, &dataStr},
					&dataWhich, nil, nil),
				ArgRuleOptional("opt", TypeObject, &opt, emptyOpt),
			)
			switch dataWhich {
			case 0:
				data = dataBytes.Value()
			case 1:
				data = []byte(dataStr.Value())
			default:
				c.RaiseRuntimeError("unexpected data type")
			}
			algo := this.GetMember("__algo", c).ToGoValue()
			result, err := compressCompress(c, opt, algo, data)
			if err != nil {
				c.RaiseRuntimeError("compress error: %+v", err)
			}
			return NewBytes(result)
		}).
		Method("decompress", func(c *Context, this ValueObject, args []Value) Value {
			var (
				data ValueBytes
				opt  ValueObject
			)
			EnsureFuncParams(c, "decompress", args,
				ArgRuleRequired("data", TypeBytes, &data),
				ArgRuleOptional("opt", TypeObject, &opt, emptyOpt),
			)
			algo := this.GetMember("__algo", c).ToGoValue()
			result, err := compressDecompress(c, opt, algo, data.Value())
			if err != nil {
				c.RaiseRuntimeError("decompress error: %+v", err)
			}
			return NewBytes(result)
		}).
		Build()
}

func compressNewAlgo(c *Context, name string, algo any) ValueObject {
	return NewObjectAndInit(compressAlgorithmClass, c, NewStr(name), NewGoValue(algo))
}
