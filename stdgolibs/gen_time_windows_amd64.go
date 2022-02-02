package stdgolibs

import (
	pkg "time"

	"reflect"
)

func init() {
	registerValues("time", map[string]reflect.Value{
		// Functions
		"LoadLocationFromTZData": reflect.ValueOf(pkg.LoadLocationFromTZData),
		"Parse":                  reflect.ValueOf(pkg.Parse),
		"ParseInLocation":        reflect.ValueOf(pkg.ParseInLocation),
		"ParseDuration":          reflect.ValueOf(pkg.ParseDuration),
		"Since":                  reflect.ValueOf(pkg.Since),
		"Until":                  reflect.ValueOf(pkg.Until),
		"Now":                    reflect.ValueOf(pkg.Now),
		"Unix":                   reflect.ValueOf(pkg.Unix),
		"Date":                   reflect.ValueOf(pkg.Date),
		"FixedZone":              reflect.ValueOf(pkg.FixedZone),
		"LoadLocation":           reflect.ValueOf(pkg.LoadLocation),
		"Sleep":                  reflect.ValueOf(pkg.Sleep),
		"NewTimer":               reflect.ValueOf(pkg.NewTimer),
		"After":                  reflect.ValueOf(pkg.After),
		"AfterFunc":              reflect.ValueOf(pkg.AfterFunc),
		"NewTicker":              reflect.ValueOf(pkg.NewTicker),
		"Tick":                   reflect.ValueOf(pkg.Tick),

		// Consts

		"ANSIC":       reflect.ValueOf(pkg.ANSIC),
		"UnixDate":    reflect.ValueOf(pkg.UnixDate),
		"RubyDate":    reflect.ValueOf(pkg.RubyDate),
		"RFC822":      reflect.ValueOf(pkg.RFC822),
		"RFC822Z":     reflect.ValueOf(pkg.RFC822Z),
		"RFC850":      reflect.ValueOf(pkg.RFC850),
		"RFC1123":     reflect.ValueOf(pkg.RFC1123),
		"RFC1123Z":    reflect.ValueOf(pkg.RFC1123Z),
		"RFC3339":     reflect.ValueOf(pkg.RFC3339),
		"RFC3339Nano": reflect.ValueOf(pkg.RFC3339Nano),
		"Kitchen":     reflect.ValueOf(pkg.Kitchen),
		"Stamp":       reflect.ValueOf(pkg.Stamp),
		"StampMilli":  reflect.ValueOf(pkg.StampMilli),
		"StampMicro":  reflect.ValueOf(pkg.StampMicro),
		"StampNano":   reflect.ValueOf(pkg.StampNano),
		"January":     reflect.ValueOf(pkg.January),
		"February":    reflect.ValueOf(pkg.February),
		"March":       reflect.ValueOf(pkg.March),
		"April":       reflect.ValueOf(pkg.April),
		"May":         reflect.ValueOf(pkg.May),
		"June":        reflect.ValueOf(pkg.June),
		"July":        reflect.ValueOf(pkg.July),
		"August":      reflect.ValueOf(pkg.August),
		"September":   reflect.ValueOf(pkg.September),
		"October":     reflect.ValueOf(pkg.October),
		"November":    reflect.ValueOf(pkg.November),
		"December":    reflect.ValueOf(pkg.December),
		"Sunday":      reflect.ValueOf(pkg.Sunday),
		"Monday":      reflect.ValueOf(pkg.Monday),
		"Tuesday":     reflect.ValueOf(pkg.Tuesday),
		"Wednesday":   reflect.ValueOf(pkg.Wednesday),
		"Thursday":    reflect.ValueOf(pkg.Thursday),
		"Friday":      reflect.ValueOf(pkg.Friday),
		"Saturday":    reflect.ValueOf(pkg.Saturday),
		"Nanosecond":  reflect.ValueOf(pkg.Nanosecond),
		"Microsecond": reflect.ValueOf(pkg.Microsecond),
		"Millisecond": reflect.ValueOf(pkg.Millisecond),
		"Second":      reflect.ValueOf(pkg.Second),
		"Minute":      reflect.ValueOf(pkg.Minute),
		"Hour":        reflect.ValueOf(pkg.Hour),

		// Variables

		"UTC":   reflect.ValueOf(&pkg.UTC),
		"Local": reflect.ValueOf(&pkg.Local),
	})
	registerTypes("time", map[string]reflect.Type{
		// Non interfaces

		"ParseError": reflect.TypeOf((*pkg.ParseError)(nil)).Elem(),
		"Time":       reflect.TypeOf((*pkg.Time)(nil)).Elem(),
		"Month":      reflect.TypeOf((*pkg.Month)(nil)).Elem(),
		"Weekday":    reflect.TypeOf((*pkg.Weekday)(nil)).Elem(),
		"Duration":   reflect.TypeOf((*pkg.Duration)(nil)).Elem(),
		"Location":   reflect.TypeOf((*pkg.Location)(nil)).Elem(),
		"Timer":      reflect.TypeOf((*pkg.Timer)(nil)).Elem(),
		"Ticker":     reflect.TypeOf((*pkg.Ticker)(nil)).Elem(),
	})
}
