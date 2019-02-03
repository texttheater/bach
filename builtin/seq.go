package builtin

import (
	"github.com/texttheater/bach/values"
)

func Range(inputValue values.Value, argumentValues []values.Value) values.Value {
	start := argumentValues[0].(values.NumValue)
	end := argumentValues[1].(values.NumValue)
	channel := make(chan values.Value)
	go func() {
		for i := start; i < end; i++ {
			channel <- values.NumValue(i)
		}
		close(channel)
	}()
	return values.SeqValue(channel)
}
