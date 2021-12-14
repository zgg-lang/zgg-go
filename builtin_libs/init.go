package builtin_libs

import (
	"sync"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var findLock sync.Mutex
var libs = map[string]ValueObject{}

func getLib(c *Context, name string, getter func(*Context) ValueObject) ValueObject {
	findLock.Lock()
	defer findLock.Unlock()
	lib, found := libs[name]
	if found {
		return lib
	}
	lib = getter(c)
	libs[name] = lib
	return lib
}

func FindLib(c *Context, name string) (ValueObject, bool) {
	switch name {
	case "base64":
		return getLib(c, "base64", libBase64), true
	case "concurrent":
		return getLib(c, "concurrent", libConcurrent), true
	case "cron":
		return getLib(c, "cron", libCron), true
	case "db":
		return getLib(c, "db", libDb), true
	case "dbop":
		fallthrough
	case "db/op":
		return getLib(c, "db/op", libDbOp), true
	case "dom":
		return getLib(c, "dom", libDom), true
	case "file":
		return getLib(c, "file", libFile), true
	case "go":
		return getLib(c, "go", libGo), true
	case "http":
		return getLib(c, "http", libHttp), true
	case "json":
		return getLib(c, "json", libJson), true
	case "kv":
		return getLib(c, "kv", libKv), true
	case "msgpack":
		return getLib(c, "msgpack", libMsgpack), true
	case "nsq":
		return getLib(c, "nsq", libNsq), true
	case "ptable":
		return getLib(c, "ptable", libPtable), true
	case "random":
		return getLib(c, "random", libRandom), true
	case "redis":
		return getLib(c, "redis", libRedis), true
	case "regex":
		return getLib(c, "regex", libRegex), true
	case "regex2":
		return getLib(c, "regex2", libRegex2), true
	case "sh":
		return getLib(c, "sh", libSh), true
	case "sys":
		return getLib(c, "sys", libSys), true
	case "template":
		return getLib(c, "template", libTemplate), true
	case "time":
		return getLib(c, "time", libTime), true
	case "url":
		return getLib(c, "url", libUrl), true
	default:
		return NewObject(), false
	}
}
