package builtin_libs

import (
	"math/rand"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libRandom(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("random", NewNativeFunction("random", func(c *Context, this Value, args []Value) Value {
		switch len(args) {
		case 0:
			return NewFloat(rand.Float64())
		case 1:
			n := c.MustInt(args[0], "random:end")
			if n <= 0 {
				c.OnRuntimeError("random(end): end must > 0")
			}
			return NewInt(rand.Int63n(n))
		case 2:
			m := c.MustInt(args[0], "random:began")
			n := c.MustInt(args[1], "random:end")
			if n <= m {
				c.OnRuntimeError("random(begin, end): end must > begin")
			}
			return NewInt(rand.Int63n(n-m) + m)
		default:
			c.OnRuntimeError("random: requires 0/1/2 argument(s)")
		}
		return nil
	}), nil)
	lib.SetMember("choice", NewNativeFunction("choice", func(c *Context, this Value, args []Value) Value {
		var choices ValueArray
		EnsureFuncParams(c, "random.choice", args,
			ArgRuleRequired{"choices", TypeArray, &choices},
		)
		if choices.Len() < 1 {
			c.OnRuntimeError("random.choice(arr): arr cannot be empty")
			return nil
		}
		n := rand.Intn(choices.Len())
		return choices.GetIndex(n, c)
	}), nil)
	lib.SetMember("shuffle", NewNativeFunction("shuffle", func(c *Context, this Value, args []Value) Value {
		var (
			arr ValueArray
		)
		EnsureFuncParams(c, "shuffle", args, ArgRuleRequired{"array", TypeArray, &arr})
		n := arr.Len()
		for i := n - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			if j != i {
				t := arr.GetIndex(i, c)
				arr.SetIndex(i, arr.GetIndex(j, c), c)
				arr.SetIndex(j, t, c)
			}
		}
		return arr
	}), nil)
	return lib
}
