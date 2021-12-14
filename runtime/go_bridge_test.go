package runtime

import (
	"reflect"
	"testing"
)

func TestToGoValue(t *testing.T) {
	{
		inVal := NewArrayByValues(NewInt(1), NewInt(2), NewInt(3))
		var x []int
		outVal := reflect.New(reflect.TypeOf(x))
		toGoValue(NewContext(true, true), inVal, outVal.Elem())
		t.Logf("outval: %+v", outVal.Elem().Interface())
	}
}
