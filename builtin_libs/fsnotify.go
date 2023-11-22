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
	return MakeIterator(c, func() Value {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			ev := NewObject()
			ev.SetMember("name", NewStr(event.Name), c)
			ev.SetMember("op", NewStr(event.Op.String()), c)
			return ev
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			ev := NewObject()
			ev.SetMember("name", NewStr(err.Error()), c)
			ev.SetMember("op", NewStr("ERROR"), c)
			return ev
		}
		return nil
	}, func() {
		watcher.Close()
	})
}
