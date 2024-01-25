package builtin_libs

import (
	"fmt"
	"net"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	udpConnClass  ValueType
	udpNoDeadline ValueObject
)

func libUdp(c *Context) ValueObject {
	lib := NewObject()
	udpInitConnClass()
	udpNoDeadline = NewObjectAndInit(timeDurationClass, c, NewGoValue(time.Duration(-1)))
	lib.SetMember("Conn", udpConnClass, c)
	lib.SetMember("__call__", udpConnClass, c)
	return lib
}

type (
	udpConnReserved struct {
		conn *net.UDPConn
	}
)

func udpInitConnClass() {
	udpConnClass = NewClassBuilder("Conn").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var laddr ValueStr
			EnsureFuncParams(c, "Conn.__init__", args,
				ArgRuleOptional("laddr", TypeStr, &laddr, NewStr("")),
			)
			var conn *net.UDPConn
			var err error
			if addr := laddr.Value(); addr != "" {
				var la *net.UDPAddr
				la, err = net.ResolveUDPAddr("udp", addr)
				if err == nil {
					conn, err = net.ListenUDP("udp", la)
					fmt.Println("ListenUDP", la, err)
				}
			} else {
				conn, err = net.ListenUDP("udp", nil)
			}
			if err != nil {
				c.RaiseRuntimeError("Conn.__init__: init with address [%s] error %s", laddr.Value(), err)
			}
			this.Reserved = &udpConnReserved{
				conn: conn,
			}
		}).
		Method("addr", func(c *Context, this ValueObject, args []Value) Value {
			conn := this.Reserved.(*udpConnReserved).conn
			return NewStr(conn.LocalAddr().String())
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			this.Reserved.(*udpConnReserved).conn.Close()
			return this
		}).
		Method("sendTo", func(c *Context, this ValueObject, args []Value) Value {
			var (
				content      Value
				address      ValueStr
				waitDuration timeDurationArg
				raddr        *net.UDPAddr
				err          error
			)
			EnsureFuncParams(c, "Conn.sendTo", args,
				ArgRuleRequired("content", TypeAny, &content),
				ArgRuleRequired("raddr", TypeStr, &address),
				waitDuration.Rule(c, "waitDuration", udpNoDeadline),
			)
			if raddr, err = net.ResolveUDPAddr("udp", address.Value()); err != nil {
				c.RaiseRuntimeError("Conn.sendTo: resolve address %s error %+v", address.Value(), err)
			}
			conn := this.Reserved.(*udpConnReserved).conn
			if d := waitDuration.GetDuration(c); d >= 0 {
				if err := conn.SetWriteDeadline(time.Now().Add(d)); err != nil {
					c.RaiseRuntimeError("Conn.sendTo: sendTo set write deadline error %+v", err)
				}
			}
			var pkg []byte
			switch arg := content.(type) {
			case ValueBytes:
				pkg = arg.Value()
			default:
				pkg = []byte(arg.ToString(c))
			}
			n, err := conn.WriteTo(pkg, raddr)
			if err != nil {
				if e, is := err.(net.Error); is && e.Timeout() {
					return NewInt(0)
				}
				c.RaiseRuntimeError("Conn.sendTo: send error %+v", err)
			}
			return NewInt(int64(n))
		}).
		Method("recvFrom", func(c *Context, this ValueObject, args []Value) Value {
			var (
				maxRecv      ValueInt
				waitDuration timeDurationArg
			)
			EnsureFuncParams(c, "Conn.recvFrom", args,
				ArgRuleRequired("maxRecv", TypeInt, &maxRecv),
				waitDuration.Rule(c, "waitDuration", udpNoDeadline),
			)
			conn := this.Reserved.(*udpConnReserved).conn
			buf := make([]byte, maxRecv.AsInt())
			if d := waitDuration.GetDuration(c); d >= 0 {
				if err := conn.SetReadDeadline(time.Now().Add(d)); err != nil {
					c.RaiseRuntimeError("Conn.recvFrom: recvFrom set read deadline error %+v", err)
				}
			}
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				if e, is := err.(net.Error); is && e.Timeout() {
					return NewArrayByValues(NewBytes(buf[:0]), NewStr(""))
				}
				c.RaiseRuntimeError("Conn.recv: recv error %+v", err)
			}
			return NewArrayByValues(NewBytes(buf[:n]), NewStr(addr.String()))
		}, "maxRecv", "waitDuration").
		Build()
}
