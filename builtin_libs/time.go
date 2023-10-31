package builtin_libs

import (
	"regexp"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	timeTimeClass     ValueType
	timeDurationClass ValueType
)

type timeTimeArg struct {
	o     ValueObject
	s     ValueStr
	i     ValueInt
	which int
}

func (a *timeTimeArg) Rule(name string) ArgRule {
	return ArgRuleOneOf(name,
		[]ValueType{timeTimeClass, TypeStr, TypeInt},
		[]any{&a.o, &a.s, &a.i},
		&a.which,
		nil, nil)
}

func (a *timeTimeArg) Get(c *Context) ValueObject {
	switch a.which {
	case 1:
		a.o = NewObjectAndInit(timeTimeClass, c, a.s)
	case 2:
		if ts := a.i.Value(); ts < 10000000000 {
			a.o = NewObjectAndInit(timeTimeClass, c, NewInt(ts*1000000000))
		} else if ts := a.i.Value(); ts < 10000000000000 {
			a.o = NewObjectAndInit(timeTimeClass, c, NewInt(ts*1000000))
		} else {
			a.o = NewObjectAndInit(timeTimeClass, c, a.i)
		}
	}
	return a.o
}

func (a *timeTimeArg) GetTime(c *Context) time.Time {
	return a.Get(c).Reserved.(timeTimeInfo).t
}

type timeDurationArg struct {
	o     ValueObject
	s     ValueStr
	f     ValueFloat
	which int
}

func (a *timeDurationArg) Rule(name string, dv ...ValueObject) ArgRule {
	if len(dv) > 0 {
		return ArgRuleOneOf(name,
			[]ValueType{timeDurationClass, TypeStr, TypeFloat},
			[]any{&a.o, &a.s, &a.f},
			&a.which,
			&a.o, dv[0])
	}
	return ArgRuleOneOf(name,
		[]ValueType{timeDurationClass, TypeStr, TypeFloat},
		[]any{&a.o, &a.s, &a.f},
		&a.which,
		nil, nil)
}

func (a *timeDurationArg) Get(c *Context) ValueObject {
	switch a.which {
	case 1:
		a.o = NewObjectAndInit(timeDurationClass, c, a.s)
	case 2:
		dur := time.Duration(a.f.Value() * float64(time.Second))
		a.o = NewObjectAndInit(timeDurationClass, c, NewGoValue(dur))
	}
	return a.o
}

func (a *timeDurationArg) GetDuration(c *Context) time.Duration {
	return a.Get(c).ToGoValue().(time.Duration)
}

