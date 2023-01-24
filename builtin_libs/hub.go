package builtin_libs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"io"
	"net/http"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func libHub(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("invokeByRSA", NewNativeFunction("invoke", func(c *Context, this Value, args []Value) Value {
		var (
			code     ValueStr
			url      ValueStr
			username ValueStr
			keyBs    ValueBytes
		)
		EnsureFuncParams(c, "invokeByRSA", args,
			ArgRuleRequired("code", TypeStr, &code),
			ArgRuleRequired("url", TypeStr, &url),
			ArgRuleRequired("username", TypeStr, &username),
			ArgRuleRequired("rsa", TypeBytes, &keyBs),
		)
		priv, err := x509.ParsePKCS1PrivateKey(keyBs.Value())
		if err != nil {
			c.RaiseRuntimeError("Parse privatekey error %v", err)
		}
		hash := sha256.New()
		codeBs := []byte(code.Value())
		hash.Write(codeBs)
		sign, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash.Sum(nil))
		if err != nil {
			c.RaiseRuntimeError("SignPKCS1v15 error %v", err)
		}
		req, err := http.NewRequest("POST", url.Value(), bytes.NewReader(codeBs))
		if err != nil {
			c.RaiseRuntimeError("Make request to hub error %v", err)
		}
		req.SetBasicAuth(username.Value(), hex.EncodeToString(sign))
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.RaiseRuntimeError("Request hub error %v", err)
		}
		defer rsp.Body.Close()
		rspBs, err := io.ReadAll(rsp.Body)
		if err != nil {
			c.RaiseRuntimeError("Read result from hub error %v", err)
		}
		return NewBytes(rspBs)
	}, "code", "url", "rsa"), nil)
	return lib
}
