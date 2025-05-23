package builtin_libs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os/exec"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	shModuleClass  ValueType
	shClassSession ValueType
	shClassCommand ValueType
)

func libSh(c *Context) ValueObject {
	rv := NewObjectAndInit(shModuleClass, c)
	return rv
}

func shGetCommandOutput(c *Context, cmds []*exec.Cmd) []byte {
	numCmd := len(cmds)
	if numCmd == 0 {
		return nil
	}
	var inCmd = cmds[0]
	closers := make([]io.Closer, 0, numCmd)
	defer func() {
		for _, c := range closers {
			c.Close()
		}
	}()
	for i := 1; i < numCmd; i++ {
		pipe, _ := inCmd.StdoutPipe()
		closers = append(closers, pipe)
		cmds[i].Stdin = pipe
		inCmd = cmds[i]
	}
	for _, cmd := range cmds[:numCmd-1] {
		if err := cmd.Start(); err != nil {
			c.RaiseRuntimeError("start command %s error %s", cmd, err)
		}
	}
	outs, err := cmds[numCmd-1].CombinedOutput()
	if err != nil {
		c.RaiseRuntimeError("wait cmd %s error %s", cmds[numCmd-1], err)
	}
	return outs
}

func init() {
	var shCommand ValueType
	shCommand = NewClassBuilder("Command").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			this.SetMember("__cmds", args[0], c)
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			cmds := this.GetMember("__cmds", c).ToGoValue(c).([]*exec.Cmd)
			outs := shGetCommandOutput(c, cmds)
			return NewStr(string(outs))
		}).
		Method("text", func(c *Context, this ValueObject, args []Value) Value {
			cmds := this.GetMember("__cmds", c).ToGoValue(c).([]*exec.Cmd)
			outs := shGetCommandOutput(c, cmds)
			return NewStr(string(outs))
		}).
		Method("lines", func(c *Context, this ValueObject, args []Value) Value {
			cmds := this.GetMember("__cmds", c).ToGoValue(c).([]*exec.Cmd)
			outs := shGetCommandOutput(c, cmds)
			rd := bufio.NewReader(bytes.NewReader(outs))
			rv := NewArray()
			for {
				line, err := rd.ReadString('\n')
				if len(line) > 0 {
					rv.PushBack(NewStr(line[:len(line)-1]))
				}
				if err != nil {
					break
				}
			}
			return rv
		}).
		Method("json", func(c *Context, this ValueObject, args []Value) Value {
			cmds := this.GetMember("__cmds", c).ToGoValue(c).([]*exec.Cmd)
			outs := shGetCommandOutput(c, cmds)
			var j interface{}
			if err := json.Unmarshal(outs, &j); err != nil {
				c.RaiseRuntimeError("Command.json(): decode json error %s", err)
			}
			return jsonToValue(j, c)
		}).
		Method("__bitOr__", func(c *Context, this ValueObject, args []Value) Value {
			if other, is := args[0].(ValueObject); is {
				cmds := []*exec.Cmd{}
				cmds = append(cmds, this.GetMember("__cmds", c).ToGoValue(c).([]*exec.Cmd)...)
				cmds = append(cmds, other.GetMember("__cmds", c).ToGoValue(c).([]*exec.Cmd)...)
				return NewObjectAndInit(shCommand, c, NewGoValue(cmds))
			}
			if other, is := c.GetCallable(args[0]); is {
				c.Invoke(other, nil, Args(this))
				return c.RetVal
			}
			c.RaiseRuntimeError("invalid right value type")
			return nil
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			switch args[0].ToString(c) {
			case "s":
				c.InvokeMethod(this, "text", NoArgs)
			case "l":
				c.InvokeMethod(this, "lines", NoArgs)
			case "j":
				c.InvokeMethod(this, "json", NoArgs)
			}
			return c.RetVal
		}).
		Build()
	shModuleClass = NewClassBuilder("sh").
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			name := args[0].ToString(c)
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				cmdArgs := make([]string, len(args))
				for i := range cmdArgs {
					cmdArgs[i] = args[i].ToString(c)
				}
				// cmd := exec.Command(name, cmdArgs...)
				// bs, err := cmd.CombinedOutput()
				// if err != nil {
				// 	c.RaiseRuntimeError("sh command %s error %s", name, err)
				// }
				// return NewStr(string(bs))
				cmd := exec.Command(name, cmdArgs...)
				return NewObjectAndInit(shCommand, c, NewGoValue([]*exec.Cmd{cmd}))
			})
		}).
		Build()
	shClassSession = NewClassBuilder("Session").
		Constructor(func(c *Context, this ValueObject, args []Value) {

		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			var cmd ValueStr
			EnsureFuncParams(c, "__getAttr__", args, ArgRuleRequired("name", TypeStr, &cmd))
			return NewObjectAndInit(shClassCommand, c, cmd, this)
		}).
		Build()
	shClassCommand = NewClassBuilder("Command").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				cmd     ValueStr
				session ValueObject
			)
			EnsureFuncParams(c, "Command.__init__", args,
				ArgRuleRequired("cmdName", TypeStr, &cmd),
				ArgRuleRequired("session", TypeObject, &session),
			)
		}).
		Method("__call__", func(c *Context, this ValueObject, args []Value) Value {
			return nil
		}).
		Build()
}
