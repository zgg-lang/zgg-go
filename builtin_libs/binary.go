package builtin_libs

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	binaryReaderClass ValueType
	binaryWriterClass ValueType
)

func libBinary(*Context) ValueObject {
	lib := NewObject()
	binaryReaderClass = binaryInitReaderClass()
	binaryWriterClass = binaryInitWriterClass()
	lib.SetMember("Reader", binaryReaderClass, nil)
	lib.SetMember("Writer", binaryWriterClass, nil)
	return lib
}

type (
	binaryReader interface {
		io.Reader
		io.ByteReader
	}
	binaryReaderInfo struct {
		byteOrder binary.ByteOrder
		rd        binaryReader
	}
	binaryWriterInfo struct {
		byteOrder binary.ByteOrder
		w         io.Writer
	}
)

func binaryInitReaderClass() ValueType {
	return NewClassBuilder("Reader").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			if len(args) != 1 {
				c.RaiseRuntimeError("Reader.__init__ requires 1 argument")
			}
			var rd binaryReader
			switch i := args[0].(type) {
			case ValueBytes:
				rd = bytes.NewReader(i.Value())
			case GoValue:
				switch ii := i.ToGoValue(c).(type) {
				case []byte:
					rd = bytes.NewReader(ii)
				case binaryReader:
					rd = ii
				case io.Reader:
					rd = bufio.NewReader(ii)
				}
			}
			if rd == nil {
				c.RaiseRuntimeError("Reader.__init__ invalid input type")
			}
			this.Reserved = &binaryReaderInfo{
				byteOrder: binary.LittleEndian,
				rd:        rd,
			}
		}).
		Methods([]string{"useBigEndian", "big"}, func(c *Context, this ValueObject, args []Value) Value {
			this.Reserved.(*binaryReaderInfo).byteOrder = binary.BigEndian
			return this
		}).
		Methods([]string{"useLittleEndian", "little"}, func(c *Context, this ValueObject, args []Value) Value {
			this.Reserved.(*binaryReaderInfo).byteOrder = binary.LittleEndian
			return this
		}).
		Method("i8", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			if b, err := info.rd.ReadByte(); err == io.EOF {
				return Nil()
			} else if err != nil {
				c.RaiseRuntimeError("read i8 error %+v", err)
				return nil
			} else {
				return NewInt(int64(int8(b)))
			}
		}).
		Method("u8", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			if b, err := info.rd.ReadByte(); err == io.EOF {
				return Nil()
			} else if err != nil {
				c.RaiseRuntimeError("read u8 error %+v", err)
				return nil
			} else {
				return NewInt(int64(b))
			}
		}).
		Method("i16", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			var buf [2]byte
			if _, err := io.ReadFull(info.rd, buf[:]); err != nil {
				if err == io.EOF {
					return Nil()
				}
				c.RaiseRuntimeError("read i16 error %+v", err)
			}
			return NewInt(int64(int16(info.byteOrder.Uint16(buf[:]))))
		}).
		Method("u16", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			var buf [2]byte
			if _, err := io.ReadFull(info.rd, buf[:]); err != nil {
				if err == io.EOF {
					return Nil()
				}
				c.RaiseRuntimeError("read i16 error %+v", err)
			}
			return NewInt(int64(info.byteOrder.Uint16(buf[:])))
		}).
		Method("i32", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			var buf [4]byte
			if _, err := io.ReadFull(info.rd, buf[:]); err != nil {
				if err == io.EOF {
					return Nil()
				}
				c.RaiseRuntimeError("read i32 error %+v", err)
			}
			return NewInt(int64(int32(info.byteOrder.Uint32(buf[:]))))
		}).
		Method("u32", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			var buf [4]byte
			if _, err := io.ReadFull(info.rd, buf[:]); err != nil {
				if err == io.EOF {
					return Nil()
				}
				c.RaiseRuntimeError("read i32 error %+v", err)
			}
			return NewInt(int64(info.byteOrder.Uint32(buf[:])))
		}).
		Method("i64", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			var buf [8]byte
			if _, err := io.ReadFull(info.rd, buf[:]); err != nil {
				if err == io.EOF {
					return Nil()
				}
				c.RaiseRuntimeError("read i64 error %+v", err)
			}
			return NewInt(int64(info.byteOrder.Uint64(buf[:])))
		}).
		Method("varint", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			if v, err := binary.ReadVarint(info.rd); err != nil {
				c.RaiseRuntimeError("read varint error %+v", err)
				return nil
			} else {
				return NewInt(v)
			}
		}).
		Method("uvarint", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryReaderInfo)
			if v, err := binary.ReadVarint(info.rd); err != nil {
				c.RaiseRuntimeError("read uvarint error %+v", err)
				return nil
			} else {
				return NewInt(int64(v))
			}
		}).
		Method("bytes", func(c *Context, this ValueObject, args []Value) Value {
			var n ValueInt
			EnsureFuncParams(c, "Reader.bytes", args, ArgRuleRequired("n", TypeInt, &n))
			nv := n.Value()
			if nv < 0 {
				c.RaiseRuntimeError("Reader.bytes: n must >= 0")
			}
			if nv == 0 {
				return NewBytes([]byte{})
			}
			out := make([]byte, int(nv))
			info := this.Reserved.(*binaryReaderInfo)
			if _, err := io.ReadFull(info.rd, out); err != nil {
				if err == io.EOF {
					return Nil()
				}
				c.RaiseRuntimeError("read %d bytes error %+v", nv, err)
			}
			return NewBytes(out)
		}).
		Build()
}

