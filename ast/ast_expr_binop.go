package ast

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/zgg-lang/zgg-go/runtime"
)

type BinOp struct {
	Left, Right Expr
}

func (binOp *BinOp) GetValues(c *runtime.Context) (runtime.Value, runtime.Value) {
	binOp.Left.Eval(c)
	left := ensureZgg(c.RetVal, c)
	binOp.Right.Eval(c)
	right := ensureZgg(c.RetVal, c)
	return left, right
}

func (binOp *BinOp) tryOverride(c *runtime.Context, fn string) (left runtime.Value, right runtime.Value, ret runtime.Value) {
	left, right = binOp.GetValues(c)
	if opFn, ok := left.GetMember(fn, c).(runtime.ValueCallable); ok {
		c.Invoke(opFn, left, runtime.Args(right))
		ret = c.RetVal
	}
	return
}

type ExprPlus struct {
	BinOp
}

func (n *ExprPlus) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	if _, isStr := right.(runtime.ValueStr); isStr {
		c.RetVal = runtime.NewStr(left.ToString(c) + right.ToString(c))
		return
	}
	switch val1 := left.(type) {
	case runtime.ValueInt:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			c.RetVal = runtime.NewInt(val1.Value() + val2.Value())
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(float64(val1.Value()) + val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Add(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueFloat:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			c.RetVal = runtime.NewFloat(val1.Value() + float64(val2.Value()))
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(val1.Value() + val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Add(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueBigNum:
		var big2, bigR big.Float
		switch val2 := right.(type) {
		case runtime.ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Add(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Add(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueBigNum:
			var bigR big.Float
			bigR.Add(val1.Value(), val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueStr:
		c.RetVal = runtime.NewStr(val1.Value() + right.ToString(c))
		return
	case runtime.ValueArray:
		switch val2 := right.(type) {
		case runtime.ValueArray:
			{
				rv := runtime.NewArray(val1.Len() + val2.Len())
				for i := 0; i < val1.Len(); i++ {
					rv.PushBack(val1.GetIndex(i, c))
				}
				for i := 0; i < val2.Len(); i++ {
					rv.PushBack(val2.GetIndex(i, c))
				}
				c.RetVal = rv
				return
			}
		}
	default:
		{
			if opFn, ok := val1.GetMember("__add__", c).(runtime.ValueCallable); ok {
				c.Invoke(opFn, val1, func() []runtime.Value { return []runtime.Value{right} })
				return
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot plus between %s and %s", left.Type().Name, right.Type().Name))
}

type ExprMinus struct {
	BinOp
}

func (n *ExprMinus) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	switch val1 := left.(type) {
	case runtime.ValueInt:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			c.RetVal = runtime.NewInt(val1.Value() - val2.Value())
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(float64(val1.Value()) - val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Sub(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueFloat:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			c.RetVal = runtime.NewFloat(val1.Value() - float64(val2.Value()))
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(val1.Value() - val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Sub(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueBigNum:
		var big2, bigR big.Float
		switch val2 := right.(type) {
		case runtime.ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Sub(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Sub(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueBigNum:
			var bigR big.Float
			bigR.Sub(val1.Value(), val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	default:
		{
			if opFn, ok := val1.GetMember("__sub__", c).(runtime.ValueCallable); ok {
				c.Invoke(opFn, val1, func() []runtime.Value { return []runtime.Value{right} })
				return
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot minus between %s and %s", left.Type().Name, right.Type().Name))
}

type ExprTimes struct {
	BinOp
}

func (n *ExprTimes) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	switch val1 := left.(type) {
	case runtime.ValueInt:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			c.RetVal = runtime.NewInt(val1.Value() * val2.Value())
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(float64(val1.Value()) * val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Mul(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueFloat:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			c.RetVal = runtime.NewFloat(val1.Value() * float64(val2.Value()))
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(val1.Value() * val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Mul(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueBigNum:
		var big2, bigR big.Float
		switch val2 := right.(type) {
		case runtime.ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Mul(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Mul(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueBigNum:
			var bigR big.Float
			bigR.Mul(val1.Value(), val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueStr:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			{
				var sb strings.Builder
				item := val1.Value()
				for times := int(val2.Value()); times > 0; times-- {
					sb.WriteString(item)
				}
				c.RetVal = runtime.NewStr(sb.String())
				return
			}
		}
	case runtime.ValueArray:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			{
				rv := runtime.NewArray(val1.Len() * val2.AsInt())
				for times := val2.AsInt(); times > 0; times-- {
					for i := 0; i < val1.Len(); i++ {
						rv.PushBack(val1.GetIndex(i, c))
					}
				}
				c.RetVal = rv
				return
			}
		}
	default:
		{
			if opFn, ok := val1.GetMember("__mul__", c).(runtime.ValueCallable); ok {
				c.Invoke(opFn, val1, func() []runtime.Value { return []runtime.Value{right} })
				return
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot times between %s and %s", left.Type().Name, right.Type().Name))
}

type ExprDiv struct {
	BinOp
}

func (n *ExprDiv) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	switch val1 := left.(type) {
	case runtime.ValueInt:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("division by zero")
				return
			}
			c.RetVal = runtime.NewInt(val1.Value() / val2.Value())
			return
		case runtime.ValueFloat:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return
			}
			c.RetVal = runtime.NewFloat(float64(val1.Value()) / val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(float64(val1.Value())).SetPrec(1024)
			var bigR big.Float
			bigR.Quo(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueFloat:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return
			}
			c.RetVal = runtime.NewFloat(val1.Value() / float64(val2.Value()))
			return
		case runtime.ValueFloat:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return
			}
			c.RetVal = runtime.NewFloat(val1.Value() / val2.Value())
			return
		case runtime.ValueBigNum:
			big1 := big.NewFloat(val1.Value()).SetPrec(1024)
			var bigR big.Float
			bigR.Quo(big1, val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	case runtime.ValueBigNum:
		var big2, bigR big.Float
		big2.SetPrec(1024)
		switch val2 := right.(type) {
		case runtime.ValueInt:
			big2.SetInt64(val2.Value())
			bigR.Quo(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueFloat:
			big2.SetFloat64(val2.Value())
			bigR.Quo(val1.Value(), &big2)
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		case runtime.ValueBigNum:
			var bigR big.Float
			bigR.Quo(val1.Value(), val2.Value())
			c.RetVal = runtime.NewBigNum(&bigR)
			return
		}
	default:
		{
			if opFn, ok := val1.GetMember("__div__", c).(runtime.ValueCallable); ok {
				c.Invoke(opFn, val1, func() []runtime.Value { return []runtime.Value{right} })
				return
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot div between %s and %s", left.Type().Name, right.Type().Name))
}

type ExprMod struct {
	BinOp
}

func (n *ExprMod) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	switch val1 := left.(type) {
	case runtime.ValueInt:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			if val2.Value() == 0 {
				c.RaiseRuntimeError("Div by zero!")
				return
			}
			c.RetVal = runtime.NewInt(val1.Value() % val2.Value())
			return
		}
	case runtime.ValueStr:
		switch val2 := right.(type) {
		case runtime.ValueArray:
			fargs := make([]interface{}, val2.Len())
			for i := range fargs {
				fargs[i] = val2.GetIndex(i, c).ToGoValue()
			}
			c.RetVal = runtime.NewStr(fmt.Sprintf(val1.Value(), fargs...))
		default:
			c.RetVal = runtime.NewStr(fmt.Sprintf(val1.Value(), val2.ToGoValue()))
		}
		return
	default:
		{
			if opFn, ok := val1.GetMember("__mod__", c).(runtime.ValueCallable); ok {
				c.Invoke(opFn, val1, func() []runtime.Value { return []runtime.Value{right} })
				return
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot mod between %s and %s", left.Type().Name, right.Type().Name))
}

type ExprPow struct {
	BinOp
}

func (n *ExprPow) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	switch val1 := left.(type) {
	case runtime.ValueInt:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			if val2.Value() >= 0 {
				c.RetVal = runtime.NewInt(int64(math.Pow(float64(val1.Value()), float64(val2.Value()))))
			} else {
				c.RetVal = runtime.NewFloat(math.Pow(float64(val1.Value()), float64(val2.Value())))
			}
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(math.Pow(float64(val1.Value()), val2.Value()))
			return
		}
	case runtime.ValueFloat:
		switch val2 := right.(type) {
		case runtime.ValueInt:
			c.RetVal = runtime.NewFloat(math.Pow(val1.Value(), float64(val2.Value())))
			return
		case runtime.ValueFloat:
			c.RetVal = runtime.NewFloat(float64(val1.Value()) - val2.Value())
			return
		}
	default:
		{
			if opFn, ok := val1.GetMember("__pow__", c).(runtime.ValueCallable); ok {
				c.Invoke(opFn, val1, func() []runtime.Value { return []runtime.Value{right} })
				return
			}
		}
	}
	c.RaiseRuntimeError(fmt.Sprintf("Cannot pow between %s and %s", left.Type().Name, right.Type().Name))
}

type IsAssign interface {
	itIsAnAssignExpr()
}

type ExprAssign struct {
	Pos
	Lval Lval
	Expr Expr
}

func (ExprAssign) itIsAnAssignExpr() {}

func (assign *ExprAssign) Eval(c *runtime.Context) {
	c.EnsureNotReadonly()
	assign.Expr.Eval(c)
	assign.Lval.SetValue(c, c.RetVal)
}

const (
	AssignTypeSingle = iota
	AssignTypeDeArray
	AssignTypeDeObject
)

type ExprLocalNewAssign struct {
	Pos
	Expr Expr
}

func (ExprLocalNewAssign) itIsAnAssignExpr() {}

func (assign *ExprLocalNewAssign) Eval(c *runtime.Context) {
	c.EnsureNotReadonly()
	assign.Expr.Eval(c)
	obj := c.MustObject(c.RetVal)
	obj.Iterate(func(k string, v runtime.Value) {
		c.SetLocalValue(k, v)
	})
	c.RetVal = obj
}

type ExprLocalAssign struct {
	Pos
	Names      []string
	ExpandLast bool
	Type       int
	Expr       Expr
}

func (ExprLocalAssign) itIsAnAssignExpr() {}

func (assign *ExprLocalAssign) Eval(c *runtime.Context) {
	c.EnsureNotReadonly()
	if len(assign.Names) < 1 {
		return
	}
	switch assign.Type {
	case AssignTypeSingle:
		assign.Expr.Eval(c)
		v := c.RetVal
		c.SetLocalValue(assign.Names[0], v)
	case AssignTypeDeArray:
		if arr, isArr := assign.Expr.(*ExprArray); isArr {
			i, j, n := 0, 0, len(assign.Names)
			for i < n {
				if j >= len(arr.Items) {
					c.SetLocalValue(assign.Names[i], runtime.Undefined())
					i++
					continue
				}
				item := arr.Items[j]
				j++
				item.Expr.Eval(c)
				itemVal := c.RetVal
				if item.ShouldExpand {
					itemArr, isArr := itemVal.(runtime.ValueArray)
					if !isArr {
						c.RaiseRuntimeError("expand value must be an array")
						return
					}
					for k := 0; k < itemArr.Len(); k++ {
						c.SetLocalValue(assign.Names[i], itemArr.GetIndex(k, c))
						i++
						if i >= n {
							break
						}
					}
				} else {
					c.SetLocalValue(assign.Names[i], itemVal)
					i++
				}
			}
		} else {
			assign.Expr.Eval(c)
			vArr := c.MustArray(c.RetVal)
			n := len(assign.Names)
			if assign.ExpandLast {
				n--
				totalValues := vArr.Len()
				if totalValues > n {
					lastValue := runtime.NewArray(totalValues - n)
					for i := n; i < totalValues; i++ {
						lastValue.PushBack(vArr.GetIndex(i, c))
					}
					c.SetLocalValue(assign.Names[n], lastValue)
				} else {
					c.SetLocalValue(assign.Names[n], runtime.NewArray(0))
				}
			}
			for i := 0; i < n; i++ {
				c.SetLocalValue(assign.Names[i], vArr.GetIndex(i, c))
			}
		}
	case AssignTypeDeObject:
		assign.Expr.Eval(c)
		v := c.RetVal
		for _, name := range assign.Names {
			c.SetLocalValue(name, v.GetMember(name, c))
		}
	}
}

const (
	CompareOpEQ = 1 + iota
	CompareOpNE
	CompareOpLT
	CompareOpLE
	CompareOpGT
	CompareOpGE
)

type ExprCompare struct {
	First   Expr
	Ops     []int
	Targets []Expr
}

func (expr *ExprCompare) Eval(c *runtime.Context) {
	if len(expr.Ops) != len(expr.Targets) {
		c.RaiseRuntimeError("invalid compare!")
	}
	expr.First.Eval(c)
	v1 := ensureZgg(c.RetVal, c)
	for i, op := range expr.Ops {
		expr.Targets[i].Eval(c)
		v2 := ensureZgg(c.RetVal, c)
		var overridedMethod string
		switch op {
		case CompareOpEQ:
			overridedMethod = "__eq__"
		case CompareOpNE:
			overridedMethod = "__ne__"
		case CompareOpLT:
			overridedMethod = "__lt__"
		case CompareOpLE:
			overridedMethod = "__le__"
		case CompareOpGT:
			overridedMethod = "__gt__"
		case CompareOpGE:
			overridedMethod = "__ge__"
		default:
			c.RaiseRuntimeError("invalid compare op %d", op)
		}
		if opFn, ok := v1.GetMember(overridedMethod, c).(runtime.ValueCallable); ok {
			c.Invoke(opFn, v1, runtime.Args(v2))
			if !c.RetVal.IsTrue() {
				return
			}
			v1 = v2
			continue
		}
		comp := v1.CompareTo(v2, c)
		isTrue := false
		switch op {
		case CompareOpEQ:
			isTrue = comp == runtime.CompareResultEqual
		case CompareOpNE:
			isTrue = (comp & runtime.CompareResultEqual) == 0
		case CompareOpLT:
			isTrue = comp == runtime.CompareResultLess
		case CompareOpLE:
			isTrue = (comp & (runtime.CompareResultLess | runtime.CompareResultEqual)) != 0
		case CompareOpGT:
			isTrue = comp == runtime.CompareResultGreater
		case CompareOpGE:
			isTrue = (comp & (runtime.CompareResultGreater | runtime.CompareResultEqual)) != 0
		default:
			c.RaiseRuntimeError("invalid compare op %d", op)
		}
		if !isTrue {
			c.RetVal = runtime.NewBool(false)
			return
		}
		v1 = v2
	}
	c.RetVal = runtime.NewBool(true)
}

type ExprEqual struct {
	BinOp
}

func (expr *ExprEqual) Eval(c *runtime.Context) {
	c.RetVal = runtime.NewBool(c.ValuesEqual(expr.GetValues(c)))
}

type ExprNotEqual struct {
	BinOp
}

func (expr *ExprNotEqual) Eval(c *runtime.Context) {
	c.RetVal = runtime.NewBool(c.ValuesNotEqual(expr.GetValues(c)))
}

type ExprGreaterThen struct {
	BinOp
}

func (expr *ExprGreaterThen) Eval(c *runtime.Context) {
	c.RetVal = runtime.NewBool(c.ValuesGreater(expr.GetValues(c)))
}

type ExprGreaterEqual struct {
	BinOp
}

func (expr *ExprGreaterEqual) Eval(c *runtime.Context) {
	c.RetVal = runtime.NewBool(c.ValuesGreaterEqual(expr.GetValues(c)))
}

type ExprLessThen struct {
	BinOp
}

func (expr *ExprLessThen) Eval(c *runtime.Context) {
	c.RetVal = runtime.NewBool(c.ValuesLess(expr.GetValues(c)))
}

type ExprLessEqual struct {
	BinOp
}

func (expr *ExprLessEqual) Eval(c *runtime.Context) {
	c.RetVal = runtime.NewBool(c.ValuesLessEqual(expr.GetValues(c)))
}

type ExprLogicNot struct {
	Expr Expr
}

func (expr *ExprLogicNot) Eval(c *runtime.Context) {
	expr.Expr.Eval(c)
	left := c.RetVal
	if opFn, ok := left.GetMember("__true__", c).(runtime.ValueCallable); ok {
		c.Invoke(opFn, left, runtime.NoArgs)
		return
	}
	c.RetVal = runtime.NewBool(!c.RetVal.IsTrue())
}

type ExprLogicAnd struct {
	BinOp
}

func (expr *ExprLogicAnd) Eval(c *runtime.Context) {
	expr.BinOp.Left.Eval(c)
	left := c.RetVal
	if opFn, ok := left.GetMember("__and__", c).(runtime.ValueCallable); ok {
		expr.BinOp.Right.Eval(c)
		c.Invoke(opFn, left, runtime.Args(c.RetVal))
		return
	}
	if !c.ReturnTrue() {
		return
	}
	expr.BinOp.Right.Eval(c)
}

type ExprLogicOr struct {
	BinOp
}

func (expr *ExprLogicOr) Eval(c *runtime.Context) {
	expr.BinOp.Left.Eval(c)
	left := c.RetVal
	if opFn, ok := left.GetMember("__or__", c).(runtime.ValueCallable); ok {
		expr.BinOp.Right.Eval(c)
		c.Invoke(opFn, left, runtime.Args(c.RetVal))
		return
	}
	if c.ReturnTrue() {
		return
	}
	expr.BinOp.Right.Eval(c)
}

type ExprFallback struct {
	BinOp
}

func (expr *ExprFallback) Eval(c *runtime.Context) {
	defer func() {
		var exc runtime.Exception
		var ok bool
		if err := recover(); err != nil {
			if exc, ok = err.(runtime.Exception); ok {
				c.PushStack()
				defer c.PopStack()
				c.SetLocalValue("__err__", runtime.ExceptionToValue(exc, c))
				expr.BinOp.Right.Eval(c)
			} else {
				panic(err)
			}
		} else {
			switch c.RetVal.(type) {
			case runtime.ValueNil, runtime.ValueUndefined:
				c.PushStack()
				defer c.PopStack()
				c.SetLocalValue("__err__", runtime.Nil())
				expr.BinOp.Right.Eval(c)
			}
		}
	}()
	expr.BinOp.Left.Eval(c)
}

type ExprBitShl struct {
	BinOp
}

func (expr *ExprBitShl) Eval(c *runtime.Context) {
	left, right, overrideRet := expr.tryOverride(c, "__shl__")
	if overrideRet != nil {
		c.RetVal = overrideRet
	} else {
		c.RetVal = runtime.NewInt(c.MustInt(left) << c.MustInt(right))
	}
}

type ExprBitShr struct {
	BinOp
}

func (expr *ExprBitShr) Eval(c *runtime.Context) {
	left, right, overrideRet := expr.tryOverride(c, "__shr__")
	if overrideRet != nil {
		c.RetVal = overrideRet
	} else {
		c.RetVal = runtime.NewInt(c.MustInt(left) >> c.MustInt(right))
	}
}

type ExprBitAnd struct {
	BinOp
}

func (expr *ExprBitAnd) Eval(c *runtime.Context) {
	left, right, overrideRet := expr.tryOverride(c, "__bitAnd__")
	if overrideRet != nil {
		c.RetVal = overrideRet
	} else {
		c.RetVal = runtime.NewInt(c.MustInt(left) & c.MustInt(right))
	}
}

type ExprBitOr struct {
	BinOp
}

func (expr *ExprBitOr) Eval(c *runtime.Context) {
	left, right := expr.GetValues(c)
	left, right, overrideRet := expr.tryOverride(c, "__bitOr__")
	if overrideRet != nil {
		c.RetVal = overrideRet
	} else {
		c.RetVal = runtime.NewInt(c.MustInt(left) | c.MustInt(right))
	}
}

type ExprBitXor struct {
	BinOp
}

func (expr *ExprBitXor) Eval(c *runtime.Context) {
	left, right, overrideRet := expr.tryOverride(c, "__bitXor__")
	if overrideRet != nil {
		c.RetVal = overrideRet
	} else {
		c.RetVal = runtime.NewInt(c.MustInt(left) ^ c.MustInt(right))
	}
}

type ExprIsType struct {
	BinOp
}

func (expr *ExprIsType) Eval(c *runtime.Context) {
	left, right, overrideRet := expr.tryOverride(c, "__is__")
	if overrideRet != nil {
		c.RetVal = overrideRet
	} else {
		rightType, isType := runtime.Unbound(right).(runtime.ValueType)
		if !isType {
			if c.IsDebug {
				c.RaiseRuntimeError("is expression: right operand is not a Type, but a %s", reflect.TypeOf(right))
			} else {
				c.RaiseRuntimeError("is expression: right operand is not a Type")
			}
		}
		c.RetVal = runtime.NewBool(left.Type().IsSubOf(rightType))
	}
}
