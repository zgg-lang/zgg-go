package builtin_libs

import (
	"fmt"
	"math"
	"strings"

	"github.com/samber/lo"
	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	mathVecClass ValueType
)

func libMath(*Context) ValueObject {
	r := NewObject()
	r.SetMember("abs", NewNativeFunction("abs", func(c *Context, _ Value, args []Value) Value {
		var (
			xi ValueInt
			xf ValueFloat
			xt int
		)
		EnsureFuncParams(c, "math.abs", args,
			ArgRuleOneOf("x", []ValueType{TypeInt, TypeFloat}, []any{&xi, &xf}, &xt, nil, nil),
		)
		switch xt {
		case 0: // int
			x := xi.Value()
			return NewInt(lo.Ternary(x >= 0, x, -x))
		case 1: // float
			return NewFloat(math.Abs(xf.Value()))
		}
		return nil
	}), nil)
	r.SetMember("gcd", NewNativeFunction("gcd", func(c *Context, _ Value, args []Value) Value {
		return NewInt(mathGCD(c, "gcd", args))
	}), nil)
	r.SetMember("lcm", NewNativeFunction("lcm", func(c *Context, _ Value, args []Value) Value {
		gcd := mathGCD(c, "lcm", args)
		r := int64(1)
		for i, a := range args {
			r *= a.(ValueInt).Value()
			if i == 0 {
				r /= gcd
			}
		}
		return NewInt(r)
	}), nil)
	r.SetMember("reduction", NewNativeFunction("reduction", func(c *Context, _ Value, args []Value) Value {
		gcd := mathGCD(c, "reduction", args)
		rv := NewArray(len(args))
		for _, a := range args {
			rv.PushBack(NewInt(a.(ValueInt).Value() / gcd))
		}
		return rv
	}), nil)
	r.SetMember("P", NewNativeFunction("P", func(c *Context, _ Value, args []Value) Value {
		var an, ar ValueInt
		EnsureFuncParams(c, "P", args, ArgRuleRequired("n", TypeInt, &an), ArgRuleRequired("r", TypeInt, &ar))
		n, r := an.AsInt(), ar.AsInt()
		if r <= 0 || n <= r {
			c.RaiseRuntimeError("invalid arguments! n %d r %d", n, r)
		}
		rv := 1
		for i := n - r + 1; i <= n; i++ {
			rv *= i
		}
		return NewInt(int64(rv))

	}, "n", "r"), nil)
	r.SetMember("C", NewNativeFunction("C", func(c *Context, _ Value, args []Value) Value {
		var an, ar ValueInt
		EnsureFuncParams(c, "C", args, ArgRuleRequired("n", TypeInt, &an), ArgRuleRequired("r", TypeInt, &ar))
		n, r := an.AsInt(), ar.AsInt()
		if r <= 0 || n <= r {
			c.RaiseRuntimeError("invalid arguments! n %d r %d", n, r)
		}
		if r+r < n {
			r = n - r
		}
		rv := 1
		for i := r + 1; i <= n; i++ {
			rv *= i
		}
		for i := n - r; i > 1; i-- {
			rv /= i
		}
		return NewInt(int64(rv))

	}, "n", "r"), nil)
	for name, fnv := range mathPortFuncs {
		switch fn := fnv.(type) {
		case func(float64) float64:
			r.SetMember(name, mathPortFuncF_F(name, fn), nil)
		case func(float64) (float64, float64):
			r.SetMember(name, mathPortFuncF_FF(name, fn), nil)
		case func(float64, float64) float64:
			r.SetMember(name, mathPortFuncFF_F(name, fn), nil)
		case func(float64, float64, float64) float64:
			r.SetMember(name, mathPortFuncFFF_F(name, fn), nil)
		}
	}
	for name, v := range mathPortConsts {
		switch vv := v.(type) {
		case float64:
			r.SetMember(name, NewFloat(vv), nil)
		}
	}
	r.SetMember("Vec", mathVecClass, nil)
	return r
}

func mathPortFuncF_F(name string, fn func(float64) float64) ValueCallable {
	return NewNativeFunction(name, func(ctx *Context, _ Value, args []Value) Value {
		var a ValueFloat
		EnsureFuncParams(ctx, "math."+name, args,
			ArgRuleRequired("a", TypeFloat, &a),
		)
		return NewFloat(fn(a.Value()))
	})
}

func mathPortFuncFF_F(name string, fn func(float64, float64) float64) ValueCallable {
	return NewNativeFunction(name, func(ctx *Context, _ Value, args []Value) Value {
		var a, b ValueFloat
		EnsureFuncParams(ctx, "math."+name, args,
			ArgRuleRequired("a", TypeFloat, &a),
			ArgRuleRequired("b", TypeFloat, &b),
		)
		return NewFloat(fn(a.Value(), b.Value()))
	})
}

func mathPortFuncFFF_F(name string, fn func(float64, float64, float64) float64) ValueCallable {
	return NewNativeFunction(name, func(ctx *Context, _ Value, args []Value) Value {
		var a, b, c ValueFloat
		EnsureFuncParams(ctx, "math."+name, args,
			ArgRuleRequired("a", TypeFloat, &a),
			ArgRuleRequired("b", TypeFloat, &b),
			ArgRuleRequired("c", TypeFloat, &c),
		)
		return NewFloat(fn(a.Value(), b.Value(), c.Value()))
	})
}

func mathPortFuncF_FF(name string, fn func(float64) (float64, float64)) ValueCallable {
	return NewNativeFunction(name, func(ctx *Context, _ Value, args []Value) Value {
		var a ValueFloat
		EnsureFuncParams(ctx, "math."+name, args,
			ArgRuleRequired("a", TypeFloat, &a),
		)
		r1, r2 := fn(a.Value())
		return NewArrayByValues(NewFloat(r1), NewFloat(r2))
	})
}

