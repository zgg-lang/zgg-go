package builtin_libs

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	. "github.com/zgg-lang/zgg-go/runtime"
)

func libPsutil(*Context) ValueObject {
	rv := NewObject()
	rv.SetMember("cpu", libPsutilModCPU, nil)
	return rv
}

var libPsutilModCPU = func() ValueObject {
	mod := NewObject()
	mod.SetMember("percent", NewNativeFunction("cpu.percent", func(c *Context, _ Value, args []Value) Value {
		var (
			intervalStr      ValueStr
			intervalDuration ValueObject
			intervalBy       int
			perCPU           ValueBool
			interval         time.Duration
		)
		EnsureFuncParams(c, "psutil.cpu", args,
			ArgRuleOneOf(
				"interval",
				[]ValueType{TypeStr, timeDurationClass},
				[]any{&intervalStr, &intervalDuration},
				&intervalBy,
				&intervalStr,
				NewStr("1s"),
			),
			ArgRuleOptional("perCPU", TypeBool, &perCPU, NewBool(false)),
		)
		if intervalBy == 1 {
			interval = intervalDuration.GetMember("__du", c).ToGoValue(c).(time.Duration)
		} else {
			var err error
			interval, err = time.ParseDuration(intervalStr.Value())
			if err != nil {
				c.RaiseRuntimeError("invalid duration %s", intervalStr.Value())
			}
		}
		percents, err := cpu.Percent(interval, perCPU.Value())
		if err != nil {
			c.RaiseRuntimeError("get cpu percent error %v", err)
		}
		if perCPU.Value() {
			rv := NewArray(len(percents))
			for _, val := range percents {
				rv.PushBack(NewFloat(val))
			}
			return rv
		} else {
			return NewFloat(percents[0])
		}
	}, "interval", "perCPU"), nil)
	logicalCount, _ := cpu.Counts(true)
	physicalCount, _ := cpu.Counts(true)
	mod.SetMember("logicalCount", NewInt(int64(logicalCount)), nil)
	mod.SetMember("physicalCount", NewInt(int64(physicalCount)), nil)
	return mod
}()
