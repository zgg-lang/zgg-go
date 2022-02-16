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
	Identifier1       string
	Identifier2       string
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
			c.OnRuntimeError("__iter__ should return a callable value")
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
			if s.Identifier2 == "" {
				c.ForceSetLocalValue(s.Identifier1, value)
			} else {
				c.ForceSetLocalValue(s.Identifier1, runtime.NewInt(int64(i)))
				c.ForceSetLocalValue(s.Identifier2, value)
			}
			if !execLoopBody(s.Label, s.Exec, c) {
				break
			}
		}
		return
	}
	switch v := iteratable.(type) {
	case runtime.CanLen:
		if obj, ok := iteratable.(runtime.ValueObject); ok {
			obj.Each(func(key string, value runtime.Value) bool {
				if s.Identifier2 == "" {
					c.ForceSetLocalValue(s.Identifier1, value)
				} else {
					c.ForceSetLocalValue(s.Identifier1, runtime.NewStr(key))
					c.ForceSetLocalValue(s.Identifier2, value)
				}
				if !execLoopBody(s.Label, s.Exec, c) {
					return false
				}
				return true
			})
		} else {
			l := v.Len()
			for i := 0; i < l; i++ {
				value := iteratable.GetIndex(i, c)
				if s.Identifier2 == "" {
					c.ForceSetLocalValue(s.Identifier1, value)
				} else {
					c.ForceSetLocalValue(s.Identifier1, runtime.NewInt(int64(i)))
					c.ForceSetLocalValue(s.Identifier2, value)
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
			if s.Identifier2 == "" {
				c.ForceSetLocalValue(s.Identifier1, value)
			} else {
				c.ForceSetLocalValue(s.Identifier1, value)
				c.ForceSetLocalValue(s.Identifier2, value)
			}
			if !execLoopBody(s.Label, s.Exec, c) {
				break
			}
		}
	default:
		c.OnRuntimeError("value is not iterable")
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
			c.OnRuntimeError("for in range must begin with an integer")
		}
		s.RangeEnd.Eval(c)
		endVal := c.RetVal
		endInt, ok := endVal.(runtime.ValueInt)
		if !ok {
			c.OnRuntimeError("for in range must end with an integer")
		}
		cur := curInt.AsInt()
		end := endInt.AsInt()
		if s.RangeIncludingEnd {
			end++
		}
		for i := -1; cur < end; cur++ {
			i++
			if s.Identifier2 == "" {
				c.ForceSetLocalValue(s.Identifier1, runtime.NewInt(int64(cur)))
			} else {
				c.ForceSetLocalValue(s.Identifier1, runtime.NewInt(int64(i)))
				c.ForceSetLocalValue(s.Identifier2, runtime.NewInt(int64(cur)))
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
		c.OnRuntimeError("export must be in module top block")
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
		c.OnRuntimeError("export must be in module top block")
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
				c.OnRuntimeError("base class %s is not a type", c.RetVal.ToString(c))
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
			c.OnRuntimeError("expand object is not acceptable in class defination")
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
			c.OnRuntimeError("expand object is not acceptable in class defination")
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
			s.Finally.Eval(c)
		}
		c.RetVal = ret
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
		c.OnRuntimeError("Assertion fail! " + c.RetVal.ToString(c))
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
		c.OnRuntimeError("export must be in module top block")
		return
	}
	s.Type.Eval(c)
	t, ok := c.RetVal.(runtime.ValueType)
	if !ok {
		c.OnRuntimeError("extending a non-type value")
	}
	if len(s.Name) != len(s.Func) {
		c.OnRuntimeError("bug! extend len(name) %d != len(func) %d", len(s.Name), len(s.Func))
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
