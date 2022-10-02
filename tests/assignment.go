package tests

import (
	"testing"
)

func TestAssignment(t *testing.T) {
	TestProgramStr(
		`1 +1 =a 3 *2 +a`,
		`Num`,
		`8`,
		nil,
		t,
	)
	TestProgramStr(
		`1 +1 ==2 =p 1 +1 ==1 =q p ==q not`,
		`Bool`,
		`true`,
		nil,
		t,
	)
}
