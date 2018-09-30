package contexts

import (
	"github.com/texttheater/bach/types"
)

// A Context consists of an input type and a number of argument types. Together
// with a name, the context determines the function 
type Context struct {
	InputType types.Type
}
