package builtin_libs

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	. "github.com/zgg-lang/zgg-go/runtime"

	"github.com/gorilla/websocket"
)

var (
	httpRequestContextClass ValueType
	websocketContextClass   ValueType
	websocketClientClass    ValueType
	httpRequestClass        ValueType
	httpResponseClass       ValueType
	httpFormFileClass       ValueType
)

func libHttp(*Context) ValueObject {
	lib := NewObject()
	// Client
	lib.SetMember("get", httpGet, nil)
	lib.SetMember("getJson", httpGetJson, nil)
	lib.SetMember("postForm", httpPostForm, nil)
	lib.SetMember("postJson", httpPostJson, nil)
	lib.SetMember("Request", httpRequestClass, nil)
	lib.SetMember("escape", NewNativeFunction("escape", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			return NewStr("")
		}
		if as, ok := args[0].(ValueObject); ok {
			rv := url.Values{}
			as.Iterate(func(k string, v Value) {
				if vs, ok := v.(ValueArray); ok {
					for i := 0; i < vs.Len(); i++ {
						item := vs.GetIndex(i, c)
						rv.Add(k, item.ToString(c))
					}
				} else {
					rv.Add(k, v.ToString(c))
				}
			})
			return NewStr(rv.Encode())
		}
		a := args[0].ToString(c)
		return NewStr(url.QueryEscape(a))
	}), nil)
	lib.SetMember("unescape", NewNativeFunction("unescape", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			return NewStr("")
		}
		a := args[0].ToString(c)
		rv, err := url.QueryUnescape(a)
		if err != nil {
			c.OnRuntimeError("unescape %s error: %s", a, err)
		}
		return NewStr(rv)
	}), nil)
	// Server
	lib.SetMember("createServer", httpCreateServer, nil)
	lib.SetMember("serve", NewNativeFunction("serve", func(c *Context, this Value, args []Value) Value {
		if len(args) != 2 {
			c.OnRuntimeError("http: serve requires 2 arguments")
			return nil
		}
		addr := c.MustStr(args[0], "http.serve(addr, handleFunc): addr")
		handleFunc := c.MustCallable(args[1], "http.serve(addr, handleFunc): function")
		http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newC := c.Clone()
			w.Header().Set("server", "zgg simple server")
			rv := NewGoValue(r)
			wv := NewGoValue(w)
			ctx := NewObjectAndInit(httpRequestContextClass, newC, wv, rv)
			newC.Invoke(handleFunc, nil, Args(ctx))
		}))
		return nil
	}), nil)
	return lib
}

var httpGet = NewNativeFunction("http.get", func(c *Context, thisArg Value, args []Value) Value {
	var (
		url     ValueStr
		headers ValueObject
	)
	EnsureFuncParams(c, "http.get", args,
		ArgRuleRequired{"url", TypeStr, &url},
		ArgRuleOptional{"headers", TypeObject, &headers, NewObject()},
	)
	request, err := http.NewRequest("GET", url.Value(), nil)
	headers.Each(func(key string, value Value) bool {
		request.Header.Add(key, value.ToString(c))
		return true
	})
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic("http.get: " + err.Error())
	}
	defer rsp.Body.Close()
	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic("http.get: " + err.Error())
	}
	return NewBytes(bs)
}, "url", "headers")

var httpGetJson = NewNativeFunction("getJson", func(c *Context, thisArg Value, args []Value) Value {
	c.AssertArgNum(len(args), 1, 2, "http.getJson")
	url := c.MustStr(args[0])
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.OnRuntimeError("http.getJson: make request error %s", err)
		return nil
	}
	if len(args) > 1 {
		headers := c.MustObject(args[1], "http.getJson::headers")
		headers.Each(func(key string, value Value) bool {
			request.Header.Add(key, value.ToString(c))
			return true
		})
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		c.OnRuntimeError("http.getJson: " + err.Error())
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	var j interface{}
	if err := json.Unmarshal(respBytes, &j); err != nil {
		c.OnRuntimeError("http.getJson: " + err.Error())
	}
	return jsonToValue(j, c)
}, "url", "headers")

