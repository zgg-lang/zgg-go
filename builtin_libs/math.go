package builtin_libs

import (
	"math"

	"github.com/samber/lo"
	. "github.com/zgg-lang/zgg-go/runtime"
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
