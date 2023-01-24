package builtin_libs

import (
	"image"
	"image/color"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	. "github.com/zgg-lang/zgg-go/runtime"
)

func libDrawing(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("Canvas", drawingCanvasClass, nil)
	lib.SetMember("from", NewNativeFunction("from", func(c *Context, this Value, args []Value) Value {
		var (
			s  GoValue
			gc *gg.Context
		)
		EnsureFuncParams(c, "fromImage", args, ArgRuleRequired("srcImage", TypeGoValue, &s))
		switch v := s.ToGoValue().(type) {
		case *image.RGBA:
			gc = gg.NewContextForRGBA(v)
		case image.Image:
			gc = gg.NewContextForImage(v)
		default:
			c.RaiseRuntimeError("srcImage is invalid")
		}
		rv := NewObject(drawingCanvasClass)
		rv.SetMember("__dc", NewGoValue(gc), c)
		return rv
	}, "srcImage"), nil)
	lib.SetMember("Font", drawingFontClass, c)
	return lib
}

var (
	drawingColorPattern3 = regexp.MustCompile(`[0-9a-fA-F]{3}$`)
	drawingColorPattern4 = regexp.MustCompile(`#[0-9a-fA-F]{3}$`)
	drawingColorPattern6 = regexp.MustCompile(`[0-9a-fA-F]{6}$`)
	drawingColorPattern7 = regexp.MustCompile(`#[0-9a-fA-F]{6}$`)
)

func drawingGetHexFromString(c string, begin, end int) uint8 {
	v, _ := strconv.ParseUint(c[begin:end], 16, 8)
	return uint8(v)
}

func drawingParseColor(c string) (out color.Color, ok bool) {
	var r, g, b, a uint8
	switch len(c) {
	case 3:
		if !drawingColorPattern3.MatchString(c) {
			return
		}
		r = drawingGetHexFromString(c, 0, 1) * 0x11
		g = drawingGetHexFromString(c, 1, 2) * 0x11
		b = drawingGetHexFromString(c, 2, 3) * 0x11
		a = 0xff
	case 4:
		if !drawingColorPattern4.MatchString(c) {
			return
		}
		r = drawingGetHexFromString(c, 1, 2) * 0x11
		g = drawingGetHexFromString(c, 2, 3) * 0x11
		b = drawingGetHexFromString(c, 3, 4) * 0x11
		a = 0xff
	case 6:
		if !drawingColorPattern6.MatchString(c) {
			return
		}
		r = drawingGetHexFromString(c, 0, 2)
		g = drawingGetHexFromString(c, 2, 4)
		b = drawingGetHexFromString(c, 4, 6)
		a = 0xff
	case 7:
		if !drawingColorPattern7.MatchString(c) {
			return
		}
		r = drawingGetHexFromString(c, 1, 3)
		g = drawingGetHexFromString(c, 3, 5)
		b = drawingGetHexFromString(c, 5, 7)
		a = 0xff
	}
	ok = true
	out = color.RGBA{R: r, G: g, B: b, A: a}
	return
}

func drawingMustParseColor(c *Context, cs string) color.Color {
	rv, is := drawingParseColor(cs)
	if !is {
		c.RaiseRuntimeError("invalid color %s", cs)
	}
	return rv
}

func drawingUseColor(c *Context, dc *gg.Context, cs ValueStr, f func()) {
	if css := cs.Value(); css != "" {
		colors := strings.Split(css, ";")
		var stroke, fill string
		if len(colors) == 1 {
			stroke = colors[0]
			fill = colors[0]
		} else {
			stroke = colors[0]
			fill = colors[1]
		}
		defer func() {
			if stroke != "" {
				dc.SetStrokeStyle(gg.NewSolidPattern(drawingMustParseColor(c, stroke)))
				dc.Stroke()
			}
			if fill != "" {
				dc.SetFillStyle(gg.NewSolidPattern(drawingMustParseColor(c, fill)))
				dc.Fill()
			}
		}()
	}
	f()
}