var httpPostForm = NewNativeFunction("postForm", func(c *Context, this Value, args []Value) Value {
	c.AssertArgNum(len(args), 2, 3, "http.postForm")
	postUrl := c.MustStr(args[0], "http.postForm::url")
	formArgs := c.MustObject(args[1], "http.postForm::form")
	form := url.Values{}
	formArgs.Each(func(key string, value Value) bool {
		if values, ok := value.(ValueArray); ok {
			for i := 0; i < values.Len(); i++ {
				elem := values.GetIndex(i, c)
				form.Add(key, elem.ToString(c))
			}
		} else {
			form.Add(key, value.ToString(c))
		}
		return true
	})
	bodyReader := strings.NewReader(form.Encode())
	request, err := http.NewRequest("POST", postUrl, bodyReader)
	if err != nil {
		c.OnRuntimeError("http.postForm: make request error %s", err)
		return nil
	}
	if len(args) > 2 {
		headers := c.MustObject(args[2], "http.postForm::headers")
		headers.Each(func(key string, value Value) bool {
			request.Header.Add(key, value.ToString(c))
			return true
		})
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		c.OnRuntimeError("http.postForm: request error %s", err)
		return nil
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.OnRuntimeError("http.postForm: read response error %s", err)
		return nil
	}
	return NewBytes(content)
}, "url", "form", "headers")

var httpPostJson = NewNativeFunction("postJson", func(c *Context, this Value, args []Value) Value {
	c.AssertArgNum(len(args), 2, 3, "http.postJson")
	postUrl := c.MustStr(args[0], "http.postJson::url")
	content := args[1]
	contentBytes, err := json.Marshal(content.ToGoValue())
	if err != nil {
		c.OnRuntimeError("http.postJson::content encode to json error %s", err)
		return nil
	}
	bodyReader := bytes.NewReader(contentBytes)
	request, err := http.NewRequest("POST", postUrl, bodyReader)
	if err != nil {
		c.OnRuntimeError("http.postJson: make request error %s", err)
		return nil
	}
	if len(args) > 2 {
		headers := c.MustObject(args[2], "http.postJson::headers")
		headers.Each(func(key string, value Value) bool {
			request.Header.Add(key, value.ToString(c))
			return true
		})
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		c.OnRuntimeError("http.postJson: request error %s", err)
		return nil
	}
	defer resp.Body.Close()
	rspContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.OnRuntimeError("http.postJson: read response error %s", err)
		return nil
	}
	return NewBytes(rspContent)
}, "url", "content", "headers")

// type httpRequestContext struct {
// 	r *http.Request
// 	w http.ResponseWriter
// 	Method string
// }

// func httpNewRequestContext(w http.ResponseWriter, r *http.Request) *httpRequestContext {
// 	return &httpRequestContext{
// 		r: r,
// 		w: w,
// 	}
// }

// func (c *httpRequestContext) write

