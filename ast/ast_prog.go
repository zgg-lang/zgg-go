package ast

import (
	"fmt"

	"github.com/zgg-lang/zgg-go/runtime"
)

type Module struct {
	Block *Block
}

func (m *Module) Eval(c *runtime.Context) {
	defer c.RunDefers()
	m.Block.Eval(c)
	c.RetVal = c.ExportValue
}

type Stmt interface {
	Node
	Position
}

type BlockType int

const (
	BlockTypeNormal BlockType = iota
	BlockTypeModuleTop
	BlockTypeFuncTop
	BlockTypeLoopTop
)

type Block struct {
	Pos
	Type  BlockType
	Stmts []Stmt
	// LocalNames []string
}

func (m *Block) Eval(c *runtime.Context) {
	c.PushStack()
	defer c.PopStack()
	for _, e := range m.Stmts {
		c.SetPosition(e.Position())
		e.Eval(c)
		if c.Breaking || c.Continuing || c.Returned {
			break
		}
	}
}

func execLoopBody(label string, block *Block, c *runtime.Context) bool {
	c.PushStack()
	defer c.PopStack()
	block.Eval(c)
	if c.Breaking {
		if label == c.BreakingLabel {
			c.Breaking = false
			c.BreakingLabel = ""
		}
		return false
	}
	if c.Continuing {
		if label == c.ContinuingLabel {
			c.Continuing = false
			c.ContinuingLabel = ""
			return true
		} else {
			return false
		}
	}
	return !c.Breaking && !c.Returned
}

type StmtFor struct {
	Pos
	Label string
	Init  Expr
	Check Expr
	Next  Expr
	Exec  *Block
}

func (s *StmtFor) Eval(c *runtime.Context) {
	c.PushStack()
	defer c.PopStack()
	s.Init.Eval(c)
	for {
		s.Check.Eval(c)
		if !c.ReturnTrue() {
			break
		}
		if !execLoopBody(s.Label, s.Exec, c) {
			break
		}
		s.Next.Eval(c)
	}
}

type StmtForEach struct {
	Pos
	Label             string
	IdIndex           string
	IdValue           string
	Iteratable        Expr
	RangeBegin        Expr
	RangeEnd          Expr
	RangeIncludingEnd bool
	Exec              *Block
}

func (s *StmtForEach) evalWithIterable(c *runtime.Context) {
	s.Iteratable.Eval(c)
	iteratable := c.RetVal
	getIter := iteratable.GetMember("__iter__", c)
	if c.IsCallable(getIter) {
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
			c.ForceSetLocalValue(s.IdValue, value)
			if id := s.IdIndex; id != "" {
				c.ForceSetLocalValue(id, runtime.NewInt(int64(i)))
			}
			if !execLoopBody(s.Label, s.Exec, c) {
				break
			}
		}
		return
	}
	switch v := iteratable.(type) {
	case runtime.CanLen:
		switch vv := iteratable.(type) {
		case runtime.ValueObject:
			vv.Each(func(key string, value runtime.Value) bool {
				c.ForceSetLocalValue(s.IdValue, value)
				if id := s.IdIndex; id != "" {
					c.ForceSetLocalValue(id, runtime.NewStr(key))
				}
				if !execLoopBody(s.Label, s.Exec, c) {
					return false
				}
				return true
			})
		case runtime.ValueMap:
			vv.Each(func(key, value runtime.Value) bool {
				c.ForceSetLocalValue(s.IdValue, value)
				if id := s.IdIndex; id != "" {
					c.ForceSetLocalValue(id, key)
				}
				if !execLoopBody(s.Label, s.Exec, c) {
					return false
				}
				return true
			})
		default:
			l := v.Len()
			for i := 0; i < l; i++ {
				value := iteratable.GetIndex(i, c)
				c.ForceSetLocalValue(s.IdValue, value)
				if s.IdIndex != "" {
					c.ForceSetLocalValue(s.IdIndex, runtime.NewInt(int64(i)))
				}
				if !execLoopBody(s.Label, s.Exec, c) {
					break
				}
			}
		}
	case runtime.ValueInt:
		l := v.AsInt()
		for i := 0; i < l; i++ {
			value := runtime.NewInt(int64(i))
			c.ForceSetLocalValue(s.IdValue, value)
			if s.IdIndex != "" {
				c.ForceSetLocalValue(s.IdIndex, value)
			}
			if !execLoopBody(s.Label, s.Exec, c) {
				break
			}
		}
	default:
		c.RaiseRuntimeError("value is not iterable")
	}
}

