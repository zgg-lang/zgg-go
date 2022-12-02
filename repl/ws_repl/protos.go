package ws_repl

type (
	payload struct {
		Type    string      `json:"type"`
		Content string      `json:"content"`
		Data    interface{} `json:"data,omitempty"`
	}
)

const (
	readCode    = "INPUT"
	readHint    = "HINT"
	writeReturn = "RETURN"
	writeTable  = "TABLE"
	writeStdout = "STDOUT"
	writeStderr = "STDERR"
)