var httpCreateServer = NewNativeFunction("createServer", func(c *Context, thisArg Value, args []Value) Value {
	serverName := "zgg http server"
	if len(args) > 0 {
		opts := c.MustObject(args[0], "http.createServer options")
		if _serverName, ok := opts.GetMember("serverName", c).(ValueStr); ok {
			serverName = _serverName.Value()
		}
	}
	svr := http.NewServeMux()
	rv := NewObject()
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  65536,
		WriteBufferSize: 65536,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rv.SetMember("route", NewNativeFunction("route", func(c *Context, thisArg Value, args []Value) Value {
		if len(args) < 2 {
			c.OnRuntimeError("http.server: route requires at least 2 arguments")
			return nil
		}
		path := c.MustStr(args[0], "http.server.route(path, handleFunc): path")
		handleFunc := c.MustCallable(args[1], "http.server.route(path, handleFunc): function")
		svr.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			newC := c.Clone()
			w.Header().Set("server", serverName)
			rv := NewGoValue(r)
			wv := NewGoValue(w)
			ctx := NewObjectAndInit(httpRequestContextClass, newC, wv, rv)
			newC.Invoke(handleFunc, nil, Args(ctx))
			if callable, ok := ctx.GetMember("close", c).(ValueCallable); ok {
				c.Invoke(callable, ctx, NoArgs)
			}
		})
		return thisArg
	}, "path", "handleFunc"), nil)
	rv.SetMember("routeWebsocket", NewNativeFunction("routeWebsocket", func(c *Context, thisArgs Value, args []Value) Value {
		var (
			routePath  ValueStr
			handleFunc ValueCallable
		)
		EnsureFuncParams(c, "routeWebsocket", args,
			ArgRuleRequired{"path", TypeStr, &routePath},
			ArgRuleRequired{"handleFunc", TypeFunc, &handleFunc},
		)
		svr.HandleFunc(routePath.Value(), func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			defer conn.Close()
			newC := c.Clone()
			w.Header().Set("server", serverName)
			ctx := NewObjectAndInit(websocketContextClass, newC, NewGoValue(conn), NewGoValue(w), NewGoValue(r))
			newC.Invoke(handleFunc, nil, Args(ctx))
		})
		return thisArg
	}, "path", "handleFunc"), nil)
	rv.SetMember("serve", NewNativeFunction("serve", func(c *Context, thisArgs Value, args []Value) Value {
		if len(args) != 1 {
			c.OnRuntimeError("http.server: serve requires 1 argument")
			return nil
		}
		addr := c.MustStr(args[0], "http.server.serve(listenAddr): listenAddr")
		if err := http.ListenAndServe(addr, svr); err != nil {
			c.OnRuntimeError("http.server.serve fail: " + err.Error())
			return nil
		}
		return Undefined()
	}), nil)
	return rv
})