func (s *StmtForEach) Eval(c *runtime.Context) {
	c.PushStack()
	defer c.PopStack()
	if s.Iteratable != nil {
		s.evalWithIterable(c)
	} else {
		s.RangeBegin.Eval(c)
		curVal := c.RetVal
		curInt, ok := curVal.(runtime.ValueInt)
		if !ok {
			c.RaiseRuntimeError("for in range must begin with an integer")
		}
		s.RangeEnd.Eval(c)
		endVal := c.RetVal
		endInt, ok := endVal.(runtime.ValueInt)
		if !ok {
			c.RaiseRuntimeError("for in range must end with an integer")
		}
		cur := curInt.AsInt()
		end := endInt.AsInt()
		if s.RangeIncludingEnd {
			end++
		}
		for i := -1; cur < end; cur++ {
			i++
			c.ForceSetLocalValue(s.IdValue, runtime.NewInt(int64(cur)))
			if s.IdIndex != "" {
				c.ForceSetLocalValue(s.IdIndex, runtime.NewInt(int64(i)))
			}
			if !execLoopBody(s.Label, s.Exec, c) {
				break
			}
		}
	}
}

type StmtDoWhile struct {
	Pos
	Label string
	Check Expr
	Exec  *Block
}

func (s *StmtDoWhile) Eval(c *runtime.Context) {
	for {
		if !execLoopBody(s.Label, s.Exec, c) {
			break
		}
		s.Check.Eval(c)
		if !c.ReturnTrue() {
			break
		}
	}
}

type StmtWhile struct {
	Pos
	Label string
	Check Expr
	Exec  *Block
}

func (s *StmtWhile) Eval(c *runtime.Context) {
	for {
		s.Check.Eval(c)
		if !c.ReturnTrue() {
			break
		}
		if !execLoopBody(s.Label, s.Exec, c) {
			break
		}
	}
}

type StmtBreak struct {
	Pos
	Level   int
	ToLabel string
}

func (s *StmtBreak) Eval(c *runtime.Context) {
	c.Breaking = true
	c.BreakingLabel = s.ToLabel
}

type StmtContinue struct {
	Pos
	Level   int
	ToLabel string
}

func (s *StmtContinue) Eval(c *runtime.Context) {
	c.Continuing = true
	c.ContinuingLabel = s.ToLabel
}

type IfCase struct {
	Check      Expr
	Assignment Expr
	Do         *Block
}

type StmtIf struct {
	Pos
	Cases  []*IfCase
	ElseDo *Block
}

func (s *StmtIf) Eval(c *runtime.Context) {
	stacks := 0
	defer func() {
		for ; stacks > 0; stacks-- {
			c.PopStack()
		}
	}()
	for _, ifCase := range s.Cases {
		if ifCase.Assignment != nil {
			c.PushStack()
			ifCase.Assignment.Eval(c)
			stacks++
		}
		ifCase.Check.Eval(c)
		if !c.ReturnTrue() {
			continue
		}
		ifCase.Do.Eval(c)
		return
	}
	if s.ElseDo != nil {
		s.ElseDo.Eval(c)
	}
}

type SwitchCase struct {
	Condition   ValueCondition
	Code        *Block
	Fallthrough bool
}

type StmtSwitch struct {
	Pos
	Val     Expr
	Cases   []SwitchCase
	Default *Block
}

func (s *StmtSwitch) Eval(c *runtime.Context) {
	s.Val.Eval(c)
	val := c.RetVal
	for _, case_ := range s.Cases {
		if case_.Condition.IsMatch(c, val) {
			case_.Code.Eval(c)
			if !case_.Fallthrough {
				return
			}
		}
	}
	if s.Default != nil {
		s.Default.Eval(c)
	}
}

type StmtReturn struct {
	Pos
	Value Expr
}

func (s *StmtReturn) Eval(c *runtime.Context) {
	if s.Value != nil {
		s.Value.Eval(c)
	}
	c.Returned = true
}

type StmtExport struct {
	Pos
	Name string
	Expr Expr
}

func (s *StmtExport) Eval(c *runtime.Context) {
	if !c.IsModuleTop() {
		c.RaiseRuntimeError("export must be in module top block")
		return
	}
	s.Expr.Eval(c)
	c.ExportValue.SetMember(s.Name, c.RetVal, c)
}

type StmtClassDefine struct {
	Pos
	Exported bool
	Name     string
	Bases    []Expr
	Body     *ExprObject
	Static   *ExprObject
}

