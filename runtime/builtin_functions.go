package runtime

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ValueBuiltinFunction struct {
	ValueBase
	name string
	args []string
	body func(*Context, Value, []Value) Value
}

func NewNativeFunction(name string, body func(*Context, Value, []Value) Value, args ...string) *ValueBuiltinFunction {
	return &ValueBuiltinFunction{
		name: name,
		body: body,
		args: args,
	}
}

func (f *ValueBuiltinFunction) GoType() reflect.Type {
	return reflect.TypeOf(f.body)
}

func (f *ValueBuiltinFunction) GetIndex(int, *Context) Value {
	return constUndefined
}

func (f *ValueBuiltinFunction) GetMember(name string, c *Context) Value {
	switch name {
	case "__name__":
		if f.name != "" {
			return NewStr(f.name)
		}
		return NewStr("anoymous")
	case "__args__":
		{
			args := NewArray(len(f.args))
			for _, arg := range f.args {
				args.PushBack(NewStr(arg))
			}
			return args
		}
	}
	return getCallableMember(f, name, c)
}

func (f *ValueBuiltinFunction) Type() ValueType {
	return TypeFunc
}

func (f *ValueBuiltinFunction) IsTrue() bool {
	return true
}

func (f *ValueBuiltinFunction) ToString(*Context) string {
	return fmt.Sprintf("<function %s>", f.name)
}

func (f *ValueBuiltinFunction) ToGoValue() interface{} {
	return f.body
}

func (f *ValueBuiltinFunction) CompareTo(other Value, c *Context) CompareResult {
	return CompareResultNotEqual
}

func (f *ValueBuiltinFunction) GetName() string {
	return f.name
}

func (f *ValueBuiltinFunction) GetArgNames() []string {
	return f.args
}

func (f *ValueBuiltinFunction) GetRefs() []string {
	return []string{}
}

func (f *ValueBuiltinFunction) Invoke(c *Context, thisArg Value, args []Value) {
	c.PushFuncStack(f.name)
	defer c.PopStack()
	c.RetVal = f.body(c, thisArg, args)
}

func buildJson(v Value, c *Context) interface{} {
	switch val := v.(type) {
	case ValueInt:
		return val.ToGoValue()
	case ValueFloat:
		return val.ToGoValue()
	case ValueBool:
		return val.ToGoValue()
	case ValueObject:
		{
			rv := make(map[string]interface{})
			val.Each(func(k string, v Value) bool {
				if !strings.HasPrefix(k, "__") || !strings.HasSuffix(k, "__") {
					rv[k] = buildJson(v, c)
				}
				return true
			})
			return rv
		}
	case ValueArray:
		{
			rv := make([]interface{}, val.Len())
			for i := range rv {
				rv[i] = buildJson(val.GetIndex(i, c), c)
			}
			return rv
		}
	}
	return v.ToString(c)
}

func getInt(v Value) (int, bool) {
	if iv, ok := v.(ValueInt); ok {
		return int(iv.Value()), true
	}
	return 0, false
}

func getArgInt(args []Value, n int) (int, bool) {
	if len(args) < n+1 {
		return 0, false
	}
	return getInt(args[n])
}

func getArgStr(args []Value, n int) (string, bool) {
	if len(args) < n+1 {
		return "", false
	}
	if v, ok := args[n].(ValueStr); !ok {
		return "", false
	} else {
		return v.Value(), true
	}
}

func mustGetArgStr(c *Context, name string, args []Value, n int) string {
	if s, ok := getArgStr(args, n); ok {
		return s
	}
	c.OnRuntimeError("get string argument at position %d fail", n)
	return ""
}

func mustGetArgInt(c *Context, name string, args []Value, n int) int {
	if s, ok := getArgInt(args, n); ok {
		return s
	}
	c.OnRuntimeError("get string argument at position %d fail", n)
	return -1
}

