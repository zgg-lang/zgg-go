package builtin_libs

import (
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"

	nsq "github.com/nsqio/go-nsq"
)

func libNsq(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Consumer", nsqConsumerClass, c)
	return lib
}

var nsqConsumerClass = NewClassBuilder("nsq.Consumer").
	Constructor(func(c *Context, this ValueObject, args []Value) {
		var (
			dType ValueStr
			addr  ValueStr
		)
		EnsureFuncParams(c, "nsq.Consumer.__init__", args,
			ArgRuleRequired{"dType", TypeStr, &dType},
			ArgRuleRequired{"addr", TypeStr, &addr},
		)
		if dt := dType.Value(); dt != "nsqd" && dt != "nsqlookupd" {
			c.RaiseRuntimeError("nsq.Consumer.__init__ invalid dType %s", dt)
		}
		this.SetMember("_dType", dType, c)
		this.SetMember("_addr", addr, c)
		this.SetMember("_consumers", NewArray(), c)
	}).
	Method("on", func(c *Context, this ValueObject, args []Value) Value {
		dType := c.MustStr(this.GetMember("_dType", c))
		addr := c.MustStr(this.GetMember("_addr", c))
		var (
			topic    ValueStr
			channel  ValueStr
			callback ValueCallable
		)
		EnsureFuncParams(c, "nsq.Consumer.on", args,
			ArgRuleRequired{"topic", TypeStr, &topic},
			ArgRuleRequired{"channel", TypeStr, &channel},
			ArgRuleRequired{"callback", TypeFunc, &callback},
		)
		conf := nsq.NewConfig()
		conf.LookupdPollInterval = 60 * time.Second
		consumer, err := nsq.NewConsumer(topic.Value(), channel.Value(), conf)
		if err != nil {
			c.RaiseRuntimeError("create nsq consumer error %s", err)
		}
		consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
			newC := c.Clone()
			newC.Invoke(callback, nil, Args(
				NewBytes(msg.Body),
				NewGoValue(msg),
			))
			return nil
		}))
		switch dType {
		case "nsqd":
			err = consumer.ConnectToNSQD(addr)
		case "nsqlookupd":
			err = consumer.ConnectToNSQLookupd(addr)
		default:
			c.RaiseRuntimeError("invalid nsq dtype %s", dType)
		}
		if err != nil {
			c.RaiseRuntimeError("nsq consumer connect to lookupd err %s", err)
		}
		consumers := c.MustArray(this.GetMember("_consumers", c))
		consumers.PushBack(NewGoValue(consumer))
		return this
	}).
	Method("closeAll", func(c *Context, this ValueObject, args []Value) Value {
		consumers := c.MustArray(this.GetMember("_consumers", c))
		for i := 0; i < consumers.Len(); i++ {
			c := consumers.GetIndex(i, c).ToGoValue().(*nsq.Consumer)
			c.Stop()
		}
		return Undefined()
	}).
	Build()
