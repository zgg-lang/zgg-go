package builtin_libs

import (
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
	shModuleClass = NewClassBuilder("sh").
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			name := args[0].ToString(c)
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				cmdArgs := make([]string, len(args))
				for i := range cmdArgs {
					cmdArgs[i] = args[i].ToString(c)
				}
				cmd := exec.Command(name, cmdArgs...)
				bs, err := cmd.CombinedOutput()
				if err != nil {
					c.OnRuntimeError("sh command %s error %s", name, err)
				}
				return NewStr(string(bs))
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