var builtinFunctions = map[string]ValueCallable{
	"println": &ValueBuiltinFunction{
		name: "println",
		body: func(c *Context, thisArg Value, args []Value) Value {
			printArgs := make([]interface{}, len(args))
			for i, arg := range args {
				printArgs[i] = arg.ToString(c)
			}
			fmt.Fprintln(c.Stdout, printArgs...)
			return NewInt(0)
		},
	},
	"printf": &ValueBuiltinFunction{
		name: "printf",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				c.OnRuntimeError("printf requires at least 1 argument")
				return nil
			}
			printFmt := args[0].ToString(c)
			args = args[1:]
			printArgs := make([]interface{}, len(args))
			for i, arg := range args {
				printArgs[i] = arg.ToGoValue()
			}
			fmt.Fprintf(c.Stdout, printFmt, printArgs...)
			return NewInt(0)
		},
	},
	"sprintf": &ValueBuiltinFunction{
		name: "sprintf",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				c.OnRuntimeError("printf requires at least 1 argument")
				return nil
			}
			printFmt := args[0].ToString(c)
			args = args[1:]
			printArgs := make([]interface{}, len(args))
			for i, arg := range args {
				printArgs[i] = arg.ToGoValue()
			}
			return NewStr(fmt.Sprintf(printFmt, printArgs...))
		},
	},
	"int": &ValueBuiltinFunction{
		name: "int",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				return NewInt(0)
			}
			switch v := args[0].(type) {
			case ValueInt:
				return args[0]
			case ValueFloat:
				return NewInt(int64(v.Value()))
			case ValueStr:
				{
					vs := v.Value()
					var vi int64
					var err error
					if strings.HasPrefix(vs, "0x") || strings.HasPrefix(vs, "0X") {
						vi, err = strconv.ParseInt(vs[2:], 16, 64)
					} else if strings.HasPrefix(vs, "0b") || strings.HasPrefix(vs, "0b") {
						vi, err = strconv.ParseInt(vs[2:], 2, 64)
					} else if len(vs) > 1 && strings.HasPrefix(vs, "0") || strings.HasPrefix(vs, "0") {
						vi, err = strconv.ParseInt(vs[1:], 8, 64)
					} else {
						vi, err = strconv.ParseInt(vs, 10, 64)
					}
					if err == nil {
						return NewInt(vi)
					}
				}
			case ValueBool:
				if v.Value() {
					return NewInt(1)
				}
			}
			return NewInt(0)
		},
	},
	"float": &ValueBuiltinFunction{
		name: "float",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				return NewFloat(0)
			}
			switch v := args[0].(type) {
			case ValueInt:
				return NewFloat(float64(v.Value()))
			case ValueFloat:
				return args[0]
			case ValueStr:
				{
					vs := v.Value()
					vf, err := strconv.ParseFloat(vs, 64)
					if err == nil {
						return NewFloat(vf)
					}
				}
			case ValueBool:
				if v.Value() {
					return NewFloat(1)
				}
			}
			return NewFloat(0)
		},
	},
	"str": &ValueBuiltinFunction{
		name: "str",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				return NewStr("")
			}
			return NewStr(args[0].ToString(c))
		},
	},
	"bytes": &ValueBuiltinFunction{
		name: "bytes",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				return NewBytes([]byte{})
			}
			return NewBytes([]byte(args[0].ToString(c)))
		},
	},
	"typeName": &ValueBuiltinFunction{
		name: "typeName",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				return NewStr("undefined")
			}
			return NewStr(args[0].Type().ToString(c))
		},
	},
	"type": &ValueBuiltinFunction{
		name: "type",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) < 1 {
				return TypeUndefined
			}
			return args[0].Type()
		},
	},
	"len": &ValueBuiltinFunction{
		name: "len",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				return NewInt(-1)
			}
			if v, ok := args[0].(CanLen); ok {
				return NewInt(int64(v.Len()))
			}
			c.OnRuntimeError("len function's argument cannot be %s", args[0].Type().Name)
			return nil
		},
	},
	"range": &ValueBuiltinFunction{
		name: "range",
		body: func(c *Context, thisArg Value, args []Value) Value {
			begin, end, step := 0, -1, 1
			ok := true
			switch len(args) {
			case 3:
				if step, ok = getInt(args[2]); !ok {
					c.OnRuntimeError("range arg 2 must be an integer")
					return nil
				}
				if step == 0 {
					c.OnRuntimeError("range argument step cannot be 0")
					return nil
				}
				fallthrough
			case 2:
				if begin, ok = getInt(args[0]); !ok {
					c.OnRuntimeError("range arg 0 must be an integer")
					return nil
				}
				if end, ok = getInt(args[1]); !ok {
					c.OnRuntimeError("range arg 1 must be an integer")
					return nil
				}
			case 1:
				if end, ok = getInt(args[0]); !ok {
					c.OnRuntimeError("range arg 0 must be an integer")
					return nil
				}
			}
			if step < 0 {
				if begin <= end {
					c.OnRuntimeError("range when step < 0, begin must be greater than end")
					return nil
				}
			} else {
				if begin >= end {
					c.OnRuntimeError("range when step > 0, begin must be less than end")
					return nil
				}
			}
			rv := NewArray()
			for i := begin; i < end; i += step {
				rv.PushBack(NewInt(int64(i)))
			}
			return rv
		},
	},
	"seq": NewNativeFunction("seq", func(c *Context, this Value, args []Value) Value {
		if len(args) != 2 {
			c.OnRuntimeError("seq requires 2 arguments")
			return nil
		}
		next := args[0]
		last := args[1]
		rv := NewArray()
		rv.PushBack(next)
		for {
			nextFn := next.GetMember("__next__", c)
			if !c.IsCallable(nextFn) {
				c.OnRuntimeError("not all the items in seq has __next__ method")
			}
			c.Invoke(nextFn, next, Args())
			next = c.RetVal
			rv.PushBack(next)
			if c.ValuesEqual(next, last) {
				break
			}
		}
		return rv
	}),
	"sha1": &ValueBuiltinFunction{
		name: "sha1",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				c.OnRuntimeError("sha1 requires only one argument")
				return nil
			}
			var bs []byte
			switch arg := args[0].(type) {
			case ValueStr:
				bs = []byte(arg.Value())
			case ValueBytes:
				bs = arg.Value()
			default:
				c.OnRuntimeError("sha1 module arg must be a string or bytes")
				return nil
			}
			res := sha1.Sum(bs)
			return NewStr(hex.EncodeToString(res[:]))
		},
	},
	"sha256": &ValueBuiltinFunction{
		name: "sha256",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				c.OnRuntimeError("sha256 requires only one argument")
				return nil
			}
			var bs []byte
			switch arg := args[0].(type) {
			case ValueStr:
				bs = []byte(arg.Value())
			case ValueBytes:
				bs = arg.Value()
			default:
				c.OnRuntimeError("sha256 module arg must be a string or bytes")
				return nil
			}
			res := sha256.Sum256(bs)
			return NewStr(hex.EncodeToString(res[:]))
		},
	},
	"md5": &ValueBuiltinFunction{
		name: "md5",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				c.OnRuntimeError("md5 requires only one argument")
				return nil
			}
			var bs []byte
			switch arg := args[0].(type) {
			case ValueStr:
				bs = []byte(arg.Value())
			case ValueBytes:
				bs = arg.Value()
			default:
				c.OnRuntimeError("md5 module arg must be a string or bytes")
				return nil
			}
			res := md5.Sum(bs)
			return NewStr(hex.EncodeToString(res[:]))
		},
	},
	"import": NewNativeFunction("import", func(c *Context, thisArg Value, args []Value) Value {
		importType := "script"
		forceReload := false
		var modName string
		switch len(args) {
		case 3:
			importType = c.MustStr(args[2], "import importType")
			fallthrough
		case 2:
			forceReload = c.MustBool(args[1], "import forceReload")
			fallthrough
		case 1:
			modName = c.MustStr(args[0], "import modName")
		default:
			c.OnRuntimeError("import requires only one or two argument(s)")
			return nil
		}
		return c.ImportModule(modName, forceReload, importType)
	}, "name", "force", "type"),
	"isUndefined": &ValueBuiltinFunction{
		name: "isUndefined",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				c.OnRuntimeError("isUndefined requires only one argument")
				return nil
			}
			_, rv := args[0].(ValueUndefined)
			return NewBool(rv)
		},
	},
	"isCallable": &ValueBuiltinFunction{
		name: "isCallable",
		body: func(c *Context, thisArg Value, args []Value) Value {
			if len(args) != 1 {
				c.OnRuntimeError("isCallable requires only one argument")
				return nil
			}
			return NewBool(c.IsCallable(args[0]))
		},
	},
	"isArray": NewNativeFunction("isArray", func(c *Context, thisArg Value, args []Value) Value {
		if len(args) == 1 {
			if _, ok := args[0].(ValueArray); ok {
				return NewBool(true)
			}
		}
		return NewBool(false)
	}),
	"isObject": NewNativeFunction("isObject", func(c *Context, thisArg Value, args []Value) Value {
		if len(args) == 1 {
			if _, ok := args[0].(ValueObject); ok {
				return NewBool(true)
			}
		}
		return NewBool(false)
	}),
	"assertError": NewNativeFunction("assertError", func(c *Context, this Value, args []Value) Value {
		if len(args) == 0 {
			return constUndefined
		}
		if arr, isArr := args[0].(ValueArray); isArr && arr.Len() > 0 {
			l := arr.Len()
			if l > 0 {
				lastVal := arr.GetIndex(l-1, c)
				if err, isErr := lastVal.ToGoValue().(error); isErr {
					c.OnRuntimeError("assertError failed: %s", err)
					return nil
				}
				l--
			}
			if l == 0 {
				return constUndefined
			} else if l == 1 {
				return arr.GetIndex(0, c)
			} else {
				rv := NewArray(l)
				for i := 0; i < l; i++ {
					rv.PushBack(arr.GetIndex(i, c))
				}
				return rv
			}
		}
		if err, isErr := args[0].ToGoValue().(error); isErr {
			c.OnRuntimeError("assertError failed: %s", err)
			return nil
		}
		return args[0]
	}),
	"rand": &ValueBuiltinFunction{
		name: "rand",
		body: func(c *Context, thisArg Value, args []Value) Value {
			switch len(args) {
			case 0:
				return NewFloat(rand.Float64())
			case 1:
				if maxVal, isInt := args[0].(ValueInt); isInt {
					max := int(maxVal.Value())
					if max <= 0 {
						c.OnRuntimeError(fmt.Sprintf("rand(n): expected n > 0, got %d", max))
						return nil
					}
					return NewInt(int64(rand.Intn(max)))
				} else {
					c.OnRuntimeError("rand(n): n must be an integer")
					return nil
				}
			case 2:
				{
					var min, max int
					if minVal, isInt := args[0].(ValueInt); isInt {
						min = int(minVal.Value())
					} else {
						c.OnRuntimeError("rand(m, n): m must be an integer")
						return nil
					}
					if maxVal, isInt := args[1].(ValueInt); isInt {
						max = int(maxVal.Value())
					} else {
						c.OnRuntimeError("rand(m, n): n must be an integer")
						return nil
					}
					if max <= min {
						c.OnRuntimeError(fmt.Sprintf("rand(m, n): expected n > m, got m=%d, n=%d", min, max))
						return nil
					}
					rv := rand.Intn(max-min) + min
					return NewInt(int64(rv))
				}
			default:
				c.OnRuntimeError("rand requires 0 or 1 or 2 argument")
				return nil
			}
		},
	},
	"eval": NewNativeFunction("eval", func(c *Context, this Value, args []Value) Value {
		var code string
		evalCtx := c
		switch len(args) {
		case 2:
			if sandbox, ok := args[1].(ValueObject); ok {
				evalCtx = NewContext(false, c.IsDebug, c.CanEval)
				evalCtx.ImportFunc = c.ImportFunc
				evalCtx.Stdin = c.Stdin
				evalCtx.Stdout = c.Stdout
				evalCtx.Stderr = c.Stderr
				sandbox.Iterate(func(k string, v Value) {
					evalCtx.ForceSetLocalValue(k, v)
				})
			} else {
				c.OnRuntimeError("eval(code, [sandbox]): sandbox must be an object")
			}
			fallthrough
		case 1:
			code = args[0].ToString(c)
		default:
			c.OnRuntimeError("eval requires 1 argument")
		}
		return evalCtx.Eval(code, false)
	}),
	"bind": NewNativeFunction("bind", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			c.OnRuntimeError("bind requires at least 1 argument")
		}
		f := c.MustCallable(args[0])
		args = args[1:]
		return NewNativeFunction(f.GetName(), func(c *Context, this Value, args2 []Value) Value {
			c.Invoke(f, nil, func() []Value {
				if len(args2) == 0 {
					return args
				}
				r := make([]Value, len(args)+len(args2))
				p := 0
				for _, a := range args {
					r[p] = a
					p++
				}
				for _, a := range args2 {
					r[p] = a
					p++
				}
				return r
			})
			return c.RetVal
		})
	}),
	"max": NewNativeFunction("max", func(c *Context, this Value, args []Value) Value {
		n := len(args)
		if n == 0 {
			return constUndefined
		}
		r := args[0]
		for i := 1; i < n; i++ {
			if c.ValuesGreater(args[i], r) {
				r = args[i]
			}
		}
		return r
	}),
	"min": NewNativeFunction("min", func(c *Context, this Value, args []Value) Value {
		n := len(args)
		if n == 0 {
			return constUndefined
		}
		r := args[0]
		for i := 1; i < n; i++ {
			if c.ValuesLess(args[i], r) {
				r = args[i]
			}
		}
		return r
	}),
	"input": NewNativeFunction("input", func(c *Context, this Value, args []Value) Value {
		var prompt ValueStr
		EnsureFuncParams(c, "input", args, ArgRuleOptional{"prompt", TypeStr, &prompt, NewStr("")})
		fmt.Fprint(c.Stdout, prompt.Value())
		scanner := bufio.NewScanner(c.Stdin)
		if scanner.Scan() {
			return NewStr(scanner.Text())
		}
		return constUndefined
	}, "prompt"),
	"log": (func() ValueObject {
		rv := NewObject()
		const (
			DEBUG = iota
			INFO
			WARN
			ERROR
			ASSERT
		)
		var doLog = func(level int, c *Context, args []Value) {
			if len(args) < 1 {
				return
			}
			var minLevel int
			if v, isInt := rv.GetMember("logLevel", c).(ValueInt); isInt {
				minLevel = v.AsInt()
			}
			if level < minLevel {
				return
			}
			handler := rv.GetMember("handler", c)
			if c.IsCallable(handler) {
				var stackArgs [10]Value
				var a []Value
				n := len(args)
				if n < 10 {
					a = stackArgs[:n+1]
				} else {
					a = make([]Value, n+1)
				}
				a[0] = NewInt(int64(level))
				for i, arg := range args {
					a[i+1] = arg
				}
				c.Invoke(handler.(ValueCallable), nil, Args(a...))
				if c.RetVal.IsTrue() {
					return
				}
			}
			var b strings.Builder
			b.WriteString(time.Now().Format("2006-01-02 15:04:05"))
			switch level {
			case DEBUG:
				b.WriteString("|DBG|")
			case INFO:
				b.WriteString("|INF|")
			case WARN:
				b.WriteString("|WRN|")
			case ERROR:
				b.WriteString("|ERR|")
			case ASSERT:
				b.WriteString("|AST|")
			}
			for i, arg := range args {
				if i > 0 {
					b.WriteRune(' ')
				}
				b.WriteString(arg.ToString(c))
			}
			b.WriteRune('\n')
			var logWriter io.Writer
			if w, ok := rv.GetMember("writer", c).ToGoValue().(io.Writer); ok {
				logWriter = w
			} else {
				logWriter = c.Stdout
			}
			io.WriteString(logWriter, b.String())
		}
		var bindLevel = func(level int) ValueCallable {
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				doLog(level, c, args)
				return constUndefined
			})
		}
		rv.SetMember("DEBUG", NewInt(DEBUG), nil)
		rv.SetMember("INFO", NewInt(INFO), nil)
		rv.SetMember("WARN", NewInt(WARN), nil)
		rv.SetMember("ERROR", NewInt(ERROR), nil)
		rv.SetMember("ASSERT", NewInt(ASSERT), nil)
		rv.SetMember("__call__", bindLevel(DEBUG), nil)
		rv.SetMember("debug", bindLevel(DEBUG), nil)
		rv.SetMember("info", bindLevel(INFO), nil)
		rv.SetMember("warn", bindLevel(WARN), nil)
		rv.SetMember("error", bindLevel(ERROR), nil)
		rv.SetMember("assert", bindLevel(ASSERT), nil)
		return rv
	})(),
}

func init() {
	rand.Seed(time.Now().UnixNano() ^ int64(os.Getpid()))
}
