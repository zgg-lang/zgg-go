package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/zgg-lang/zgg-go/builtin_libs"
	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/repl"
	"github.com/zgg-lang/zgg-go/runtime"

	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
)

var (
	BUILD_TIME = "N/A"
	BUILD_HASH = "N/A"
)

func runRepl(isDebug bool) {
	repl.ReplLoop(repl.NewConsoleReplContext(isDebug, true), !isDebug)
}

func runFile(name string, inFile io.Reader, stdout, stderr io.Writer, dir string, args []string, isDebug bool) {
	srcBytes, err := io.ReadAll(inFile)
	if err != nil {
		panic(err)
	}
	srcText := string(srcBytes)
	t, errs := parser.ParseFromString(name, srcText, !isDebug)
	if n := len(errs); n > 0 {
		for i, e := range errs {
			if i >= 5 {
				fmt.Printf("%d more error(s) ...\n", n-i)
				break
			}
			fmt.Println(e.String())
		}
	} else if t == nil {
		fmt.Println("parse codes fail")
	} else {
		c := runtime.NewContext(true, isDebug, os.Getenv("CAN_EVAL") != "")
		c.Path = dir
		c.IsDebug = isDebug
		c.Args = args
		c.ImportFunc = parser.SimpleImport
		c.Stdout = stdout
		c.Stderr = stderr
		func() {
			defer func() {
				if !isDebug {
					e := recover()
					if e == nil {
						return
					}
					switch err := e.(type) {
					case runtime.Exception:
						fmt.Fprint(c.Stderr, err.MessageWithStack())
					default:
						fmt.Fprintln(c.Stderr, err)
					}
				}
			}()
			t.Eval(c)
		}()
	}
}

func updateZgg() int {
	log := func(tag, s string, v ...interface{}) {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf(tag+"|"+now+"|"+s+"\n", v...)
	}
	installSciprtUrl := "https://www.zgg-lang.org/install/install.sh"
	if s := os.Getenv("INSTALL_SH_URL"); s != "" {
		installSciprtUrl = s
	}
	log("INF", "正在从%s加载ZGG一键安装脚本", installSciprtUrl)
	resp, err := http.Get(installSciprtUrl)
	if err != nil {
		log("ERR", "一键安装脚本加载失败: %s", err)
		return 1
	}
	log("INF", "一键安装脚本加载完成")
	defer resp.Body.Close()
	c := exec.Command("/bin/sh")
	stdin, _ := c.StdinPipe()
	stdout, _ := c.StdoutPipe()
	defer stdout.Close()
	stderr, _ := c.StderrPipe()
	defer stderr.Close()
	if err := c.Start(); err != nil {
		log("ERR", "启动一键安装脚本失败: %s", err)
		return 1
	}
	io.Copy(stdin, resp.Body)
	stdin.Close()
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	if err := c.Wait(); err != nil {
		log("ERR", "运行一键安装脚本失败: %s", err)
		return 1
	}
	log("INF", "一键安装脚本执行完成")
	return 0
}

func main() {
	isDebug := os.Getenv("DEBUG") != ""
	if isDebug {
		fmt.Println("[INFO]Running in debug mode")
		prf, _ := os.Create("/tmp/pprof_result.pprof")
		defer prf.Close()
		pprof.StartCPUProfile(prf)
		defer func() {
			pprof.StopCPUProfile()
		}()
	}
	numArgs := len(os.Args)
	if numArgs > 1 {
		switch os.Args[1] {
		case "-c":
			code := ""
			if numArgs > 2 {
				code = strings.Join(os.Args[2:], " ")
			}
			runFile("", strings.NewReader(code), os.Stdout, os.Stderr, ".", []string{}, true)
		case "-m":
			if numArgs > 2 {
				moduleName := os.Args[2]
				filename := parser.GetModulePath(nil, moduleName)
				if filename == "" {
					if _, found := builtin_libs.StdLibMap[moduleName]; !found {
						return
					}

				} else if f, err := os.Open(filename); err == nil {
					defer f.Close()
					runFile(filename, f, os.Stdout, os.Stderr, filepath.Dir(filename), os.Args[3:], isDebug)
				} else {
					panic(err)
				}
			} else {
				fmt.Printf("expected module name after -m\n")
			}
		case "--update":
			os.Exit(updateZgg())
		case "--info":
			fmt.Printf("BUILD_TIME : %s\nBUILD_HASH : %s\n", BUILD_TIME, BUILD_HASH)
		case "stdin":
			runFile("input", os.Stdin, os.Stdout, os.Stderr, ".", os.Args[2:], isDebug)
		case "hub":
			runHub(os.Args[2:])
		case "deps":
			runDeps(os.Args[2:])
		case "add":
			runAddDep(os.Args[2:])
		case "ws":
			runWebsocket(isDebug, os.Args[2:])
		default:
			if f, err := os.Open(os.Args[1]); err == nil {
				defer f.Close()
				runFile(os.Args[1], f, os.Stdout, os.Stderr, filepath.Dir(os.Args[1]), os.Args[2:], isDebug)
			} else {
				panic(err)
			}
		}
	} else {
		runRepl(isDebug)
	}
}
