package builtin_libs

import (
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

func init() {
	var shCommand ValueType
	shCommand = NewClassBuilder("Command").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			this.SetMember("__cmds", args[0], c)
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			cmds := this.GetMember("__cmds", c).ToGoValue().([]*exec.Cmd)
			numCmd := len(cmds)
			if numCmd == 0 {
				return NewStr("")
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
					c.OnRuntimeError("start command %s error %s", cmd, err)
				}
			}
			outs, err := cmds[numCmd-1].CombinedOutput()
			if err != nil {
				c.OnRuntimeError("wait cmd %s error %s", cmds[numCmd-1], err)
			}
			return NewStr(string(outs))
		}).
		Method("__bitOr__", func(c *Context, this ValueObject, args []Value) Value {
			other := c.MustObject(args[0])
			cmds := append([]*exec.Cmd{}, this.GetMember("__cmds", c).ToGoValue().([]*exec.Cmd)...)
			cmds = append(cmds, other.GetMember("__cmds", c).ToGoValue().([]*exec.Cmd)...)
			return NewObjectAndInit(shCommand, c, NewGoValue(cmds))
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
				// 	c.OnRuntimeError("sh command %s error %s", name, err)
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
			EnsureFuncParams(c, "__getAttr__", args, ArgRuleRequired{"name", TypeStr, &cmd})
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
				ArgRuleRequired{"cmdName", TypeStr, &cmd},
				ArgRuleRequired{"session", TypeObject, &session},
			)
		}).
		Method("__call__", func(c *Context, this ValueObject, args []Value) Value {
			return nil
		}).
		Build()
}
