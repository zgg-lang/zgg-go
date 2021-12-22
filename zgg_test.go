package zgg

import "testing"

func TestRunCode(t *testing.T) {
	RunCode(`
		println(1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9)
	`)
}

func TestRunEval(t *testing.T) {
	res, err := Eval("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9")
	if err != nil {
		t.Errorf("eval error %s", err)
	} else {
		t.Logf("result is %+v", res)
	}
}

func TestRunner(t *testing.T) {
	runner := NewRunner().
		Var("a", 10).
		Var("b", 11)
	runner.Run("return a * b")
}
