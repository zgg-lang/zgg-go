package builtin_libs

import (
	"sync"

	. "github.com/zgg-lang/zgg-go/runtime"
)

type LibInfo struct {
	name   string
	getter func(*Context) ValueObject
}

var (
	findLock  sync.Mutex
	libs      = map[string]ValueObject{}
	StdLibMap = map[string]LibInfo{
		"base64":     LibInfo{name: "base64", getter: libBase64},
		"concurrent": LibInfo{name: "concurrent", getter: libConcurrent},
		"cron":       LibInfo{name: "cron", getter: libCron},
		"db":         LibInfo{name: "db", getter: libDb},
		"dbop":       LibInfo{name: "db/op", getter: libDbOp},
		"db/op":      LibInfo{name: "db/op", getter: libDbOp},
		"dom":        LibInfo{name: "dom", getter: libDom},
		"drawing":    LibInfo{name: "drawing", getter: libDrawing},
		"etree":      LibInfo{name: "etree", getter: libEtree},
		"file":       LibInfo{name: "file", getter: libFile},
		"go":         LibInfo{name: "go", getter: libGo},
		"graph":      LibInfo{name: "graph", getter: libGraph},
		"http":       LibInfo{name: "http", getter: libHttp},
		"json":       LibInfo{name: "json", getter: libJson},
		"kv":         LibInfo{name: "kv", getter: libKv},
		"msgpack":    LibInfo{name: "msgpack", getter: libMsgpack},
		"math":       LibInfo{name: "math", getter: libMath},
		"nsq":        LibInfo{name: "nsq", getter: libNsq},
		"psutil":     LibInfo{name: "psutil", getter: libPsutil},
		"ptable":     LibInfo{name: "ptable", getter: libPtable},
		"random":     LibInfo{name: "random", getter: libRandom},
		"redis":      LibInfo{name: "redis", getter: libRedis},
		"regex":      LibInfo{name: "regex", getter: libRegex},
		"regex2":     LibInfo{name: "regex2", getter: libRegex2},
		"sh":         LibInfo{name: "sh", getter: libSh},
		"sys":        LibInfo{name: "sys", getter: libSys},
		"template":   LibInfo{name: "template", getter: libTemplate},
		"time":       LibInfo{name: "time", getter: libTime},
		"url":        LibInfo{name: "url", getter: libUrl},
		"yaml":       LibInfo{name: "yaml", getter: libYaml},
	}
)

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
	info, found := StdLibMap[name]
	if !found {
		return NewObject(), false
	}
	return getLib(c, name, info.getter), true
}
