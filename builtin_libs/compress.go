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
	println(len(data))
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

type compressAlgoInfo struct {
	Name ValueStr
	Algo any
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
			a := algo.ToGoValue(c)
			_, canCompress := a.(compressCompressor)
			_, canDecompress := a.(compressDecompressor)
			if !canCompress && !canDecompress {
				c.RaiseRuntimeError("algo is neither a Compressor nor a Decompressor")
			}
			this.Reserved = compressAlgoInfo{
				Name: name,
				Algo: a,
			}
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			var field ValueStr
			EnsureFuncParams(c, "Time.__getAttr__", args, ArgRuleRequired("field", TypeStr, &field))
			switch field.Value() {
			case "name":
				return this.Reserved.(compressAlgoInfo).Name
			}
			return Undefined()
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			return this.Reserved.(compressAlgoInfo).Name
		}).
		Method("__call__", func(c *Context, this ValueObject, args []Value) Value {
			return c.InvokeMethod(this, "compress", Args(args...))
		}).
		Method("compress", func(c *Context, this ValueObject, args []Value) Value {
			var (
				dataBytes ValueBytes
				dataStr   ValueStr
				dataWhich int
				data      []byte
				opt       ValueObject
				optLevel  ValueInt
				optWhich  int
			)
			EnsureFuncParams(c, "compress", args,
				ArgRuleOneOf("data",
					[]ValueType{TypeBytes, TypeStr},
					[]any{&dataBytes, &dataStr},
					&dataWhich, nil, nil),
				ArgRuleOneOf("opt",
					[]ValueType{TypeObject, TypeInt},
					[]any{&opt, &optLevel},
					&optWhich, &opt, emptyOpt),
			)
			switch dataWhich {
			case 0:
				data = dataBytes.Value()
			case 1:
				data = []byte(dataStr.Value())
			default:
				c.RaiseRuntimeError("unexpected data type")
			}
			switch optWhich {
			case 1:
				opt = NewObject()
				opt.SetMember("level", optLevel, c)
			}
			algo := this.Reserved.(compressAlgoInfo).Algo
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
			algo := this.Reserved.(compressAlgoInfo).Algo
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
