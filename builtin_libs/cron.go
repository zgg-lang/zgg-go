package builtin_libs

import (
	. "github.com/zgg-lang/zgg-go/runtime"

	cron "github.com/robfig/cron/v3"
)

var (
	cronAppClass ValueType
)

func libCron(c *Context) ValueObject {
	rv := NewObject()
	rv.SetMember("App", cronGetAppClass(), c)
	return rv
}

func cronGetAppClass() ValueType {
	if cronAppClass == nil {
		cronAppClass = NewClassBuilder("App").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				app := cron.New()
				this.SetMember("__app", NewGoValue(app), c)
			}).
			Method("add", func(c *Context, this ValueObject, args []Value) Value {
				var (
					spec     ValueStr
					callback ValueCallable
				)
				EnsureFuncParams(c, "App.add", args,
					ArgRuleRequired{"spec", TypeStr, &spec},
					ArgRuleRequired{"callback", TypeCallable, &callback},
				)
				app := this.GetMember("__app", c).ToGoValue().(*cron.Cron)
				newContext := c.Clone()
				id, err := app.AddFunc(spec.Value(), func() {
					defer newContext.Recover()
					newContext.Invoke(callback, nil, Args(this))
				})
				if err != nil {
					c.RaiseRuntimeError("CronApp.add: add callback error: %s", err)
				}
				return NewInt(int64(id))
			}).
			Method("remove", func(c *Context, this ValueObject, args []Value) Value {
				var jobId ValueInt
				EnsureFuncParams(c, "App.remove", args, ArgRuleRequired{"jobId", TypeInt, &jobId})
				app := this.GetMember("__app", c).ToGoValue().(*cron.Cron)
				app.Remove(cron.EntryID(jobId.AsInt()))
				return this
			}).
			Method("start", func(c *Context, this ValueObject, args []Value) Value {
				app := this.GetMember("__app", c).ToGoValue().(*cron.Cron)
				app.Start()
				return this
			}).
			Method("stop", func(c *Context, this ValueObject, args []Value) Value {
				app := this.GetMember("__app", c).ToGoValue().(*cron.Cron)
				app.Stop()
				return this
			}).
			Build()
	}
	return cronAppClass
}
