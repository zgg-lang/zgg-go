package builtin_libs

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/samber/lo"
	. "github.com/zgg-lang/zgg-go/runtime"

	redis "github.com/redis/go-redis/v9"
)

var (
	redisClientClass      ValueType
	redisPipeSessionClass ValueType
)

func libRedis(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("open", NewNativeFunction("open", func(c *Context, this Value, args []Value) Value {
		var (
			redisAddr  ValueStr
			redisAddrs ValueArray
			addrBy     int
		)
		EnsureFuncParams(c, "redis.open", args,
			ArgRuleOneOf(
				"addr",
				[]ValueType{TypeStr, TypeArray},
				[]any{&redisAddr, &redisAddrs},
				&addrBy,
				&redisAddr,
				NewStr("127.0.0.1:6379"),
			),
		)
		opts := &redis.UniversalOptions{}
		for i := 1; i < len(args); i++ {
			switch v := args[i].(type) {
			case ValueInt:
				opts.DB = v.AsInt()
			case ValueObject:
				if val, ok := v.GetMember("username", c).(ValueStr); ok {
					opts.Username = val.Value()
				}
				if val, ok := v.GetMember("password", c).(ValueStr); ok {
					opts.Password = val.Value()
				}
				if val, ok := v.GetMember("database", c).(ValueInt); ok {
					opts.DB = val.AsInt()
				}
				switch val := v.GetMember("readTimeout", c).(type) {
				case ValueFloat:
					opts.ReadTimeout = time.Duration(val.Value()) * time.Second
				case ValueInt:
					opts.ReadTimeout = time.Duration(val.Value()) * time.Second
				}
				switch val := v.GetMember("connTimeout", c).(type) {
				case ValueFloat:
					opts.DialTimeout = time.Duration(val.Value()) * time.Second
				case ValueInt:
					opts.DialTimeout = time.Duration(val.Value()) * time.Second
				}
				switch val := v.GetMember("writeTimeout", c).(type) {
				case ValueFloat:
					opts.WriteTimeout = time.Duration(val.Value()) * time.Second
				case ValueInt:
					opts.WriteTimeout = time.Duration(val.Value()) * time.Second
				}
			}
		}
		var client redis.Cmdable
		if addrBy == 1 {
			opts.Addrs = lo.Map(*redisAddrs.Values, func(v Value, _ int) string {
				return v.ToString(c)
			})
			client = redis.NewClusterClient(opts.Cluster())
		} else {
			opts.Addrs = []string{redisAddr.Value()}
			client = redis.NewClient(opts.Simple())
		}
		return NewObjectAndInit(redisClientClass, c, NewGoValue(client))
	}), nil)
	lib.SetMember("RedisClient", redisClientClass, nil)
	lib.SetMember("RedisPipeSession", redisPipeSessionClass, nil)
	return lib
}

func initRedisPipeSessionClass() ValueType {
	rv := NewClassBuilder("RedisPipeSession").
		Constructor(func(c *Context, thisObj ValueObject, args []Value) {
			thisObj.SetMember("_cmds", NewGoValue([][]any{}), c)
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			cmd := c.MustStr(args[0])
			if cmd == strings.ToUpper(cmd) {
				cmds := this.GetMember("_cmds", c).ToGoValue().([][]any)
				session := this
				return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
					defer func() {
						session.SetMember("_cmds", NewGoValue(cmds), c)
					}()
					cmdArgs := make([]any, 1+len(args))
					cmdArgs[0] = cmd
					for i, arg := range args {
						cmdArgs[i+1] = arg.ToGoValue()
					}
					cmds = append(cmds, cmdArgs)
					return Undefined()
				})
			} else {
				return Undefined()
			}
		}).
		Build()
	return rv
}

