package ws_repl

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/repl"
	"github.com/zgg-lang/zgg-go/runtime"
)

type WebsocketReplContext struct {
	c     *runtime.Context
	conn  *websocket.Conn
	wlock sync.Mutex
}

func New(isMain, isDebug, canEval bool, conn *websocket.Conn) *WebsocketReplContext {
	c := runtime.NewContext(isMain, isDebug, canEval)
	c.ImportFunc = parser.SimpleImport
	r := &WebsocketReplContext{
		c:    c,
		conn: conn,
	}
	c.Stdout = newOutputPipe(r.onStdout)
	c.Stderr = newOutputPipe(r.onStderr)
	return r
}

func (c *WebsocketReplContext) write(v interface{}) {
	c.wlock.Lock()
	defer c.wlock.Unlock()
	c.conn.WriteJSON(v)
}

func (c *WebsocketReplContext) onStdout(bs []byte) {
	c.write(payload{
		Type:    writeStdout,
		Content: string(bs),
	})
}

func (c *WebsocketReplContext) onStderr(bs []byte) {
	c.write(payload{
		Type:    writeStderr,
		Content: string(bs),
	})
}

func (c *WebsocketReplContext) Context() *runtime.Context { return c.c }

func (c *WebsocketReplContext) ReadAction(shouldRecover bool) repl.ReplAction {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return nil
	}
	var input payload
	if err := json.Unmarshal(msg, &input); err != nil {
		return nil
	}
	switch input.Type {
	}
	compiled, err := repl.ParseInputCode(input.Content, shouldRecover)
	return repl.ReplRunCode{Compiled: compiled, Err: err}
}

func (c *WebsocketReplContext) writePtable(content string, obj runtime.Value) bool {
	cc := c.Context()
	cc.InvokeMethod(obj, "toArray", runtime.Args(runtime.NewBool(true)))
	arr, is := c.Context().RetVal.(runtime.ValueArray)
	if !is {
		return false
	}
	data := make([][]string, arr.Len())
	for i := 0; i < arr.Len(); i++ {
		item, is := arr.GetIndex(i, cc).(runtime.ValueArray)
		if !is {
			return false
		}
		row := make([]string, item.Len())
		for j := range row {
			row[j] = item.GetIndex(j, cc).ToString(cc)
		}
		data[i] = row
	}
	c.write(payload{
		Type:    writeTable,
		Content: content,
		Data:    data,
	})
	return true
}

func (c *WebsocketReplContext) WriteResult(result interface{}) {
	if result == nil {
		c.write(payload{
			Type:    writeReturnNothing,
			Content: "",
		})
		return
	}
	var content string
	switch rv := result.(type) {
	case runtime.Value:
		content = rv.ToString(c.Context())
		if rv.Type().Name == "PTable" {
			if c.writePtable(content, rv) {
				return
			}
		}
	case string:
		content = rv
	default:
		content = fmt.Sprint(rv)
	}
	c.write(payload{
		Type:    writeReturn,
		Content: content,
	})
}

func (c *WebsocketReplContext) WriteException(e runtime.Exception) {
	c.write(payload{
		Type:    writeException,
		Content: e.MessageWithStack(),
		Data: M{
			"error": e.GetMessage(),
			"stack": lo.Map(e.GetStack(), func(stack runtime.Stack, _ int) M {
				return M{"filename": stack.FileName, "line": stack.Line, "function": stack.Function}
			}),
		},
	})
}

func (c *WebsocketReplContext) OnEnter() {
	c.write(payload{
		Type:    writeStdout,
		Content: "Welcome to ZGG Web Repl",
	})
}

func (*WebsocketReplContext) OnExit() {
}
