package builtin

import (
	"fmt"
	"log"

	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initFmt() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		// for Num fmt:f Str
		expressions.SimpleFuncer(
			types.Num{},
			"fmtF",
			nil,
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := float64(inputValue.(states.NumValue))
				outputStr := fmt.Sprintf("%f", inputNum)
				return states.StrValue(outputStr), nil
			},
		),
		// for Num fmt:f(Num) Str
		expressions.SimpleFuncer(
			types.Num{},
			"fmtF",
			[]types.Type{
				types.Num{},
			},
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := float64(inputValue.(states.NumValue))
				precision := int(argumentValues[0].(states.NumValue))
				var outputStr string
				if precision < 0 {
					outputStr = fmt.Sprintf("%f", inputNum)
				} else {
					outputStr = fmt.Sprintf("%."+fmt.Sprint(precision)+"f", inputNum)
				}
				return states.StrValue(outputStr), nil
			},
		),
		// for Num fmt:f(Num) Str
		expressions.SimpleFuncer(
			types.Num{},
			"fmtF",
			[]types.Type{
				types.Num{},
				types.Num{},
			},
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := float64(inputValue.(states.NumValue))
				precision := int(argumentValues[0].(states.NumValue))
				width := int(argumentValues[1].(states.NumValue))
				log.Println(precision)
				log.Println(width)
				formatString := "%"
				if width > 0 {
					formatString += fmt.Sprint(width)
				}
				if precision >= 0 {
					formatString += "."
					formatString += fmt.Sprint(precision)
				}
				formatString += "f"
				log.Println(formatString)
				outputString := fmt.Sprintf(formatString, inputNum)
				return states.StrValue(outputString), nil
			},
		),
	})
}
