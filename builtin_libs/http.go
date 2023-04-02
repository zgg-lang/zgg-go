package builtin_libs

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

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

const (
	httpUnixPrefix = "unix://"
)

func libHttp(*Context) ValueObject {
	lib := NewObject()
	// Client
	lib.SetMember("get", httpGet, nil)
	lib.SetMember("getJson", httpGetJson, nil)
	lib.SetMember("postForm", httpPostForm, nil)
	lib.SetMember("postMultipartForm", httpPostMultipartForm, nil)
	lib.SetMember("postJson", httpPostJson, nil)
	lib.SetMember("Request", httpRequestClass, nil)
	lib.SetMember("WebsocketClient", websocketClientClass, nil)
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
			c.RaiseRuntimeError("unescape %s error: %s", a, err)
		}
		return NewStr(rv)
	}), nil)
	// Server
	lib.SetMember("createServer", httpCreateServer, nil)
	lib.SetMember("serve", NewNativeFunction("serve", func(c *Context, this Value, args []Value) Value {
		if len(args) != 2 {
			c.RaiseRuntimeError("http: serve requires 2 arguments")
			return nil
		}
		addr := c.MustStr(args[0], "http.serve(addr, handleFunc): addr")
		handleFunc := c.MustCallable(args[1], "http.serve(addr, handleFunc): function")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newC := c.Clone()
			w.Header().Set("server", "zgg simple server")
			rv := NewGoValue(r)
			wv := NewGoValue(w)
			ctx := NewObjectAndInit(httpRequestContextClass, newC, wv, rv)
			newC.Invoke(handleFunc, nil, Args(ctx))
		})
		if strings.HasPrefix(addr, httpUnixPrefix) {
			unixAddr, err := net.ResolveUnixAddr("unix", addr[len(httpUnixPrefix):])
			if err != nil {
				c.RaiseRuntimeError("resolve unix addr %s error %s", addr, err)
			}
			listener, err := net.ListenUnix("unix", unixAddr)
			if err != nil {
				c.RaiseRuntimeError("listen %s error %s", addr, err)
			}
			defer listener.Close()
			listener.SetUnlinkOnClose(true)
			if err := http.Serve(listener, handler); err != nil {
				c.RaiseRuntimeError("http serve on %s error %s", addr, err)
			}
		} else {
			if err := http.ListenAndServe(addr, handler); err != nil {
				c.RaiseRuntimeError("http serve on %s error %s", addr, err)
			}
		}
		return nil
	}), nil)
	// 命令行功能
	lib.SetMember("static", NewNativeFunction("static", func(c *Context, this Value, args []Value) Value {
		addr := c.Args[0]
		fs := http.FileServer(http.Dir("."))
		http.Handle("/", fs)
		http.ListenAndServe(addr, nil)
		return Undefined()
	}), nil)
	return lib
}

var httpGet = NewNativeFunction("http.get", func(c *Context, thisArg Value, args []Value) Value {
	var (
		url     ValueStr
		headers ValueObject
	)
	EnsureFuncParams(c, "http.get", args,
		ArgRuleRequired("url", TypeStr, &url),
		ArgRuleOptional("headers", TypeObject, &headers, NewObject()),
	)
	request, err := http.NewRequest("GET", url.Value(), nil)
	if err != nil {
		c.RaiseRuntimeError("http.get: make reqeust failed %s", err)
	}
	headers.Each(func(key string, value Value) bool {
		request.Header.Add(key, value.ToString(c))
		return true
	})
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		c.RaiseRuntimeError("http.get: request failed %s", err)
	}
	defer rsp.Body.Close()
	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		c.RaiseRuntimeError("http.get: read response failed %s", err)
	}
	return NewBytes(bs)
}, "url", "headers")

var httpGetJson = NewNativeFunction("getJson", func(c *Context, thisArg Value, args []Value) Value {
	c.AssertArgNum(len(args), 1, 2, "http.getJson")
	url := c.MustStr(args[0])
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.RaiseRuntimeError("http.getJson: make request error %s", err)
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
		c.RaiseRuntimeError("http.getJson: " + err.Error())
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	var j interface{}
	if err := json.Unmarshal(respBytes, &j); err != nil {
		c.RaiseRuntimeError("http.getJson: " + err.Error())
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
		c.RaiseRuntimeError("http.postForm: make request error %s", err)
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
		c.RaiseRuntimeError("http.postForm: request error %s", err)
		return nil
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.RaiseRuntimeError("http.postForm: read response error %s", err)
		return nil
	}
	return NewBytes(content)
}, "url", "form", "headers")

