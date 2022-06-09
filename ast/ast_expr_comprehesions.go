package ast

import "github.com/zgg-lang/zgg-go/runtime"

type basicComprehension struct {
	ValueName         string
	IndexerName       string
	Iterable          Expr
	RangeBegin        Expr
	RangeEnd          Expr
	RangeIncludingEnd bool
	FilterExpr        Expr
}

type SetResult = func(c *runtime.Context)

func (e *basicComprehension) tryPushItem(c *runtime.Context, index runtime.Value, value runtime.Value, setResult SetResult) {
	c.ForceSetLocalValue(e.ValueName, value)
	if e.IndexerName != "" {
		c.ForceSetLocalValue(e.IndexerName, index)
	}
	if e.FilterExpr != nil {
		e.FilterExpr.Eval(c)
		if !c.RetVal.IsTrue() {
			return
		}
	}
	setResult(c)
}

func (e *basicComprehension) evalWithIterable(c *runtime.Context, setResult SetResult) {
	e.Iterable.Eval(c)
	iteratable := c.RetVal
	getIter := iteratable.GetMember("__iter__", c)
	if c.IsCallable(getIter) {
		e.evalWithCustomIterable(c, getIter.(runtime.ValueCallable), setResult)
	} else {
		e.evalWithBuiltinIterable(c, iteratable, setResult)
	}
}

func (e *basicComprehension) evalWithCustomIterable(c *runtime.Context, getIter runtime.ValueCallable, setResult SetResult) {
	c.Invoke(getIter.(runtime.ValueCallable), nil, runtime.NoArgs)
	iter := c.RetVal
	if !c.IsCallable(iter) {
		c.RaiseRuntimeError("__iter__ should return a callable value")
	}
	iterFn := iter.(runtime.ValueCallable)
	for i := 0; ; i++ {
		c.Invoke(iterFn, nil, runtime.NoArgs)
		var value runtime.Value
		if retArr, isArray := c.RetVal.(runtime.ValueArray); !isArray || retArr.Len() != 2 || !retArr.GetIndex(1, c).IsTrue() {
			break
		} else {
			value = retArr.GetIndex(0, c)
		}
		e.tryPushItem(c,
			runtime.NewInt(int64(i)),
			value,
			setResult)
	}
}

func (e *basicComprehension) evalWithBuiltinIterable(c *runtime.Context, iteratable runtime.Value, setResult SetResult) {
	switch it := iteratable.(type) {
	case runtime.ValueInt:
		for i := 0; i < it.AsInt(); i++ {
			v := runtime.NewInt(int64(i))
			e.tryPushItem(c, v, v, setResult)
		}
	case runtime.CanLen:
		if o, ok := iteratable.(runtime.ValueObject); ok {
			o.Each(func(key string, value runtime.Value) bool {
				e.tryPushItem(c,
					runtime.NewStr(key),
					value,
					setResult)
				return true
			})
		} else {
			l := it.Len()
			for i := 0; i < l; i++ {
				e.tryPushItem(c,
					runtime.NewInt(int64(i)),
					iteratable.GetIndex(i, c),
					setResult)
			}
		}
	}
}

func (e *basicComprehension) evalWithRange(c *runtime.Context, setResult SetResult) {
	e.RangeBegin.Eval(c)
	current := c.RetVal
	curInt, isInt := current.(runtime.ValueInt)
	if !isInt {
		c.RaiseRuntimeError("array comprehension: range begin must be an integer")
	}
	e.RangeEnd.Eval(c)
	end := c.RetVal
	endInt, isInt := end.(runtime.ValueInt)
	if !isInt {
		c.RaiseRuntimeError("array comprehension: range end must be an integer")
	}
	i := -1
	ci, ei := curInt.Value(), endInt.Value()
	for ; ci < ei; ci++ {
		i++
		e.tryPushItem(c,
			runtime.NewInt(int64(i)),
			runtime.NewInt(ci),
			setResult,
		)
	}
	if e.RangeIncludingEnd {
		e.tryPushItem(c,
			runtime.NewInt(int64(i+1)),
			end,
			setResult,
		)
	}
}

func (e *basicComprehension) eval(c *runtime.Context, setResult SetResult) {
	c.PushStack()
	defer c.PopStack()
	if e.Iterable != nil {
		e.evalWithIterable(c, setResult)
	} else {
		e.evalWithRange(c, setResult)
	}
}

type ArrayComprehension struct {
	basicComprehension
	ItemExpr Expr
}

func (e *ArrayComprehension) Eval(c *runtime.Context) {
	rv := runtime.NewArray()
	e.eval(c, func(c *runtime.Context) {
		e.ItemExpr.Eval(c)
		rv.PushBack(c.RetVal)
	})
	c.RetVal = rv
}

type ObjectComprehension struct {
	basicComprehension
	KeyExpr   Expr
	ValueExpr Expr
}

func (e *ObjectComprehension) Eval(c *runtime.Context) {
	rv := runtime.NewObject()
	e.eval(c, func(c *runtime.Context) {
		e.KeyExpr.Eval(c)
		key := c.RetVal.ToString(c)
		e.ValueExpr.Eval(c)
		rv.SetMember(key, c.RetVal, c)
	})
	c.RetVal = rv
}
