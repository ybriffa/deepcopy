package deepcopy

import (
	"reflect"
	"testing"
)

func foo() {
}

type testStruct struct {
	Str  string
	Num  int
	Dec  float64
	Bool bool
	Nil  interface{}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		obj interface{}
	}{
		{
			obj: 42,
		},
		{
			obj: map[string]interface{}{"str": (*string)(nil), "n": 0, "nil": nil},
		},
		{
			obj: []string{"a", "b", "c"},
		},
		{
			obj: &testStruct{
				Str:  "test",
				Num:  42,
				Dec:  42.42,
				Bool: true,
				Nil:  nil,
			},
		},
	}

	for idx, test := range tests {
		cp := Copy(test.obj)

		t.Log(cp)

		if &cp == &test.obj {
			t.Fatalf("[test #%d] unexpected equal pointer", idx)
		}

		if !reflect.DeepEqual(test.obj, cp) {
			t.Fatalf("[test #%d] unequal objects", idx)
		}
	}
}