func initHttpRequestContextClass() ValueType {
	className := "http.RequestContext"
	return NewClassBuilder("RequestContext").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			// w := args[0].ToGoValue().(http.ResponseWriter)
			r := args[1].ToGoValue().(*http.Request)
			this.SetMember("_w", args[0], c)
			this.SetMember("_r", args[1], c)
			this.SetMember("method", NewStr(r.Method), c)
			this.SetMember("path", NewStr(r.URL.Path), c)
			this.SetMember("querystr", NewStr(r.URL.RawQuery), c)
		}).
		Method("getBody", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				c.OnRuntimeError("%s.getData error %s", className, err)
				return nil
			}
			return NewBytes(body)
		}).
		Method("query", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var argName ValueStr
			EnsureFuncParams(c, className+".query", args,
				ArgRuleRequired{"name", TypeStr, &argName},
			)
			if arg, found := r.URL.Query()[argName.Value()]; found && len(arg) > 0 {
				return NewStr(arg[0])
			}
			return Nil()
		}).
		Method("getHeader", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var argName ValueStr
			EnsureFuncParams(c, className+".getHeader", args,
				ArgRuleRequired{"name", TypeStr, &argName},
			)
			return NewStr(r.Header.Get(argName.Value()))
		}).
		Method("getHeaders", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var argName ValueStr
			EnsureFuncParams(c, className+".getHeaders", args,
				ArgRuleRequired{"name", TypeStr, &argName},
			)
			headers := r.Header.Values(argName.Value())
			rv := NewArray(len(headers))
			for _, h := range headers {
				rv.PushBack(NewStr(h))
			}
			return rv
		}).
		Method("queryAll", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var argName ValueStr
			EnsureFuncParams(c, className+".queryAll", args,
				ArgRuleRequired{"name", TypeStr, &argName},
			)
			if arg, found := r.URL.Query()[argName.Value()]; found {
				rv := NewArray(len(arg))
				for _, a := range arg {
					rv.PushBack(NewStr(a))
				}
				return rv
			}
			return NewArray()
		}).
		Method("parseMultipartForm", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var maxSize ValueInt
			EnsureFuncParams(c, className+".parseMultipartForm", args,
				ArgRuleOptional{"maxSize", TypeInt, &maxSize, NewInt(10 * 1024 * 1024)},
			)
			r.ParseMultipartForm(maxSize.Value())
			return Nil()
		}).
		Method("file", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var (
				name  ValueStr
				index ValueInt
			)
			EnsureFuncParams(c, className+".file", args,
				ArgRuleRequired{"name", TypeStr, &name},
				ArgRuleOptional{"index", TypeInt, &index, NewInt(0)},
			)
			if r.MultipartForm == nil {
				c.OnRuntimeError("RequestContext.file: cannot get file without parseMultipartForm")
			}
			fhs := r.MultipartForm.File[name.Value()]
			if len(fhs) <= index.AsInt() {
				return Nil()
			}
			rv := NewObjectAndInit(httpFormFileClass, c, NewGoValue(fhs[index.AsInt()]))
			return rv
		}).
		Method("files", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var (
				name  ValueStr
				index ValueInt
			)
			EnsureFuncParams(c, className+".file", args,
				ArgRuleRequired{"name", TypeStr, &name},
				ArgRuleOptional{"index", TypeInt, &index, NewInt(0)},
			)
			if r.MultipartForm == nil {
				c.OnRuntimeError("RequestContext.file: cannot get file without parseMultipartForm")
			}
			fhs := r.MultipartForm.File[name.Value()]
			rv := NewArray(len(fhs))
			for _, fh := range fhs {
				rfh := NewObjectAndInit(httpFormFileClass, c, NewGoValue(fh))
				rv.PushBack(rfh)
			}
			return rv
		}).
		Method("form", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			r.ParseForm()
			var argName ValueStr
			EnsureFuncParams(c, className+".form", args,
				ArgRuleRequired{"name", TypeStr, &argName},
			)
			if arg, found := r.Form[argName.Value()]; found && len(arg) > 0 {
				return NewStr(arg[0])
			}
			return Nil()
		}).
		Method("formAll", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var argName ValueStr
			EnsureFuncParams(c, className+".formAll", args,
				ArgRuleRequired{"name", TypeStr, &argName},
			)
			if arg, found := r.Form[argName.Value()]; found {
				rv := NewArray(len(arg))
				for _, a := range arg {
					rv.PushBack(NewStr(a))
				}
				return rv
			}
			return NewArray()
		}).
		Method("getBodyStr", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				c.OnRuntimeError("%s.getData error %s", className, err)
				return nil
			}
			return NewStr(string(body))
		}).
		Method("addHeader", func(c *Context, this ValueObject, args []Value) Value {
			var (
				key ValueStr
				val ValueStr
			)
			EnsureFuncParams(c, "http.RequestContext.addHeader", args,
				ArgRuleRequired{"key", TypeStr, &key},
				ArgRuleRequired{"val", TypeStr, &val},
			)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			w.Header().Add(key.Value(), val.Value())
			return Undefined()
		}).
		Method("setHeader", func(c *Context, this ValueObject, args []Value) Value {
			var (
				key ValueStr
				val ValueStr
			)
			EnsureFuncParams(c, "http.RequestContext.setHeader", args,
				ArgRuleRequired{"key", TypeStr, &key},
				ArgRuleRequired{"val", TypeStr, &val},
			)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			w.Header().Set(key.Value(), val.Value())
			return Undefined()
		}).
		Method("write", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) < 1 {
				c.OnRuntimeError("http.RequestContext.write requires at least one argument")
				return nil
			}
			statusCode := c.MustInt(args[0])
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			w.WriteHeader(int(statusCode))
			for i := 1; i < len(args); i++ {
				arg := args[i]
				if reader, ok := arg.ToGoValue().(io.Reader); ok {
					io.Copy(w, reader)
					continue
				}
				switch argVal := arg.(type) {
				case ValueBytes:
					{
						bs := argVal.Value()
						written := 0
						for written < len(bs) {
							n, err := w.Write(bs[written:])
							if err != nil {
								break
							}
							written += n
						}
					}
				default:
					io.WriteString(w, arg.ToString(c))
				}
			}
			return Undefined()
		}).
		Method("writeJson", func(c *Context, this ValueObject, args []Value) Value {
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			statusCode := 200
			contentAt := 0
			switch len(args) {
			case 2:
				statusCode = int(c.MustInt(args[0]))
				contentAt++
				fallthrough
			case 1:
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(statusCode)
				if err := json.NewEncoder(w).Encode(args[contentAt].ToGoValue()); err != nil {
					c.OnRuntimeError("writeJson: encode to json error %s", err)
				}
			default:
				c.OnRuntimeError("writeJson usage: writeJson([statusCode,] contentValue)")
			}
			return Undefined()
		}).
		Method("sendFile", func(c *Context, this ValueObject, args []Value) Value {
			var (
				filename ValueStr
			)
			EnsureFuncParams(c, "http.RequestContext.sendFile", args,
				ArgRuleRequired{"filename", TypeStr, &filename},
			)
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			http.ServeFile(w, r, filename.Value())
			return Undefined()
		}).
		Method("redirect", func(c *Context, this ValueObject, args []Value) Value {
			var (
				url  ValueStr
				code ValueInt
			)
			switch len(args) {
			case 1:
				EnsureFuncParams(c, "http.RequestContext.redirect", args,
					ArgRuleRequired{"url", TypeStr, &url},
				)
				code = NewInt(302)
			default:
				EnsureFuncParams(c, "http.RequestContext.redirect", args,
					ArgRuleRequired{"url", TypeStr, &url},
					ArgRuleRequired{"code", TypeInt, &code},
				)
			}
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			http.Redirect(w, r, url.Value(), code.AsInt())
			return Undefined()
		}).
		Method("forward", func(c *Context, this ValueObject, args []Value) Value {
			var (
				target  ValueStr
				options ValueObject
			)
			EnsureFuncParams(c, "http.RequestContext.forward", args,
				ArgRuleRequired{"target", TypeStr, &target},
				ArgRuleOptional{"options", TypeObject, &options, NewObject()},
			)
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			var (
				targetUrl = target.Value()
				host      = targetUrl
			)
			if sp := strings.Index(targetUrl, "://"); sp >= 0 {
				host = host[sp+3:]
			} else {
				targetUrl = r.URL.Scheme + "://" + targetUrl
			}
			targetUrl += r.URL.Path
			if r.URL.RawQuery != "" {
				targetUrl += "?" + r.URL.RawQuery
			}
			reqBody := r.Body
			if reqBody != nil {
				defer reqBody.Close()
			}
			request, err := http.NewRequest(r.Method, targetUrl, reqBody)
			if err != nil {
				c.OnRuntimeError("make forward request error: %s", err)
			}
			for field, values := range r.Header {
				if strings.ToLower(field) == "host" {
					request.Header.Set(field, host)
				} else {
					for _, value := range values {
						request.Header.Add(field, value)
					}
				}
			}
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				c.OnRuntimeError("do forward request error: %s", err)
			}
			defer response.Body.Close()
			respHeader := w.Header()
			for field, values := range response.Header {
				for i, value := range values {
					if i == 0 {
						respHeader.Set(field, value)
					} else {
						respHeader.Add(field, value)
					}
				}
			}
			w.WriteHeader(response.StatusCode)
			io.Copy(w, response.Body)
			return Undefined()
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			r.Body.Close()
			return Undefined()
		}).
		Build()
}

