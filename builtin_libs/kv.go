package builtin_libs

import (
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/zgg-lang/zgg-go/runtime"

	"github.com/gomodule/redigo/redis"
)

var (
	kvAdapterSchemeMap = map[string]ValueType{}
	kvManagerClass     = getKvManagerClass()
)

func libKv(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Manager", kvManagerClass, c)
	initKvFileSystemAdapter()
	initKvRedisAdapter()
	initKvRedisHashAdapter()
	return lib
}

func getKvManagerClass() ValueType {
	return NewClassBuilder("Manager").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var adapter Value
			EnsureFuncParams(c, "Manager.__init__", args,
				ArgRuleRequired{"scheme", TypeAny, &adapter},
			)
			if adapterScheme, ok := adapter.(ValueStr); ok {
				u, err := url.Parse(adapterScheme.Value())
				if err != nil {
					c.OnRuntimeError("parse adapter scheme error %s", err)
				}
				adapterClass, found := kvAdapterSchemeMap[u.Scheme]
				if !found {
					c.OnRuntimeError("unknown adapter scheme %s", u.Scheme)
				}
				adapter = NewObjectAndInit(adapterClass, c, NewGoValue(u))
			}
			this.SetMember("__adapter", adapter, c)
			if onOpen := adapter.GetMember("onOpen", c); c.IsCallable(onOpen) {
				c.Invoke(onOpen, this, NoArgs)
			}
		}).
		Method("get", func(c *Context, this ValueObject, args []Value) Value {
			adapter := c.MustObject(this.GetMember("__adapter", c))
			c.InvokeMethod(adapter, "get", Args(args...))
			return c.RetVal
		}).
		Method("set", func(c *Context, this ValueObject, args []Value) Value {
			adapter := c.MustObject(this.GetMember("__adapter", c))
			c.InvokeMethod(adapter, "set", Args(args...))
			return this
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			adapter := c.MustObject(this.GetMember("__adapter", c))
			if close := adapter.GetMember("close", c); c.IsCallable(close) {
				c.Invoke(close, adapter, NoArgs)
			}
			return this
		}).
		Build()
}

func initKvFileSystemAdapter() {
	t := NewClassBuilder("FileSystemAdapter").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			u := args[0].ToGoValue().(*url.URL)
			q := u.Query()
			this.SetMember("__root", NewStr(u.Path), c)
			this.SetMember("__prefix", NewStr(q.Get("prefix")), c)
		}).
		Method("get", func(c *Context, this ValueObject, args []Value) Value {
			var (
				root   = c.MustStr(this.GetMember("__root", c))
				prefix = c.MustStr(this.GetMember("__prefix", c))
				key    ValueStr
			)
			EnsureFuncParams(c, "FileSystemAdapter.get", args, ArgRuleRequired{"key", TypeStr, &key})
			filename := filepath.Join(root, prefix+key.Value())
			bs, err := ioutil.ReadFile(filename)
			if err != nil {
				if os.IsNotExist(err) {
					return Nil()
				}
				c.OnRuntimeError("FileSystemAdapter.get: read file %s error %s", filename, err)
			}
			return NewBytes(bs)
		}).
		Method("set", func(c *Context, this ValueObject, args []Value) Value {
			var (
				root   = c.MustStr(this.GetMember("__root", c))
				prefix = c.MustStr(this.GetMember("__prefix", c))
				key    ValueStr
				data   Value
			)
			EnsureFuncParams(c, "FileSystemAdapter.set", args,
				ArgRuleRequired{"key", TypeStr, &key},
				ArgRuleRequired{"value", TypeAny, &data},
			)
			filename := filepath.Join(root, prefix+key.Value())
			var bs []byte
			if b, ok := data.(ValueBytes); ok {
				bs = b.Value()
			} else {
				bs = []byte(data.ToString(c))
			}
			if err := ioutil.WriteFile(filename, bs, 0600); err != nil {
				c.OnRuntimeError("FileSystemAdapter.set: write file %s error %s", filename, err)
			}
			return Undefined()
		}).
		Build()
	kvAdapterSchemeMap["fs"] = t
}

