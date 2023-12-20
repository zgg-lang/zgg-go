package builtin_libs

import (
	"os"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	ip2regionSearcherClass ValueType
)

func libIp2region(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Searcher", ip2regionSearcherClass, nil)
	return lib
}

type ip2regionSearcherReserved struct {
	searcher *xdb.Searcher
}

func ip2regionInitSearcherClass() {
	ip2regionSearcherClass = NewClassBuilder("Searcher").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				dbfile ValueStr
			)
			EnsureFuncParams(c, "Searcher.__init__", args,
				ArgRuleRequired("dbfile", TypeStr, &dbfile),
			)
			content, err := os.ReadFile(dbfile.Value())
			if err != nil {
				c.RaiseRuntimeError("open ip2region xdb path %s error %+v", dbfile.Value(), err)
			}
			searcher, err := xdb.NewWithBuffer(content)
			if err != nil {
				c.RaiseRuntimeError("open ip2region xdb path %s error %+v", dbfile.Value(), err)
			}
			this.Reserved = &ip2regionSearcherReserved{
				searcher: searcher,
			}
		}).
		Method("search", func(c *Context, this ValueObject, args []Value) Value {
			var (
				ipStr ValueStr
				ipInt ValueInt
				by    int
				r     string
				err   error
				s     = this.Reserved.(*ip2regionSearcherReserved).searcher
			)
			EnsureFuncParams(c, "Searcher.search", args,
				ArgRuleOneOf("ip",
					[]ValueType{TypeInt, TypeStr},
					[]any{&ipInt, &ipStr},
					&by, nil, nil))
			switch by {
			case 0:
				r, err = s.Search(uint32(ipInt.Value()))
			case 1:
				r, err = s.SearchByStr(ipStr.Value())
			default:
				c.RaiseRuntimeError("invalid ip type")
			}
			if err != nil {
				c.RaiseRuntimeError("search error %+v", err)
			}
			//中国|0|广东省|广州市|电信
			fields := strings.Split(r, "|")
			rv := NewObject()
			switch len(fields) {
			case 5:
				rv.SetMember("carrier", NewStr(fields[4]), c)
				fallthrough
			case 4:
				rv.SetMember("city", NewStr(fields[3]), c)
				fallthrough
			case 3:
				rv.SetMember("province", NewStr(fields[2]), c)
				fallthrough
			case 2:
				rv.SetMember("region", NewStr(fields[1]), c)
				fallthrough
			case 1:
				rv.SetMember("country", NewStr(fields[0]), c)
			}
			return rv
		}).
		Build()
}

func init() {
	ip2regionInitSearcherClass()
}