var mathPortFuncs = map[string]any{
	"acos":            math.Acos,
	"acosh":           math.Acosh,
	"asin":            math.Asin,
	"asinh":           math.Asinh,
	"atan":            math.Atan,
	"atan2":           math.Atan2,
	"atanh":           math.Atanh,
	"cbrt":            math.Cbrt,
	"ceil":            math.Ceil,
	"copysign":        math.Copysign,
	"cos":             math.Cos,
	"cosh":            math.Cosh,
	"dim":             math.Dim,
	"erf":             math.Erf,
	"erfc":            math.Erfc,
	"erfcinv":         math.Erfcinv,
	"erfinv":          math.Erfinv,
	"exp":             math.Exp,
	"exp2":            math.Exp2,
	"expm1":           math.Expm1,
	"fma":             math.FMA,
	"float32bits":     math.Float32bits,
	"float32frombits": math.Float32frombits,
	"float64bits":     math.Float64bits,
	"float64frombits": math.Float64frombits,
	"floor":           math.Floor,
	"frexp":           math.Frexp,
	"gamma":           math.Gamma,
	"hypot":           math.Hypot,
	"ilogb":           math.Ilogb,
	"inf":             math.Inf,
	"isInf":           math.IsInf,
	"isNaN":           math.IsNaN,
	"j0":              math.J0,
	"j1":              math.J1,
	"jn":              math.Jn,
	"ldexp":           math.Ldexp,
	"lgamma":          math.Lgamma,
	"log":             math.Log,
	"log10":           math.Log10,
	"log1p":           math.Log1p,
	"log2":            math.Log2,
	"logb":            math.Logb,
	"max":             math.Max,
	"min":             math.Min,
	"mod":             math.Mod,
	"modf":            math.Modf,
	"naN":             math.NaN,
	"nextafter":       math.Nextafter,
	"nextafter32":     math.Nextafter32,
	"pow":             math.Pow,
	"pow10":           math.Pow10,
	"remainder":       math.Remainder,
	"round":           math.Round,
	"roundToEven":     math.RoundToEven,
	"signbit":         math.Signbit,
	"sin":             math.Sin,
	"sincos":          math.Sincos,
	"sinh":            math.Sinh,
	"sqrt":            math.Sqrt,
	"tan":             math.Tan,
	"tanh":            math.Tanh,
	"trunc":           math.Trunc,
	"y0":              math.Y0,
	"y1":              math.Y1,
	"yn":              math.Yn,
}

var mathPortConsts = map[string]any{
	"PI": math.Pi,
	"E":  math.E,
}

func mathGCD2(a, b int64) int64 {
	for b > 0 {
		t := a % b
		a = b
		b = t
	}
	return a
}

func mathMustGetPositiveInt(c *Context, a Value) int64 {
	if ai, is := a.(ValueInt); !is {
		c.RaiseRuntimeError("argument must be an integer")
		return -1
	} else if vi := ai.Value(); vi <= 0 {
		c.RaiseRuntimeError("argument must be larger than 0")
		return -1
	} else {
		return vi
	}
}

func mathGCD(c *Context, f string, args []Value) int64 {
	n := len(args)
	if n < 2 {
		c.RaiseRuntimeError("math.%s: requires at lease 2 arguments", f)
		return -1
	}
	a := mathMustGetPositiveInt(c, args[0])
	b := mathMustGetPositiveInt(c, args[1])
	r := mathGCD2(a, b)
	for i := 2; i < n; i++ {
		a = r
		b = mathMustGetPositiveInt(c, args[i])
		r = mathGCD2(a, b)
	}
	return r
}

type mathVecInfo struct {
	components []float64
}

func mathInitVecClass() ValueType {
	initVec := func(components []float64, vec ValueObject) {
		vec.Reserved = &mathVecInfo{
			components: components,
		}
	}
	checkBinArgs := func(c *Context, name string, this ValueObject, args []Value) (a, b []float64) {
		var other ValueObject
		EnsureFuncParams(c, "Vec."+name, args, ArgRuleRequired("other", mathVecClass, &other))
		a = this.Reserved.(*mathVecInfo).components
		b = other.Reserved.(*mathVecInfo).components
		if len(a) != len(b) {
			c.RaiseRuntimeError("Vec.%s: size not equal", name)
		}
		return
	}
	return NewClassBuilder("Vec").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			c.AssertArgNum(len(args), 1, 99999, "Vec.__init__")
			components := lo.Map(args, func(a Value, _ int) float64 {
				return c.MustFloat(a)
			})
			initVec(components, this)
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			var b strings.Builder
			components := this.Reserved.(*mathVecInfo).components
			b.WriteRune('(')
			fmt.Fprint(&b, components[0])
			for i := 1; i < len(components); i++ {
				b.WriteString(", ")
				fmt.Fprint(&b, components[i])
			}
			b.WriteRune(')')
			return NewStr(b.String())
		}).
		Method("__add__", func(c *Context, this ValueObject, args []Value) Value {
			a, b := checkBinArgs(c, "__add__", this, args)
			result := NewObject(mathVecClass)
			initVec(lo.Map(a, func(c1 float64, i int) float64 { return c1 + b[i] }), result)
			return result
		}).
		Method("__sub__", func(c *Context, this ValueObject, args []Value) Value {
			a, b := checkBinArgs(c, "__sub__", this, args)
			result := NewObject(mathVecClass)
			initVec(lo.Map(a, func(c1 float64, i int) float64 { return c1 - b[i] }), result)
			return result
		}).
		Build()
}

func init() {
	mathVecClass = mathInitVecClass()
}
