package ws_repl

type (
	payload struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}
)

const (
	writeReturn = "RETURN"
	writeStdout = "STDOUT"
	writeStderr = "STDERR"
)