func httpGetMultipartForm(c *Context, form ValueObject) ([]byte, string) {
	var buf bytes.Buffer
	formWriter := multipart.NewWriter(&buf)
	form.Iterate(func(key string, val Value) {
		if rd, ok := val.ToGoValue().(io.Reader); ok {
			filename := ""
			if file, ok := rd.(*os.File); ok {
				filename = filepath.Base(file.Name())
			}
			if field, err := formWriter.CreateFormFile(key, filename); err == nil {
				io.Copy(field, rd)
			} else {
				c.RaiseRuntimeError("create form file %s error %s", key, err)
			}
		} else if callable, ok := c.GetCallable(val); ok {
			c.Invoke(callable, nil, NoArgs)
			var (
				file    = c.RetVal
				content []byte
				name    = ""
			)
			if arr, ok := file.(ValueArray); ok {
				if arr.Len() > 1 {
					name = arr.GetIndex(1, c).ToString(c)
				}
				file = arr.GetIndex(0, c)
			}
			if bs, ok := file.(ValueBytes); ok {
				content = bs.Value()
			} else {
				content = []byte(file.ToString(c))
			}
			if field, err := formWriter.CreateFormFile(key, name); err == nil {
				field.Write(content)
			} else {
				c.RaiseRuntimeError("create form file %s error %s", key, err)
			}
		} else {
			if field, err := formWriter.CreateFormField(key); err == nil {
				if bs, ok := val.(ValueBytes); ok {
					field.Write(bs.Value())
				} else {
					field.Write([]byte(val.ToString(c)))
				}
			}
		}
	})
	formWriter.Close()
	return buf.Bytes(), formWriter.FormDataContentType()
}

var httpPostMultipartForm = NewNativeFunction("postMultipartForm", func(c *Context, this Value, args []Value) Value {
	var (
		url     ValueStr
		form    ValueObject
		headers ValueObject
	)
	EnsureFuncParams(c, "http.postMultiPartForm", args,
		ArgRuleRequired("url", TypeStr, &url),
		ArgRuleRequired("form", TypeObject, &form),
		ArgRuleOptional("headers", TypeObject, &headers, NewObject()),
	)
	formBs, contentType := httpGetMultipartForm(c, form)
	request, err := http.NewRequest("POST", url.Value(), bytes.NewReader(formBs))
	if err != nil {
		c.RaiseRuntimeError("postMultipartForm url %s new request error %s", url.Value(), err)
	}
	request.Header.Set("Content-Type", contentType)
	headers.Iterate(func(key string, value Value) {
		request.Header.Add(key, value.ToString(c))
	})
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		c.RaiseRuntimeError("postMultipartForm url %s do request error %s", url.Value(), err)
	}
	defer resp.Body.Close()
	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		c.RaiseRuntimeError("postMultipartForm read response error %s", err)
	}
	return NewBytes(ret)
}, "url", "form", "headers")

