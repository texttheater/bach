package values

import (
	"github.com/texttheater/bach/types"
)

type SeqValue struct {
	ElementType types.Type
	Channel     chan Value
}

func (v SeqValue) String() string {
	return "<seq>"
}

func (v SeqValue) Out() string {
	return v.String()
}

func (v SeqValue) Iter() <-chan Value {
	// TODO safeguard against iterating twice?
	return v.Channel
}