func initWebsocketContextClass() ValueType {
	// className := "http.WebsocketContext"
	return NewClassBuilder("WebsocketContext").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			this.SetMember("_conn", args[0], c)
			this.SetMember("w", args[1], c)
			this.SetMember("r", args[2], c)
		}).
		Method("read", func(c *Context, this ValueObject, args []Value) Value {
			conn := this.GetMember("_conn", c).ToGoValue().(*websocket.Conn)
			mt, pkg, err := conn.ReadMessage()
			if err != nil {
				c.OnRuntimeError("websocket read message error %s", err)
			}
			switch mt {
			case websocket.TextMessage:
				return NewStr(string(pkg))
			case websocket.BinaryMessage:
				return NewBytes(pkg)
			}
			return nil
		}).
		Method("write", func(c *Context, this ValueObject, args []Value) Value {
			conn := this.GetMember("_conn", c).ToGoValue().(*websocket.Conn)
			var err error
			for _, arg := range args {
				switch pkg := arg.(type) {
				case ValueBytes:
					err = conn.WriteMessage(websocket.BinaryMessage, pkg.Value())
				case ValueObject:
					err = conn.WriteJSON(pkg.ToGoValue())
				default:
					err = conn.WriteMessage(websocket.TextMessage, []byte(pkg.ToString(c)))
				}
				if err != nil {
					c.OnRuntimeError("websocket write message error %s", err)
				}
			}
			return Undefined()
		}).
		Build()

}

