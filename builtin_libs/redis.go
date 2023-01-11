package builtin_libs

import (
	"reflect"
	"strings"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"

	"github.com/gomodule/redigo/redis"
)

var (
	redisClientClass      ValueType
	redisPipeSessionClass ValueType
)

func libRedis(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("open", NewNativeFunction("open", func(c *Context, this Value, args []Value) Value {
		var redisAddr ValueStr
		EnsureFuncParams(c, "redis.open", args,
			ArgRuleOptional("addr", TypeStr, &redisAddr, NewStr("127.0.0.1:6379")),
		)
		opts := make([]redis.DialOption, 0)
		for i := 1; i < len(args); i++ {
			switch v := args[i].(type) {
			case ValueInt:
				opts = append(opts, redis.DialDatabase(v.AsInt()))
			case ValueObject:
				if val, ok := v.GetMember("username", c).(ValueStr); ok {
					opts = append(opts, redis.DialUsername(val.Value()))
				}
				if val, ok := v.GetMember("password", c).(ValueStr); ok {
					opts = append(opts, redis.DialPassword(val.Value()))
				}
				if val, ok := v.GetMember("database", c).(ValueInt); ok {
					opts = append(opts, redis.DialDatabase(val.AsInt()))
				}
				switch val := v.GetMember("readTimeout", c).(type) {
				case ValueFloat:
					opts = append(opts, redis.DialReadTimeout(time.Duration(val.Value())*time.Second))
				case ValueInt:
					opts = append(opts, redis.DialReadTimeout(time.Duration(val.Value())*time.Second))
				}
				switch val := v.GetMember("connTimeout", c).(type) {
				case ValueFloat:
					opts = append(opts, redis.DialConnectTimeout(time.Duration(val.Value())*time.Second))
				case ValueInt:
					opts = append(opts, redis.DialConnectTimeout(time.Duration(val.Value())*time.Second))
				}
				switch val := v.GetMember("writeTimeout", c).(type) {
				case ValueFloat:
					opts = append(opts, redis.DialWriteTimeout(time.Duration(val.Value())*time.Second))
				case ValueInt:
					opts = append(opts, redis.DialWriteTimeout(time.Duration(val.Value())*time.Second))
				}
			}
		}
		conn, err := redis.Dial("tcp", redisAddr.Value(), opts...)
		if err != nil {
			c.RaiseRuntimeError("redis.open error: %s", err)
			return nil
		}
		return NewObjectAndInit(redisClientClass, c, NewGoValue(conn))
	}), nil)
	lib.SetMember("RedisClient", redisClientClass, nil)
	lib.SetMember("RedisPipeSession", redisPipeSessionClass, nil)
	return lib
}

type redisPipeCmd struct {
	cmd  string
	args []interface{}
}

