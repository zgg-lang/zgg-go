package builtin_libs

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libSys(c *Context) ValueObject {
	lib := NewObject()
	{
		args := NewArray(len(c.Args))
		for _, arg := range c.Args {
			args.PushBack(NewStr(arg))
		}
		lib.SetMember("args", args, nil)
	}
	{
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(os.Getenv("ZGG_RUN_ENV")), &m); err == nil {
			lib.SetMember("runEnv", jsonToValue(m, c), c)
		}
	}
	lib.SetMember("env", NewNativeFunction("sys.env", func(c *Context, this Value, args []Value) Value {
		switch len(args) {
		case 0:
			rv := NewObject()
			for _, item := range os.Environ() {
				kv := strings.SplitN(item, "=", 2)
				if len(kv) == 2 {
					rv.SetMember(kv[0], NewStr(kv[1]), c)
				}
			}
			return rv
		case 1:
			return NewStr(os.Getenv(args[0].ToString(c)))
		case 2:
			os.Setenv(args[0].ToString(c), args[1].ToString(c))
			return args[1]
		}
		c.OnRuntimeError("sys.env: invalid parameters")
		return nil
	}), nil)
	lib.SetMember("Command", sysCommandClass, c)
	lib.SetMember("createTempFile", NewNativeFunction("sys.createTempFile", func(c *Context, this Value, args []Value) Value {
		file, err := ioutil.TempFile("", "")
		if err != nil {
			c.OnRuntimeError("sys.createTempFile: create fail %s", err.Error())
			return nil
		}
		file.Close()
		res := NewObject()
		res.SetMember("name", NewStr(file.Name()), c)
		res.SetMember("drop", NewNativeFunction("drop", func(c *Context, this Value, args []Value) Value {
			os.Remove(file.Name())
			return Undefined()
		}), c)
		return res
	}), nil)
	lib.SetMember("getResult", NewNativeFunction("sys.getResult", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			c.OnRuntimeError("sys.getResult requires at least 1 argument")
		}
		name := args[0].ToString(c)
		cmdArgs := make([]string, len(args)-1)
		for i := range cmdArgs {
			cmdArgs[i] = args[i+1].ToString(c)
		}
		cmd := exec.Command(name, cmdArgs...)
		bs, err := cmd.CombinedOutput()
		if err != nil {
			c.OnRuntimeError("sys.getResult: command %s error %s", name, err)
		}
		return NewStr(string(bs))
	}), nil)
	return lib
}

