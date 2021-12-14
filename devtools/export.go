package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

type exportInfo struct {
	Name          string
	ImportPath    string
	Exported      bool
	Funcs         []string
	Consts        []string
	Vars          []string
	NonInterfaces []string
	TypeMapping   map[string]string
	symbols       map[string]bool
}

func exportValue(info *exportInfo, s ast.Spec, toArr []string) []string {
	if spec, ok := s.(*ast.ValueSpec); ok {
		for _, nameIdent := range spec.Names {
			name := nameIdent.Name
			if !ast.IsExported(name) {
				continue
			}
			if info.symbols[name] {
				continue
			}
			info.symbols[name] = true
			toArr = append(toArr, name)
		}
	}
	return toArr
}

func exportProcessSrcFile(info *exportInfo, srcFile *ast.File) error {
	for _, d := range srcFile.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			if decl.Recv == nil { // 只导出函数，不导出方法
				if name := decl.Name.Name; ast.IsExported(name) && !info.symbols[name] {
					info.Funcs = append(info.Funcs, name)
					info.symbols[name] = true
				}
			}
		case *ast.GenDecl:
			switch decl.Tok {
			case token.CONST:
				for _, spec := range decl.Specs {
					info.Consts = exportValue(info, spec, info.Consts)
				}
			case token.VAR:
				for _, spec := range decl.Specs {
					info.Vars = exportValue(info, spec, info.Vars)
				}
			case token.TYPE:
				for _, spec := range decl.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if name := ts.Name.Name; ast.IsExported(name) && !info.symbols[name] {
						if _, ok := ts.Type.(*ast.InterfaceType); !ok {
							info.NonInterfaces = append(info.NonInterfaces, name)
							info.symbols[name] = true
						}
					}
				}
			}
		}
	}
	return nil
}

func dirExists(filepath string) (bool, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	} else if !fi.IsDir() {
		return false, nil
	}
	return true, nil
}

func exportPkg(c *cli.Context, root, pkg string) (*exportInfo, error) {
	fs := token.NewFileSet()
	dir := path.Join(root, "src", pkg)
	if exists, err := dirExists(dir); err != nil {
		return nil, err
	} else if !exists {
		return nil, nil
	}
	getFilesCmd := exec.Command(c.String("go"), "list", "-f", "{{range $i,$f := .GoFiles}}{{if $i}},{{end}}{{$f}}{{end}}")
	getFilesCmd.Dir = dir
	srcFiles := map[string]bool{}
	if output, err := getFilesCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("go list error: %s", err)
	} else {
		fmt.Fprintln(os.Stderr, string(output))
		for _, f := range strings.Split(string(output), ",") {
			if f != "" {
				srcFiles[strings.Trim(f, " \n\r\t")] = true
			}
		}
	}
	fileFilter := func(fi os.FileInfo) bool {
		return srcFiles[fi.Name()]
	}
	pkgs, err := parser.ParseDir(fs, dir, fileFilter, 0)
	if err != nil {
		return nil, err
	}
	info := &exportInfo{ImportPath: pkg, symbols: map[string]bool{}, TypeMapping: map[string]string{}}
	for _, pkg := range pkgs {
		if pkg.Name == "main" {
			continue
		}
		info.Name = pkg.Name
		for filename, srcFile := range pkg.Files {
			if !strings.HasSuffix(filename, ".go") {
				continue
			}
			if strings.HasSuffix(filename, "_test.go") {
				continue
			}
			if err := exportProcessSrcFile(info, srcFile); err != nil {
				return nil, err
			}
		}
		break
	}
	return info, nil
}

func exportGoPkg(c *cli.Context) error {
	pkgName := c.Args().First()
	gopath := os.Getenv("GOPATH")
	if root := c.String("root"); root != "" {
		gopath = root
	}
	paths := strings.Split(gopath, ":")
	var info *exportInfo
	var err error
	for _, srcPath := range paths {
		if info, err = exportPkg(c, srcPath, pkgName); err != nil && !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "导出错误：", err)
			return err
		} else if info != nil {
			break
		}
	}
	if info == nil {
		err := fmt.Errorf("找不到指定包:%s", pkgName)
		fmt.Fprintln(os.Stderr, "导出错误：", err)
		return err
	}
	info.Exported = len(info.symbols) > 0
	shouldFormat := false
	tpl := c.String("gotemplate")
	if tpl != "" {
		shouldFormat = true
	} else {
		tpl = c.String("template")
	}
	if st := c.String("spectypes"); st != "" {
		for _, mapping := range strings.Split(st, ";") {
			kv := strings.Split(mapping, ":")
			if len(kv) == 2 {
				info.TypeMapping[kv[0]] = kv[1]
			}
		}
	}
	if tpl != "" {
		tplContent, err := ioutil.ReadFile(tpl)
		if err != nil {
			fmt.Fprintln(os.Stderr, "导出错误：查找模板时发生错误：", err)
			return err
		}
		if t, err := template.New("code").Parse(string(tplContent)); err == nil {
			buf := bytes.NewBuffer(nil)
			if err := t.Execute(buf, info); err != nil {
				fmt.Fprintln(os.Stderr, "导出错误：渲染模板时发生错误：", err)
				return err
			}
			if shouldFormat {
				if out, err := format.Source(buf.Bytes()); err != nil {
					fmt.Fprintln(os.Stderr, "导出错误：格式化输出代码时发生错误：", err)
					return err
				} else {
					fmt.Print(string(out))
				}
			} else {
				fmt.Print(buf.String())
			}
		} else {
			fmt.Fprintln(os.Stderr, "导出错误：解析模板时发生错误：", err)
			return err
		}
	} else {
		o, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(o))
	}
	return nil
}

func init() {
	commands = append(commands, &cli.Command{
		Name:  "export",
		Usage: "导出go模块符号",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "root",
				Usage: "搜索目录，分号分隔。默认为GOPATH",
			},
			&cli.StringFlag{
				Name:  "go",
				Value: "go",
				Usage: "go命令路径，默认根据PATH查找",
			},
			&cli.StringFlag{
				Name:  "gotemplate",
				Usage: "导出GO代码模板模板",
			},
			&cli.StringFlag{
				Name:  "template",
				Usage: "导出模板，不指定时导出json",
			},
			&cli.StringFlag{
				Name:  "spectypes",
				Usage: "指定符号类型",
			},
		},
		Action: exportGoPkg,
	})
}
