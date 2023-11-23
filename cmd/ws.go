package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zgg-lang/zgg-go"
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

func checkAuthByScript(authScript string, w http.ResponseWriter, r *http.Request, event string) bool {
	if authScript != "" {
		if bs, err := os.ReadFile(authScript); err != nil {
			log(ERROR, "read authScript %s fail: %+v", authScript, err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		} else if authHandler, err := zgg.CompileCode(string(bs)); err != nil {
			log(ERROR, "compile authScript %s fail: %+v", authScript, err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		} else if r, err := zgg.RunCode(authHandler, zgg.Var{"event", zgg.Val{event}},
			zgg.Var{"request", r}, zgg.Var{"responseWriter", w}); err != nil {
			log(ERROR, "execute auth script error: %+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		} else if pass, is := r["pass"].(bool); !is || !pass {
			return false
		}
	}
	return true
}

//go:embed wsindex.html
var wsIndex string

func runWebsocket(isDebug bool, args []string) {
	var listen, path, authScript, indexPath string
	fs := flag.NewFlagSet("zgg ws", flag.ExitOnError)
	fs.StringVar(&authScript, "auth", "", "指定鉴权处理脚本路径，留空为不鉴权")
	fs.StringVar(&indexPath, "index", "", "指定首页HTML文件路径，留空为内置首页")
	fs.Parse(args)
	addr := fs.Arg(0)
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
		if !checkAuthByScript(authScript, w, r, "index") {
			log(INFO, "index blocked by authscript")
			return
		}
		if indexPath != "" {
			http.ServeFile(w, r, indexPath)
		} else {
			http.ServeContent(w, r, "index.html", time.Time{}, strings.NewReader(wsIndex))
		}
	})
	http.HandleFunc(path+"/session", func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthByScript(authScript, w, r, "connect") {
			log(INFO, "ws blocked by authscript")
			return
		}
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
