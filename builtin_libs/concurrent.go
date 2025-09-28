package builtin_libs

import (
	"context"
	"sync"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func concurrentStartFunc(c *Context, parentCtx context.Context, this Value, args []Value) Value {
	if len(args) < 1 {
		c.RaiseRuntimeError("concurrent.start: requires at least 1 argument")
		return nil
	}
	callee, isCallable := c.GetCallable(args[0])
	if !isCallable {
		c.RaiseRuntimeError("concurrent.start: argument 0 must callable")
		return nil
	}
	c.StartThread(parentCtx, callee, nil, args[1:])
	return c.RetVal
}

func libConcurrent(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("start", NewNativeFunction("start", func(c *Context, this Value, args []Value) Value {
		return concurrentStartFunc(c, context.Background(), this, args)
	}), nil)
	lib.SetMember("startChild", NewNativeFunction("start", func(c *Context, this Value, args []Value) Value {
		return concurrentStartFunc(c, c.Ctx, this, args)
	}), nil)
	lib.SetMember("all", NewNativeFunction("all", func(c *Context, this Value, args []Value) Value {
		rv := make([]Value, len(args))
		awaits := make([]func() Value, len(args))
		for i, arg := range args {
			callee, isCallable := c.GetCallable(arg)
			if !isCallable {
				c.RaiseRuntimeError("concurrent.all: argument %d must callable", i)
				return nil
			}
			awaits[i] = c.StartThread(c.Ctx, callee, nil, []Value{})
		}
		for i, await := range awaits {
			rv[i] = await()
		}
		return NewArrayByValues(rv...)
	}), nil)
	{
		objMutex := NewClassBuilder("Mutex").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				lock := new(sync.Mutex)
				this.Reserved = lock
			}).
			Method("lock", func(c *Context, this ValueObject, args []Value) Value {
				this.Reserved.(*sync.Mutex).Lock()
				return Undefined()
			}).
			Method("unlock", func(c *Context, this ValueObject, args []Value) Value {
				this.Reserved.(*sync.Mutex).Unlock()
				return Undefined()
			}).
			Method("run", func(c *Context, this ValueObject, args []Value) Value {
				switch len(args) {
				case 1:
					if !c.IsCallable(args[0]) {
						c.RaiseRuntimeError("run: first argument must callable")
						return nil
					}
				default:
					c.RaiseRuntimeError("run: requires 1 argument(s)")
					return nil
				}
				lockVal := this.Reserved.(*sync.Mutex)
				lockVal.Lock()
				defer lockVal.Unlock()
				c.Invoke(args[0], nil, NoArgs)
				rv := c.RetVal
				return rv
			}).
			Build()
		lib.SetMember("Mutex", objMutex, nil)
	}
	{
		noTimeout := NewObjectAndInit(timeDurationClass, c, NewGoValue(time.Duration(-1)))
		objChan := NewClassBuilder("Chan").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				var n ValueInt
				EnsureFuncParams(c, "Chan.__init__", args,
					ArgRuleOptional("n", TypeInt, &n, NewInt(1)),
				)
				ch := make(chan Value, n.AsInt())
				this.Reserved = ch
			}).
			Method("send", func(c *Context, this ValueObject, args []Value) Value {
				ch := this.Reserved.(chan Value)
				ch <- args[0]
				return Undefined()
			}).
			Method("recv", func(c *Context, this ValueObject, args []Value) Value {
				var n timeDurationArg
				EnsureFuncParams(c, "Chan.recv", args,
					n.Rule(c, "timeout", noTimeout),
				)
				ch := this.Reserved.(chan Value)
				timeout := n.GetDuration(c)
				if timeout < 0 {
					return <-ch
				}
				select {
				case <-time.After(timeout):
					return Undefined()
				case v := <-ch:
					return v
				}
			}).
			Build()
		lib.SetMember("Chan", objChan, nil)
	}
	{
		objLimiter := NewClassBuilder("Limiter").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				if len(args) > 1 {
					c.RaiseRuntimeError("Excutor.__init__: requires 0 or 1 argument(s)")
				}
				max := 1
				if len(args) == 1 {
					max = int(c.MustInt(args[0], "Limiter.__init__(maxConcurrent): maxConcurrent"))
					if max < 1 {
						c.RaiseRuntimeError("Limiter.__init__(maxConcurrent): maxConcurrent must > 0")
					}
				}
				info := &concurrentLimiterInfo{
					ch: make(chan struct{}, max),
				}
				this.Reserved = info
			}).
			Method("run", func(c *Context, this ValueObject, args []Value) Value {
				if len(args) < 1 {
					c.RaiseRuntimeError("run: requires 1 argument(s)")
					return nil
				}
				if !c.IsCallable(args[0]) {
					c.RaiseRuntimeError("run: first argument must callable")
					return nil
				}
				info := this.Reserved.(*concurrentLimiterInfo)
				info.ch <- struct{}{}
				info.wg.Add(1)
				joinFunc := c.StartThread(c.Ctx, args[0], nil, args[1:])
				rv := c.RetVal
				go func() {
					defer func() {
						defer c.Recover()
						info.wg.Done()
						<-info.ch
					}()
					joinFunc()
				}()
				return rv
			}).
			Method("wait", func(c *Context, this ValueObject, args []Value) Value {
				this.Reserved.(*concurrentLimiterInfo).wg.Wait()
				return Undefined()
			}).
			Build()
		lib.SetMember("Limiter", objLimiter, nil)
	}
	return lib
}

type concurrentLimiterInfo struct {
	ch chan struct{}
	wg sync.WaitGroup
}
