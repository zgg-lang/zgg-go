package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ziipin-server/zplog"
)

func runHub(args []string) {
	var (
		addr    string
		rsaRoot string
		secret  string
	)
	flagset := flag.NewFlagSet("hub", flag.ExitOnError)
	flagset.StringVar(&addr, "addr", ":40000", "http listening address")
	flagset.StringVar(&rsaRoot, "rsa", "", "rsa public keys' root")
	flagset.StringVar(&secret, "secret", "", "secret")
	flagset.Parse(args)
	if rsaRoot != "" {
		hubAuthRequest = hubGetAuthByRSA(rsaRoot)
	} else if secret != "" {
		hubAuthRequest = hubGetAuthBySecret(secret)
	} else {
		hubAuthRequest = hubDefaultAuthRequest
	}
	fmt.Printf("Start serving on %s...\n", addr)
	http.ListenAndServe(addr, http.HandlerFunc(hubHandleRequest))
}

var hubAuthRequest func(*http.Request) []byte

func hubDefaultAuthRequest(r *http.Request) []byte {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	return bs
}

var (
	hubRsaPublicKeyType = reflect.TypeOf((*rsa.PublicKey)(nil))
)

func hubGetPub(pubfile string) (*rsa.PublicKey, error) {
	pubBs, err := ioutil.ReadFile(pubfile)
	if err != nil {
		return nil, err
	}
	rsaPub, err := x509.ParsePKCS1PublicKey(pubBs)
	if err != nil {
		return nil, err
	}
	return rsaPub, nil
}

func hubGetAuthByRSA(rsaRoot string) func(*http.Request) []byte {
	return func(r *http.Request) []byte {
		username, signatureHex, ok := r.BasicAuth()
		if !ok {
			zplog.LogError("AuthByRSA: get basicauth failed")
			return nil
		}
		if strings.ContainsRune(username, '.') {
			zplog.LogError("AuthByRSA: invalid username %s", username)
			return nil
		}
		signature, err := hex.DecodeString(signatureHex)
		if err != nil {
			zplog.LogError("AuthByRSA: decode signature [%s] error %s", signatureHex, err)
			return nil
		}
		rsaPub, err := hubGetPub(filepath.Join(rsaRoot, username+".pub"))
		if err != nil {
			zplog.LogError("AuthByRSA: parse pub file error %s", err)
			return nil
		}
		code, err := ioutil.ReadAll(r.Body)
		if err != nil {
			zplog.LogError("AuthByRSA: read code error %s", err)
			return nil
		}
		hash := sha256.New()
		hash.Write(code)
		if err := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hash.Sum(nil), signature); err != nil {
			zplog.LogError("AuthByRSA: verify signature error %s", err)
			return nil
		}
		return code
	}
}

func hubGetAuthBySecret(secret string) func(*http.Request) []byte {
	return func(r *http.Request) []byte {
		if secret != r.Header.Get("X-ZGG-SECRET") {
			return nil
		}
		bs, _ := ioutil.ReadAll(r.Body)
		return bs
	}
}

func hubHandleRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	code := hubAuthRequest(r)
	if code == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var (
		qs   = r.URL.Query()
		args = qs["args"]
	)
	if args == nil {
		args = []string{}
	}
	runFile("", bytes.NewReader(code), w, w, ".", args, false)
}
