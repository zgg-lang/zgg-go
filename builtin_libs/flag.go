package builtin_libs

import (
	"flag"
	"time"

	"github.com/samber/lo"
	. "github.com/zgg-lang/zgg-go/runtime"
)

var flagParser ValueType

func libFlag(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Parser", flagGetParserClass(), nil)
	return lib
}

type flagParserInfo struct {
	name string
	fset *flag.FlagSet
	opts map[string]any
}

func flagProcessSingleVar[T any](c *Context, fname string, this ValueObject, args []Value) Value {
	var name, usage ValueStr
	EnsureFuncParams(c, "Parser."+fname, args,
		ArgRuleRequired("name", TypeStr, &name),
		ArgRuleOptional("usage", TypeStr, &usage, NewStr("")),
	)
	info := this.Reserved.(*flagParserInfo)
	var r any = new(T)
	switch rr := r.(type) {
	case *int:
		info.fset.IntVar(rr, name.Value(), 0, usage.Value())
	case *string:
		info.fset.StringVar(rr, name.Value(), "", usage.Value())
	case *float64:
		info.fset.Float64Var(rr, name.Value(), 0, usage.Value())
	case *bool:
		info.fset.BoolVar(rr, name.Value(), false, usage.Value())
	case *time.Duration:
		info.fset.DurationVar(rr, name.Value(), 0, usage.Value())
	}
	info.opts[name.Value()] = r
	return this
}

func flagGetParserClass() ValueType {
	if flagParser == nil {
		flagParser = NewClassBuilder("Parser").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				var name ValueStr
				EnsureFuncParams(c, "Parser.__init__", args,
					ArgRuleOptional("name", TypeStr, &name, NewStr("")),
				)
				info := &flagParserInfo{
					name: name.Value(),
					opts: make(map[string]any),
				}
				info.fset = flag.NewFlagSet(info.name, flag.ContinueOnError)
				this.Reserved = info
			}).
			Method("int", func(c *Context, this ValueObject, args []Value) Value {
				return flagProcessSingleVar[int](c, "int", this, args)
			}).
			Method("str", func(c *Context, this ValueObject, args []Value) Value {
				return flagProcessSingleVar[string](c, "str", this, args)
			}).
			Method("bool", func(c *Context, this ValueObject, args []Value) Value {
				return flagProcessSingleVar[bool](c, "bool", this, args)
			}).
			Method("float", func(c *Context, this ValueObject, args []Value) Value {
				return flagProcessSingleVar[float64](c, "float", this, args)
			}).
			Method("duration", func(c *Context, this ValueObject, args []Value) Value {
				return flagProcessSingleVar[time.Duration](c, "duration", this, args)
			}).
			Method("parse", func(c *Context, this ValueObject, args []Value) Value {
				info := this.Reserved.(*flagParserInfo)
				err := info.fset.Parse(lo.Map(args, func(a Value, _ int) string { return a.ToString(c) }))
				if err == flag.ErrHelp {
					return NewArrayByValues(NewBool(true))
				} else if err != nil {
					c.RaiseRuntimeError("parse argumetns error: %+v", err)
				}
				retOpts := NewObject()
				for k, v := range info.opts {
					switch vv := v.(type) {
					case *int:
						if vv != nil {
							retOpts.SetMember(k, NewInt(int64(*vv)), c)
						}
					case *string:
						if vv != nil {
							retOpts.SetMember(k, NewStr(*vv), c)
						}
					case *bool:
						if vv != nil {
							retOpts.SetMember(k, NewBool(*vv), c)
						}
					case *float64:
						if vv != nil {
							retOpts.SetMember(k, NewFloat(*vv), c)
						}
					case *time.Duration:
						if vv != nil {
							retOpts.SetMember(k, NewObjectAndInit(timeDurationClass, c, NewGoValue(*vv)), c)
						}
					}
				}
				retArgs := NewArrayByValues(lo.Map(info.fset.Args(), func(a string, _ int) Value {
					return NewStr(a)
				})...)
				return NewArrayByValues(NewBool(false), retOpts, retArgs)
			}).
			Build()
	}
	return flagParser
}
