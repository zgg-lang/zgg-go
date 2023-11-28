package ws_repl

type (
	payload struct {
		Type    string      `json:"type"`
		Content string      `json:"content"`
		Data    interface{} `json:"data,omitempty"`
	}
	M = map[string]any
)

const (
	readCode           = "INPUT"
	readHint           = "HINT"
	writeException     = "EXCEPTION"
	writeReturn        = "RETURN"
	writeReturnNothing = "RETURN_NOTHING"
	writeTable         = "TABLE"
	writeStdout        = "STDOUT"
	writeStderr        = "STDERR"
)
