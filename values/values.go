package values

import (
	"github.com/texttheater/bach/types"
)

type Value interface {
	String() string
	Out() string
	Iter() <-chan Value
	Inhabits(types.Type) bool
}

func inhabits(v Value, t types.UnionType) bool {
	for _, disjunct := range t {
		if v.Inhabits(disjunct) {
			return true
		}
	}
	return false
}
