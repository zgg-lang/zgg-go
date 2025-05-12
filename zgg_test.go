package zgg

import (
	"strings"
	"testing"

	"github.com/zgg-lang/zgg-go/parser"
	"github.com/zgg-lang/zgg-go/runtime"
)

func TestRunCode(t *testing.T) {
	r, e := RunCode(`
		println(1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9)
		export result := 123321
	`)
	t.Log("err:", e)
	t.Log("result:", r["result"].(int64))
}

func TestRunEval(t *testing.T) {
	for i := 1; i < 10; i++ {
		res, err := Eval(`
			fib := n => when n {
				1, 2 -> 1
				else -> fib(n-1) + fib(n-2)
			}
			fib(v)
		`, Var{"v", Val{i}})
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

func TestCustomImport(t *testing.T) {
	myImport := func(c *runtime.Context, name string, code string, importType string, reloadIfNewer bool) (runtime.Value, int64, bool) {
		if name == "sys" {
			return nil, 0, false
		}
		return parser.SimpleImport(c, name, code, importType, reloadIfNewer)
	}
	exported, err := RunCode(`a := @sys.getResult('ls')`, ImportFunc(myImport))
	t.Log(exported, err)
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
	prog, _ := CompileCode("1+2+3+a")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunCode(prog, Var{"a", Val{i}})
	}
}

func BenchmarkCalcWithPrecalc(b *testing.B) {
	parser.CanCalcInCompileTime = true
	code, _ := NewRunner().CompileExpr("1+2+3*3-3**234.2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eval(code)
	}
}

func BenchmarkCalcWithoutPrecalc(b *testing.B) {
	parser.CanCalcInCompileTime = false
	code, _ := NewRunner().CompileExpr("1+2+3*3-3**234.2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eval(code)
	}
}
