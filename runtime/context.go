package runtime

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type deferCall struct {
	callee   Value
	args     []Value
	optional bool
}

type contextFrame struct {
	level     int
	loopLevel int
	variables *sync.Map
	parent    *contextFrame
	funcName  string
	funcLevel int
	filename  string
	lineNum   int
	defers    []deferCall
}

func newContextFrame(parent *contextFrame) *contextFrame {
	f := &contextFrame{
		parent:    parent,
		level:     0,
		funcName:  "<root>",
		funcLevel: 0,
		variables: new(sync.Map),
	}
	if parent != nil {
		f.level = parent.level + 1
		f.funcName = parent.funcName
		f.funcLevel = parent.funcLevel
		f.lineNum = parent.lineNum
		f.filename = parent.filename
	}
	return f
}

func (s *contextFrame) resetAsRoot() {
	s.parent = nil
	s.level = 0
	s.loopLevel = 0
	s.funcName = "<root>"
	s.funcLevel = 0
	s.variables = new(sync.Map)
	s.filename = ""
	s.lineNum = 0
	s.defers = nil
}

func (s *contextFrame) findValue(name string) Value {
	for stack := s; stack != nil; stack = stack.parent {
		if value, valueFound := stack.variables.Load(name); valueFound {
			return value.(Value)
		}
	}
	return nil
}

func (frame *contextFrame) addDefer(callee Value, args []Value, optional bool) {
	// fmt.Println("addDefer at frame", frame.level)
	frame.defers = append(frame.defers, deferCall{callee: callee, args: args, optional: optional})
}

func (frame *contextFrame) clone() *contextFrame {
	return &contextFrame{
		level:     frame.level,
		loopLevel: frame.loopLevel,
		variables: frame.variables,
		parent:    frame.parent,
		funcName:  frame.funcName,
		funcLevel: frame.funcLevel,
		filename:  frame.filename,
		lineNum:   frame.lineNum,
		defers:    frame.defers,
	}
}

type ModuleInfo struct {
	Value      Value
	ModifyTime int64
}

type Context struct {
	lock            sync.Mutex
	main            bool
	IsDebug         bool
	debugLogger     *log.Logger
	Path            string
	ImportPaths     []string
	Args            []string
	RetVal          Value
	Breaking        bool
	BreakingLabel   string
	Stdout          io.Writer
	Stderr          io.Writer
	Stdin           io.Reader
	Continuing      bool
	ContinuingLabel string
	Returned        bool
	ExportValue     ValueObject
	ImportFunc      func(*Context, string, string, string, int64) (Value, int64, bool)
	modules         *sync.Map
	curFrame        *contextFrame
	funcRootFrame   *contextFrame
	rootFrame       *contextFrame
	local           ValueObject
	builtins        *sync.Map
	CanEval         bool
}

func GetImportPaths() []string {
	importPaths := []string{".", "./zgg_modules"}
	if wd, err := os.Getwd(); err == nil {
		importPaths[0] = wd
		importPaths[1] = filepath.Join(wd, "zgg_modules")
	}
	if zggPath := os.Getenv("ZGGPATH"); zggPath != "" {
		importPaths = append(importPaths, strings.Split(zggPath, ":")...)
	}
	return importPaths
}

func NewContext(isMain bool, isDebug, canEval bool) *Context {
	f := newContextFrame(nil)
	ctx := &Context{
		RetVal:          Undefined(),
		main:            isMain,
		IsDebug:         isDebug,
		debugLogger:     log.New(os.Stderr, "DBG", log.Ldate|log.Lshortfile),
		Path:            ".",
		ImportPaths:     GetImportPaths(),
		Breaking:        false,
		BreakingLabel:   "",
		Continuing:      false,
		ContinuingLabel: "",
		Returned:        false,
		ExportValue:     NewObject(),
		curFrame:        f,
		funcRootFrame:   f,
		rootFrame:       f,
		builtins:        new(sync.Map),
		local:           NewObject(),
		CanEval:         canEval,
	}
	ctx.modules = new(sync.Map)
	ctx.Stdin = os.Stdin
	ctx.Stdout = os.Stdout
	ctx.Stderr = os.Stderr
	for name, value := range builtins {
		ctx.builtins.Store(name, value)
	}
	ctx.builtins.Store("isMain", NewBool(isMain))
	return ctx
}

func (c *Context) Reset(resetVars ...bool) {
	f := c.curFrame
	f.resetAsRoot()
	c.RetVal = Undefined()
	c.Path = "."
	c.Breaking = false
	c.BreakingLabel = ""
	c.Continuing = false
	c.ContinuingLabel = ""
	c.Returned = false
	c.ExportValue = NewObject()
	c.curFrame = f
	c.funcRootFrame = f
	c.rootFrame = f
	c.local = NewObject()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
}

