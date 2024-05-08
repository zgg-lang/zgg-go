package builtin_libs

import (
	"fmt"
	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	dbFilterHelperClass          ValueType
	dbFilterHelperFieldClass     ValueType
	dbFilterHelperConditionClass ValueType
)

func init() {
	dbFilterHelperClass = NewClassBuilder("FilterHelper").
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			return NewObjectAndInit(dbFilterHelperFieldClass, c, args...)
		}).
		Build()
	dbFilterHelperFieldClass = func() ValueType {
		b := NewClassBuilder("FilterHelperField").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				var name ValueStr
				EnsureFuncParams(c, "FilterHelperField.__init__", args, ArgRuleRequired("name", TypeStr, &name))
				this.SetMember("name", name, c)
				this.SetMember("args", NewArrayByValues(args...), c)
			}).
			Method("encode", func(c *Context, this ValueObject, args []Value) Value {
				rargs := NewArray()
				return nil
			})
		for _, op := range []string{"eq", "ne", "gt", "ge", "lt", "le", "and", "or"} {
			(func(op string) {
				fn := fmt.Sprintf("__%s__", op)
				b.Method(fn, func(c *Context, this ValueObject, args []Value) Value {
					c.AssertArgNum(len(args), 1, 1, fn)
					return NewObjectAndInit(dbFilterHelperConditionClass, c, NewStr(op), this.GetMember("name", c), args[0])
				})
			})(op)
		}
		return b.Build()
	}()
	dbFilterHelperConditionClass = NewClassBuilder("FilterHelperCondition").
		Build()
}
