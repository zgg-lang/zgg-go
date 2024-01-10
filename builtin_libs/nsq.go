package builtin_libs

import (
	"io"
	"log"
	"time"

	. "github.com/zgg-lang/zgg-go/runtime"

	nsq "github.com/nsqio/go-nsq"
)

var nsqDefaultLogger = log.New(io.Discard, "", 0)

func libNsq(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Consumer", nsqConsumerClass, c)
	lib.SetMember("Producer", nsqProducerClass, c)
	lib.SetMember("publish", NewNativeFunction("nsq.publish", func(c *Context, this Value, args []Value) Value {
		var (
			addr  ValueStr
			topic ValueStr
			data  Value
		)
		EnsureFuncParams(c, "publish", args,
			ArgRuleRequired("addr", TypeStr, &addr),
			ArgRuleRequired("topic", TypeStr, &topic),
			ArgRuleRequired("data", TypeAny, &data),
		)
		p, err := nsq.NewProducer(addr.Value(), nsq.NewConfig())
		if err != nil {
			c.RaiseRuntimeError("nsq.publish new producer error %+v", err)
		}
		defer p.Stop()
		p.SetLogger(nsqDefaultLogger, nsq.LogLevelMax)
		var payload []byte
		switch v := data.(type) {
		case ValueBytes:
			payload = v.Value()
		case ValueStr:
			payload = []byte(v.Value())
		default:
			var err error
			payload, err = jsonMarshal(data.ToGoValue(c))
			if err != nil {
				c.RaiseRuntimeError("nsq.Producer.publish marshal data error %+v", err)
			}
		}
		if err := p.Publish(topic.Value(), payload); err != nil {
			c.RaiseRuntimeError("nsq.publish publish error %+v", err)
		}
		return Undefined()
	}), c)
	return lib
}

var (
	nsqConsumerClass = NewClassBuilder("nsq.Consumer").
				Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				dType ValueStr
				addr  ValueStr
			)
			EnsureFuncParams(c, "nsq.Consumer.__init__", args,
				ArgRuleRequired("dType", TypeStr, &dType),
				ArgRuleRequired("addr", TypeStr, &addr),
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
				ArgRuleRequired("topic", TypeStr, &topic),
				ArgRuleRequired("channel", TypeStr, &channel),
				ArgRuleRequired("callback", TypeFunc, &callback),
			)
			conf := nsq.NewConfig()
			conf.LookupdPollInterval = 60 * time.Second
			consumer, err := nsq.NewConsumer(topic.Value(), channel.Value(), conf)
			if err != nil {
				c.RaiseRuntimeError("create nsq consumer error %s", err)
			}
			consumer.SetLogger(nsqDefaultLogger, nsq.LogLevelMax)
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
				c := consumers.GetIndex(i, c).ToGoValue(c).(*nsq.Consumer)
				c.Stop()
			}
			return Undefined()
		}).
		Build()

	nsqProducerClass = NewClassBuilder("nsq.Producer").
				Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				addr ValueStr
			)
			EnsureFuncParams(c, "nsq.Consumer.__init__", args,
				ArgRuleRequired("addr", TypeStr, &addr),
			)
			p, err := nsq.NewProducer(addr.Value(), nsq.NewConfig())
			if err != nil {
				c.RaiseRuntimeError("Init nsq producer error %+v", err)
			}
			p.SetLogger(nsqDefaultLogger, nsq.LogLevelMax)
			this.SetMember("_addr", addr, c)
			this.SetMember("_producer", NewGoValue(p), c)
		}).
		Method("publish", func(c *Context, this ValueObject, args []Value) Value {
			var (
				topic ValueStr
				data  Value
			)
			EnsureFuncParams(c, "nsq.Producer.publish", args,
				ArgRuleRequired("topic", TypeStr, &topic),
				ArgRuleRequired("data", TypeAny, &data),
			)
			p := this.GetMember("_producer", c).ToGoValue(c).(*nsq.Producer)
			var payload []byte
			switch v := data.(type) {
			case ValueBytes:
				payload = v.Value()
			case ValueStr:
				payload = []byte(v.Value())
			default:
				var err error
				payload, err = jsonMarshal(data.ToGoValue(c))
				if err != nil {
					c.RaiseRuntimeError("nsq.Producer.publish marshal data error %+v", err)
				}
			}
			if err := p.Publish(topic.Value(), payload); err != nil {
				c.RaiseRuntimeError("nsq.Producer.publish publish error %+v", err)
			}
			return this
		}).
		Method("ping", func(c *Context, this ValueObject, args []Value) Value {
			p := this.GetMember("_producer", c).ToGoValue(c).(*nsq.Producer)
			if err := p.Ping(); err != nil {
				c.RaiseRuntimeError("nsq.Producer.ping ping error %+v", err)
			}
			return this
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			p := this.GetMember("_producer", c).ToGoValue(c).(*nsq.Producer)
			p.Stop()
			return Undefined()
		}).
		Build()
)