func libTime(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Time", timeTimeClass, nil)
	lib.SetMember("__call__", timeTimeClass, nil)
	lib.SetMember("Duration", timeDurationClass, nil)
	lib.SetMember("time", NewNativeFunction("time", func(c *Context, this Value, args []Value) Value {
		return NewInt(time.Now().Unix())
	}), nil)
	lib.SetMember("now", NewNativeFunction("now", func(c *Context, this Value, args []Value) Value {
		var timeType ValueStr
		EnsureFuncParams(c, "now", args, ArgRuleOptional("timeType", TypeStr, &timeType, NewStr("")))
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
			c.RaiseRuntimeError("Invalid time type %s", asType)
		}
		now := NewObjectAndInit(timeTimeClass, c, NewInt(nowTs-nowTs%mod))
		if asType != "" {
			now.SetMember("__as", timeType, c)
		}
		return now
	}), nil)
	lib.SetMember("fromUnix", NewNativeFunction("fromUnix", func(c *Context, this Value, args []Value) Value {
		var ts ValueInt
		EnsureFuncParams(c, "time.fromUnix", args, ArgRuleRequired("unixTimestamp", TypeInt, &ts))
		return NewObjectAndInit(timeTimeClass, c, NewInt(ts.Value()*1e9))
	}, "timestamp"), nil)
	lib.SetMember("fromGoTime", NewNativeFunction("fromGoTime", func(c *Context, this Value, args []Value) Value {
		var gt GoValue
		EnsureFuncParams(c, "time.fromGoTime", args, ArgRuleRequired("time", TypeGoValue, &gt))
		if _, ok := gt.ToGoValue().(time.Time); !ok {
			c.RaiseRuntimeError("Not a time.Time!")
		}
		return NewObjectAndInit(timeTimeClass, c, gt)
	}, "time"), nil)
	lib.SetMember("sleep", NewNativeFunction("sleep", func(c *Context, this Value, args []Value) Value {
		var sleepDur timeDurationArg
		EnsureFuncParams(c, "time.sleep", args, sleepDur.Rule("duration"))
		time.Sleep(sleepDur.Get(c).Reserved.(time.Duration))
		return Undefined()
	}), nil)
	lib.SetMember("since", NewNativeFunction("since", func(c *Context, this Value, args []Value) Value {
		var from timeTimeArg
		EnsureFuncParams(c, "since", args, from.Rule("from"))
		du := time.Since(from.GetTime(c))
		return NewObjectAndInit(timeDurationClass, c, NewGoValue(du))
	}), nil)
	oneDay := NewObjectAndInit(timeDurationClass, c, NewGoValue(24*time.Hour))
	_iter := func(c *Context, args []Value, canCallback bool) (begin, end, step ValueObject, callback ValueCallable) {
		var (
			beginArg timeTimeArg
			endArg   timeTimeArg
			stepArg  timeDurationArg
			rules    = []ArgRule{
				beginArg.Rule("begin"),
				endArg.Rule("end"),
				stepArg.Rule("step", oneDay),
			}
		)
		if canCallback {
			rules = append(rules, ArgRuleOptional("callback", TypeCallable, &callback, nil))
		}
		EnsureFuncParams(c, "iter", args, rules...)
		begin = beginArg.Get(c)
		end = endArg.Get(c)
		step = stepArg.Get(c)
		return
	}
	lib.SetMember("iter", NewNativeFunction("iter", func(c *Context, this Value, args []Value) Value {
		begin, end, step, callback := _iter(c, args, true)
		if callback != nil {
			current := begin
			for c.ValuesLess(current, end) {
				c.Invoke(callback, nil, Args(current))
				current = c.InvokeMethod(current, "__add__", Args(step)).(ValueObject)
			}
			return Undefined()
		}
		rv := NewObject()
		rv.SetMember("__iter__", NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
			next := begin
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				cur := next
				if !c.ValuesLess(cur, end) {
					return NewArrayByValues(Undefined(), NewBool(false))
				}
				next = c.InvokeMethod(cur, "__add__", Args(step)).(ValueObject)
				return NewArrayByValues(cur, NewBool(true))
			})
		}), c)
		return rv
	}), nil)
	lib.SetMember("list", NewNativeFunction("list", func(c *Context, this Value, args []Value) Value {
		begin, end, step, _ := _iter(c, args, false)
		current := begin
		rv := NewArray()
		for c.ValuesLess(current, end) {
			rv.PushBack(current)
			current = c.InvokeMethod(current, "__add__", Args(step)).(ValueObject)
		}
		return rv
	}), c)
	lib.SetMember("timeit", NewNativeFunction("timeit", func(c *Context, this Value, args []Value) (rv Value) {
		var callable ValueCallable
		EnsureFuncParams(c, "timeit", args,
			ArgRuleRequired("callable", TypeCallable, &callable),
		)
		t0 := time.Now()
		defer func() {
			rv = NewInt(time.Now().UnixNano() - t0.UnixNano())
		}()
		c.Invoke(callable, nil, NoArgs)
		return
	}), nil)
	lib.SetMember("as", NewNativeFunction("as", func(c *Context, this Value, args []Value) (rv Value) {
		var as ValueStr
		EnsureFuncParams(c, "as", args, ArgRuleRequired("as", TypeStr, &as))
		info := timeTimeInfo{
			t:  time.Now(),
			as: as.Value(),
		}
		return NewObjectAndInit(timeTimeClass, c, NewGoValue(info))
	}), nil)
	return lib
}

type timeTimeInfo struct {
	t  time.Time
	as string
}

