package builtin_libs

import (
	"net"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	tcpConnClass ValueType
)

func libTcp(*Context) ValueObject {
	lib := NewObject()
	tcpInitConnClass()
	lib.SetMember("connect", NewNativeFunction("connect", func(c *Context, _ Value, args []Value) Value {
		var remoteAddr, localAddr ValueStr
		EnsureFuncParams(c, "connect", args,
			ArgRuleRequired("remoteAddr", TypeStr, &remoteAddr),
			ArgRuleOptional("localAddr", TypeStr, &localAddr, NewStr("")),
		)
		raddr, err := net.ResolveTCPAddr("tcp", remoteAddr.Value())
		if err != nil {
			c.RaiseRuntimeError("resolve remote addr %s error %+v", remoteAddr.Value(), err)
		}
		var laddr *net.TCPAddr
		if la := localAddr.Value(); la != "" {
			if laddr, err = net.ResolveTCPAddr("tcp", la); err != nil {
				c.RaiseRuntimeError("resolve local addr %s error %+v", la, err)
			}
		}
		conn, err := net.DialTCP("tcp", laddr, raddr)
		if err != nil {
			c.RaiseRuntimeError("connecting to %s error %+v", remoteAddr.Value(), err)
		}
		return NewObjectAndInit(tcpConnClass, c, NewGoValue(conn))
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
			conn := this.Reserved.(*tcpConnReserved).conn
			totalSent := 0
			for _, a := range args {
				var pkg []byte
				switch arg := a.(type) {
				case ValueBytes:
					pkg = arg.Value()
				default:
					pkg = []byte(arg.ToString(c))
				}
				n, err := conn.Write(pkg)
				if err != nil {
					c.RaiseRuntimeError("TcpConn.send: send error %+v", err)
				}
				totalSent += n
				if n < len(pkg) {
					break
				}
			}
			return NewInt(int64(totalSent))
		}).
		Method("recv", func(c *Context, this ValueObject, args []Value) Value {
			var maxRecv ValueInt
			EnsureFuncParams(c, "TcpConn.recv", args,
				ArgRuleRequired("maxRecv", TypeInt, &maxRecv))
			conn := this.Reserved.(*tcpConnReserved).conn
			buf := make([]byte, maxRecv.AsInt())
			n, err := conn.Read(buf)
			if err != nil {
				c.RaiseRuntimeError("TcpConn.recv: recv error %+v", err)
			}
			return NewBytes(buf[:n])
		}).
		Build()
}