var sysCommandClass = func() ValueType {
	captureOutput := func(rd io.ReadCloser, callable ValueCallable, c *Context) {
		defer rd.Close()
		var buf [64 * 1024]byte
		for {
			n, err := rd.Read(buf[:])
			if err != nil {
				break
			}
			c.Invoke(callable, Undefined(), Args(NewBytes(buf[:n])))
		}
	}
	return NewClassBuilder("sys.Command").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			if len(args) < 1 {
				c.OnRuntimeError("sys.Command.__init__(cmd, ...args) requires at lease 1 argument(s)")
				return
			}
			name := args[0].ToString(c)
			cmdArgs := make([]string, len(args)-1)
			for i := range cmdArgs {
				cmdArgs[i] = args[i+1].ToString(c)
			}
			cmd := exec.Command(name, cmdArgs...)
			this.SetMember("_name", args[0], c)
			this.SetMember("_cmd", NewGoValue(cmd), c)
			// this.SetMember("_stdin", NewGoValue(stdin), c)
			// this.SetMember("_stdout", NewGoValue(stdout), c)
			// this.SetMember("_stderr", NewGoValue(stderr), c)
			this.SetMember("_started", NewBool(false), c)
		}).
		Method("start", func(c *Context, this ValueObject, args []Value) Value {
			name := this.GetMember("_name", c)
			cmd := this.GetMember("_cmd", c).ToGoValue().(*exec.Cmd)
			stdin, err := cmd.StdinPipe()
			if err != nil {
				c.OnRuntimeError("sys.Command,input get stdin pipe error %s", err)
			}
			if err := cmd.Start(); err != nil {
				c.OnRuntimeError("sys.Command.start: command %s start fail %s", name.ToString(c), err)
				return nil
			}
			this.SetMember("_input", NewGoValue(stdin), c)
			this.SetMember("_started", NewBool(true), c)
			return this
		}).
		Method("input", func(c *Context, this ValueObject, args []Value) Value {
			stdin := this.GetMember("_input", c).ToGoValue().(io.WriteCloser)
			for _, arg := range args {
				var bs []byte
				switch a := arg.(type) {
				case ValueBytes:
					bs = a.Value()
				default:
					bs = []byte(a.ToString(c))
				}
				for i := 0; i < len(bs); {
					c.DebugLog("stdin %+v, bs %+v", stdin, bs)
					n, err := stdin.Write(bs[i:])
					if err != nil {
						c.OnRuntimeError("sys.Command.input write arg into stdin error %s", err)
					}
					i += n
				}
			}
			return Undefined()
		}).
		Method("onStdout", func(c *Context, this ValueObject, args []Value) Value {
			var callback ValueCallable
			EnsureFuncParams(c, "sys.Command.onStdout", args, ArgRuleRequired{"callback", TypeFunc, &callback})
			stdoutCallback := this.GetMember("_onStdout", c)
			if _, isUndefined := stdoutCallback.(ValueUndefined); !isUndefined {
				c.OnRuntimeError("sys.Command.onStdout: stdoutCallback already set")
				return nil
			}
			cmd := this.GetMember("_cmd", c).ToGoValue().(*exec.Cmd)
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				c.OnRuntimeError("sys.Command.onStdout: get stdout pipe error %s", err)
			}
			go captureOutput(stdout, callback, c)
			this.SetMember("_onStdout", callback, c)
			return this
		}).
		Method("onStderr", func(c *Context, this ValueObject, args []Value) Value {
			var callback ValueCallable
			EnsureFuncParams(c, "sys.Command.onStderr", args, ArgRuleRequired{"callback", TypeFunc, &callback})
			stderrCallback := this.GetMember("_onStderr", c)
			if _, isUndefined := stderrCallback.(ValueUndefined); !isUndefined {
				c.OnRuntimeError("sys.Command.onStderr: stderrCallback already set")
				return nil
			}
			cmd := this.GetMember("_cmd", c).ToGoValue().(*exec.Cmd)
			stderr, err := cmd.StderrPipe()
			if err != nil {
				c.OnRuntimeError("sys.Command.onStderr: get stderr pipe error %s", err)
			}
			go captureOutput(stderr, callback, c)
			this.SetMember("_onStderr", callback, c)
			return this
		}).
		Method("wait", func(c *Context, this ValueObject, args []Value) Value {
			name := this.GetMember("_name", c)
			cmd := this.GetMember("_cmd", c).ToGoValue().(*exec.Cmd)
			started := c.MustBool(this.GetMember("_started", c))
			if !started {
				if err := cmd.Start(); err != nil {
					c.OnRuntimeError("sys.Command.start: command %s start fail %s", name.ToString(c), err)
					return nil
				}
				this.SetMember("_started", NewBool(true), c)
			}
			if err := cmd.Wait(); err != nil {
				if ee, ok := err.(*exec.ExitError); ok {
					return NewInt(int64(ee.ExitCode()))
				}
				c.OnRuntimeError("sys.Command.wait: command %s error %s", name.ToString(c), err)
				return nil
			}
			return NewInt(0)
		}).
		Method("waitOutput", func(c *Context, this ValueObject, args []Value) Value {
			name := this.GetMember("_name", c)
			cmd := this.GetMember("_cmd", c).ToGoValue().(*exec.Cmd)
			bs, err := cmd.CombinedOutput()
			if err != nil {
				c.OnRuntimeError("sys.Command.waitOutput: command %s error %s", name.ToString(c), err)
				return nil
			}
			return NewBytes(bs)
		}).
		Method("kill", func(c *Context, this ValueObject, args []Value) Value {
			cmd := this.GetMember("_cmd", c).ToGoValue().(*exec.Cmd)
			cmd.Process.Kill()
			return Undefined()
		}).
		Build()
}()
