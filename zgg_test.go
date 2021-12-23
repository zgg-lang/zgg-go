package zgg

import (
	"strings"
	"testing"
)

func TestRunCode(t *testing.T) {
	RunCode(`
		println(1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9)
	`)
}

func TestRunEval(t *testing.T) {
	for i := 1; i < 10; i++ {
		res, err := Eval(`
			fib := n => when n {
				1, 2 -> 1
				else -> fib(n-1) + fib(n-2)
			}
			fib(v)
		`, Var{"v", Val(i)})
		if err != nil {
			t.Fatalf("eval error %s", err)
		}
		t.Logf("fib(%d) = %d", i, res.(int64))
	}
}

func TestRunner(t *testing.T) {
	var outbuf strings.Builder
	runner := NewRunner().
		Stdout(&outbuf).
		Var("a", 10).
		Var("b", 11)
	runner.Run("println(a * b)")
	t.Log(outbuf.String())
	outbuf.Reset()

	runner.Var("a", 1000)
	runner.Run("println(a + b)")
	t.Log(outbuf.String())
}

func BenchmarkWithoutPrecompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRunner().Eval("1+2+3")
	}
}

func BenchmarkWithoutPrecompileReuseRunner(b *testing.B) {
	runner := NewRunner()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner.Eval("1+2+3")
	}
}

func BenchmarkWithPrecompile(b *testing.B) {
	code, _ := NewRunner().CompileExpr("1+2+3")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewRunner().Eval(code)
	}
}

func BenchmarkWithPrecompileAndReuseRunner(b *testing.B) {
	runner := NewRunner()
	code, _ := runner.CompileExpr("1+2+3")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner.Eval(code)
	}
}

func BenchmarkRunWithPool(b *testing.B) {
	prog, _ := NewRunner().CompileCode("1+2+3+a")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunCode(prog, Var{"a", Val(i)})
	}
}
