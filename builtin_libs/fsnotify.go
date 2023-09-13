package builtin_libs

import (
	"github.com/fsnotify/fsnotify"
	. "github.com/zgg-lang/zgg-go/runtime"
)

func libFsnotify(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("watch", NewNativeFunction("watch", fsnotifyWatch), nil)
	return lib
}

func fsnotifyWatch(c *Context, this Value, args []Value) Value {
	if len(args) < 1 {
		c.RaiseRuntimeError("required at least 1 watch path")
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		c.RaiseRuntimeError("create fsnotify watcher error %+v", err)
	}
	for _, a := range args {
		watcher.Add(a.ToString(c))
	}
	r := NewObject()
	r.SetMember("__iter__", fsnotifyMakeIter(watcher), c)
	r.SetMember("close", NewNativeFunction("", func(*Context, Value, []Value) Value {
		watcher.Close()
		return Undefined()
	}), c)
	return r
}

func fsnotifyMakeIter(watcher *fsnotify.Watcher) Value {
	return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
		return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return NewArrayByValues(Undefined(), NewBool(false))
				}
				ev := NewObject()
				ev.SetMember("name", NewStr(event.Name), c)
				ev.SetMember("op", NewStr(event.Op.String()), c)
				return NewArrayByValues(ev, NewBool(true))
			case err, ok := <-watcher.Errors:
				if !ok {
					return NewArrayByValues(Undefined(), NewBool(false))
				}
				ev := NewObject()
				ev.SetMember("name", NewStr(err.Error()), c)
				ev.SetMember("op", NewStr("ERROR"), c)
				return NewArrayByValues(ev, NewBool(true))
			}
		})
	})
}