var (
	drawingFontClass = NewClassBuilder("Font").
				Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				filename ValueStr
				points   ValueFloat
			)
			EnsureFuncParams(c, "Font.__init__", args,
				ArgRuleRequired("filename", TypeStr, &filename),
				ArgRuleRequired("points", TypeFloat, &points),
			)
			ff, err := gg.LoadFontFace(filename.Value(), points.Value())
			if err != nil {
				c.RaiseRuntimeError("load fontface file %s error %s", filename.Value(), err)
			}
			this.SetMember("__ff", NewGoValue(ff), c)
		}).
		Build()
	drawingCanvasClass = NewClassBuilder("Canvas").
				Constructor(func(c *Context, this ValueObject, args []Value) {
			var width, height ValueInt
			var bgColor ValueStr
			EnsureFuncParams(c, "Canvas.__init__", args,
				ArgRuleRequired("width", TypeInt, &width),
				ArgRuleRequired("height", TypeInt, &height),
				ArgRuleOptional("bgColor", TypeStr, &bgColor, NewStr("")),
			)
			dc := gg.NewContext(width.AsInt(), height.AsInt())
			if bgc := bgColor.Value(); bgc != "" {
				w := float64(width.AsInt())
				h := float64(height.AsInt())
				dc.DrawRectangle(0, 0, w, h)
				dc.SetColor(drawingMustParseColor(c, bgc))
				dc.Fill()
			}
			this.SetMember("__dc", NewGoValue(dc), c)
		}).
		Method("set", func(c *Context, this ValueObject, args []Value) Value {
			var opts ValueObject
			EnsureFuncParams(c, "Canvas.set", args,
				ArgRuleRequired("options", TypeObject, &opts),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			if lineCap, is := opts.GetMember("lineCap", c).ToGoValue().(int); is {
				dc.SetLineCap(gg.LineCap(lineCap))
			}
			switch lineWidth := opts.GetMember("lineWidth", c).ToGoValue().(type) {
			case int64:
				dc.SetLineWidth(float64(lineWidth))
			case float64:
				dc.SetLineWidth(lineWidth)
			}
			if color, is := opts.GetMember("color", c).ToGoValue().(string); is {
				dc.SetColor(drawingMustParseColor(c, color))
			}
			if color, is := opts.GetMember("strokeColor", c).ToGoValue().(string); is {
				dc.SetStrokeStyle(gg.NewSolidPattern(drawingMustParseColor(c, color)))
			}
			if color, is := opts.GetMember("fillColor", c).ToGoValue().(string); is {
				dc.SetFillStyle(gg.NewSolidPattern(drawingMustParseColor(c, color)))
			}
			return this
		}, "options").
		Method("moveTo", func(c *Context, this ValueObject, args []Value) Value {
			var x, y ValueFloat
			EnsureFuncParams(c, "Canvas.moveTo", args,
				ArgRuleRequired("x", TypeFloat, &x),
				ArgRuleRequired("y", TypeFloat, &y),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			dc.MoveTo(x.Value(), y.Value())
			return this
		}).
		Method("lineTo", func(c *Context, this ValueObject, args []Value) Value {
			var x, y ValueFloat
			EnsureFuncParams(c, "Canvas.lineTo", args,
				ArgRuleRequired("x", TypeFloat, &x),
				ArgRuleRequired("y", TypeFloat, &y),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			dc.LineTo(x.Value(), y.Value())
			return this
		}).
		Method("pixel", func(c *Context, this ValueObject, args []Value) Value {
			var x, y ValueInt
			var cl ValueStr
			EnsureFuncParams(c, "Canvas.pixel", args,
				ArgRuleRequired("x", TypeInt, &x),
				ArgRuleRequired("y", TypeInt, &y),
				ArgRuleOptional("c", TypeStr, &cl, NewStr("")),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			drawingUseColor(c, dc, cl, func() {
				dc.SetPixel(x.AsInt(), y.AsInt())
			})
			return this
		}).
		Method("line", func(c *Context, this ValueObject, args []Value) Value {
			var x1, y1, x2, y2 ValueFloat
			var cl ValueStr
			EnsureFuncParams(c, "Canvas.line", args,
				ArgRuleRequired("x1", TypeFloat, &x1),
				ArgRuleRequired("y1", TypeFloat, &y1),
				ArgRuleRequired("x2", TypeFloat, &x2),
				ArgRuleRequired("y2", TypeFloat, &y2),
				ArgRuleOptional("c", TypeStr, &cl, NewStr("")),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			drawingUseColor(c, dc, cl, func() {
				dc.DrawLine(x1.Value(), y1.Value(), x2.Value(), y2.Value())
			})
			return this
		}).
		Method("rect", func(c *Context, this ValueObject, args []Value) Value {
			var x, y, w, h ValueFloat
			var cl ValueStr
			EnsureFuncParams(c, "Canvas.line", args,
				ArgRuleRequired("x", TypeFloat, &x),
				ArgRuleRequired("y", TypeFloat, &y),
				ArgRuleRequired("w", TypeFloat, &w),
				ArgRuleRequired("h", TypeFloat, &h),
				ArgRuleOptional("c", TypeStr, &cl, NewStr("")),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			drawingUseColor(c, dc, cl, func() {
				dc.DrawRectangle(x.Value(), y.Value(), w.Value(), h.Value())
			})
			return this
		}).
		Method("circle", func(c *Context, this ValueObject, args []Value) Value {
			var x, y, r ValueFloat
			var cl ValueStr
			EnsureFuncParams(c, "Canvas.circle", args,
				ArgRuleRequired("x", TypeFloat, &x),
				ArgRuleRequired("y", TypeFloat, &y),
				ArgRuleRequired("r", TypeFloat, &r),
				ArgRuleOptional("c", TypeStr, &cl, NewStr("")),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			drawingUseColor(c, dc, cl, func() {
				dc.DrawCircle(x.Value(), y.Value(), r.Value())
			})
			return this
		}).
		Method("arc", func(c *Context, this ValueObject, args []Value) Value {
			var x, y, r, begin, end ValueFloat
			var cl ValueStr
			EnsureFuncParams(c, "Canvas.circle", args,
				ArgRuleRequired("x", TypeFloat, &x),
				ArgRuleRequired("y", TypeFloat, &y),
				ArgRuleRequired("r", TypeFloat, &r),
				ArgRuleRequired("begin", TypeFloat, &begin),
				ArgRuleRequired("end", TypeFloat, &end),
				ArgRuleOptional("c", TypeStr, &cl, NewStr("")),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			drawingUseColor(c, dc, cl, func() {
				dc.DrawArc(x.Value(), y.Value(), r.Value(), begin.Value(), end.Value())
			})
			return this
		}).
		Method("loadFont", func(c *Context, this ValueObject, args []Value) Value {
			var (
				fontPath ValueStr
				points   ValueFloat
			)
			EnsureFuncParams(c, "Canvas.loadFont", args,
				ArgRuleRequired("fontPath", TypeStr, &fontPath),
				ArgRuleRequired("points", TypeFloat, &points),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			if err := dc.LoadFontFace(fontPath.Value(), points.Value()); err != nil {
				c.RaiseRuntimeError("load font face error %s", err)
			}
			return this
		}).
		Method("useFont", func(c *Context, this ValueObject, args []Value) Value {
			var font ValueObject
			EnsureFuncParams(c, "Canvas.useFont", args, ArgRuleRequired("font", TypeObject, &font))
			ff, is := font.GetMember("__ff", c).ToGoValue().(font.Face)
			if !is {
				c.RaiseRuntimeError("Cannot get font face from argument")
			}
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			dc.SetFontFace(ff)
			return this
		}).
		Method("text", func(c *Context, this ValueObject, args []Value) Value {
			var (
				s  ValueStr
				x  ValueFloat
				y  ValueFloat
				cl ValueStr
			)
			EnsureFuncParams(c, "Canvas.text", args,
				ArgRuleRequired("s", TypeStr, &s),
				ArgRuleRequired("x", TypeFloat, &x),
				ArgRuleRequired("y", TypeFloat, &y),
				ArgRuleOptional("c", TypeStr, &cl, NewStr("")),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			drawingUseColor(c, dc, cl, func() {
				dc.DrawString(s.Value(), x.Value(), y.Value())
			})
			return this
		}).
		Method("measureText", func(c *Context, this ValueObject, args []Value) Value {
			var (
				s ValueStr
			)
			EnsureFuncParams(c, "Canvas.measureText", args,
				ArgRuleRequired("s", TypeStr, &s),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			w, h := dc.MeasureString(s.Value())
			return NewArrayByValues(NewFloat(w), NewFloat(h))
		}).
		Method("stroke", func(c *Context, this ValueObject, args []Value) Value {
			var cl ValueStr
			EnsureFuncParams(c, "Canvas.stroke", args,
				ArgRuleRequired("c", TypeStr, &cl),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			dc.SetColor(drawingMustParseColor(c, cl.Value()))
			dc.Stroke()
			return this
		}).
		Method("fill", func(c *Context, this ValueObject, args []Value) Value {
			var cl ValueStr
			EnsureFuncParams(c, "Canvas.fill", args,
				ArgRuleRequired("c", TypeStr, &cl),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			dc.SetColor(drawingMustParseColor(c, cl.Value()))
			dc.Fill()
			return this
		}).
		Method("save", func(c *Context, this ValueObject, args []Value) Value {
			var filename ValueStr
			EnsureFuncParams(c, "Canvas.save", args,
				ArgRuleRequired("filename", TypeStr, &filename),
			)
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			dc.SavePNG(filename.Value())
			return this
		}).
		Method("show", func(c *Context, this ValueObject, args []Value) Value {
			f, err := os.CreateTemp("", "*.png")
			if err != nil {
				c.RaiseRuntimeError("create temp file error %s", err)
			}
			dc := this.GetMember("__dc", c).ToGoValue().(*gg.Context)
			dc.EncodePNG(f)
			if err := f.Close(); err != nil {
				c.RaiseRuntimeError("close temp file error %s", err)
			}
			if err := exec.Command("open", f.Name()).Run(); err != nil {
				c.RaiseRuntimeError("run open commmand error %s", err)
			}
			return NewStr(f.Name())
		}).
		Build()
)
