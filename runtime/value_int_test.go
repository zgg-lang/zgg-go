package runtime

import "testing"

func TestNewInt(t *testing.T) {
	{
		v := NewInt(10)
		if v.Value() != 10 {
			t.Errorf("v.Value() returns %d not %d", v.Value(), 10)
		}
	}
}
