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
	return a.Get(c).GetMember("__t", c).ToGoValue().(time.Time)
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
		time.Sleep(sleepDur.Get(c).GetMember("__du", c).ToGoValue().(time.Duration))
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
	return lib
}

func timeInittimeTimeClass() {
	layoutP := regexp.MustCompile("%.")
	timeTimeClass = NewClassBuilder("Time").
		Constructor(func(c *Context, thisObj ValueObject, args []Value) {
			// var ts int64
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
			thisObj.SetMember("__t", NewGoValue(_t), c)
			if as != "" {
				thisObj.SetMember("__as", NewStr(as), c)
			}
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			var field ValueStr
			EnsureFuncParams(c, "Time.__getAttr__", args, ArgRuleRequired("field", TypeStr, &field))
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
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
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
			rt := t.Add(d)
			rv := NewObjectAndInit(timeTimeClass, c, NewGoValue(rt))
			rv.SetMember("__as", this.GetMember("__as", c), c)
			return rv
		}).
		Method("addDays", func(c *Context, this ValueObject, args []Value) Value {
			var days ValueInt
			EnsureFuncParams(c, "Time.addDays", args,
				ArgRuleRequired("days", TypeInt, &days),
			)
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
			rv := NewObjectAndInit(timeTimeClass, c, NewGoValue(t.AddDate(0, 0, days.AsInt())))
			rv.SetMember("__as", this.GetMember("__as", c), c)
			return rv
		}).
		Method("addHours", func(c *Context, this ValueObject, args []Value) Value {
			var hours ValueInt
			EnsureFuncParams(c, "Time.addHours", args,
				ArgRuleRequired("days", TypeInt, &hours),
			)
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
			rv := NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(time.Duration(hours.AsInt())*time.Hour)))
			rv.SetMember("__as", this.GetMember("__as", c), c)
			return rv
		}).
		Method("as", func(c *Context, this ValueObject, args []Value) Value {
			var asType ValueStr
			EnsureFuncParams(c, "Time.as", args, ArgRuleRequired("asType", TypeStr, &asType))
			this.SetMember("__as", asType, c)
			return this
		}).
		Method("__next__", func(c *Context, this ValueObject, args []Value) Value {
			if as, ok := this.GetMember("__as", c).(ValueStr); ok {
				var r Value
				t := this.GetMember("__t", c).ToGoValue().(time.Time)
				switch as.Value() {
				case "day":
					r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.AddDate(0, 0, 1)))
				case "hour":
					r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(time.Hour)))
				case "minute":
					r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(time.Minute)))
				case "second":
					r = NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(time.Second)))
				}
				if r != nil {
					r.(ValueObject).SetMember("__as", as, c)
					return r
				}
			}
			c.RaiseRuntimeError("Time object cannot get __next__ without specialized time type")
			return nil
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
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
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
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
			switch diff := args[0].(type) {
			case ValueObject:
				if diff.Type().TypeId != timeDurationClass.TypeId {
					c.RaiseRuntimeError("invalid duration class")
				}
				r := NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(diff.GetMember("__du", c).ToGoValue().(time.Duration))))
				return r
			case ValueInt:
				if as, ok := this.GetMember("__as", c).(ValueStr); ok {
					var r Value
					d := diff.AsInt()
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
					c.RaiseRuntimeError("Time object cannot get __next__ without specialized time type")
				}
			}
			return nil
		}).
		Method("__sub__", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) != 1 {
				c.RaiseRuntimeError("__sub__ requires one arugment!")
			}
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
			switch diff := args[0].(type) {
			case ValueObject:
				switch diff.Type().TypeId {
				case timeTimeClass.TypeId:
					t2 := diff.GetMember("__t", c).ToGoValue().(time.Time)
					r := NewObjectAndInit(timeDurationClass, c, NewGoValue(t.Sub(t2)))
					return r
				case timeDurationClass.TypeId:
					du := diff.GetMember("__du", c).ToGoValue().(time.Duration)
					r := NewObjectAndInit(timeTimeClass, c, NewGoValue(t.Add(-du)))
					return r
				default:
					c.RaiseRuntimeError("invalid duration class")
				}
			case ValueInt:
				if as, ok := this.GetMember("__as", c).(ValueStr); ok {
					var r Value
					t := this.GetMember("__t", c).ToGoValue().(time.Time)
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
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
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
			t := this.GetMember("__t", c).ToGoValue().(time.Time)
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
			t1 := this.GetMember("__t", c).ToGoValue().(time.Time)
			t2 := args[0].GetMember("__t", c).ToGoValue().(time.Time)
			return NewBool(t1.Before(t2))
		}).
		Method("__le__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.GetMember("__t", c).ToGoValue().(time.Time)
			t2 := args[0].GetMember("__t", c).ToGoValue().(time.Time)
			return NewBool(!t1.After(t2))
		}).
		Method("__gt__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.GetMember("__t", c).ToGoValue().(time.Time)
			t2 := args[0].GetMember("__t", c).ToGoValue().(time.Time)
			return NewBool(t1.After(t2))
		}).
		Method("__ge__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.GetMember("__t", c).ToGoValue().(time.Time)
			t2 := args[0].GetMember("__t", c).ToGoValue().(time.Time)
			return NewBool(!t1.Before(t2))
		}).
		Method("__eq__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.GetMember("__t", c).ToGoValue().(time.Time)
			t2 := args[0].GetMember("__t", c).ToGoValue().(time.Time)
			return NewBool(!t1.Before(t2) && !t1.After(t2))
		}).
		Method("__ne__", func(c *Context, this ValueObject, args []Value) Value {
			t1 := this.GetMember("__t", c).ToGoValue().(time.Time)
			t2 := args[0].GetMember("__t", c).ToGoValue().(time.Time)
			return NewBool(t1.Before(t2) || t1.After(t2))
		}).
		Method("toGoTime", func(c *Context, this ValueObject, args []Value) Value {
			return this.GetMember("__t", c)
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
			d2 = d2a.GetMember("__du", c).ToGoValue().(time.Duration)
		}
		return d2
	}
	compareDurations := func(c *Context, this ValueObject, args []Value) (gt bool, lt bool) {
		d1 := this.GetMember("__du", c).ToGoValue().(time.Duration)
		d2 := getOther(c, args)
		return d1 < d2, d1 > d2
	}
	timeDurationClass = NewClassBuilder("Duration").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			switch len(args) {
			case 1:
				switch dv := args[0].(type) {
				case ValueInt:
					du := time.Duration(dv.Value()) * time.Second
					this.SetMember("__du", NewGoValue(du), c)
					return
				case ValueFloat:
					du := time.Duration(dv.Value() * float64(time.Second))
					this.SetMember("__du", NewGoValue(du), c)
					return
				case ValueStr:
					du, err := time.ParseDuration(dv.Value())
					if err != nil {
						c.RaiseRuntimeError("invalid duration string %s", dv.Value())
					}
					this.SetMember("__du", NewGoValue(du), c)
					return
				case GoValue:
					if du, is := dv.ToGoValue().(time.Duration); is {
						this.SetMember("__du", NewGoValue(du), c)
						return
					}
				}
			}
			c.RaiseRuntimeError("Duration.__init__: invalid duration argumenet")
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			du := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewStr(du.String())
		}).
		// To floats
		Methods([]string{"nanoseconds", "ns"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewFloat(float64(d.Nanoseconds()))
		}).
		Methods([]string{"milliseconds", "ms"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewFloat(float64(d.Milliseconds()))
		}).
		Methods([]string{"seconds", "s"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewFloat(d.Seconds())
		}).
		Methods([]string{"minutes", "m"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewFloat(d.Minutes())
		}).
		Methods([]string{"hours", "h"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewFloat(d.Hours())
		}).
		Methods([]string{"days", "d"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewFloat(d.Hours() / 24)
		}).
		Methods([]string{"weeks", "w"}, func(c *Context, this ValueObject, args []Value) Value {
			d := this.GetMember("__du", c).ToGoValue().(time.Duration)
			return NewFloat(d.Hours() / 24 / 7)
		}).
		// Comparation
		Method("__lt__", func(c *Context, this ValueObject, args []Value) Value {
			lt, _ := compareDurations(c, this, args)
			return NewBool(lt)
		}).
		Method("__le__", func(c *Context, this ValueObject, args []Value) Value {
			_, gt := compareDurations(c, this, args)
			return NewBool(!gt)
		}).
		Method("__gt__", func(c *Context, this ValueObject, args []Value) Value {
			_, gt := compareDurations(c, this, args)
			return NewBool(gt)
		}).
		Method("__ge__", func(c *Context, this ValueObject, args []Value) Value {
			lt, _ := compareDurations(c, this, args)
			return NewBool(!lt)
		}).
		Method("__eq__", func(c *Context, this ValueObject, args []Value) Value {
			lt, gt := compareDurations(c, this, args)
			return NewBool(!(lt || gt))
		}).
		Method("__ne__", func(c *Context, this ValueObject, args []Value) Value {
			lt, gt := compareDurations(c, this, args)
			return NewBool(lt || gt)
		}).
		// Add & sub
		Method("__add__", func(c *Context, this ValueObject, args []Value) Value {
			d1 := this.GetMember("__du", c).ToGoValue().(time.Duration)
			d2 := getOther(c, args)
			du := time.Duration(int64(d1) + int64(d2))
			return NewObjectAndInit(timeDurationClass, c, NewGoValue(du))
		}).
		Method("__sub__", func(c *Context, this ValueObject, args []Value) Value {
			d1 := this.GetMember("__du", c).ToGoValue().(time.Duration)
			d2 := getOther(c, args)
			du := time.Duration(int64(d1) - int64(d2))
			return NewObjectAndInit(timeDurationClass, c, NewGoValue(du))
		}).
		Build()
}

func init() {
	timeInittimeTimeClass()
	timeInitDurationClass()
}