func (c *Context) DebugLog(msg string, args ...interface{}) {
	if c.IsDebug {
		c.debugLogger.Output(2, fmt.Sprintf(msg+"\n", args...))
	}
}

func (c *Context) IsInFunc() bool {
	return c.funcRootFrame != c.rootFrame
}

func (c *Context) IsModuleTop() bool {
	return c.curFrame == c.rootFrame || (c.curFrame.parent != nil &&
		c.curFrame.funcLevel == c.curFrame.parent.funcLevel &&
		c.curFrame.parent == c.rootFrame)
}

func (c *Context) FindValue(name string) (Value, bool) {
	if name == "_" {
		return nil, false
	}
	v := c.curFrame.findValue(name)
	if v != nil {
		return v, true
	}
	builtin, builtinFound := c.builtins.Load(name)
	if builtinFound {
		return builtin.(Value), true
	}
	switch name {
	case "local":
		return c.local, true
	}
	return nil, false
}

func (c *Context) ModifyValue(name string, value Value) {
	if name == "_" {
		return
	}
	for s := c.curFrame; s != nil; s = s.parent {
		if _, found := (*s.variables).Load(name); found {
			(*s.variables).Store(name, value)
			return
		}
	}
	c.OnRuntimeError("variable %s not exists", name)
}

func (c *Context) SetLocalValue(name string, value Value) {
	if name == "_" {
		return
	}
	if _, found := (*c.curFrame.variables).Load(name); found {
		c.OnRuntimeError(fmt.Sprintf("variable %s redefined", name))
	} else {
		(*c.curFrame.variables).Store(name, value)
	}
}

func (c *Context) ForceSetLocalValue(name string, value Value) {
	if name != "_" {
		(*c.curFrame.variables).Store(name, value)
	}
}

func (c *Context) PushStack() {
	nextFrame := newContextFrame(c.curFrame)
	// fmt.Println("PushStack", c.curFrame.level, nextFrame.level)
	c.curFrame = nextFrame
}

func (c *Context) PushFuncStack(funcName string) {
	nextFrame := newContextFrame(c.curFrame)
	if funcName != "" {
		nextFrame.funcName = funcName
		nextFrame.funcLevel = c.curFrame.funcLevel + 1
		c.funcRootFrame = nextFrame
	}
	// fmt.Println("PushFuncStack", c.curFrame.level, nextFrame.level)
	c.curFrame = nextFrame
}

func (c *Context) AddDefer(callee Value, args []Value, optional bool) {
	c.funcRootFrame.addDefer(callee, args, optional)
}

func (c *Context) AddBlockDefer(callee Value, args []Value, optional bool) {
	c.curFrame.addDefer(callee, args, optional)
}

func (c *Context) RunDefers() {
	frame := c.curFrame
	retVal := c.RetVal
	if n := len(frame.defers); n > 0 {
		for i := n - 1; i >= 0; i-- {
			call := frame.defers[i]
			if !c.Invoke(call.callee, nil, func() []Value { return call.args }) && !call.optional {
				c.OnRuntimeError("defer call not callable")
				return
			}
		}
	}
	c.RetVal = retVal
}

func (c *Context) PopStack() {
	frame := c.curFrame
	// fmt.Println("on pop stack defers", frame.funcName, "line", frame.lineNum, "level", frame.level, "defers", len(frame.defers))
	c.RunDefers()
	c.curFrame = frame.parent
	if frame.funcLevel != c.curFrame.funcLevel {
		c.Returned = false
		funcRoot := c.curFrame.parent
		if funcRoot != nil {
			for funcRoot.parent != nil && funcRoot.parent.funcLevel == funcRoot.funcLevel {
				funcRoot = funcRoot.parent
			}
		}
		c.funcRootFrame = funcRoot
	}
	if c.curFrame == nil {
		panic("Context pop stack in bottom frame")
	}
}

func (c *Context) OnRuntimeError(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	e := &RuntimeError{
		Message: msg,
	}
	frame := c.curFrame
	for frame != nil {
		level := frame.funcLevel
		stack := Stack{
			FileName: frame.filename,
			Line:     frame.lineNum,
			Function: frame.funcName,
		}
		e.Stack = append(e.Stack, stack)
		frame = frame.parent
		for frame != nil && frame.funcLevel == level {
			frame = frame.parent
		}
	}
	panic(e)
}

func (c *Context) PrintStack() {
	frame := c.curFrame
	for frame != nil {
		level := frame.funcLevel
		fmt.Printf("%s line %d (%s)\n",
			frame.filename,
			frame.lineNum,
			frame.funcName,
		)
		frame = frame.parent
		for frame != nil && frame.funcLevel == level {
			frame = frame.parent
		}
	}
}

func (c *Context) SetPosition(filename string, lineNum int) {
	c.curFrame.filename = filename
	c.curFrame.lineNum = lineNum
}