var httpPostJson = NewNativeFunction("postJson", func(c *Context, this Value, args []Value) Value {
	c.AssertArgNum(len(args), 2, 3, "http.postJson")
	postUrl := c.MustStr(args[0], "http.postJson::url")
	content := args[1]
	contentBytes, err := json.Marshal(content.ToGoValue())
	if err != nil {
		c.RaiseRuntimeError("http.postJson::content encode to json error %s", err)
		return nil
	}
	bodyReader := bytes.NewReader(contentBytes)
	request, err := http.NewRequest("POST", postUrl, bodyReader)
	if err != nil {
		c.RaiseRuntimeError("http.postJson: make request error %s", err)
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
		c.RaiseRuntimeError("http.postJson: request error %s", err)
		return nil
	}
	defer resp.Body.Close()
	rspContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.RaiseRuntimeError("http.postJson: read response error %s", err)
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
			c.RaiseRuntimeError("http.server: route requires at least 2 arguments")
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
			if callable, ok := newC.GetCallable(ctx.GetMember("close", c)); ok {
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
			ArgRuleRequired("path", TypeStr, &routePath),
			ArgRuleRequired("handleFunc", TypeFunc, &handleFunc),
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
		var (
			addrStr ValueStr
			options ValueObject
		)
		EnsureFuncParams(c, "serve", args,
			ArgRuleRequired("address", TypeStr, &addrStr),
			ArgRuleOptional("options", TypeObject, &options, NewObject()),
		)
		addr := addrStr.Value()
		if strings.HasPrefix(addr, httpUnixPrefix) {
			unixAddr, err := net.ResolveUnixAddr("unix", addr[len(httpUnixPrefix):])
			if err != nil {
				c.RaiseRuntimeError("resolve unix addr %s error %s", addr, err)
			}
			listener, err := net.ListenUnix("unix", unixAddr)
			if err != nil {
				c.RaiseRuntimeError("listen %s error %s", addr, err)
			}
			defer listener.Close()
			listener.SetUnlinkOnClose(true)
			if err := http.Serve(listener, svr); err != nil {
				c.RaiseRuntimeError("http serve on %s error %s", addr, err)
			}
		} else {
			var certFile, keyFile string
			useTls := false
			if cf, ok := options.GetMember("certFile", c).(ValueStr); ok {
				if kf, ok := options.GetMember("keyFile", c).(ValueStr); ok {
					useTls = true
					certFile = cf.Value()
					keyFile = kf.Value()
				}
			}
			var err error
			if useTls {
				err = http.ListenAndServeTLS(addr, certFile, keyFile, svr)
			} else {
				err = http.ListenAndServe(addr, svr)
			}
			if err != nil {
				c.RaiseRuntimeError("http.server.serve fail: %s", err)
				return nil
			}
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
				c.RaiseRuntimeError("%s.getData error %s", className, err)
				return nil
			}
			return NewBytes(body)
		}).
		Method("query", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var argName ValueStr
			EnsureFuncParams(c, className+".query", args,
				ArgRuleRequired("name", TypeStr, &argName),
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
				ArgRuleRequired("name", TypeStr, &argName),
			)
			return NewStr(r.Header.Get(argName.Value()))
		}).
		Method("getHeaders", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			var argName ValueStr
			EnsureFuncParams(c, className+".getHeaders", args,
				ArgRuleRequired("name", TypeStr, &argName),
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
				ArgRuleRequired("name", TypeStr, &argName),
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
				ArgRuleOptional("maxSize", TypeInt, &maxSize, NewInt(10*1024*1024)),
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
				ArgRuleRequired("name", TypeStr, &name),
				ArgRuleOptional("index", TypeInt, &index, NewInt(0)),
			)
			if r.MultipartForm == nil {
				c.RaiseRuntimeError("RequestContext.file: cannot get file without parseMultipartForm")
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
				ArgRuleRequired("name", TypeStr, &name),
				ArgRuleOptional("index", TypeInt, &index, NewInt(0)),
			)
			if r.MultipartForm == nil {
				c.RaiseRuntimeError("RequestContext.file: cannot get file without parseMultipartForm")
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
				ArgRuleRequired("name", TypeStr, &argName),
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
				ArgRuleRequired("name", TypeStr, &argName),
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
				c.RaiseRuntimeError("%s.getData error %s", className, err)
				return nil
			}
			return NewStr(string(body))
		}).
		Method("cookie", func(c *Context, this ValueObject, args []Value) Value {
			var (
				name ValueStr
			)
			EnsureFuncParams(c, className+".setCookie", args,
				ArgRuleRequired("name", TypeStr, &name),
			)
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			cookie, err := r.Cookie(name.Value())
			if err == http.ErrNoCookie {
				return Nil()
			} else if err != nil {
				c.RaiseRuntimeError("get cookie error %+v", err)
			}
			return NewStr(cookie.Value)
		}).
		Method("cookies", func(c *Context, this ValueObject, args []Value) Value {
			r := this.GetMember("_r", c).ToGoValue().(*http.Request)
			cookies := r.Cookies()
			rv := NewObject()
			for _, cookie := range cookies {
				rv.SetMember(cookie.Name, NewStr(cookie.Value), c)
			}
			return rv
		}).
		Method("addHeader", func(c *Context, this ValueObject, args []Value) Value {
			var (
				key ValueStr
				val ValueStr
			)
			EnsureFuncParams(c, "http.RequestContext.addHeader", args,
				ArgRuleRequired("key", TypeStr, &key),
				ArgRuleRequired("val", TypeStr, &val),
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
				ArgRuleRequired("key", TypeStr, &key),
				ArgRuleRequired("val", TypeStr, &val),
			)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			w.Header().Set(key.Value(), val.Value())
			return Undefined()
		}).
		Method("setCookie", func(c *Context, this ValueObject, args []Value) Value {
			var (
				name    ValueStr
				value   ValueStr
				options ValueObject
			)
			EnsureFuncParams(c, className+".setCookie", args,
				ArgRuleRequired("name", TypeStr, &name),
				ArgRuleRequired("value", TypeStr, &value),
				ArgRuleOptional("options", TypeObject, &options, NewObject()),
			)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			var cookie http.Cookie
			cookie.Name = name.Value()
			cookie.Value = value.Value()
			if v, ok := options.GetMember("path", c).(ValueStr); ok {
				cookie.Path = v.Value()
			}
			if v, ok := options.GetMember("domain", c).(ValueStr); ok {
				cookie.Domain = v.Value()
			}
			if v, ok := options.GetMember("maxAge", c).(ValueInt); ok {
				cookie.MaxAge = v.AsInt()
			}
			if v, ok := options.GetMember("secure", c).(ValueBool); ok {
				cookie.Secure = v.Value()
			}
			if v, ok := options.GetMember("httpOnly", c).(ValueBool); ok {
				cookie.HttpOnly = v.Value()
			}
			http.SetCookie(w, &cookie)
			return Undefined()
		}).
		Method("delCookie", func(c *Context, this ValueObject, args []Value) Value {
			var (
				name ValueStr
			)
			EnsureFuncParams(c, className+".setCookie", args,
				ArgRuleRequired("name", TypeStr, &name),
			)
			w := this.GetMember("_w", c).ToGoValue().(http.ResponseWriter)
			http.SetCookie(w, &http.Cookie{Name: name.Value(), MaxAge: -1})
			return Undefined()
		}).
		Method("write", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) < 1 {
				c.RaiseRuntimeError("http.RequestContext.write requires at least one argument")
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
					c.RaiseRuntimeError("writeJson: encode to json error %s", err)
				}
			default:
				c.RaiseRuntimeError("writeJson usage: writeJson([statusCode,] contentValue)")
			}
			return Undefined()
		}).
		Method("sendFile", func(c *Context, this ValueObject, args []Value) Value {
			var (
				filename ValueStr
			)
			EnsureFuncParams(c, "http.RequestContext.sendFile", args,
				ArgRuleRequired("filename", TypeStr, &filename),
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
			EnsureFuncParams(c, "http.RequestContext.redirect", args,
				ArgRuleRequired("url", TypeStr, &url),
				ArgRuleOptional("code", TypeInt, &code, NewInt(200)),
			)
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
				ArgRuleRequired("target", TypeStr, &target),
				ArgRuleOptional("options", TypeObject, &options, NewObject()),
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
			reqPath := r.URL.Path
			if r.URL.RawQuery != "" {
				reqPath += "?" + r.URL.RawQuery
			}
			if rewriteFunc := options.GetMember("rewrite", c); c.IsCallable(rewriteFunc) {
				c.Invoke(rewriteFunc, nil, Args(NewStr(reqPath)))
				reqPath = c.RetVal.ToString(c)
			}
			targetUrl += reqPath
			reqBody := r.Body
			if reqBody != nil {
				defer reqBody.Close()
			}
			request, err := http.NewRequest(r.Method, targetUrl, reqBody)
			if err != nil {
				c.RaiseRuntimeError("make forward request error: %s", err)
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
			request.Host = r.Host
			var client http.Client
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
			response, err := client.Do(request)
			if err != nil {
				c.RaiseRuntimeError("do forward request error: %s", err)
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
				c.RaiseRuntimeError("websocket read message error %s", err)
			}
			switch mt {
			case websocket.TextMessage:
				return NewStr(string(pkg))
			case websocket.BinaryMessage:
				return NewBytes(pkg)
			}
			return nil
		}).
		Method("readJson", func(c *Context, this ValueObject, args []Value) Value {
			conn := this.GetMember("_conn", c).ToGoValue().(*websocket.Conn)
			_, pkg, err := conn.ReadMessage()
			if err != nil {
				c.RaiseRuntimeError("websocket read message error %s", err)
			}
			var j interface{}
			if err := json.Unmarshal(pkg, &j); err != nil {
				c.RaiseRuntimeError("json.decode error %v", err)
				return nil
			}
			rv := jsonToValue(j, c)
			return rv
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
					c.RaiseRuntimeError("websocket write message error %s", err)
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
					ArgRuleRequired("url", TypeStr, &url),
				)
				this.SetMember("method", NewStr("GET"), c)
			} else {
				EnsureFuncParams(c, "Request.__init__", args,
					ArgRuleRequired("method", TypeStr, &method),
					ArgRuleRequired("url", TypeStr, &url),
				)
				this.SetMember("method", method, c)
			}
			this.SetMember("url", url, c)
			this.SetMember("__headers", NewArray(), c)
			this.SetMember("__body", NewArray(), c)
		}).
		Method("certFile", func(c *Context, this ValueObject, args []Value) Value {
			var (
				certFile ValueStr
				keyFile  ValueStr
				caFile   ValueStr
			)
			EnsureFuncParams(c, "http.Request.certFile", args,
				ArgRuleRequired("certFile", TypeStr, &certFile),
				ArgRuleRequired("keyFile", TypeStr, &keyFile),
				ArgRuleOptional("caFile", TypeStr, &caFile, NewStr("")),
			)
			this.SetMember("__certs", NewArrayByValues(certFile, keyFile, caFile), c)
			return this
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
			EnsureFuncParams(c, "Request.form", args, ArgRuleRequired("form", TypeObject, &formArg))
			headers := this.GetMember("__headers", c).(ValueArray)
			headers.PushBack(NewArrayByValues(NewStr("Content-Type"), NewStr("application/x-www-form-urlencoded")))
			body := this.GetMember("__body", c).(ValueArray)
			form := url.Values{}
			formArg.Iterate(func(key string, val Value) {
				form.Set(key, val.ToString(c))
			})
			body.PushBack(NewStr(form.Encode()))
			return this
		}, "form").
		Method("multipartForm", func(c *Context, this ValueObject, args []Value) Value {
			var formArg ValueObject
			EnsureFuncParams(c, "Request.multipartForm", args, ArgRuleRequired("form", TypeObject, &formArg))
			formBs, contentType := httpGetMultipartForm(c, formArg)
			body := this.GetMember("__body", c).(ValueArray)
			body.PushBack(NewBytes(formBs))
			headers := this.GetMember("__headers", c).(ValueArray)
			headers.PushBack(NewArrayByValues(NewStr("Content-Type"), NewStr(contentType)))
			return this
		}, "form").
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
			c.RaiseRuntimeError("Request.header: invalid argument(s)")
			return nil
		}).
		Method("host", func(c *Context, this ValueObject, args []Value) Value {
			var host ValueStr
			EnsureFuncParams(c, "Request.host", args, ArgRuleRequired("host", TypeStr, &host))
			this.SetMember("__host", host, c)
			return this
		}, "host").
		Method("hosts", func(c *Context, this ValueObject, args []Value) Value {
			var hosts ValueObject
			EnsureFuncParams(c, "Request.hosts", args, ArgRuleRequired("hosts", TypeObject, &hosts))
			this.SetMember("__hosts", hosts, c)
			return this
		}, "hosts").
		Method("followRedirect", func(c *Context, this ValueObject, args []Value) Value {
			var shouldFollow ValueBool
			EnsureFuncParams(c, "Request.followRedirect", args,
				ArgRuleRequired("shouldFollow", TypeBool, &shouldFollow))
			this.SetMember("__shouldFollowRedirect", shouldFollow, c)
			return this
		}, "shouldFollow").
		Method("tlsConfig", func(c *Context, this ValueObject, args []Value) Value {
			var tlsConfig ValueObject
			EnsureFuncParams(c, "Request.tlsConfig", args, ArgRuleRequired("config", TypeObject, &tlsConfig))
			this.SetMember("__tlsConfig", tlsConfig, c)
			return this
		}, "config").
		Method("timeout", func(c *Context, this ValueObject, args []Value) Value {
			var timeout ValueFloat
			EnsureFuncParams(c, "Request.timeout", args, ArgRuleRequired("timeout", TypeFloat, &timeout))
			this.SetMember("__timeout", timeout, c)
			return this
		}, "timeout").
		Method("useClient", func(c *Context, this ValueObject, args []Value) Value {
			var client GoValue
			EnsureFuncParams(c, "Request.useClient", args,
				ArgRuleRequired("goHttpClient", TypeGoValue, &client))
			if _, ok := client.ToGoValue().(*http.Client); !ok {
				c.RaiseRuntimeError("Request.useClient: argument must be a *http.Client")
			}
			this.SetMember("__goHttpClient", client, c)
			return this
		}, "goHttpClient").
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
			req, err := http.NewRequest(strings.ToUpper(method), url, reqBody)
			if err != nil {
				c.RaiseRuntimeError("Request.build: make request error %s", err)
			}
			headers := this.GetMember("__headers", c).(ValueArray)
			for i := 0; i < headers.Len(); i++ {
				item := headers.GetIndex(i, c)
				key := item.GetIndex(0, c)
				val := item.GetIndex(1, c)
				req.Header.Add(key.ToString(c), val.ToString(c))
			}
			var httpClient *http.Client = nil
			if hc, ok := this.GetMember("__goHttpClient", c).ToGoValue().(*http.Client); ok {
				httpClient = hc
			} else {
				if certs, ok := this.GetMember("__certs", c).(ValueArray); ok && certs.Len() == 3 {
					certFile := certs.GetIndex(0, c).ToString(c)
					keyFile := certs.GetIndex(1, c).ToString(c)
					caFile := certs.GetIndex(2, c).ToString(c)
					cert, err := tls.LoadX509KeyPair(certFile, keyFile)
					if err != nil {
						c.RaiseRuntimeError("http.Request.call: load key pair error: %s", err)
					}
					tlsConfig := &tls.Config{
						Certificates: []tls.Certificate{cert},
					}
					if caFile != "" {
						caCertPool := x509.NewCertPool()
						if caCert, err := ioutil.ReadFile(caFile); err != nil {
							c.RaiseRuntimeError("http.Request.call: load ca cert error: %s", err)
						} else {
							caCertPool.AppendCertsFromPEM(caCert)
							tlsConfig.RootCAs = caCertPool
						}
					}
					tlsConfig.BuildNameToCertificate()
					transport := &http.Transport{TLSClientConfig: tlsConfig}
					if httpClient == nil {
						httpClient = &http.Client{Transport: transport}
					} else {
						httpClient.Transport = transport
					}
				}
				if s, ok := this.GetMember("__shouldFollowRedirect", c).(ValueBool); ok && !s.Value() {
					if httpClient == nil {
						httpClient = &http.Client{}
					}
					httpClient.CheckRedirect = func(*http.Request, []*http.Request) error {
						return http.ErrUseLastResponse
					}
				}
				if hosts, ok := this.GetMember("__hosts", c).(ValueObject); ok {
					if httpClient == nil {
						httpClient = &http.Client{}
					}
					hostsMap := hosts.ToGoValue().(map[string]interface{})
					var dialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
						if mapped, exists := hostsMap[addr]; exists {
							addr = fmt.Sprint(mapped)
						}
						var d net.Dialer
						return d.DialContext(ctx, network, addr)
					}
					if httpClient.Transport == nil {
						httpClient.Transport = &http.Transport{
							DialContext: dialContext,
						}
					} else if tr, ok := httpClient.Transport.(*http.Transport); ok {
						if tr.DialContext != nil {
							c.RaiseRuntimeError("DailContext already exists")
						}
						tr.DialContext = dialContext
					} else {
						c.RaiseRuntimeError("Cannot set DialContext")
					}
				}
				if tlsConfig, ok := this.GetMember("__tlsConfig", c).(ValueObject); ok {
					if httpClient == nil {
						httpClient = &http.Client{}
					}
					var transport *http.Transport
					if httpClient.Transport == nil {
						transport = &http.Transport{}
						httpClient.Transport = transport
					} else if tr, ok := httpClient.Transport.(*http.Transport); ok {
						transport = tr
					} else {
						c.RaiseRuntimeError("Get transport failed")
					}
					if transport.TLSClientConfig == nil {
						transport.TLSClientConfig = &tls.Config{}
					}
					if jsonBs, err := json.Marshal(tlsConfig.ToGoValue()); err != nil {
						c.RaiseRuntimeError("set tls config error: %+v", err)
					} else if err := json.Unmarshal(jsonBs, transport.TLSClientConfig); err != nil {
						c.RaiseRuntimeError("set tls config error: %+v", err)
					}
				}
				if timeout, ok := this.GetMember("__timeout", c).(ValueFloat); ok {
					if httpClient == nil {
						httpClient = &http.Client{}
					}
					httpClient.Timeout = time.Duration(timeout.Value() * float64(time.Second))
				}
				if httpClient == nil {
					httpClient = http.DefaultClient
				}
			}
			if host, ok := this.GetMember("__host", c).(ValueStr); ok {
				req.Host = host.Value()
			}
			resp, err := httpClient.Do(req)
			if err != nil {
				c.RaiseRuntimeError("http.Request.call: do request error %s", err)
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
			EnsureFuncParams(c, "http.Response.header", args, ArgRuleRequired("header", TypeStr, &h))
			return NewStr(resp.Header.Get(h.Value()))
		}).
		Method("bytes", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.RaiseRuntimeError("http.Response.bytes: read body error %s", err)
			}
			return NewBytes(bytes)
		}).
		Method("chunk", func(c *Context, this ValueObject, args []Value) Value {
			var chunkSize ValueInt
			EnsureFuncParams(c, "http.Response.chunk", args, ArgRuleOptional("chunkSize", TypeInt, &chunkSize, NewInt(512*1024)))
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			var stackBuf [512 * 1024]byte
			var buf []byte
			if s := chunkSize.AsInt(); s <= len(stackBuf) {
				buf = stackBuf[:s]
			} else {
				buf = make([]byte, s)
			}
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				n, err := resp.Body.Read(buf)
				if err == io.EOF {
					return NewArrayByValues(NewBytes(buf[:n]), NewBool(false))
				}
				if err != nil {
					c.RaiseRuntimeError("Read chunk error: %s", err)
				}
				return NewArrayByValues(NewBytes(buf[:n]), NewBool(true))
			})
		}, "chunkSize").
		Method("iterChunk", func(c *Context, this ValueObject, args []Value) Value {
			var chunkSize ValueInt
			EnsureFuncParams(c, "http.Response.iterChunk", args, ArgRuleOptional("chunkSize", TypeInt, &chunkSize, NewInt(512*1024)))
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			var stackBuf [512 * 1024]byte
			var buf []byte
			if s := chunkSize.AsInt(); s <= len(stackBuf) {
				buf = stackBuf[:s]
			} else {
				buf = make([]byte, s)
			}
			rv := NewObject()
			rv.SetMember("__iter__", NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
					n, err := resp.Body.Read(buf)
					if err == io.EOF {
						return NewArrayByValues(NewBytes(buf[:n]), NewBool(false))
					}
					if err != nil {
						c.RaiseRuntimeError("Read chunk error: %s", err)
					}
					return NewArrayByValues(NewBytes(buf[:n]), NewBool(true))
				})
			}), c)
			return rv
		}, "chunkSize").
		Method("text", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.RaiseRuntimeError("http.Response.text: read body error %s", err)
			}
			return NewStr(string(bytes))
		}).
		Method("json", func(c *Context, this ValueObject, args []Value) Value {
			resp := this.GetMember("__resp", c).ToGoValue().(*http.Response)
			var o interface{}
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&o); err != nil {
				c.RaiseRuntimeError("http.Response.json: decode body error %s", err)
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
				c.RaiseRuntimeError("FormFile.bytes: open file error %s", err)
			}
			defer file.Close()
			fileBs, err := ioutil.ReadAll(file)
			if err != nil {
				c.RaiseRuntimeError("FormFile.bytes: read file error %s", err)
			}
			bs := NewBytes(fileBs)
			this.SetMember("_bytes", bs, c)
			return bs
		}).
		Build()
}