func binaryInitWriterClass() ValueType {
	writeAll := func(w io.Writer, bs []byte) error {
		o, s := 0, len(bs)
		for o < s {
			n, err := w.Write(bs[o:])
			if err != nil {
				return err
			}
			o += n
		}
		return nil
	}
	return NewClassBuilder("Writer").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				capacity ValueInt
				writer   GoValue
				selected int
			)
			EnsureFuncParams(c, "Writer.__init__", args,
				ArgRuleOneOf("v",
					[]ValueType{TypeInt, TypeGoValue},
					[]any{&capacity, &writer},
					&selected,
					&capacity,
					NewInt(0),
				),
			)
			var w io.Writer
			switch selected {
			case 2:
				if wv, is := writer.ToGoValue(c).(io.Writer); !is {
					c.RaiseRuntimeError("Writer.__init__: not a writer")
				} else {
					w = wv
				}
			default:
				w = bytes.NewBuffer(make([]byte, capacity.AsInt()))
			}
			this.Reserved = &binaryWriterInfo{
				byteOrder: binary.LittleEndian,
				w:         w,
			}
		}).
		Methods([]string{"useBigEndian", "big"}, func(c *Context, this ValueObject, args []Value) Value {
			this.Reserved.(*binaryWriterInfo).byteOrder = binary.BigEndian
			return this
		}).
		Methods([]string{"useLittleEndian", "little"}, func(c *Context, this ValueObject, args []Value) Value {
			this.Reserved.(*binaryWriterInfo).byteOrder = binary.LittleEndian
			return this
		}).
		Method("i8", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "Writer.i8", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			buf := [1]byte{byte(int8(v.AsInt()))}
			if e := writeAll(info.w, buf[:]); e != nil {
				c.RaiseRuntimeError("Writer.i8: write error %+v", e)
			}
			return this
		}).
		Method("u8", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "Writer.u8", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			buf := [1]byte{byte(uint8(v.AsInt()))}
			if e := writeAll(info.w, buf[:]); e != nil {
				c.RaiseRuntimeError("Writer.u8: write error %+v", e)
			}
			return this
		}).
		Method("i16", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "Writer.i16", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			var buf [2]byte
			info.byteOrder.PutUint16(buf[:], uint16(int16(v.Value())))
			if e := writeAll(info.w, buf[:]); e != nil {
				c.RaiseRuntimeError("Writer.i16: write error %+v", e)
			}
			return this
		}).
		Method("u16", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "Writer.u16", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			var buf [2]byte
			info.byteOrder.PutUint16(buf[:], uint16(v.Value()))
			if e := writeAll(info.w, buf[:]); e != nil {
				c.RaiseRuntimeError("Writer.u16: write error %+v", e)
			}
			return this
		}).
		Method("i32", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "Writer.i32", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			var buf [4]byte
			info.byteOrder.PutUint32(buf[:], uint32(int32(v.Value())))
			if e := writeAll(info.w, buf[:]); e != nil {
				c.RaiseRuntimeError("Writer.i32: write error %+v", e)
			}
			return this
		}).
		Method("u32", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "Writer.u32", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			var buf [4]byte
			info.byteOrder.PutUint32(buf[:], uint32(v.Value()))
			if e := writeAll(info.w, buf[:]); e != nil {
				c.RaiseRuntimeError("Writer.u32: write error %+v", e)
			}
			return this
		}).
		Method("i64", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "Writer.i64", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			var buf [8]byte
			info.byteOrder.PutUint64(buf[:], uint64(v.Value()))
			if e := writeAll(info.w, buf[:]); e != nil {
				c.RaiseRuntimeError("Writer.i64: write error %+v", e)
			}
			return this
		}).
		Method("varint", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			var buf [binary.MaxVarintLen64]byte
			EnsureFuncParams(c, "Writer.varint", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			n := binary.PutVarint(buf[:], v.Value())
			if e := writeAll(info.w, buf[:n]); e != nil {
				c.RaiseRuntimeError("Writer.varint: write error %+v", e)
			}
			return this
		}).
		Method("uvarint", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			var buf [binary.MaxVarintLen64]byte
			EnsureFuncParams(c, "Writer.uvarint", args, ArgRuleRequired("v", TypeInt, &v))
			info := this.Reserved.(*binaryWriterInfo)
			n := binary.PutUvarint(buf[:], uint64(v.Value()))
			if e := writeAll(info.w, buf[:n]); e != nil {
				c.RaiseRuntimeError("Writer.uvarint: write error %+v", e)
			}
			return this
		}).
		Method("bytes", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueBytes
			EnsureFuncParams(c, "Writer.bytes", args, ArgRuleRequired("v", TypeBytes, &v))
			info := this.Reserved.(*binaryWriterInfo)
			if e := writeAll(info.w, v.Value()); e != nil {
				c.RaiseRuntimeError("Writer.bytes: write error %+v", e)
			}
			return this
		}).
		Method("toBytes", func(c *Context, this ValueObject, args []Value) Value {
			info := this.Reserved.(*binaryWriterInfo)
			if b, is := info.w.(*bytes.Buffer); !is {
				c.RaiseRuntimeError("Writer.toBytes: not a byte buffer")
				return nil
			} else {
				return NewBytes(b.Bytes())
			}
		}).
		Build()
}