func timeInittimeTimeClass() {
	layoutP := regexp.MustCompile("%.")
	mustTime := func(c *Context, name string, args []Value) ValueObject {
		var t ValueObject
		EnsureFuncParams(c, name, args, ArgRuleRequired("other", timeTimeClass, &t))
		return t
	}
	timeTimeClass = NewClassBuilder("Time").
		Constructor(func(c *Context, thisObj ValueObject, args []Value) {
			var _t time.Time
			var as string
			switch len(args) {
			case 0:
				_t = time.Now()
			case 1:
				switch v := args[0].(type) {
				case ValueInt:
					ts := v.Value()
					_t = time.Unix(ts/1e9, ts%1e9)
				case GoValue:
					switch gv := v.ToGoValue().(type) {
					case time.Time:
						_t = gv
					case timeTimeInfo:
						_t = gv.t
						as = gv.as
					default:
						c.RaiseRuntimeError("Time.__init__: invalid arg %v", gv)
					}
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
							c.RaiseRuntimeError("Time.__init__: invalid time str %s", v.Value())
						}
						t, err := time.Parse(layout, v.Value())
						if err != nil {
							c.RaiseRuntimeError("Time.__init__: parse time error %s", err)
						}
						// ts = t.UnixNano()
						_t = t
					}
				}
			case 2:
				{
					var timeStr, layout ValueStr
					EnsureFuncParams(c, "Time.__init__", args,
						ArgRuleRequired("timeStr", TypeStr, &timeStr),
						ArgRuleRequired("layout", TypeStr, &layout),
					)
					t, err := time.Parse(layout.Value(), timeStr.Value())
					if err != nil {
						c.RaiseRuntimeError("Time.__init__: parse time error %s", err)
					}
					_t = t
				}
			case 3:
				{
					var year, month, day ValueInt
					EnsureFuncParams(c, "Time.__init__", args,
						ArgRuleRequired("year", TypeInt, &year),
						ArgRuleRequired("month", TypeInt, &month),
						ArgRuleRequired("day", TypeInt, &day),
					)
					_t = time.Date(year.AsInt(), time.Month(month.AsInt()), day.AsInt(), 0, 0, 0, 0, time.Local)
				}
			case 6:
				{
					var year, month, day, hour, minute, second ValueInt
					EnsureFuncParams(c, "Time.__init__", args,
						ArgRuleRequired("year", TypeInt, &year),
						ArgRuleRequired("month", TypeInt, &month),
						ArgRuleRequired("day", TypeInt, &day),
						ArgRuleRequired("hour", TypeInt, &hour),
						ArgRuleRequired("minute", TypeInt, &minute),
						ArgRuleRequired("second", TypeInt, &second),
					)
					_t = time.Date(
						year.AsInt(),
						time.Month(month.AsInt()),
						day.AsInt(),
						hour.AsInt(),
						minute.AsInt(),
						second.AsInt(),
						0,
						time.Local,
					)
				}
			default:
				c.RaiseRuntimeError("Time.__init__: invalid args")
			}
			if as != "" {
				switch as {
				case "day":
					_t = time.Date(_t.Year(), _t.Month(), _t.Day(), 0, 0, 0, 0, _t.Location())
				case "hour":
					_t = time.Date(_t.Year(), _t.Month(), _t.Day(), _t.Hour(), 0, 0, 0, _t.Location())
				case "minute":
					_t = time.Date(_t.Year(), _t.Month(), _t.Day(), _t.Hour(), _t.Minute(), 0, 0, _t.Location())
				case "second":
					_t = time.Date(_t.Year(), _t.Month(), _t.Day(), _t.Hour(), _t.Minute(), _t.Second(), 0, _t.Location())
				default:
					c.RaiseRuntimeError("Invalid specialized time type %s", as)
				}
			}
			thisObj.Reserved = timeTimeInfo{
				t:  _t,
				as: as,
			}
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			var field ValueStr
			EnsureFuncParams(c, "Time.__getAttr__", args, ArgRuleRequired("field", TypeStr, &field))
			t := this.Reserved.(timeTimeInfo).t
			switch field.Value() {
			case "unix":
				return NewInt(t.Unix())
			case "unixNano":
				return NewInt(t.UnixNano())
			case "year":
				return NewInt(int64(t.Year()))
			case "month":
				return NewInt(int64(t.Month()))
			case "day":
				return NewInt(int64(t.Day()))
			case "hour":
				return NewInt(int64(t.Hour()))
			case "minute":
				return NewInt(int64(t.Minute()))
			case "second":
				return NewInt(int64(t.Second()))
			case "weekday":
				return NewInt(int64(t.Weekday()))
			}
			return Undefined()
		}).
		Method("add", func(c *Context, this ValueObject, args []Value) Value {
			var duration ValueStr
			EnsureFuncParams(c, "Time.add", args, ArgRuleRequired("duration", TypeStr, &duration))
			d, err := time.ParseDuration(duration.Value())
			if err != nil {
				c.RaiseRuntimeError("Invalid duration %s", duration.Value())
			}
			info := this.Reserved.(timeTimeInfo)
			info.t = info.t.Add(d)
			return NewObjectAndInit(timeTimeClass, c, NewGoValue(info))
		}).
		Method("addDays", func(c *Context, this ValueObject, args []Value) Value {
			var days ValueInt
			EnsureFuncParams(c, "Time.addDays", args,
				ArgRuleRequired("days", TypeInt, &days),
			)
			info := this.Reserved.(timeTimeInfo)
			info.t = info.t.AddDate(0, 0, days.AsInt())
			return NewObjectAndInit(timeTimeClass, c, NewGoValue(info))
		}).
		Method("addHours", func(c *Context, this ValueObject, args []Value) Value {
			var hours ValueInt
			EnsureFuncParams(c, "Time.addHours", args,
				ArgRuleRequired("days", TypeInt, &hours),
			)
			info := this.Reserved.(timeTimeInfo)
			info.t = info.t.Add(time.Duration(hours.AsInt()) * time.Hour)
			return NewObjectAndInit(timeTimeClass, c, NewGoValue(info))
		}).
		Method("as", func(c *Context, this ValueObject, args []Value) Value {
			var asType ValueStr
			EnsureFuncParams(c, "Time.as", args, ArgRuleRequired("asType", TypeStr, &asType))
			info := this.Reserved.(timeTimeInfo)
			info.as = asType.Value()
			if info.as != "" {
				t := info.t
				switch info.as {
				case "day":
					t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
				case "hour":
					t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
				case "minute":
					t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
				case "second":
					t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
				default:
					c.RaiseRuntimeError("Invalid specialized time type %s", info.as)
				}
				info.t = t
			}
			this.Reserved = info
			return this
		}).
		Method("__next__", func(c *Context, this ValueObject, args []Value) Value {
			var (
				info = this.Reserved.(timeTimeInfo)
				t    = info.t
			)
			switch info.as {
			case "day":
				info.t = t.AddDate(0, 0, 1)
			case "hour":
				info.t = t.Add(time.Hour)
			case "minute":
				info.t = t.Add(time.Minute)
			case "second":
				info.t = t.Add(time.Second)
			default:
				c.RaiseRuntimeError("Time object cannot get __next__ without specialized time type")
			}
			return NewObjectAndInit(timeTimeClass, c, NewGoValue(info))
		}).
		Method("format", func(c *Context, this ValueObject, args []Value) Value {
			var (
				layout   ValueStr
				timezone ValueStr
			)
			EnsureFuncParams(c, "Time.format", args,
				ArgRuleRequired("layout", TypeStr, &layout),
				ArgRuleOptional("timezone", TypeStr, &timezone, NewStr("")),
			)
			t := this.Reserved.(timeTimeInfo).t
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
					c.RaiseRuntimeError("Invalid timezone %s", tz)
				}
				t = t.In(loc)
			}
			return NewStr(t.Format(layoutStr))
		}).
		Method("__add__", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) != 1 {
				c.RaiseRuntimeError("__add__ requires one arugment!")
			}
			info := this.Reserved.(timeTimeInfo)
			newInfo := info
			t := info.t
			switch diff := args[0].(type) {
			case ValueObject:
				if diff.Type().TypeId != timeDurationClass.TypeId {
					c.RaiseRuntimeError("invalid duration class")
				}
				newInfo.t = t.Add(diff.Reserved.(time.Duration))
				return NewObjectAndInit(timeTimeClass, c, NewGoValue(newInfo))
			case ValueInt:
				d := diff.AsInt()
				switch info.as {
				case "day":
					newInfo.t = t.AddDate(0, 0, d)
				case "hour":
					newInfo.t = t.Add(time.Hour * time.Duration(d))
				case "minute":
					newInfo.t = t.Add(time.Minute * time.Duration(d))
				case "second":
					newInfo.t = t.Add(time.Second * time.Duration(d))
				default:
					c.RaiseRuntimeError("Time object cannot get __next__ without specialized time type")
				}
				return NewObjectAndInit(timeTimeClass, c, NewGoValue(newInfo))
			}
			return nil
		}).
		Method("__sub__", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) != 1 {
				c.RaiseRuntimeError("__sub__ requires one arugment!")
			}
			t := this.Reserved.(timeTimeInfo).t
			switch diff := args[0].(type) {
			case ValueObject:
				switch diff.Type().TypeId {
				case timeTimeClass.TypeId:
					t2 := diff.Reserved.(timeTimeInfo).t
					r := NewObjectAndInit(timeDurationClass, c, NewGoValue(t.Sub(t2)))
					return r
				case timeDurationClass.TypeId:
					du := diff.Reserved.(time.Duration)
					r := NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(-du)))
					return r
				default:
					c.RaiseRuntimeError("invalid duration class")
				}
			case ValueInt:
				if as, ok := this.GetMember("__as", c).(ValueStr); ok {
					var r Value
					t := this.Reserved.(timeTimeInfo).t
					d := -diff.AsInt()
					switch as.Value() {
					case "day":
						r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.AddDate(0, 0, d)))
					case "hour":
						r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(time.Hour*time.Duration(d))))
					case "minute":
						r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(time.Minute*time.Duration(d))))
					case "second":
						r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(time.Second*time.Duration(d))))
					}
					if r != nil {
						r.(ValueObject).SetMember("__as", as, c)
						return r
					}
				}
				c.RaiseRuntimeError("Time object cannot get __next__ without specialized time type")
			}
			return nil
		}).
		Method("timetuple", func(c *Context, this ValueObject, args []Value) Value {
			t := this.Reserved.(timeTimeInfo).t
			return NewArrayByValues(
				NewInt(int64(t.Year())),
				NewInt(int64(t.Month())),
				NewInt(int64(t.Day())),
				NewInt(int64(t.Hour())),
				NewInt(int64(t.Minute())),
				NewInt(int64(t.Second())),
				NewStr(t.Location().String()),
			)
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			var (
				info = this.Reserved.(timeTimeInfo)
				t    = this.Reserved.(timeTimeInfo).t
			)
			switch info.as {
			case "day":
				return NewStr(t.Format("2006-01-02"))
			case "hour":
				return NewStr(t.Format("2006-01-02 15"))
			case "minute":
				return NewStr(t.Format("2006-01-02 15:04"))
			case "second":
				return NewStr(t.Format("2006-01-02 15:04:05"))
			}
			return NewStr(t.Format("Time(2006-01-02 15:04:05)"))
		}).
		Method("__lt__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.Reserved.(timeTimeInfo).t
			t2 := mustTime(c, "__lt__", args).Reserved.(timeTimeInfo).t
			return NewBool(t1.Before(t2))
		}).
		Method("__le__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.Reserved.(timeTimeInfo).t
			t2 := mustTime(c, "__le__", args).Reserved.(timeTimeInfo).t
			return NewBool(!t1.After(t2))
		}).
		Method("__gt__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.Reserved.(timeTimeInfo).t
			t2 := mustTime(c, "__gt__", args).Reserved.(timeTimeInfo).t
			return NewBool(t1.After(t2))
		}).
		Method("__ge__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.Reserved.(timeTimeInfo).t
			t2 := mustTime(c, "__ge__", args).Reserved.(timeTimeInfo).t
			return NewBool(!t1.Before(t2))
		}).
		Method("__eq__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.Reserved.(timeTimeInfo).t
			t2 := mustTime(c, "__eq__", args).Reserved.(timeTimeInfo).t
			return NewBool(!t1.Before(t2) && !t1.After(t2))
		}).
		Method("__ne__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.Reserved.(timeTimeInfo).t
			t2 := mustTime(c, "__ne__", args).Reserved.(timeTimeInfo).t
			return NewBool(t1.Before(t2) || t1.After(t2))
		}).
		Method("toGoTime", func(c *Context, this ValueObject, args []Value) Value {
			return NewGoValue(this.Reserved.(timeTimeInfo).t)
		}).
		Build()
}

