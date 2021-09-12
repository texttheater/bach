package tests

import (
	"testing"
)

func TestConditionals(t *testing.T) {
	TestProgramStr(`if true then 2 else 3 ok`, `Num`, `2`, nil, t)
	TestProgramStr(`for Num def heart Bool as if <3 then true else false ok ok 2 heart`, `Bool`, `true`, nil, t)
	TestProgramStr(`for Num def heart Bool as if <3 then true else false ok ok 4 heart`, `Bool`, `false`, nil, t)
	TestProgramStr(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 -1 expand`, `Num`, `-2`, nil, t)
	TestProgramStr(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 1 expand`, `Num`, `2`, nil, t)
	TestProgramStr(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 expand`, `Num`, `0`, nil, t)
}