func initRedisClientClass() ValueType {
	rv := NewClassBuilder("RedisClient").
		Constructor(func(c *Context, thisObj ValueObject, args []Value) {
			thisObj.SetMember("_conn", args[0], c)
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			conn := this.GetMember("_conn", c).ToGoValue().(redis.UniversalClient)
			if err := conn.Close(); err != nil {
				c.RaiseRuntimeError("RedisClient.close fail on close: %s", err)
			}
			return Undefined()
		}).
		Method("exec", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) < 1 {
				c.RaiseRuntimeError("RedisClient.exec requires at least 1 argument")
				return nil
			}
			var (
				ctx  = context.Background()
				conn = this.GetMember("_conn", c).ToGoValue().(redis.Cmdable)
			)
			if cmds, ok := args[0].(ValueArray); ok {
				var (
					pipeline = conn.Pipeline()
					n        = cmds.Len()
					resCmds  = make([]*redis.Cmd, n)
				)
				for i := 0; i < n; i++ {
					var (
						cmdArgs = []interface{}{}
					)
					cmd := cmds.GetIndex(i, c)
					if cmdItems, ok := cmd.(ValueArray); ok {
						if n2 := cmdItems.Len(); n2 > 0 {
							cmdArgs = append(make([]any, 0, n2), cmdItems.GetIndex(0, c).ToString(c))
							for i := 1; i < n2; i++ {
								cmdArgs = append(cmdArgs, cmdItems.GetIndex(i, c).ToGoValue())
							}
						} else {
							cmdArgs = []any{"PING"}
						}
					} else {
						cmdArgs = []any{cmd.ToString(c)}
					}
					resCmds[i] = pipeline.Do(ctx, cmdArgs...)
				}
				if _, err := pipeline.Exec(ctx); err != nil {
					c.RaiseRuntimeError("redis.exec: exec piped command error: %s", err)
				}
				rv := NewArray(len(resCmds))
				for _, rc := range resCmds {
					if r, err := rc.Result(); err != nil {
						c.RaiseRuntimeError("redis.exec: receive piped command error: %s", err)
					} else if r == nil {
						rv.PushBack(Nil())
					} else {
						v := reflect.ValueOf(r)
						rv.PushBack(FromGoValue(v, c))
					}
				}
				return rv
			} else {
				return redisDo(c, ctx, conn, lo.Map(args, func(a Value, _ int) any { return a.ToGoValue() })...)
			}
		}).
		Method("pipe", func(c *Context, this ValueObject, args []Value) Value {
			action := c.MustCallable(args[0])
			session := NewObjectAndInit(redisPipeSessionClass, c)
			c.Invoke(action, nil, Args(session))
			var (
				ctx  = context.Background()
				cmds = session.GetMember("_cmds", c).ToGoValue().([][]any)
				conn = this.GetMember("_conn", c).ToGoValue().(redis.Cmdable)
				pipe = conn.Pipeline()
				rets = make([]*redis.Cmd, len(cmds))
			)
			for i, cmd := range cmds {
				rets[i] = pipe.Do(ctx, cmd...)
			}
			pipe.Exec(ctx)
			rv := NewArray(len(cmds))
			for _, cmd := range rets {
				if rsp, err := cmd.Result(); err != nil {
					if redisIsNil(err) {
						rv.PushBack(Nil())
					} else {
						c.RaiseRuntimeError("redis.exec: receive piped command error: %s", err)
					}
				} else if rsp == nil {
					rv.PushBack(Nil())
				} else {
					v := reflect.ValueOf(rsp)
					rv.PushBack(FromGoValue(v, c))
				}
			}
			return rv
		}).
		Method("eachShared", func(c *Context, this ValueObject, args []Value) Value {
			conn, iscc := this.GetMember("_conn", c).ToGoValue().(*redis.ClusterClient)
			if !iscc {
				c.RaiseRuntimeError("not a cluster client")
			}
			conn.ForEachShard(func(c *redis.Client) {
				sc := NewObjectAndInit(redisClientClass, c, NewGoValue(c))
			})
			return rv
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			cmd := c.MustStr(args[0])
			if cmd == strings.ToUpper(cmd) {
				conn := this.GetMember("_conn", c).ToGoValue().(redis.Cmdable)
				return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
					cmdArgs := make([]interface{}, 1+len(args))
					cmdArgs[0] = cmd
					for i, arg := range args {
						cmdArgs[i+1] = arg.ToGoValue()
					}
					return redisDo(c, context.Background(), conn, cmdArgs...)
				})
			} else {
				return Undefined()
			}
		}).
		Build()
	return rv
}

func redisIsNil(err error) bool {
	return err != nil && err.Error() == "redis: nil"
}

func redisDo(c *Context, ctx context.Context, conn redis.Cmdable, cmdArgs ...any) Value {
	var (
		result any
		err    error
	)
	if uc, is := conn.(redis.UniversalClient); is {
		result, err = uc.Do(ctx, cmdArgs...).Result()
	} else {
		var (
			pipe = conn.Pipeline()
			cmd  = pipe.Do(ctx, cmdArgs...)
		)
		if _, err = pipe.Exec(ctx); err == nil {
			result, err = cmd.Result()
		}
	}
	if err != nil {
		if redisIsNil(err) {
			return Nil()
		}
		c.RaiseRuntimeError("RedisClient.exec pipe error %+v", reflect.TypeOf(err))
		return nil
	} else if result == nil {
		return Nil()
	} else {
		v := reflect.ValueOf(result)
		return FromGoValue(v, c)
	}
}

func init() {
	redisClientClass = initRedisClientClass()
	redisPipeSessionClass = initRedisPipeSessionClass()
}