func initKvRedisAdapter() {
	t := NewClassBuilder("RedisAdapter").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			u := args[0].ToGoValue().(*url.URL)
			q := u.Query()
			var host string
			if h, p, err := net.SplitHostPort(u.Host); err != nil {
				c.OnRuntimeError("RedisAdapter.__init__: invalid url %s: %s", u, err)
			} else {
				if p == "" {
					p = "6379"
				}
				if h == "" {
					h = "localhost"
				}
				host = h + ":" + p
			}
			opts := []redis.DialOption{}
			if db, err := strconv.Atoi(q.Get("db")); err == nil && db != 0 {
				opts = append(opts, redis.DialDatabase(db))
			}
			if ui := u.User; ui != nil {
				opts = append(opts, redis.DialUsername(ui.Username()))
				if p, has := ui.Password(); has {
					opts = append(opts, redis.DialPassword(p))
				}
			}
			maxIdle := 10
			if mi, err := strconv.Atoi(q.Get("maxIdle")); err == nil {
				maxIdle = mi
			}
			pool := redis.NewPool(func() (redis.Conn, error) {
				return redis.Dial("tcp", host, opts...)
			}, maxIdle)
			this.SetMember("__pool", NewGoValue(pool), c)
			this.SetMember("__prefix", NewStr(q.Get("prefix")), c)
			if ttl, err := strconv.Atoi(q.Get("ttl")); err == nil && ttl > 0 {
				this.SetMember("__ttl", NewInt(int64(ttl)), c)
			} else {
				this.SetMember("__ttl", NewInt(0), c)
			}
		}).
		Method("get", func(c *Context, this ValueObject, args []Value) Value {
			var (
				pool   = this.GetMember("__pool", c).ToGoValue().(*redis.Pool)
				prefix = c.MustStr(this.GetMember("__prefix", c))
				key    ValueStr
			)
			EnsureFuncParams(c, "RedisAdapter.get", args, ArgRuleRequired{"key", TypeStr, &key})
			conn := pool.Get()
			defer conn.Close()
			bs, err := redis.Bytes(conn.Do("GET", prefix+key.Value()))
			if err != nil {
				if err == redis.ErrNil {
					return Nil()
				}
				c.OnRuntimeError("RedisAdapter.get: %s", err)
			}
			return NewBytes(bs)
		}).
		Method("set", func(c *Context, this ValueObject, args []Value) Value {
			var (
				pool   = this.GetMember("__pool", c).ToGoValue().(*redis.Pool)
				prefix = c.MustStr(this.GetMember("__prefix", c))
				ttl    = c.MustInt(this.GetMember("__ttl", c))
				key    ValueStr
				data   Value
			)
			EnsureFuncParams(c, "RedisAdapter.set", args,
				ArgRuleRequired{"key", TypeStr, &key},
				ArgRuleRequired{"value", TypeAny, &data},
			)
			var bs []byte
			if b, ok := data.(ValueBytes); ok {
				bs = b.Value()
			} else {
				bs = []byte(data.ToString(c))
			}
			conn := pool.Get()
			defer conn.Close()
			redisKey := prefix + key.Value()
			var err error
			if ttl > 0 {
				_, err = conn.Do("SET", redisKey, bs, "EX", ttl)
			} else {
				_, err = conn.Do("SET", redisKey, bs)
			}
			if err != nil {
				c.OnRuntimeError("RedisAdapter.set: %s", err)
			}
			return Undefined()
		}).
		Build()
	kvAdapterSchemeMap["redis"] = t
}

func initKvRedisHashAdapter() {
	t := NewClassBuilder("RedisHashAdapter").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			u := args[0].ToGoValue().(*url.URL)
			q := u.Query()
			hashkey := u.Path
			if strings.HasPrefix(hashkey, "/") {
				hashkey = hashkey[1:]
			}
			if hashkey == "" {
				c.OnRuntimeError("RedisHashAdapter.__init__: hashkey cannot be empty")
			}
			var host string
			if h, p, err := net.SplitHostPort(u.Host); err != nil {
				c.OnRuntimeError("RedisHashAdapter.__init__: invalid url %s: %s", u, err)
			} else {
				if p == "" {
					p = "6379"
				}
				if h == "" {
					h = "localhost"
				}
				host = h + ":" + p
			}
			opts := []redis.DialOption{}
			if db, err := strconv.Atoi(q.Get("db")); err == nil && db != 0 {
				opts = append(opts, redis.DialDatabase(db))
			}
			if ui := u.User; ui != nil {
				opts = append(opts, redis.DialUsername(ui.Username()))
				if p, has := ui.Password(); has {
					opts = append(opts, redis.DialPassword(p))
				}
			}
			maxIdle := 10
			if mi, err := strconv.Atoi(q.Get("maxIdle")); err == nil {
				maxIdle = mi
			}
			pool := redis.NewPool(func() (redis.Conn, error) {
				return redis.Dial("tcp", host, opts...)
			}, maxIdle)
			this.SetMember("__pool", NewGoValue(pool), c)
			this.SetMember("__hashkey", NewStr(hashkey), c)
			this.SetMember("__prefix", NewStr(q.Get("prefix")), c)
		}).
		Method("get", func(c *Context, this ValueObject, args []Value) Value {
			var (
				pool    = this.GetMember("__pool", c).ToGoValue().(*redis.Pool)
				hashkey = this.GetMember("__hashkey", c)
				prefix  = c.MustStr(this.GetMember("__prefix", c))
				key     ValueStr
			)
			EnsureFuncParams(c, "RedisHashAdapter.get", args, ArgRuleRequired{"key", TypeStr, &key})
			conn := pool.Get()
			defer conn.Close()
			bs, err := redis.Bytes(conn.Do("HGET", hashkey.ToString(c), prefix+key.Value()))
			if err != nil {
				if err == redis.ErrNil {
					return Nil()
				}
				c.OnRuntimeError("RedisHashAdapter.get: %s", err)
			}
			return NewBytes(bs)
		}).
		Method("set", func(c *Context, this ValueObject, args []Value) Value {
			var (
				pool    = this.GetMember("__pool", c).ToGoValue().(*redis.Pool)
				hashkey = this.GetMember("__hashkey", c)
				prefix  = c.MustStr(this.GetMember("__prefix", c))
				key     ValueStr
				data    Value
			)
			EnsureFuncParams(c, "RedisHashAdapter.set", args,
				ArgRuleRequired{"key", TypeStr, &key},
				ArgRuleRequired{"value", TypeAny, &data},
			)
			var bs []byte
			if b, ok := data.(ValueBytes); ok {
				bs = b.Value()
			} else {
				bs = []byte(data.ToString(c))
			}
			conn := pool.Get()
			defer conn.Close()
			if _, err := conn.Do("HSET", hashkey.ToString(c), prefix+key.Value(), bs); err != nil {
				c.OnRuntimeError("RedisHashAdapter.set: %s", err)
			}
			return Undefined()
		}).
		Build()
	kvAdapterSchemeMap["redishash"] = t
}
