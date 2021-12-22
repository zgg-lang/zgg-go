package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"

	"github.com/urfave/cli"
)

func genPair(c *cli.Context) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	pub := priv.Public().(*rsa.PublicKey)
	privName := c.Args().First()
	if err := ioutil.WriteFile(privName, x509.MarshalPKCS1PrivateKey(priv), 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(privName+".pub", x509.MarshalPKCS1PublicKey(pub), 0644); err != nil {
		return err
	}
	return nil
}

func init() {
	commands = append(commands, cli.Command{
		Name:  "genpair",
		Usage: "生成请求hub的密钥对",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "文件名前缀，生成的私钥为${prefix}，公钥为${prefix}.pub",
			},
		},
		Action: genPair,
	})
}
