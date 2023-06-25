package builtin_libs

import (
	"io"

	"github.com/beevik/etree"
	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	etreeDocument ValueType
	etreeElement  ValueType
)

func libEtree(*Context) ValueObject {
	lib := NewObject()
	lib.SetMember("fromFile", NewNativeFunction("fromFile", func(c *Context, this Value, args []Value) Value {
		var filename ValueStr
		EnsureFuncParams(c, "etree.fromFile", args, ArgRuleRequired("filename", TypeStr, &filename))
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(filename.Value()); err != nil {
			c.RaiseRuntimeError("read from file error %+v", err)
		}
		return NewObjectAndInit(etreeDocument, c, NewGoValue(doc))
	}, "filename"), nil)
	lib.SetMember("fromString", NewNativeFunction("fromString", func(c *Context, this Value, args []Value) Value {
		var xmlStr ValueStr
		EnsureFuncParams(c, "etree.fromString", args, ArgRuleRequired("xmlStr", TypeStr, &xmlStr))
		doc := etree.NewDocument()
		if err := doc.ReadFromString(xmlStr.Value()); err != nil {
			c.RaiseRuntimeError("read from string error %+v", err)
		}
		return NewObjectAndInit(etreeDocument, c, NewGoValue(doc))
	}, "xmlStr"), nil)
	lib.SetMember("fromBytes", NewNativeFunction("fromBytes", func(c *Context, this Value, args []Value) Value {
		var xmlBytes ValueBytes
		EnsureFuncParams(c, "etree.fromBytes", args, ArgRuleRequired("xmlBytes", TypeBytes, &xmlBytes))
		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(xmlBytes.Value()); err != nil {
			c.RaiseRuntimeError("read from Bytes error %+v", err)
		}
		return NewObjectAndInit(etreeDocument, c, NewGoValue(doc))
	}, "xmlBytes"), nil)
	lib.SetMember("new", NewNativeFunction("new", func(c *Context, this Value, args []Value) Value {
		var rootTag ValueStr
		EnsureFuncParams(c, "etree.new", args, ArgRuleRequired("rootTag", TypeStr, &rootTag))
		doc := etree.NewDocument()
		doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
		doc.SetRoot(etree.NewElement(rootTag.Value()))
		return NewObjectAndInit(etreeDocument, c, NewGoValue(doc))
	}), nil)
	return lib
}

func etreeInitDocument() {
	etreeDocument = NewClassBuilder("Document").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var doc GoValue
			EnsureFuncParams(c, "Document.__init__", args, ArgRuleRequired("doc", TypeGoValue, &doc))
			this.SetMember("__doc", doc, c)
		}).
		Method("encode", func(c *Context, this ValueObject, args []Value) Value {
			var indent ValueInt
			EnsureFuncParams(c, "Document.encode", args, ArgRuleOptional("indent", TypeInt, &indent, NewInt(0)))
			doc := this.GetMember("__doc", c).ToGoValue().(*etree.Document)
			if ind := indent.AsInt(); ind > 0 {
				doc.Indent(ind)
				defer doc.Unindent()
			}
			s, err := doc.WriteToString()
			if err != nil {
				c.RaiseRuntimeError("document encoding error %+v", err)
			}
			return NewStr(s)
		}).
		Method("write", func(c *Context, this ValueObject, args []Value) Value {
			var (
				writer GoValue
				indent ValueInt
			)
			EnsureFuncParams(c, "Document.write", args,
				ArgRuleRequired("writer", TypeGoValue, &writer),
				ArgRuleOptional("indent", TypeInt, &indent, NewInt(0)),
			)
			doc := this.GetMember("__doc", c).ToGoValue().(*etree.Document)
			if ind := indent.AsInt(); ind > 0 {
				doc.Indent(ind)
				defer doc.Unindent()
			}
			if w, is := writer.ToGoValue().(io.Writer); !is {
				c.RaiseRuntimeError("document write target is not a writer")
			} else if _, err := doc.WriteTo(w); err != nil {
				c.RaiseRuntimeError("document write error %+v", err)
			}
			return this
		}).
		Method("root", func(c *Context, this ValueObject, args []Value) Value {
			doc := this.GetMember("__doc", c).ToGoValue().(*etree.Document)
			el := doc.Root()
			if el == nil {
				return Nil()
			}
			return NewObjectAndInit(etreeElement, c, NewGoValue(el))
		}).
		Method("add", func(c *Context, this ValueObject, args []Value) Value {
			root := c.InvokeMethod(this, "root", NoArgs)
			if nv, is := root.(ValueNil); is {
				return nv
			}
			return c.InvokeMethod(root, "add", Args(args...))
		}).
		Method("find", func(c *Context, this ValueObject, args []Value) Value {
			var path ValueStr
			EnsureFuncParams(c, "Document.find", args, ArgRuleRequired("path", TypeStr, &path))
			doc := this.GetMember("__doc", c).ToGoValue().(*etree.Document)
			el := doc.FindElement(path.Value())
			if el == nil {
				return Nil()
			}
			return NewObjectAndInit(etreeElement, c, NewGoValue(el))
		}).
		Method("findAll", func(c *Context, this ValueObject, args []Value) Value {
			var path ValueStr
			EnsureFuncParams(c, "Document.findAll", args, ArgRuleRequired("path", TypeStr, &path))
			doc := this.GetMember("__doc", c).ToGoValue().(*etree.Document)
			els := doc.FindElements(path.Value())
			rv := NewArray(len(els))
			for _, el := range els {
				v := NewObjectAndInit(etreeElement, c, NewGoValue(el))
				rv.PushBack(v)
			}
			return rv
		}).
		Build()
}

