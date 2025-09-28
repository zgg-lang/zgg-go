package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func reqHub(c *cli.Context) error {
	var (
		keyPath = c.String("key")
		hubURL  = c.Args().First()
	)
	code, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	keyBs, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}
	priv, err := x509.ParsePKCS1PrivateKey(keyBs)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write(code)
	sign, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return err
	}
	signature := hex.EncodeToString(sign)
	req, err := http.NewRequestWithContext(context.Background(), "POST", hubURL, bytes.NewReader(code))
	if err != nil {
		return err
	}
	username := filepath.Base(keyPath)
	req.SetBasicAuth(username, signature)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	io.Copy(os.Stdout, rsp.Body)
	return nil
}

func init() {
	commands = append(commands, cli.Command{
		Name:  "reqhub",
		Usage: "请求远端hub运行代码",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "key",
				Usage: "私钥文件路径",
			},
		},
		Action: reqHub,
	})
}
