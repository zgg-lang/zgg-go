package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type depsRequirement struct {
	Name string
	URL  string
}

func depsLog(verbose bool, tag, msg string, args ...interface{}) {
	if !verbose && tag == "VER" {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05.000")
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	fmt.Println(now + "|" + tag + "|" + msg)
}

func runDeps(args []string) {
	var (
		depFile string
		verbose bool
	)
	flagset := flag.NewFlagSet("deps", flag.ExitOnError)
	flagset.StringVar(&depFile, "f", "zggdeps.txt", "依赖文件")
	flagset.BoolVar(&verbose, "v", false, "show detail logs")
	flagset.Parse(args)
	reqs, err := depsGetRequirements(depFile, verbose)
	if err != nil {
		panic(err)
	}
	for _, r := range reqs {
		depsLog(verbose, "VER", "正在从%s下载依赖项%s的内容...", r.URL, r.Name)
		z, err := depsGetZip(r.URL)
		if err != nil {
			panic(err)
		}
		root := filepath.Join("zgg_modules", r.Name)
		for _, f := range z.File {
			if err := depsSaveFile(f, root); err != nil {
				panic(err)
			}
		}
	}
	depsLog(verbose, "INF", "完成")
}

func depsGetRequirements(filename string, verbose bool) ([]depsRequirement, error) {
	depsLog(verbose, "VER", "正在从%s读取依赖项", filename)
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		depsLog(verbose, "ERR", "加载依赖文件%s失败: %s", filename, err)
		return nil, err
	}
	lines := strings.Split(string(bs), "\n")
	rv := make([]depsRequirement, 0, len(lines))
	for _, line := range lines {
		if p := strings.Index(line, "//"); p >= 0 {
			line = line[:p]
		}
		line = strings.Trim(line, " \t\r")
		if line == "" {
			continue
		}
		req, err := depsParseRequirement(line, verbose)
		if err != nil {
			return nil, err
		}
		rv = append(rv, req)
	}
	return rv, nil
}

func depsParseRequirement(s string, verbose bool) (depsRequirement, error) {
	depsLog(verbose, "VER", "正在解析依赖%s...", s)
	parts := strings.SplitN(s, "@", 2)
	url := parts[0]
	tag := ""
	if len(parts) > 1 {
		tag = parts[1]
	}
	name := url
	if !strings.ContainsRune(url, '/') {
		url = "github.com/zgg-libs/" + url
	}
	if strings.HasPrefix(url, "github.com/") {
		if tag != "" {
			url = fmt.Sprintf("https://%s/archive/refs/tags/%s.zip", url, tag)
		} else {
			url = fmt.Sprintf("https://%s/archive/refs/heads/main.zip", url)
		}
	} else {
		return depsRequirement{}, fmt.Errorf("未支持的依赖描述%s", s)
	}
	return depsRequirement{
		Name: name,
		URL:  url,
	}, nil
}

func depsGetZip(url string) (*zip.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败 %d", resp.StatusCode)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return zip.NewReader(bytes.NewReader(bs), int64(len(bs)))
}

func depsSaveFile(f *zip.File, root string) error {
	fName := f.Name
	if p := strings.IndexRune(fName, '/'); p >= 0 {
		fName = fName[p+1:]
	}
	dst := filepath.Join(root, fName)
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return fmt.Errorf("创建目录%s失败: %s", dst, err)
		}
		return nil
	}
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目录%s失败: %s", dstDir, err)
	}
	rd, err := f.Open()
	if err != nil {
		return fmt.Errorf("从zip包读取文件%s失败: %s", f.Name, err)
	}
	defer rd.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("从zip包读取文件%s失败: %s", f.Name, err)
	}
	defer dstFile.Close()
	io.Copy(dstFile, rd)
	return nil
}