func initHttpWebsocketClientClass() ValueType {
	return NewClassBuilder("WebsocketClient").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				url ValueStr
			)
			EnsureFuncParams(c, "http.WebSocket", args,
				ArgRuleRequired("url", TypeStr, &url),
			)
			this.SetMember("__url", url, c)
		}).
		Method("connect", func(c *Context, this ValueObject, args []Value) Value {
			var argHeaders ValueObject
			EnsureFuncParams(c, "http.WebSocketClient.connect", args, ArgRuleOptional("headers", TypeObject, &argHeaders, NewObject()))
			conn, ok := this.GetMember("__conn", c).ToGoValue().(*websocket.Conn)
			if ok && conn != nil {
				conn.Close()
			}
			var err error
			headers := http.Header{}
			argHeaders.Iterate(func(key string, val Value) {
				if valArr, ok := val.(ValueArray); ok {
					for i := 0; i < valArr.Len(); i++ {
						headers.Add(key, valArr.GetIndex(i, c).ToString(c))
					}
				} else {
					headers.Add(key, val.ToString(c))
				}
			})
			conn, _, err = websocket.DefaultDialer.Dial(this.GetMember("__url", c).ToString(c), headers)
			this.SetMember("__conn", NewGoValue(conn), c)
			if err != nil {
				c.RaiseRuntimeError("websocket connect error: %s", err)
			}
			return Undefined()
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			return this
		}).
		Method("read", func(c *Context, this ValueObject, args []Value) Value {
			conn, ok := this.GetMember("__conn", c).ToGoValue().(*websocket.Conn)
			if !ok || conn == nil {
				c.RaiseRuntimeError("websocket read error: no connection")
			}
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				c.RaiseRuntimeError("websocket read error: %s", err)
			}
			if msgType == websocket.TextMessage {
				return NewStr(string(msg))
			} else {
				return NewBytes(msg)
			}
		}).
		Method("readJson", func(c *Context, this ValueObject, args []Value) Value {
			conn, ok := this.GetMember("__conn", c).ToGoValue().(*websocket.Conn)
			if !ok || conn == nil {
				c.RaiseRuntimeError("websocket readJson error: no connection")
			}
			_, msg, err := conn.ReadMessage()
			if err != nil {
				c.RaiseRuntimeError("websocket readJson error: %s", err)
			}
			var j interface{}
			if err := json.Unmarshal(msg, &j); err != nil {
				c.RaiseRuntimeError("websocket readJson error: %s", err)
			}
			return jsonToValue(j, c)
		}).
		Method("write", func(c *Context, this ValueObject, args []Value) Value {
			conn, ok := this.GetMember("__conn", c).ToGoValue().(*websocket.Conn)
			if !ok || conn == nil {
				c.RaiseRuntimeError("websocket write error: no connection")
			}
			var err error
			for _, arg := range args {
				if bs, ok := arg.(ValueBytes); ok {
					err = conn.WriteMessage(websocket.BinaryMessage, bs.Value())
				} else {
					err = conn.WriteMessage(websocket.TextMessage, []byte(arg.ToString(c)))
				}
				if err != nil {
					c.RaiseRuntimeError("websocket write error: %s", err)
				}
			}
			return Undefined()
		}).
		Method("writeJson", func(c *Context, this ValueObject, args []Value) Value {
			var val Value
			EnsureFuncParams(c, "writeJson", args, ArgRuleRequired("value", TypeAny, &val))
			bs, err := jsonMarshal(val.ToGoValue())
			if err != nil {
				c.RaiseRuntimeError("websocket writeJson error: %s", err)
			}
			return c.InvokeMethod(this, "write", Args(NewBytes(bs)))
		}).
		Build()
}

func init() {
	httpRequestContextClass = initHttpRequestContextClass()
	websocketContextClass = initWebsocketContextClass()
	websocketClientClass = initHttpWebsocketClientClass()
	httpRequestClass = initHttpRequestClass()
	httpResponseClass = initHttpResponseClass()
	httpFormFileClass = initHttpFormFileClass()
}
