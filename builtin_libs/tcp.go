package builtin_libs

import (
	"net"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	tcpConnClass  ValueType
	tcpNoDeadline ValueObject
)

func libTcp(c *Context) ValueObject {
	lib := NewObject()
	tcpNoDeadline = NewObjectAndInit(timeDurationClass, c, NewGoValue(time.Duration(-1)))
	tcpInitConnClass()
	lib.SetMember("connect", NewNativeFunction("connect", func(c *Context, _ Value, args []Value) Value {
		var (
			remoteAddr   ValueStr
			waitDuration timeDurationArg
			conn         net.Conn
			err          error
		)
		EnsureFuncParams(c, "connect", args,
			ArgRuleRequired("remoteAddr", TypeStr, &remoteAddr),
			waitDuration.Rule(c, "waitDuration", tcpNoDeadline),
		)
		if d := waitDuration.GetDuration(c); d >= 0 {
			conn, err = net.DialTimeout("tcp", remoteAddr.Value(), d)
		} else {
			conn, err = net.Dial("tcp", remoteAddr.Value())
		}
		if err != nil {
			c.RaiseRuntimeError("connecting to %s error %+v", remoteAddr.Value(), err)
		} else if tcpc, is := conn.(*net.TCPConn); !is {
			c.RaiseRuntimeError("connecting to %s error invalid conn type", remoteAddr.Value())
		} else {
			return NewObjectAndInit(tcpConnClass, c, NewGoValue(tcpc))
		}
		return nil
	}), nil)
	lib.SetMember("serve", NewNativeFunction("serve", func(c *Context, _ Value, args []Value) Value {
		var (
			argAddr    ValueStr
			handleConn ValueCallable
		)
		EnsureFuncParams(c, "serve", args,
			ArgRuleRequired("addr", TypeStr, &argAddr),
			ArgRuleRequired("handleConn", TypeCallable, &handleConn),
		)
		addr, err := net.ResolveTCPAddr("tcp", argAddr.Value())
		if err != nil {
			c.RaiseRuntimeError("resolve listen addr %s error %+v", argAddr.Value(), err)
		}
		ln, err := net.ListenTCP("tcp", addr)
		if err != nil {
			c.RaiseRuntimeError("listen at %s error %+v", argAddr.Value(), err)
		}
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				break
			}
			go func(c *Context) {
				defer recover()
				cliConn := NewObjectAndInit(tcpConnClass, c, NewGoValue(conn))
				c.Invoke(handleConn, nil, Args(cliConn))
			}(c.Clone())
		}
		return Undefined()
	}), nil)
	return lib
}

type (
	tcpConnReserved struct {
		conn *net.TCPConn
	}
)

func tcpInitConnClass() {
	tcpConnClass = NewClassBuilder("TcpConn").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var conn GoValue
			EnsureFuncParams(c, "TcpConn.__init__", args,
				ArgRuleRequired("conn", TypeGoValue, &conn),
			)
			this.Reserved = &tcpConnReserved{
				conn: conn.ToGoValue(c).(*net.TCPConn),
			}
		}).
		Method("addr", func(c *Context, this ValueObject, args []Value) Value {
			conn := this.Reserved.(*tcpConnReserved).conn
			rv := NewObject()
			rv.SetMember("remote", NewStr(conn.RemoteAddr().String()), c)
			rv.SetMember("local", NewStr(conn.LocalAddr().String()), c)
			return rv
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			this.Reserved.(*tcpConnReserved).conn.Close()
			return this
		}).
		Method("send", func(c *Context, this ValueObject, args []Value) Value {
			var (
				content      Value
				waitDuration timeDurationArg
			)
			EnsureFuncParams(c, "TcpConn.send", args,
				ArgRuleRequired("content", TypeAny, &content),
				waitDuration.Rule(c, "waitDuration", tcpNoDeadline),
			)
			conn := this.Reserved.(*tcpConnReserved).conn
			if d := waitDuration.GetDuration(c); d >= 0 {
				if err := conn.SetWriteDeadline(time.Now().Add(d)); err != nil {
					c.RaiseRuntimeError("TcpConn.send: send set write deadline error %+v", err)
				}
			}
			var pkg []byte
			switch arg := content.(type) {
			case ValueBytes:
				pkg = arg.Value()
			default:
				pkg = []byte(arg.ToString(c))
			}
			n, err := conn.Write(pkg)
			if err != nil {
				if e, is := err.(net.Error); is && e.Timeout() {
					return NewInt(0)
				}
				c.RaiseRuntimeError("TcpConn.send: send error %+v", err)
			}
			return NewInt(int64(n))
		}).
		Method("recv", func(c *Context, this ValueObject, args []Value) Value {
			var (
				maxRecv      ValueInt
				waitDuration timeDurationArg
			)
			EnsureFuncParams(c, "TcpConn.recv", args,
				ArgRuleRequired("maxRecv", TypeInt, &maxRecv),
				waitDuration.Rule(c, "waitDuration", tcpNoDeadline),
			)
			conn := this.Reserved.(*tcpConnReserved).conn
			buf := make([]byte, maxRecv.AsInt())
			if d := waitDuration.GetDuration(c); d >= 0 {
				if err := conn.SetReadDeadline(time.Now().Add(d)); err != nil {
					c.RaiseRuntimeError("TcpConn.recv: recv set read deadline error %+v", err)
				}
			}
			n, err := conn.Read(buf)
			if err != nil {
				if e, is := err.(net.Error); is && e.Timeout() {
					return NewBytes(buf[:0])
				}
				c.RaiseRuntimeError("TcpConn.recv: recv error %+v", err)
			}
			return NewBytes(buf[:n])
		}, "maxRecv", "waitDuration").
		Build()
}
