package builtin_libs

import (
	"errors"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/samber/lo"
	. "github.com/zgg-lang/zgg-go/runtime"
)

func libPath(c *Context) ValueObject {
	lib := NewObject()
	pathSetBridge(lib, c, "abs", filepath.Abs, "path")
	pathSetBridge(lib, c, "base", filepath.Base, "path")
	pathSetBridge(lib, c, "clean", filepath.Clean, "path")
	pathSetBridge(lib, c, "dir", filepath.Dir, "path")
	pathSetBridge(lib, c, "evalSymlinks", filepath.EvalSymlinks, "path")
	pathSetBridge(lib, c, "ext", filepath.Ext, "path")
	pathSetBridge(lib, c, "fromSlash", filepath.FromSlash, "path")
	lib.SetMember("glob", NewNativeFunction("path.glob", func(c *Context, _ Value, args []Value) Value {
		var pattern ValueStr
		EnsureFuncParams(c, "path.glob", args, ArgRuleRequired("pattern", TypeStr, &pattern))
		matches, err := filepath.Glob(pattern.Value())
		if err != nil {
			c.RaiseRuntimeError("glob error %+v", err)
		}
		r := NewArray(len(matches))
		for _, m := range matches {
			r.PushBack(NewStr(m))
		}
		return r
	}, "pattern"), nil)
	pathSetBridge(lib, c, "isAbs", filepath.IsAbs, "path")
	lib.SetMember("join", NewNativeFunction("path.join", func(c *Context, _ Value, args []Value) Value {
		elem := make([]string, len(args))
		for i, a := range args {
			elem[i] = a.ToString(c)
		}
		joined := filepath.Join(elem...)
		return NewStr(joined)
	}), nil)
	pathSetBridge(lib, c, "match", filepath.Match, "pattern", "name")
	pathSetBridge(lib, c, "rel", filepath.Rel, "basepath", "targetpath")
	pathSetBridge(lib, c, "toSlash", filepath.ToSlash, "path")
	pathSetBridge(lib, c, "volumeName", filepath.VolumeName, "path")
	lib.SetMember("walk", NewNativeFunction("path.walk", func(c *Context, _ Value, args []Value) Value {
		var (
			root ValueStr
			cb   ValueCallable
		)
		switch len(args) {
		case 1:
			EnsureFuncParams(c, "path.walk", args, ArgRuleRequired("root", TypeStr, &root))
			ch := make(chan Value)
			walking := true
			go func(root string) {
				defer func() {
					walking = false
					close(ch)
				}()
				filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !walking {
						return errors.New("end")
					}
					item := NewObject()
					item.SetMember("path", NewStr(path), nil)
					item.SetMember("name", NewStr(d.Name()), nil)
					item.SetMember("isDir", NewBool(d.IsDir()), nil)
					// 避免死锁，加上send超时
					t := time.NewTicker(5 * time.Second)
					defer t.Stop()
					for walking {
						select {
						case ch <- item:
							return nil
						case <-t.C:
						}
					}
					return nil
				})
			}(root.Value())
			return MakeIterator(c,
				func() Value {
					item, ok := <-ch
					return lo.Ternary(ok, item, nil)
				},
				func() {
					walking = false
				},
			)
		default:
			EnsureFuncParams(c, "path.walk", args,
				ArgRuleRequired("root", TypeStr, &root),
				ArgRuleRequired("callback", TypeCallable, &cb),
			)
			filepath.WalkDir(root.Value(), func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				item := NewObject()
				item.SetMember("path", NewStr(path), nil)
				item.SetMember("name", NewStr(d.Name()), nil)
				item.SetMember("isDir", NewBool(d.IsDir()), nil)
				c.Invoke(cb, nil, Args(item))
				return nil
			})
			return Undefined()
		}
	}, "root"), nil)
	return lib
}

func pathSetBridge(lib ValueObject, c *Context, name string, f any, argNames ...string) {
	fullname := "path." + name
	var nf func(*Context, Value, []Value) Value
	switch fv := f.(type) {
	case func(string) string:
		nf = func(c *Context, _ Value, args []Value) Value {
			var a ValueStr
			EnsureFuncParams(c, fullname, args, ArgRuleRequired(argNames[0], TypeStr, &a))
			return NewStr(fv(a.Value()))
		}
	case func(string) (string, error):
		nf = func(c *Context, _ Value, args []Value) Value {
			var a ValueStr
			EnsureFuncParams(c, fullname, args, ArgRuleRequired(argNames[0], TypeStr, &a))
			r, err := fv(a.Value())
			if err != nil {
				c.RaiseRuntimeError("%s error %+v", fullname, err)
			}
			return NewStr(r)
		}
	case func(string) bool:
		nf = func(c *Context, _ Value, args []Value) Value {
			var a ValueStr
			EnsureFuncParams(c, fullname, args, ArgRuleRequired(argNames[0], TypeStr, &a))
			return NewBool(fv(a.Value()))
		}
	case func(string, string) (bool, error):
		nf = func(c *Context, _ Value, args []Value) Value {
			var a, b ValueStr
			EnsureFuncParams(c, fullname, args,
				ArgRuleRequired(argNames[0], TypeStr, &a),
				ArgRuleRequired(argNames[1], TypeStr, &b),
			)
			r, err := fv(a.Value(), b.Value())
			if err != nil {
				c.RaiseRuntimeError("%s error %+v", fullname, err)
			}
			return NewBool(r)
		}
	case func(string, string) (string, error):
		nf = func(c *Context, _ Value, args []Value) Value {
			var a, b ValueStr
			EnsureFuncParams(c, fullname, args,
				ArgRuleRequired(argNames[0], TypeStr, &a),
				ArgRuleRequired(argNames[1], TypeStr, &b),
			)
			r, err := fv(a.Value(), b.Value())
			if err != nil {
				c.RaiseRuntimeError("%s error %+v", fullname, err)
			}
			return NewStr(r)
		}
	default:
		panic("invalid func type " + fullname)
	}
	lib.SetMember(name, NewNativeFunction(fullname, nf, argNames...), c)
}
