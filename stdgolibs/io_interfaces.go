package stdgolibs

import (
	"io"

	. "github.com/zgg-lang/zgg-go/runtime"
)

// 万能interface代理
type ioDelegate interfaceDelegate

func (d *ioDelegate) ReadByte() (b byte, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "ReadByte",
		[]interface{}{},
		[]interface{}{&b, &err},
	)
	return
}

func (d *ioDelegate) UnreadByte() (err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "UnreadByte",
		[]interface{}{},
		[]interface{}{&err},
	)
	return
}

func (d *ioDelegate) WriteByte(c byte) (err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "WriteByte",
		[]interface{}{c},
		[]interface{}{&err},
	)
	return
}

func (d *ioDelegate) Close() (err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "Close",
		[]interface{}{},
		[]interface{}{&err},
	)
	return
}

func (d *ioDelegate) Read(p []byte) (n int, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "Read",
		[]interface{}{p},
		[]interface{}{&n, &err},
	)
	return
}

func (d *ioDelegate) ReadAt(p []byte, off int64) (n int, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "ReadAt",
		[]interface{}{p, off},
		[]interface{}{&n, &err},
	)
	return
}

func (d *ioDelegate) ReadFrom(r io.Reader) (n int64, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "ReadAt",
		[]interface{}{r},
		[]interface{}{&n, &err},
	)
	return
}

func (d *ioDelegate) ReadRune() (r rune, size int, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "ReadRune",
		[]interface{}{},
		[]interface{}{&r, &size, &err},
	)
	return
}

func (d *ioDelegate) UnreadRune() (err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "UnreadRune",
		[]interface{}{},
		[]interface{}{&err},
	)
	return
}

func (d *ioDelegate) Seek(offset int64, whence int) (n int64, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "Seek",
		[]interface{}{offset, whence},
		[]interface{}{&n, &err},
	)
	return
}

func (d *ioDelegate) WriteString(s string) (n int, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "WriteString",
		[]interface{}{s},
		[]interface{}{&n, &err},
	)
	return
}

func (d *ioDelegate) Write(p []byte) (n int, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "Write",
		[]interface{}{p},
		[]interface{}{&n, &err},
	)
	return
}

func (d *ioDelegate) WriteAt(p []byte, off int64) (n int, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "WriteAt",
		[]interface{}{p, off},
		[]interface{}{&n, &err},
	)
	return
}

func (d *ioDelegate) WriteTo(w io.Writer) (n int64, err error) {
	interfaceMethodBridge(
		d.c, d.obj, d, "WriteTo",
		[]interface{}{w},
		[]interface{}{&n, &err},
	)
	return
}

var (
	_ io.ByteReader      = &ioDelegate{}
	_ io.ByteScanner     = &ioDelegate{}
	_ io.ByteWriter      = &ioDelegate{}
	_ io.Closer          = &ioDelegate{}
	_ io.ReadCloser      = &ioDelegate{}
	_ io.ReadSeekCloser  = &ioDelegate{}
	_ io.ReadSeeker      = &ioDelegate{}
	_ io.ReadWriteCloser = &ioDelegate{}
	_ io.ReadWriteSeeker = &ioDelegate{}
	_ io.ReadWriter      = &ioDelegate{}
	_ io.Reader          = &ioDelegate{}
	_ io.ReaderAt        = &ioDelegate{}
	_ io.ReaderFrom      = &ioDelegate{}
	_ io.RuneReader      = &ioDelegate{}
	_ io.RuneScanner     = &ioDelegate{}
	_ io.Seeker          = &ioDelegate{}
	_ io.StringWriter    = &ioDelegate{}
	_ io.WriteCloser     = &ioDelegate{}
	_ io.WriteSeeker     = &ioDelegate{}
	_ io.Writer          = &ioDelegate{}
	_ io.WriterAt        = &ioDelegate{}
	_ io.WriterTo        = &ioDelegate{}
)

func init() {
	registerFuncs("io", map[string]*ValueBuiltinFunction{
		"Delegate": NewNativeFunction("io.Delegate", func(c *Context, this Value, args []Value) Value {
			var obj ValueObject
			EnsureFuncParams(c, "Delegate", args, ArgRuleRequired("delegate", TypeObject, &obj))
			return NewGoValue(&ioDelegate{c: c, obj: obj})
		}),
	})
}
