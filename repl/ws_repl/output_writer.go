package ws_repl

import (
	"bytes"
	"strings"
)

type outputPipe struct {
	buf   *bytes.Buffer
	flush func([]byte)
}

func newOutputPipe(flush func([]byte)) *outputPipe {
	return &outputPipe{
		buf:   bytes.NewBuffer(nil),
		flush: flush,
	}
}

func (op *outputPipe) Write(p []byte) (n int, err error) {
	n, err = op.buf.Write(p)
	if err != nil {
		return
	}
	content := string(op.buf.Bytes())
	lastNewline := strings.LastIndexAny(content, "\n")
	if lastNewline == -1 {
		return
	}
	lastNewline++
	toFlush := []byte(content[:lastNewline])
	toBuffer := []byte(content[lastNewline:])
	op.buf.Reset()
	if len(toBuffer) > 0 {
		op.buf.Write(toBuffer)
	}
	if len(toFlush) > 0 {
		op.flush(toFlush)
	}
	return
}
