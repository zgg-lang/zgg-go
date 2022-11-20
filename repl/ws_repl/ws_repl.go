package ws_repl

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zgg-lang/zgg-go/parser"
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

func (c *WebsocketReplContext) ReadCode(newCode bool, initCode string) (string, bool) {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return "", false
	}
	var input payload
	if err := json.Unmarshal(msg, &input); err != nil {
		return "", false
	}
	return input.Content, true
}

func (c *WebsocketReplContext) WriteResult(result string) {
	c.write(payload{
		Type:    writeReturn,
		Content: result,
	})
}

func (c *WebsocketReplContext) OnEnter() {
	c.write(payload{
		Type:    writeStdout,
		Content: "Welcome to ZGG Repl",
	})
}

func (*WebsocketReplContext) OnExit() {
}

