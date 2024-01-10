package builtin_libs

import (
	"sync"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libConcurrent(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("start", NewNativeFunction("start", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			c.RaiseRuntimeError("concurrent.start: requires at least 1 argument")
			return nil
		}
		callee, isCallable := c.GetCallable(args[0])
		if !isCallable {
			c.RaiseRuntimeError("concurrent.start: argument 0 must callable")
			return nil
		}
		c.StartThread(callee, nil, args[1:])
		return c.RetVal
	}), nil)
	lib.SetMember("all", NewNativeFunction("all", func(c *Context, this Value, args []Value) Value {
		rv := make([]Value, len(args))
		for i, arg := range args {
			callee, isCallable := c.GetCallable(arg)
			if !isCallable {
				c.RaiseRuntimeError("concurrent.all: argument %d must callable", i)
				return nil
			}
			c.StartThread(callee, nil, []Value{})
			rv[i] = c.RetVal
		}
		for i, t := range rv {
			await, ok := c.GetCallable(t.GetMember("await", c))
			if ok {
				c.Invoke(await, nil, NoArgs)
				rv[i] = c.RetVal
			}
		}
		return NewArrayByValues(rv...)
	}), nil)
	{
		objMutex := NewClassBuilder("Mutex").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				lock := new(sync.Mutex)
				this.SetMember("__lock", NewGoValue(lock), c)
			}).
			Method("lock", func(c *Context, this ValueObject, args []Value) Value {
				lockVal := this.GetMember("__lock", c)
				lockVal.ToGoValue(c).(*sync.Mutex).Lock()
				return Undefined()
			}).
			Method("unlock", func(c *Context, this ValueObject, args []Value) Value {
				lockVal := this.GetMember("__lock", c)
				lockVal.ToGoValue(c).(*sync.Mutex).Unlock()
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
				lockVal := this.GetMember("__lock", c).ToGoValue(c).(*sync.Mutex)
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
		objChan := NewClassBuilder("Chan").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				var n ValueInt
				EnsureFuncParams(c, "Chan.__init__", args,
					ArgRuleOptional("n", TypeInt, &n, NewInt(1)),
				)
				ch := make(chan Value, n.AsInt())
				this.SetMember("__ch", NewGoValue(&ch), c)
			}).
			Method("send", func(c *Context, this ValueObject, args []Value) Value {
				ch := this.GetMember("__ch", c).ToGoValue(c).(*chan Value)
				(*ch) <- args[0]
				return Undefined()
			}).
			Method("recv", func(c *Context, this ValueObject, args []Value) Value {
				var n ValueFloat
				EnsureFuncParams(c, "Chan.recv", args,
					ArgRuleOptional("timeout", TypeFloat, &n, NewFloat(-1)),
				)
				ch := this.GetMember("__ch", c).ToGoValue(c).(*chan Value)
				timeout := n.Value()
				if timeout < 0 {
					return <-*ch
				}
				select {
				case <-time.After(time.Duration(timeout) * time.Second):
					return Undefined()
				case v := <-*ch:
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
				ch := make(chan bool, max)
				this.SetMember("__ch", NewGoValue(&ch), c)
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
				ch := this.GetMember("__ch", c).ToGoValue(c).(*chan bool)
				*ch <- true
				joinFunc := c.StartThread(args[0], nil, []Value{})
				rv := c.RetVal
				go func() {
					joinFunc()
					<-*ch
				}()
				return rv
			}).
			Method("wait", func(c *Context, this ValueObject, args []Value) Value {
				ch := this.GetMember("__ch", c).ToGoValue(c).(*chan bool)
				for len(*ch) > 0 {
					time.Sleep(10 * time.Millisecond)
				}
				return Undefined()
			}).
			Build()
		lib.SetMember("Limiter", objLimiter, nil)
	}
	return lib
}