func (c *Context) AddModule(name string, val Value, modTime int64) {
	if _, isUndefined := val.(ValueUndefined); !isUndefined {
		c.modules.Store(name, ModuleInfo{Value: val, ModifyTime: modTime})
	} else {
		c.modules.Delete(name)
	}
}

const (
	ImportTypeScript = "script"
	ImportTypeText   = "text"
	ImportTypeBytes  = "bytes"
	ImportTypeCsv    = "csv"
	ImportTypeJson   = "json"
)

func (c *Context) ImportModule(modPath string, forceReload bool, importType string) Value {
	modInfo, found := c.modules.Load(modPath)
	modTime := int64(0)
	if found {
		info := modInfo.(ModuleInfo)
		if !forceReload || info.ModifyTime == 0 {
			return info.Value
		}
		modTime = info.ModifyTime
	}
	modVal, thisTime, success := c.ImportFunc(c, modPath, "", importType, modTime)
	if !success {
		c.OnRuntimeError("ImportError: module %s not exists", modPath)
		return constUndefined
	}
	if thisTime == modTime {
		if found {
			return modInfo.(ModuleInfo).Value
		}
	}
	if _, isUndefined := modVal.(ValueUndefined); !isUndefined {
		c.AddModule(modPath, modVal, thisTime)
	}
	return modVal
}

// Helpful functions
func (c *Context) ReturnTrue() bool {
	return c.RetVal != nil && c.RetVal.IsTrue()
}

func (c *Context) throwValueTypeError(v Value, wanted string, name []string) {
	if len(name) > 0 {
		c.OnRuntimeError("%s requires %s, got %s", wanted, v.Type().Name)
	} else {
		c.OnRuntimeError("value requires %s, got %s", wanted, v.Type().Name)
	}
}

func (c *Context) AssertArgNum(num, min, max int, funcName string) {
	if num < min || num > max {
		if min == max {
			c.OnRuntimeError("%s requires %d argument(s)", funcName, min)
		} else {
			c.OnRuntimeError("%s requires %d ~ %d argument(s)", funcName, min, max)
		}
	}
}

func (c *Context) MustInt(v Value, name ...string) int64 {
	if iv, ok := v.(ValueInt); ok {
		return iv.Value()
	}
	c.throwValueTypeError(v, "int", name)
	return -1
}

func (c *Context) MustFloat(v Value, name ...string) float64 {
	if iv, ok := v.(ValueFloat); ok {
		return iv.Value()
	}
	if iv, ok := v.(ValueInt); ok {
		return float64(iv.Value())
	}
	c.throwValueTypeError(v, "int or float", name)
	return -1
}

func (c *Context) MustBool(v Value, name ...string) bool {
	if iv, ok := v.(ValueBool); ok {
		return iv.Value()
	}
	c.throwValueTypeError(v, "bool", name)
	return false
}

func (c *Context) MustStr(v Value, name ...string) string {
	if iv, ok := v.(ValueStr); ok {
		return iv.Value()
	}
	c.throwValueTypeError(v, "str", name)
	return ""
}

func (c *Context) MustArray(v Value, name ...string) ValueArray {
	if av, ok := v.(ValueArray); ok {
		return av
	}
	c.throwValueTypeError(v, "array", name)
	return ValueArray{}
}

func (c *Context) MustObject(v Value, name ...string) ValueObject {
	if ov, ok := v.(ValueObject); ok {
		return ov
	}
	c.throwValueTypeError(v, "object", name)
	return ValueObject{}
}

func (c *Context) MustCallable(v Value, name ...string) ValueCallable {
	if ov, ok := v.(ValueCallable); ok {
		return ov
	}
	c.throwValueTypeError(v, "callable", name)
	return nil
}

func (c *Context) IsCallable(value Value) bool {
	switch vv := value.(type) {
	case ValueObject:
		return c.IsCallable(vv.GetMember("__call__", c))
	case ValueCallable:
		return true
	}
	return false
}

func (c *Context) Invoke(calleeVal Value, this Value, getArgs func() []Value) bool {
	switch callee := calleeVal.(type) {
	case ValueCallable:
		callee.Invoke(c, this, getArgs())
		return true
	case ValueType:
		{
			newObj := NewObject(callee)
			initMember := newObj.GetMember("__init__", c)
			if initFunc, isCallable := initMember.(ValueCallable); isCallable {
				c.Invoke(initFunc, newObj, getArgs)
			}
			c.RetVal = newObj
			return true
		}
	}
	return false
}

func (c *Context) InvokeMethod(this Value, method string, getArgs func() []Value) Value {
	m := c.MustCallable(this.GetMember(method, c))
	c.Invoke(m, this, getArgs)
	return c.RetVal
}

