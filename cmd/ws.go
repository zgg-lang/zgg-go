package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zgg-lang/zgg-go/repl"
	"github.com/zgg-lang/zgg-go/repl/ws_repl"
)

type logLevel string

const (
	INFO  logLevel = "INF"
	ERROR logLevel = "ERR"
)

func log(level logLevel, msg string, args ...interface{}) {
	_, file, line, _ := goruntime.Caller(1)
	file = filepath.Base(file)
	now := time.Now().Format("2006-01-02 15:04:05")
	f := fmt.Sprintf("%s|%s|%d|%s|%s\n", now, file, line, level, msg)
	fmt.Fprintf(os.Stderr, f, args...)
}

func runWebsocket(isDebug bool, addr string) {
	var listen, path string
	if p := strings.Index(addr, "/"); p >= 0 {
		listen = addr[0:p]
		path = addr[p:]
	} else {
		listen = addr
		path = "/"
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		clientAddr := r.RemoteAddr
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			log(ERROR, "Upgrade %s error: %v", clientAddr, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer conn.Close()
		log(INFO, "new connection from %s", clientAddr)
		repl.ReplLoop(ws_repl.New(true, isDebug, true, conn), !isDebug)
	})
	log(INFO, "starting serving websocket console at %s%s...", listen, path)
	http.ListenAndServe(listen, nil)
}
