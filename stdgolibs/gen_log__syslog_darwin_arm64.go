package stdgolibs

import (
	pkg "log/syslog"

	"reflect"
)

func init() {
	registerValues("log/syslog", map[string]reflect.Value{
		// Functions
		"New":       reflect.ValueOf(pkg.New),
		"Dial":      reflect.ValueOf(pkg.Dial),
		"NewLogger": reflect.ValueOf(pkg.NewLogger),

		// Consts

		"LOG_EMERG":    reflect.ValueOf(pkg.LOG_EMERG),
		"LOG_ALERT":    reflect.ValueOf(pkg.LOG_ALERT),
		"LOG_CRIT":     reflect.ValueOf(pkg.LOG_CRIT),
		"LOG_ERR":      reflect.ValueOf(pkg.LOG_ERR),
		"LOG_WARNING":  reflect.ValueOf(pkg.LOG_WARNING),
		"LOG_NOTICE":   reflect.ValueOf(pkg.LOG_NOTICE),
		"LOG_INFO":     reflect.ValueOf(pkg.LOG_INFO),
		"LOG_DEBUG":    reflect.ValueOf(pkg.LOG_DEBUG),
		"LOG_KERN":     reflect.ValueOf(pkg.LOG_KERN),
		"LOG_USER":     reflect.ValueOf(pkg.LOG_USER),
		"LOG_MAIL":     reflect.ValueOf(pkg.LOG_MAIL),
		"LOG_DAEMON":   reflect.ValueOf(pkg.LOG_DAEMON),
		"LOG_AUTH":     reflect.ValueOf(pkg.LOG_AUTH),
		"LOG_SYSLOG":   reflect.ValueOf(pkg.LOG_SYSLOG),
		"LOG_LPR":      reflect.ValueOf(pkg.LOG_LPR),
		"LOG_NEWS":     reflect.ValueOf(pkg.LOG_NEWS),
		"LOG_UUCP":     reflect.ValueOf(pkg.LOG_UUCP),
		"LOG_CRON":     reflect.ValueOf(pkg.LOG_CRON),
		"LOG_AUTHPRIV": reflect.ValueOf(pkg.LOG_AUTHPRIV),
		"LOG_FTP":      reflect.ValueOf(pkg.LOG_FTP),
		"LOG_LOCAL0":   reflect.ValueOf(pkg.LOG_LOCAL0),
		"LOG_LOCAL1":   reflect.ValueOf(pkg.LOG_LOCAL1),
		"LOG_LOCAL2":   reflect.ValueOf(pkg.LOG_LOCAL2),
		"LOG_LOCAL3":   reflect.ValueOf(pkg.LOG_LOCAL3),
		"LOG_LOCAL4":   reflect.ValueOf(pkg.LOG_LOCAL4),
		"LOG_LOCAL5":   reflect.ValueOf(pkg.LOG_LOCAL5),
		"LOG_LOCAL6":   reflect.ValueOf(pkg.LOG_LOCAL6),
		"LOG_LOCAL7":   reflect.ValueOf(pkg.LOG_LOCAL7),

		// Variables

	})
	registerTypes("log/syslog", map[string]reflect.Type{
		// Non interfaces

		"Priority": reflect.TypeOf((*pkg.Priority)(nil)).Elem(),
		"Writer":   reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
	})
}
