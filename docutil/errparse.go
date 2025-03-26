package docutil

import (
	"encoding/json"

	"github.com/texttheater/bach/errors"
)

func ParseError(input string) (error, error) {
	// special case: treat the string "null" as "no error"
	if input == "null" {
		return nil, nil
	}
	// otherwise, treat input as JSON and unmarshal it into an errors.E
	var e errors.E
	err := json.Unmarshal([]byte(input), &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