func initRedisPipeSessionClass() ValueType {
	rv := NewClassBuilder("RedisPipeSession").
		Constructor(func(c *Context, thisObj ValueObject, args []Value) {
			thisObj.SetMember("_cmds", NewGoValue([]redisPipeCmd{}), c)
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			cmd := c.MustStr(args[0])
			if cmd == strings.ToUpper(cmd) {
				cmds := this.GetMember("_cmds", c).ToGoValue().([]redisPipeCmd)
				session := this
				return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
					defer func() {
						session.SetMember("_cmds", NewGoValue(cmds), c)
					}()
					cmdArgs := make([]interface{}, len(args))
					for i, arg := range args {
						cmdArgs[i] = arg.ToGoValue()
					}
					cmds = append(cmds, redisPipeCmd{cmd: cmd, args: cmdArgs})
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
			conn := this.GetMember("_conn", c).ToGoValue().(redis.Conn)
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
			conn := this.GetMember("_conn", c).ToGoValue().(redis.Conn)
			if cmds, ok := args[0].(ValueArray); ok {
				n := cmds.Len()
				for i := 0; i < n; i++ {
					var (
						cmdName string
						cmdArgs = []interface{}{}
					)
					cmd := cmds.GetIndex(i, c)
					if cmdItems, ok := cmd.(ValueArray); ok {
						if n2 := cmdItems.Len(); n2 > 0 {
							cmdName = cmdItems.GetIndex(0, c).ToString(c)
							cmdArgs = make([]interface{}, 0, n2-1)
							for i := 1; i < n2; i++ {
								cmdArgs = append(cmdArgs, cmdItems.GetIndex(i, c).ToGoValue())
							}
						} else {
							cmdName = "PING"
							cmdArgs = []interface{}{}
						}
					} else {
						cmdName = cmd.ToString(c)
					}
					if err := conn.Send(cmdName, cmdArgs...); err != nil {
						c.RaiseRuntimeError("redis.exec: send piped command error: %s", err)
					}
				}
				if err := conn.Flush(); err != nil {
					c.RaiseRuntimeError("redis.exec: flush piped command error: %s", err)
				}
				rv := NewArray(n)
				for i := 0; i < n; i++ {
					rsp, err := conn.Receive()
					if err != nil {
						c.RaiseRuntimeError("redis.exec: receive piped command error: %s", err)
					}
					if rsp == nil {
						rv.PushBack(Nil())
					} else {
						v := reflect.ValueOf(rsp)
						rv.PushBack(FromGoValue(v, c))
					}
				}
				return rv
			} else {
				cmd := c.MustStr(args[0], "command")
				cmdArgs := make([]interface{}, len(args)-1)
				for i, arg := range args[1:] {
					cmdArgs[i] = arg.ToGoValue()
				}
				rsp, err := conn.Do(cmd, cmdArgs...)
				if err != nil {
					if err == redis.ErrNil {
						return Nil()
					}
					c.RaiseRuntimeError("RedisClient.exec error %s", err)
					return nil
				}
				if rsp == nil {
					return Nil()
				}
				v := reflect.ValueOf(rsp)
				return FromGoValue(v, c)
			}
		}).
		Method("pipe", func(c *Context, this ValueObject, args []Value) Value {
			action := c.MustCallable(args[0])
			session := NewObjectAndInit(redisPipeSessionClass, c)
			c.Invoke(action, nil, Args(session))
			cmds := session.GetMember("_cmds", c).ToGoValue().([]redisPipeCmd)
			conn := this.GetMember("_conn", c).ToGoValue().(redis.Conn)
			for _, cmd := range cmds {
				if err := conn.Send(cmd.cmd, cmd.args...); err != nil {
					c.RaiseRuntimeError("redis.pipe: send piped command error: %s", err)
				}
			}
			if err := conn.Flush(); err != nil {
				c.RaiseRuntimeError("redis.pipe: flush piped command error: %s", err)
			}
			n := len(cmds)
			rv := NewArray(n)
			for i := 0; i < n; i++ {
				rsp, err := conn.Receive()
				if err != nil {
					c.RaiseRuntimeError("redis.exec: receive piped command error: %s", err)
				}
				if rsp == nil {
					rv.PushBack(Nil())
				} else {
					v := reflect.ValueOf(rsp)
					rv.PushBack(FromGoValue(v, c))
				}
			}
			return rv
		}).
		Method("__getAttr__", func(c *Context, this ValueObject, args []Value) Value {
			cmd := c.MustStr(args[0])
			if cmd == strings.ToUpper(cmd) {
				conn := this.GetMember("_conn", c).ToGoValue().(redis.Conn)
				return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
					cmdArgs := make([]interface{}, len(args))
					for i, arg := range args {
						cmdArgs[i] = arg.ToGoValue()
					}
					rsp, err := conn.Do(cmd, cmdArgs...)
					if err != nil {
						if err == redis.ErrNil {
							return Nil()
						}
						c.RaiseRuntimeError("RedisClient.exec error %s", err)
						return nil
					}
					if rsp == nil {
						return Nil()
					}
					v := reflect.ValueOf(rsp)
					return FromGoValue(v, c)
				})
			} else {
				return Undefined()
			}
		}).
		Build()
	return rv
}

func init() {
	redisClientClass = initRedisClientClass()
	redisPipeSessionClass = initRedisPipeSessionClass()
}