func timeInitDurationClass() {
	getOther := func(c *Context, args []Value) time.Duration {
		if len(args) != 1 {
			c.RaiseRuntimeError("duration compare requries 1 argument!")
		}
		var d2 time.Duration
		if d2a := c.MustObject(args[0]); d2a.Type().TypeId != timeDurationClass.TypeId {
			c.RaiseRuntimeError("invalid other duration type")
		} else {
			d2 = d2a.Reserved.(time.Duration)
		}
		return d2
	}
	compareDurations := func(c *Context, this ValueObject, args []Value) time.Duration {
		d1 := this.Reserved.(time.Duration)
		d2 := getOther(c, args)
		return d1 - d2
	}
	timeDurationClass = NewClassBuilder("Duration").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			switch len(args) {
			case 1:
				switch dv := args[0].(type) {
				case ValueInt:
					du := time.Duration(dv.Value()) * time.Second
					this.Reserved = du
					return
				case ValueFloat:
					du := time.Duration(dv.Value() * float64(time.Second))
					this.Reserved = du
					return
				case ValueStr:
					du, err := time.ParseDuration(dv.Value())
					if err != nil {
						c.RaiseRuntimeError("invalid duration string %s", dv.Value())
					}
					this.Reserved = du
					return
				case GoValue:
					if du, is := dv.ToGoValue().(time.Duration); is {
						this.Reserved = du
						return
					}
				}
			}
			c.RaiseRuntimeError("Duration.__init__: invalid duration argumenet")
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			du := this.Reserved.(time.Duration)
			return NewStr(du.String())
		}).
		// To floats
		Methods([]string{"nanoseconds", "ns"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.Reserved.(time.Duration)
			return NewFloat(float64(d.Nanoseconds()))
		}).
		Methods([]string{"milliseconds", "ms"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.Reserved.(time.Duration)
			return NewFloat(float64(d.Milliseconds()))
		}).
		Methods([]string{"seconds", "s"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.Reserved.(time.Duration)
			return NewFloat(d.Seconds())
		}).
		Methods([]string{"minutes", "m"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.Reserved.(time.Duration)
			return NewFloat(d.Minutes())
		}).
		Methods([]string{"hours", "h"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.Reserved.(time.Duration)
			return NewFloat(d.Hours())
		}).
		Methods([]string{"days", "d"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.Reserved.(time.Duration)
			return NewFloat(d.Hours() / 24)
		}).
		Methods([]string{"weeks", "w"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.Reserved.(time.Duration)
			return NewFloat(d.Hours() / 24 / 7)
		}).
		// Comparation
		Method("__lt__", func(c *Context, this ValueObject, args []Value) Value {
			return NewBool(compareDurations(c, this, args) < 0)
		}).
		Method("__le__", func(c *Context, this ValueObject, args []Value) Value {
			return NewBool(compareDurations(c, this, args) <= 0)
		}).
		Method("__gt__", func(c *Context, this ValueObject, args []Value) Value {
			return NewBool(compareDurations(c, this, args) > 0)
		}).
		Method("__ge__", func(c *Context, this ValueObject, args []Value) Value {
			return NewBool(compareDurations(c, this, args) >= 0)
		}).
		Method("__eq__", func(c *Context, this ValueObject, args []Value) Value {
			return NewBool(compareDurations(c, this, args) == 0)
		}).
		Method("__ne__", func(c *Context, this ValueObject, args []Value) Value {
			return NewBool(compareDurations(c, this, args) != 0)
		}).
		// Add & sub
		Method("__add__", func(c *Context, this ValueObject, args []Value) Value {
			d1 := this.Reserved.(time.Duration)
			d2 := getOther(c, args)
			du := time.Duration(int64(d1) + int64(d2))
			return NewObjectAndInit(timeDurationClass, c, NewGoValue(du))
		}).
		Method("__sub__", func(c *Context, this ValueObject, args []Value) Value {
			d1 := this.Reserved.(time.Duration)
			d2 := getOther(c, args)
			du := time.Duration(int64(d1) - int64(d2))
			return NewObjectAndInit(timeDurationClass, c, NewGoValue(du))
		}).
		Method("__mul__", func(c *Context, this ValueObject, args []Value) Value {
			d1 := this.Reserved.(time.Duration)
			var times ValueFloat
			EnsureFuncParams(c, "Duration.__mul__", args, ArgRuleRequired("times", TypeFloat, &times))
			du := time.Duration(float64(d1) * times.Value())
			return NewObjectAndInit(timeDurationClass, c, NewGoValue(du))
		}).
		Method("__div__", func(c *Context, this ValueObject, args []Value) Value {
			d1 := this.Reserved.(time.Duration)
			var times ValueFloat
			EnsureFuncParams(c, "Duration.__div__", args, ArgRuleRequired("times", TypeFloat, &times))
			if t := times.Value(); t == 0 {
				c.RaiseRuntimeError("division by zero")
				return nil
			} else {
				du := time.Duration(float64(d1) / t)
				return NewObjectAndInit(timeDurationClass, c, NewGoValue(du))
			}
		}).
		Build()
}

func init() {
	timeInittimeTimeClass()
	timeInitDurationClass()
}