func initHttpRequestClass() ValueType {
	return NewClassBuilder("Request").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				method ValueStr
				url    ValueStr
			)
			if len(args) == 1 {
				EnsureFuncParams(c, "Request.__init__", args,
					ArgRuleRequired{"url", TypeStr, &url},
				)
				this.SetMember("method", NewStr("GET"), c)
			} else {
				EnsureFuncParams(c, "Request.__init__", args,
					ArgRuleRequired{"method", TypeStr, &method},
					ArgRuleRequired{"url", TypeStr, &url},
				)
				this.SetMember("method", method, c)
			}
			this.SetMember("url", url, c)
			this.SetMember("__headers", NewArray(), c)
			this.SetMember("__body", NewArray(), c)
		}).
		Method("data", func(c *Context, this ValueObject, args []Value) Value {
			body := this.GetMember("__body", c).(ValueArray)
			for _, arg := range args {
				body.PushBack(arg)
			}
			return this
		}).
		Method("form", func(c *Context, this ValueObject, args []Value) Value {
			var formArg ValueObject
			EnsureFuncParams(c, "Request.form", args, ArgRuleRequired{"form", TypeObject, &formArg})
			headers := this.GetMember("__headers", c).(ValueArray)
			headers.PushBack(NewArrayByValues(NewStr("Content-Type"), NewStr("application/x-www-form-urlencoded")))
			body := this.GetMember("__body", c).(ValueArray)
			form := url.Values{}
			formArg.Iterate(func(key string, val Value) {
				form.Set(key, val.ToString(c))
			})
			body.PushBack(NewStr(form.Encode()))
			return this
		}).
		Method("header", func(c *Context, this ValueObject, args []Value) Value {
			headers := this.GetMember("__headers", c).(ValueArray)
			switch len(args) {
			case 1:
				if argHeaders, ok := args[0].(ValueObject); ok {
					argHeaders.Iterate(func(key string, val Value) {
						headers.PushBack(NewArrayByValues(NewStr(key), val))
					})
					return this
				}
			case 2:
				headers.PushBack(NewArrayByValues(args[0], args[1]))
				return this
			}
			c.OnRuntimeError("Request.header: invalid argument(s)")
			return nil
		}).
		Method("call", func(c *Context, this ValueObject, args []Value) Value {
			body := this.GetMember("__body", c).(ValueArray)
			var reqBody io.Reader
			if bodyLen := body.Len(); bodyLen == 0 {
				reqBody = nil
			} else {
				readers := make([]io.Reader, bodyLen)
				for i := range readers {
					v := body.GetIndex(i, c)
					switch bodyVal := v.(type) {
					case ValueBytes:
						readers[i] = bytes.NewReader(bodyVal.Value())
					default:
						readers[i] = strings.NewReader(bodyVal.ToString(c))
					}
				}
				reqBody = io.MultiReader(readers...)
			}
			method := c.MustStr(this.GetMember("method", c))
			url := c.MustStr(this.GetMember("url", c))
			req, err := http.NewRequest(method, url, reqBody)
			if err != nil {
				c.OnRuntimeError("Request.build: make request error %s", err)
			}
			headers := this.GetMember("__headers", c).(ValueArray)
			for i := 0; i < headers.Len(); i++ {
				item := headers.GetIndex(i, c)
				key := item.GetIndex(0, c)
				val := item.GetIndex(1, c)
				req.Header.Add(key.ToString(c), val.ToString(c))
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				c.OnRuntimeError("http.Request.call: do request error %s", err)
			}
			return NewObjectAndInit(httpResponseClass, c, NewGoValue(resp))
			// return NewGoValue(resp)
		}).
		Build()
}

