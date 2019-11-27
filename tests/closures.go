package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestClosures(t *testing.T) {
	TestProgram(`1 =a for Any def f Num as a ok f 2 =a f`, types.NumType{}, states.NumValue(1), nil, t)
}
