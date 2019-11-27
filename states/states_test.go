package states_test

import (
	"testing"

	"github.com/texttheater/bach/states"
)

func TestArrEqual(t *testing.T) {
	elements1 := []states.Value{
		states.NumValue(1),
		states.NumValue(2),
		states.NumValue(3),
	}
	elements2 := []states.Value{
		states.NumValue(1),
		states.NumValue(2),
		states.NumValue(3),
	}
	arr1 := states.NewArrValue(elements1)
	arr2 := states.NewArrValue(elements2)
	ok, err := arr1.Equal(arr2)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !ok {
		t.Log("[1, 2, 3] == [1, 2, 3]")
		t.Fail()
	}
}