func initHttpResponseClass() ValueType {
	return NewClassBuilder("Response").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			resp := args[0].ToGoValue().(*http.Response)
			this.SetMember("__resp", args[0], c)
			this.SetMember("statusCode", NewInt(int64(resp.StatusCode)), c)
			headers := NewObject()
			for k := range resp.Header {
				headers.SetMember(k, NewStr(resp.Header.Get(k)), c)
			}
			this.SetMember("headers", headers, c)
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			resp.Body.Close()
			return this
		}).
		Method("header", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			var h ValueStr
			EnsureFuncParams(c, "http.Response.header", args, ArgRuleRequired{"header", TypeStr, &h})
			return NewStr(resp.Header.Get(h.Value()))
		}).
		Method("bytes", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.OnRuntimeError("http.Response.bytes: read body error %s", err)
			}
			return NewBytes(bytes)
		}).
		Method("chunk", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			var buf [512 * 1024]byte
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				n, err := resp.Body.Read(buf[:])
				if err == io.EOF {
					return NewArrayByValues(NewBytes(buf[:n]), NewBool(false))
				}
				if err != nil {
					c.OnRuntimeError("Read chunk error: %s", err)
				}
				return NewArrayByValues(NewBytes(buf[:n]), NewBool(true))
			})
		}).
		Method("text", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.OnRuntimeError("http.Response.text: read body error %s", err)
			}
			return NewStr(string(bytes))
		}).
		Method("json", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			var o interface{}
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&o); err != nil {
				c.OnRuntimeError("http.Response.json: decode body error %s", err)
			}
			return jsonToValue(o, c)
		}).
		Build()
}

func initHttpFormFileClass() ValueType {
	return NewClassBuilder("FormFile").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			fh := args[0].ToGoValue().(*multipart.FileHeader)
			this.SetMember("_fh", args[0], c)
			this.SetMember("name", NewStr(fh.Filename), c)
			this.SetMember("size", NewInt(fh.Size), c)
		}).
		Method("bytes", func(c *Context, this ValueObject, args []Value) Value {
			_bs := this.GetMember("_bytes", c)
			if bs, ok := _bs.(ValueBytes); ok {
				return bs
			}
			fh := this.GetMember("_fh", c).ToGoValue().(*multipart.FileHeader)
			file, err := fh.Open()
			if err != nil {
				c.OnRuntimeError("FormFile.bytes: open file error %s", err)
			}
			defer file.Close()
			fileBs, err := ioutil.ReadAll(file)
			if err != nil {
				c.OnRuntimeError("FormFile.bytes: read file error %s", err)
			}
			bs := NewBytes(fileBs)
			this.SetMember("_bytes", bs, c)
			return bs
		}).
		Build()
}

// func initHttpWebsocketClientClass() ValueType {
// 	return NewClassBuilder("WebsocketClient").
// 		Constructor(func(c *Context, this ValueObject, args []Value) {
// 			var (
// 				url ValueStr
// 			)
// 			EnsureFuncParams(c, "http.WebSocket")
// 			websocket.Dial()
// 		}).
// 		Build()
// }
func init() {
	httpRequestContextClass = initHttpRequestContextClass()
	websocketContextClass = initWebsocketContextClass()
	httpRequestClass = initHttpRequestClass()
	httpResponseClass = initHttpResponseClass()
	httpFormFileClass = initHttpFormFileClass()
}
