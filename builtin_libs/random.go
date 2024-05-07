package builtin_libs

import (
	"math/rand"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libRandom(*Context) ValueObject {
	lib := NewObject()
	//DOC random() 		返回[0, 1)范围内的一个随机浮点数
	//DOC random(end)	返回[0, end)范围内的一个随机整数，当end<=0或者end非整数时抛出异常
	//DOC random(a, b)	返回[a, b)范围内的一个随机整数，当a或者b非整数，或者a>=b时抛出异常
	lib.SetMember("random", NewNativeFunction("random", randomRandomFunc), nil)
	//DOC __call__() 		返回[0, 1)范围内的一个随机浮点数
	//DOC __call__(end)		返回[0, end)范围内的一个随机整数，当end<=0或者end非整数时抛出异常
	//DOC __call__(a, b)	返回[a, b)范围内的一个随机整数，当a或者b非整数，或者a>=b时抛出异常
	lib.SetMember("__call__", NewNativeFunction("__call__", randomRandomFunc), nil)
	//DOC choice(array) 	返回array内的一个随机元素，当array非数组或者数组为空时抛出异常
	lib.SetMember("choice", NewNativeFunction("choice", func(c *Context, this Value, args []Value) Value {
		var choices ValueArray
		EnsureFuncParams(c, "random.choice", args,
			ArgRuleRequired("choices", TypeArray, &choices),
		)
		if choices.Len() < 1 {
			c.RaiseRuntimeError("random.choice(arr): arr cannot be empty")
			return nil
		}
		n := rand.Intn(choices.Len())
		return choices.GetIndex(n, c)
	}), nil)
	//DOC shuffle(array) 	打算array内元素顺序，当array非数组时抛出异常
	lib.SetMember("shuffle", NewNativeFunction("shuffle", func(c *Context, this Value, args []Value) Value {
		var (
			arr ValueArray
		)
		EnsureFuncParams(c, "shuffle", args, ArgRuleRequired("array", TypeArray, &arr))
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

func randomRandomFunc(c *Context, this Value, args []Value) Value {
	switch len(args) {
	case 0:
		return NewFloat(rand.Float64())
	case 1:
		n := c.MustInt(args[0], "random:end")
		if n <= 0 {
			c.RaiseRuntimeError("random(end): end must > 0")
		}
		return NewInt(rand.Int63n(n))
	case 2:
		m := c.MustInt(args[0], "random:began")
		n := c.MustInt(args[1], "random:end")
		if n <= m {
			c.RaiseRuntimeError("random(begin, end): end must > begin")
		}
		return NewInt(rand.Int63n(n-m) + m)
	default:
		c.RaiseRuntimeError("random: requires 0/1/2 argument(s)")
	}
	return nil
}