func (c *Context) Clone() *Context {
	c.lock.Lock()
	defer c.lock.Unlock()
	newContext := NewContext(false, c.IsDebug, c.CanEval)
	newContext.Args = c.Args
	newContext.debugLogger = c.debugLogger
	newContext.ImportFunc = c.ImportFunc
	newContext.Stdin = c.Stdin
	newContext.Stdout = c.Stdout
	newContext.Stderr = c.Stderr
	newContext.modules = c.modules
	newContext.curFrame.lineNum = c.curFrame.lineNum
	newContext.curFrame.filename = c.curFrame.filename
	return newContext
}

func (c *Context) Recover() {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case Exception:
			io.WriteString(c.Stderr, e.MessageWithStack())
		default:
			fmt.Fprintf(c.Stderr, "error at: %s:%d\n", c.curFrame.filename, c.curFrame.lineNum)
			if c.IsDebug {
				panic(e)
			}
		}
	}
}

func (c *Context) StartThread(callee Value, this Value, args []Value) func() Value {
	newContext := c.Clone()
	done := make(chan Value, 1)
	go func() {
		defer func() {
			defer newContext.Recover()
			done <- newContext.RetVal
			close(done)
		}()
		newContext.Invoke(callee, this, func() []Value { return args })
	}()
	ret := NewObject()
	var rv Value
	var joinLock sync.Mutex
	joinFunc := func() Value {
		joinLock.Lock()
		defer joinLock.Unlock()
		ret, ok := <-done
		if ok {
			rv = ret
		}
		return rv
	}
	ret.SetMember("join", NewNativeFunction("join", func(c *Context, this Value, args []Value) Value {
		return joinFunc()
	}), c)
	ret.SetMember("await", NewNativeFunction("await", func(c *Context, this Value, args []Value) Value {
		return joinFunc()
	}), c)
	c.RetVal = ret
	return joinFunc
}

func (c *Context) CloneFrames() (cur, funcRoot, root *contextFrame) {
	var last *contextFrame = nil
	for p := c.curFrame; p != nil; p = p.parent {
		frame := p.clone()
		frame.parent = nil
		if last != nil {
			last.parent = frame
		}
		if p == c.curFrame {
			cur = frame
		} else if p == c.funcRootFrame {
			funcRoot = frame
		}
		last = frame
	}
	root = last
	return
}

func (c *Context) Eval(code string, force bool) Value {
	if !force && !c.CanEval {
		c.OnRuntimeError("eval is forbidden!")
	}
	val, _, ok := c.ImportFunc(c, "", code, ImportTypeScript, 0)
	if !ok {
		c.OnRuntimeError("eval error")
		return nil
	}
	return val
}

func (c *Context) AutoImport() {
	var roots []string
	if zggpath := os.Getenv("ZGGPATH"); zggpath != "" {
		roots = strings.Split(zggpath, ":")
	} else {
		return
	}
	for _, root := range roots {
		filename := filepath.Join(root, "_autoimport.zgg")
		if _, err := os.Stat(filename); err != nil {
			continue
		}
		val, _, ok := c.ImportFunc(c, filename, "", ImportTypeScript, 0)
		if !ok {
			continue
		}
		if exports, ok := val.(ValueObject); ok {
			exports.Iterate(func(key string, value Value) {
				c.ForceSetLocalValue(key, value)
			})
		}
	}
}

func (c *Context) valuesCompare(v1, v2 Value, fn string, expectedResults ...CompareResult) bool {
	if compFn, ok := v1.GetMember(fn, c).(ValueCallable); ok && c.IsCallable(compFn) {
		c.Invoke(compFn, v1, Args(v2))
		return c.RetVal.IsTrue()
	}
	compareResult := v1.CompareTo(v2, c)
	for _, r := range expectedResults {
		if compareResult == r {
			return true
		}
	}
	return false
}

func (c *Context) ValuesEqual(v1, v2 Value) bool {
	return c.valuesCompare(v1, v2, "__eq__", CompareResultEqual)
}

func (c *Context) ValuesNotEqual(v1, v2 Value) bool {
	return !c.ValuesEqual(v1, v2)
}

func (c *Context) ValuesLess(v1, v2 Value) bool {
	return c.valuesCompare(v1, v2, "__lt__", CompareResultLess)
}

func (c *Context) ValuesGreater(v1, v2 Value) bool {
	return c.valuesCompare(v1, v2, "__gt__", CompareResultGreater)
}

func (c *Context) ValuesLessEqual(v1, v2 Value) bool {
	return c.valuesCompare(v1, v2, "__lt__", CompareResultLess, CompareResultEqual)
}

func (c *Context) ValuesGreaterEqual(v1, v2 Value) bool {
	return c.valuesCompare(v1, v2, "__gt__", CompareResultGreater, CompareResultEqual)
}
