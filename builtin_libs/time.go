package builtin_libs

import (
	"regexp"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var timeClass ValueType

func libTime(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Time", timeClass, nil)
	lib.SetMember("time", NewNativeFunction("time", func(c *Context, this Value, args []Value) Value {
		return NewInt(time.Now().Unix())
	}), nil)
	lib.SetMember("now", NewNativeFunction("now", func(c *Context, this Value, args []Value) Value {
		var timeType ValueStr
		EnsureFuncParams(c, "now", args, ArgRuleOptional{"timeType", TypeStr, &timeType, NewStr("")})
		nowTs := time.Now().UnixNano()
		asType := timeType.Value()
		mod := int64(1)
		switch asType {
		case "day":
			mod = 86400 * 1e9
		case "hour":
			mod = 3600 * 1e9
		case "minute":
			mod = 60 * 1e9
		case "second":
			mod = 1e9
		case "":
		default:
			c.OnRuntimeError("Invalid time type %s", asType)
		}
		now := NewObjectAndInit(timeClass, c, NewInt(nowTs-nowTs%mod))
		if asType != "" {
			now.SetMember("__as", timeType, c)
		}
		return now
	}), nil)
	lib.SetMember("fromUnix", NewNativeFunction("fromUnix", func(c *Context, this Value, args []Value) Value {
		var ts ValueInt
		EnsureFuncParams(c, "time.fromUnix", args, ArgRuleRequired{"unixTimestamp", TypeInt, &ts})
		return NewObjectAndInit(timeClass, c, NewInt(ts.Value()*1e9))
	}), nil)
	lib.SetMember("sleep", NewNativeFunction("time", func(c *Context, this Value, args []Value) Value {
		if len(args) != 1 {
			c.OnRuntimeError("sleep: requires 1 argument")
			return nil
		}
		sleepSeconds := c.MustFloat(args[0])
		time.Sleep(time.Duration(sleepSeconds * float64(time.Second)))
		return Undefined()
	}), nil)
	return lib
}

func initTimeClass() {
	layoutP := regexp.MustCompile("%.")
	timeClass = NewClassBuilder("Time").
		Constructor(func(c *Context, thisObj ValueObject, args []Value) {
			var ts int64
			var as string
			switch len(args) {
			case 0:
				ts = time.Now().UnixNano()
			case 1:
				switch v := args[0].(type) {
				case ValueInt:
					ts = v.Value()
				case ValueStr:
					{
						var layout string
						switch v.Len() {
						case 8:
							layout = "20060102"
							as = "day"
						case 10:
							layout = "2006-01-02"
							as = "day"
						case 14:
							layout = "20060102150405"
							as = "second"
						case 19:
							layout = "2006-01-02 15:04:05"
							as = "second"
						default:
							c.OnRuntimeError("Time.__init__: invalid time str %s", v.Value())
						}
						t, err := time.Parse(layout, v.Value())
						if err != nil {
							c.OnRuntimeError("Time.__init__: parse time error %s", err)
						}
						ts = t.UnixNano()
					}
				}
			case 3:
				{
					var year, month, day ValueInt
					EnsureFuncParams(c, "Time.__init__", args,
						ArgRuleRequired{"year", TypeInt, &year},
						ArgRuleRequired{"month", TypeInt, &month},
						ArgRuleRequired{"day", TypeInt, &day},
					)
					ts = time.Date(year.AsInt(), time.Month(month.AsInt()), day.AsInt(), 0, 0, 0, 0, time.Local).UnixNano()
				}
			case 6:
				{
					var year, month, day, hour, minute, second ValueInt
					EnsureFuncParams(c, "Time.__init__", args,
						ArgRuleRequired{"year", TypeInt, &year},
						ArgRuleRequired{"month", TypeInt, &month},
						ArgRuleRequired{"day", TypeInt, &day},
					)
					ts = time.Date(
						year.AsInt(),
						time.Month(month.AsInt()),
						day.AsInt(),
						hour.AsInt(),
						minute.AsInt(),
						second.AsInt(),
						0,
						time.Local,
					).UnixNano()
				}
			default:
				c.OnRuntimeError("Time.__init__: invalid args")
			}
			thisObj.SetMember("unix", NewInt(ts/1e9), c)
			thisObj.SetMember("unixNano", NewInt(ts), c)
			t := time.Unix(ts/1e9, ts%1e9)
			thisObj.SetMember("_t", NewGoValue(t), c)
			thisObj.SetMember("year", NewInt(int64(t.Year())), c)
			thisObj.SetMember("month", NewInt(int64(t.Month())), c)
			thisObj.SetMember("day", NewInt(int64(t.Day())), c)
			thisObj.SetMember("hour", NewInt(int64(t.Hour())), c)
			thisObj.SetMember("minute", NewInt(int64(t.Minute())), c)
			thisObj.SetMember("second", NewInt(int64(t.Second())), c)
			if as != "" {
				thisObj.SetMember("__as", NewStr(as), c)
			}
		}).
		Method("addDays", func(c *Context, this ValueObject, args []Value) Value {
			var days ValueInt
			EnsureFuncParams(c, "Time.addDays", args,
				ArgRuleRequired{"days", TypeInt, &days},
			)
			ts := c.MustInt(this.GetMember("unixNano", c))
			nextNano := ts + days.Value()*86400*1e9
			rv := NewObjectAndInit(timeClass, c, NewInt(nextNano))
			rv.SetMember("__as", this.GetMember("__as", c), c)
			return rv
		}).
		Method("addHours", func(c *Context, this ValueObject, args []Value) Value {
			var hours ValueInt
			EnsureFuncParams(c, "Time.addHours", args,
				ArgRuleRequired{"days", TypeInt, &hours},
			)
			ts := c.MustInt(this.GetMember("unixNano", c))
			nextNano := ts + hours.Value()*3600*1e9
			rv := NewObjectAndInit(timeClass, c, NewInt(nextNano))
			return rv
		}).
		Method("as", func(c *Context, this ValueObject, args []Value) Value {
			var asType ValueStr
			EnsureFuncParams(c, "Time.as", args, ArgRuleRequired{"asType", TypeStr, &asType})
			this.SetMember("__as", asType, c)
			return this
		}).
		Method("__next__", func(c *Context, this ValueObject, args []Value) Value {
			if as, ok := this.GetMember("__as", c).(ValueStr); ok {
				var r Value
				ts := c.MustInt(this.GetMember("unixNano", c))
				switch as.Value() {
				case "day":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+86400*1e9))
				case "hour":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+3600*1e9))
				case "minute":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+60*1e9))
				case "second":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+1*1e9))
				}
				if r != nil {
					r.(ValueObject).SetMember("__as", as, c)
					return r
				}
			}
			c.OnRuntimeError("Time object cannot get __next__ without specialized time type")
			return nil
		}).
		Method("format", func(c *Context, this ValueObject, args []Value) Value {
			var (
				layout   ValueStr
				timezone ValueStr
			)
			switch len(args) {
			case 2:
				EnsureFuncParams(c, "Time.format", args,
					ArgRuleRequired{"layout", TypeStr, &layout},
					ArgRuleRequired{"timezone", TypeStr, &timezone},
				)
			default:
				EnsureFuncParams(c, "Time.format", args,
					ArgRuleRequired{"layout", TypeStr, &layout},
				)
				timezone = NewStr("")
			}
			ts := c.MustInt(this.GetMember("unixNano", c))
			t := time.Unix(ts/1e9, ts%1e9)
			layoutStr := layoutP.ReplaceAllStringFunc(layout.Value(), func(s string) string {
				switch s {
				case "%Y":
					return "2006"
				case "%y":
					return "06"
				case "%m":
					return "01"
				case "%d":
					return "02"
				case "%H":
					return "15"
				case "%M":
					return "04"
				case "%S":
					return "05"
				}
				return s
			})
			if tz := timezone.Value(); tz != "" {
				loc, err := time.LoadLocation(tz)
				if err != nil {
					c.OnRuntimeError("Invalid timezone %s", tz)
				}
				t = t.In(loc)
			}
			return NewStr(t.Format(layoutStr))
		}).
		Method("__add__", func(c *Context, this ValueObject, args []Value) Value {
			var diff ValueInt
			EnsureFuncParams(c, "Time.__add__", args, ArgRuleRequired{"add", TypeInt, &diff})
			if as, ok := this.GetMember("__as", c).(ValueStr); ok {
				var r Value
				ts := c.MustInt(this.GetMember("unixNano", c))
				switch as.Value() {
				case "day":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+diff.Value()*86400*1e9))
				case "hour":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+diff.Value()*3600*1e9))
				case "minute":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+diff.Value()*60*1e9))
				case "second":
					r = NewObjectAndInit(timeClass, c, NewInt(ts+diff.Value()*1e9))
				}
				if r != nil {
					r.(ValueObject).SetMember("__as", as, c)
					return r
				}
			}
			c.OnRuntimeError("Time object cannot get __next__ without specialized time type")
			return nil
		}).
		Method("__sub__", func(c *Context, this ValueObject, args []Value) Value {
			var diff ValueInt
			EnsureFuncParams(c, "Time.__sub__", args, ArgRuleRequired{"add", TypeInt, &diff})
			if as, ok := this.GetMember("__as", c).(ValueStr); ok {
				var r Value
				ts := c.MustInt(this.GetMember("unixNano", c))
				switch as.Value() {
				case "day":
					r = NewObjectAndInit(timeClass, c, NewInt(ts-diff.Value()*86400*1e9))
				case "hour":
					r = NewObjectAndInit(timeClass, c, NewInt(ts-diff.Value()*3600*1e9))
				case "minute":
					r = NewObjectAndInit(timeClass, c, NewInt(ts-diff.Value()*60*1e9))
				case "second":
					r = NewObjectAndInit(timeClass, c, NewInt(ts-diff.Value()*1e9))
				}
				if r != nil {
					r.(ValueObject).SetMember("__as", as, c)
					return r
				}
			}
			c.OnRuntimeError("Time object cannot get __next__ without specialized time type")
			return nil
		}).
		Method("timetuple", func(c *Context, this ValueObject, args []Value) Value {
			return NewArrayByValues(
				this.GetMember("year", c),
				this.GetMember("month", c),
				this.GetMember("day", c),
				this.GetMember("hour", c),
				this.GetMember("minute", c),
				this.GetMember("second", c),
			)
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			ts := c.MustInt(this.GetMember("unixNano", c))
			t := time.Unix(ts/1e9, ts%1e9)
			if as, ok := this.GetMember("__as", c).(ValueStr); ok {
				switch as.Value() {
				case "day":
					return NewStr(t.Format("2006-01-02"))
				case "hour":
					return NewStr(t.Format("2006-01-02 15"))
				case "minute":
					return NewStr(t.Format("2006-01-02 15:04"))
				case "second":
					return NewStr(t.Format("2006-01-02 15:04:05"))
				}
			}
			return NewStr(t.Format("Time(2006-01-02 15:04:05)"))
		}).
		Method("__lt__", func(c *Context, this ValueObject, args []Value) Value {
			ts := c.MustInt(this.GetMember("unixNano", c))
			other := c.MustInt(args[0].GetMember("unixNano", c))
			return NewBool(ts < other)
		}).
		Method("__le__", func(c *Context, this ValueObject, args []Value) Value {
			ts := c.MustInt(this.GetMember("unixNano", c))
			other := c.MustInt(args[0].GetMember("unixNano", c))
			return NewBool(ts <= other)
		}).
		Method("__gt__", func(c *Context, this ValueObject, args []Value) Value {
			ts := c.MustInt(this.GetMember("unixNano", c))
			other := c.MustInt(args[0].GetMember("unixNano", c))
			return NewBool(ts > other)
		}).
		Method("__ge__", func(c *Context, this ValueObject, args []Value) Value {
			ts := c.MustInt(this.GetMember("unixNano", c))
			other := c.MustInt(args[0].GetMember("unixNano", c))
			return NewBool(ts >= other)
		}).
		Method("__eq__", func(c *Context, this ValueObject, args []Value) Value {
			ts := c.MustInt(this.GetMember("unixNano", c))
			other := c.MustInt(args[0].GetMember("unixNano", c))
			return NewBool(ts == other)
		}).
		Method("__ne__", func(c *Context, this ValueObject, args []Value) Value {
			ts := c.MustInt(this.GetMember("unixNano", c))
			other := c.MustInt(args[0].GetMember("unixNano", c))
			return NewBool(ts != other)
		}).
		Build()
}

func init() {
	initTimeClass()
}