func (s *StmtClassDefine) Eval(c *runtime.Context) {
	if s.Exported && !c.IsModuleTop() {
		c.RaiseRuntimeError("export must be in module top block")
		return
	}
	newClass := runtime.NewType(runtime.NextTypeId(), s.Name)
	if len(s.Bases) > 0 {
		newClass.Bases = make([]runtime.ValueType, len(s.Bases))
		for i, b := range s.Bases {
			b.Eval(c)
			baseVal := c.RetVal
			if b, isBound := baseVal.(runtime.ValueBoundMethod); isBound {
				baseVal = b.Value
			}
			if baseCls, isType := baseVal.(runtime.ValueType); !isType {
				c.RaiseRuntimeError("base class %s is not a type", c.RetVal.ToString(c))
				return
			} else {
				newClass.Bases[i] = baseCls
			}
		}
	} else {
		newClass.Bases = []runtime.ValueType{runtime.TypeObject}
	}
	for _, item := range s.Body.Items {
		it, ok := item.(ExprObjectItemKV)
		if !ok {
			c.RaiseRuntimeError("expand object is not acceptable in class defination")
		}
		it.Key.Eval(c)
		key := c.RetVal.ToString(c)
		it.Value.Eval(c)
		val := c.RetVal
		if valFunc, isFunc := val.(*runtime.ValueFunc); isFunc {
			valFunc.BelongType = newClass
		}
		newClass.Members.Store(key, c.RetVal)
	}
	for _, item := range s.Static.Items {
		it, ok := item.(ExprObjectItemKV)
		if !ok {
			c.RaiseRuntimeError("expand object is not acceptable in class defination")
		}
		it.Key.Eval(c)
		key := c.RetVal.ToString(c)
		it.Value.Eval(c)
		val := c.RetVal
		if valFunc, isFunc := val.(*runtime.ValueFunc); isFunc {
			valFunc.BelongType = newClass
		}
		newClass.Statics.Store(key, c.RetVal)
	}
	c.SetLocalValue(s.Name, newClass)
	if s.Exported {
		c.ExportValue.SetMember(s.Name, newClass, c)
	}
	c.RetVal = newClass
}

type StmtDefer struct {
	Pos
	Call *ExprCall
}

func (s *StmtDefer) Eval(c *runtime.Context) {
	s.Call.Callee.Eval(c)
	callee := c.MustCallable(c.RetVal)
	args := s.Call.GetArgs(c, callee)
	c.AddDefer(callee, args, s.Call.Optional)
}

type StmtBlockDefer struct {
	Pos
	Call *ExprCall
}

func (s *StmtBlockDefer) Eval(c *runtime.Context) {
	s.Call.Callee.Eval(c)
	callee := c.MustCallable(c.RetVal)
	args := s.Call.GetArgs(c, callee)
	c.AddBlockDefer(callee, args, s.Call.Optional)
}

type StmtTry struct {
	Pos
	Try     *Block
	Catch   *Block
	ExcName string
	Finally *Block
}

func (s *StmtTry) Eval(c *runtime.Context) {
	var ret runtime.Value = runtime.Undefined()
	defer func() {
		if s.Catch != nil {
			e := recover()
			if e != nil {
				switch exc := e.(type) {
				case runtime.Exception:
					c.PushStack()
					defer c.PopStack()
					c.SetLocalValue(s.ExcName, runtime.ExceptionToValue(exc, c))
					s.Catch.Eval(c)
					ret = c.RetVal
				default:
					panic(e)
				}
			}
		}
		if s.Finally != nil {
			returned := c.Returned
			s.Finally.Eval(c)
			if returned {
				c.RetVal = ret
			}
		}
	}()
	s.Try.Eval(c)
	ret = c.RetVal
}

type StmtFallback struct {
	Pos
	Stmt     Stmt
	Fallback *Block
}

func (s *StmtFallback) Eval(c *runtime.Context) {
	if s.Fallback != nil {
		defer func() {
			e := recover()
			if e != nil {
				switch exc := e.(type) {
				case runtime.Exception:
					c.PushStack()
					defer c.PopStack()
					c.SetLocalValue("__err__", runtime.ExceptionToValue(exc, c))
					s.Fallback.Eval(c)
				default:
					panic(e)
				}
			}
		}()
	}
	s.Stmt.Eval(c)
}

type StmtAssert struct {
	Pos
	Expr    Expr
	Message Expr
}

func (s *StmtAssert) Eval(c *runtime.Context) {
	s.Expr.Eval(c)
	if !c.RetVal.IsTrue() {
		s.Message.Eval(c)
		c.RaiseRuntimeError("Assertion fail! " + c.RetVal.ToString(c))
	}
}

type StmtExtend struct {
	Pos
	Exported bool
	Type     Expr
	Name     []Expr
	Func     []Expr
}

func (s *StmtExtend) Eval(c *runtime.Context) {
	if s.Exported && !c.IsModuleTop() {
		c.RaiseRuntimeError("export must be in module top block")
		return
	}
	s.Type.Eval(c)
	t, ok := c.RetVal.(runtime.ValueType)
	if !ok {
		c.RaiseRuntimeError("extending a non-type value")
	}
	if len(s.Name) != len(s.Func) {
		c.RaiseRuntimeError("bug! extend len(name) %d != len(func) %d", len(s.Name), len(s.Func))
	}
	for i, name := range s.Name {
		name.Eval(c)
		nameVal := c.RetVal
		s.Func[i].Eval(c)
		extVal := c.RetVal
		varName := fmt.Sprintf("%d#%s", t.TypeId, nameVal.ToString(c))
		c.SetLocalValue(varName, extVal)
		if s.Exported {
			c.ExportValue.SetMember(varName, extVal, c)
		}
	}
}
