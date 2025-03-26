package builtin

import (
	"math/rand"
	"time"

	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
)

var InitialShape shapes.Shape = shapes.Shape{
	Type:  types.NullType{},
	Stack: nil,
}

func init() {
	InitialShape.Stack = InitialShape.Stack.PushAll(NullFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(IOFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(LogicFuncers)
	rand.Seed(time.Now().UnixNano())
	InitialShape.Stack = InitialShape.Stack.PushAll(MathFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(TextFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(ArrFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(ObjFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(TypeFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(ValueFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(RegexpFuncers)
	InitialShape.Stack = InitialShape.Stack.PushAll(ControlFuncers)
}
