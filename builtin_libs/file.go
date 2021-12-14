package builtin_libs

import (
	"bufio"
	"io/ioutil"
	"os"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	fileFileReaderClass ValueType
)

func libFile(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("readBytes", NewNativeFunction("file.readBytes", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.OnRuntimeError("file.readBytes requires 1 argument")
			return nil
		}
		filename := c.MustStr(args[0])
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			c.OnRuntimeError("read file error: %s", err.Error())
			return nil
		}
		return NewBytes(bytes)
	}), nil)
	lib.SetMember("readString", NewNativeFunction("file.readString", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.OnRuntimeError("file.readString requires 1 argument")
			return nil
		}
		filename := c.MustStr(args[0])
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			c.OnRuntimeError("read file error: %s", err.Error())
			return nil
		}
		return NewStr(string(bytes))
	}), nil)
	lib.SetMember("readLines", NewNativeFunction("file.readLines", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.OnRuntimeError("file.readLines requires 1 argument")
			return nil
		}
		filename := c.MustStr(args[0])
		file, err := os.Open(filename)
		if err != nil {
			c.OnRuntimeError("open file error: %s", err.Error())
			return nil
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		rv := NewArray()
		for scanner.Scan() {
			rv.PushBack(NewStr(scanner.Text()))
		}
		return rv
	}), nil)
	lib.SetMember("rewrite", NewNativeFunction("file.rewrite", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			c.OnRuntimeError("file.rewrite requires at lease 1 argument")
			return nil
		}
		filename := c.MustStr(args[0])
		file, err := os.Create(filename)
		if err != nil {
			c.OnRuntimeError("create file error: %s", err.Error())
			return nil
		}
		defer file.Close()
		total := 0
		for i := 1; i < len(args); i++ {
			var bs []byte
			switch v := args[i].(type) {
			case ValueBytes:
				bs = v.Value()
			default:
				bs = []byte(v.ToString(c))
			}
			for written := 0; written < len(bs); {
				n, err := file.Write(bs[written:])
				if err != nil {
					c.OnRuntimeError("write file error: %s", err)
					return nil
				}
				written += n
			}
			total += len(bs)
		}
		return NewInt(int64(total))
	}), nil)
	lib.SetMember("append", NewNativeFunction("file.append", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			c.OnRuntimeError("file.append requires at lease 1 argument")
			return nil
		}
		filename := c.MustStr(args[0])
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			c.OnRuntimeError("create file error: %s", err.Error())
			return nil
		}
		defer file.Close()
		total := 0
		for i := 1; i < len(args); i++ {
			var bs []byte
			switch v := args[i].(type) {
			case ValueBytes:
				bs = v.Value()
			default:
				bs = []byte(v.ToString(c))
			}
			for written := 0; written < len(bs); {
				n, err := file.Write(bs[written:])
				if err != nil {
					c.OnRuntimeError("write file error: %s", err)
					return nil
				}
				written += n
			}
			total += len(bs)
		}
		return NewInt(int64(total))
	}), nil)
	return lib
}

func initFileClass() {
	fileFileReaderClass = NewClassBuilder("FileReader").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			this.SetMember("_file", args[0], c)
		}).
		Method("nextLine", func(c *Context, this ValueObject, args []Value) Value {

			return nil
		}).
		Build()
}

func init() {

}