func etreeInitElement() {
	etreeElement = NewClassBuilder("Element").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var el GoValue
			EnsureFuncParams(c, "Element.__init__", args, ArgRuleRequired("el", TypeGoValue, &el))
			this.SetMember("__el", el, c)
		}).
		Method("find", func(c *Context, this ValueObject, args []Value) Value {
			var path ValueStr
			EnsureFuncParams(c, "Element.find", args, ArgRuleRequired("path", TypeStr, &path))
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			el := cur.FindElement(path.Value())
			if el == nil {
				return Nil()
			}
			return NewObjectAndInit(etreeElement, c, NewGoValue(el))
		}).
		Method("findAll", func(c *Context, this ValueObject, args []Value) Value {
			var path ValueStr
			EnsureFuncParams(c, "Element.findAll", args, ArgRuleRequired("path", TypeStr, &path))
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			els := cur.FindElements(path.Value())
			rv := NewArray(len(els))
			for _, el := range els {
				v := NewObjectAndInit(etreeElement, c, NewGoValue(el))
				rv.PushBack(v)
			}
			return rv
		}).
		Method("children", func(c *Context, this ValueObject, args []Value) Value {
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			els := cur.ChildElements()
			rv := NewArray(len(els))
			for _, el := range els {
				v := NewObjectAndInit(etreeElement, c, NewGoValue(el))
				rv.PushBack(v)
			}
			return rv
		}).
		Method("add", func(c *Context, this ValueObject, args []Value) Value {
			var (
				tag   ValueStr
				attrs ValueObject
			)
			EnsureFuncParams(c, "Element.add", args,
				ArgRuleRequired("tag", TypeStr, &tag),
				ArgRuleOptional("attrs", TypeObject, &attrs, NewObject()),
			)
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			child := cur.CreateElement(tag.Value())
			attrs.Iterate(func(k string, v Value) {
				child.CreateAttr(k, v.ToString(c))
			})
			return NewObjectAndInit(etreeElement, c, NewGoValue(child))
		}).
		Method("setAttr", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) == 1 {
				var attrs ValueObject
				EnsureFuncParams(c, "Element.setAttr", args,
					ArgRuleRequired("attrs", TypeObject, &attrs),
				)
				cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
				attrs.Iterate(func(k string, v Value) {
					cur.CreateAttr(k, v.ToString(c))
				})
			} else {
				var key ValueStr
				var value Value
				EnsureFuncParams(c, "Element.setAttr", args,
					ArgRuleRequired("key", TypeStr, &key),
					ArgRuleOptional("value", TypeAny, &value, Undefined()),
				)
				cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
				if IsUndefined(value) {
					cur.RemoveAttr(key.Value())
				} else {
					cur.CreateAttr(key.Value(), value.ToString(c))
				}
			}
			return this
		}).
		Method("setText", func(c *Context, this ValueObject, args []Value) Value {
			var text Value
			EnsureFuncParams(c, "Element.setText", args,
				ArgRuleRequired("text", TypeAny, &text),
			)
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			cur.SetText(text.ToString(c))
			return this
		}).
		Method("tag", func(c *Context, this ValueObject, args []Value) Value {
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			return NewStr(cur.Tag)
		}).
		Method("fullTag", func(c *Context, this ValueObject, args []Value) Value {
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			return NewStr(cur.FullTag())
		}).
		Method("attr", func(c *Context, this ValueObject, args []Value) Value {
			var key ValueStr
			EnsureFuncParams(c, "Element.attr", args,
				ArgRuleRequired("key", TypeStr, &key),
			)
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			attr := cur.SelectAttr(key.Value())
			if attr == nil {
				return Nil()
			}
			return NewStr(attr.Value)
		}).
		Method("attrs", func(c *Context, this ValueObject, args []Value) Value {
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			rv := NewObject()
			for _, attr := range cur.Attr {
				rv.SetMember(attr.FullKey(), NewStr(attr.Value), c)
			}
			return rv
		}).
		Method("text", func(c *Context, this ValueObject, args []Value) Value {
			cur := this.GetMember("__el", c).ToGoValue().(*etree.Element)
			return NewStr(cur.Text())
		}).
		Build()
}

func init() {
	etreeInitDocument()
	etreeInitElement()
}
