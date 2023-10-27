package ast

import (
	"reflect"

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
	c.ValuesPlus(left, right)
}

type ExprMinus struct {
	BinOp
}

func (n *ExprMinus) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	c.ValuesMinus(left, right)
}

type ExprTimes struct {
	BinOp
}

func (n *ExprTimes) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	c.ValuesTimes(left, right)
}

type ExprDiv struct {
	BinOp
}

func (n *ExprDiv) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	c.ValuesDiv(left, right)
}

type ExprMod struct {
	BinOp
}

func (n *ExprMod) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	c.ValuesMod(left, right)
}

type ExprPow struct {
	BinOp
}

func (n *ExprPow) Eval(c *runtime.Context) {
	left, right := n.GetValues(c)
	c.ValuesPow(left, right)
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
