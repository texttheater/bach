package builtin

import (
	"math"
	"math/bits"
	"math/rand"
	"time"

	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/golang-variadic-hypot/varhypot"
)

func initMath() {
	rand.Seed(time.Now().UnixNano())
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.SimpleFuncer(
			types.NumType{},
			"+",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum + argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"-",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum - argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"-",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(argumentValues[0].(states.NumValue))
				if math.Signbit(x) {
					return states.NumValue(math.Copysign(x, 1)), nil
				} else {
					return states.NumValue(math.Copysign(x, -1)), nil
				}
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"*",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum * argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"/",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum / argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"%",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(math.Mod(float64(inputNum), float64(argumentNum))), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"<",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum < argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			">",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum > argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"==",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum == argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"<=",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum <= argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			">=",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum >= argumentNum), nil
			},
		),
		expressions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"sum",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				iter := states.IterFromValue(inputValue)
				sum := 0.0
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.NumValue(sum), nil
					}
					sum += float64(value.(states.NumValue))
				}
			},
		),
		expressions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"avg",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				iter := states.IterFromValue(inputValue)
				sum := 0.0
				count := 0.0
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.NumValue(sum / count), nil
					}
					sum += float64(value.(states.NumValue))
					count += 1.0
				}
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"inf",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Inf(1)), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"nan",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.NaN()), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"isFinite",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(!math.IsInf(x, 0)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"isNaN",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(math.IsNaN(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"epsilon",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Nextafter(1, 2) - 1), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"largestSafeInteger",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(9007199254740991), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"largestNum",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.MaxFloat64), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"smallestSafeInteger",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(-9007199254740991), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"smallestPositiveNum",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.SmallestNonzeroFloat64), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"isInteger",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(x == float64(int(x))), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"isSafeInteger",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(x >= -9007199254740991 && x <= 9007199254740991), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"e",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.E), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"ln2",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Ln2), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"ln10",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Ln10), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"log2e",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Log2E), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"log10e",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Log10E), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"pi",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Pi), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"sqrt1_2",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(0.7071067811865476), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"sqrt2",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Sqrt2), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"abs",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Abs(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"acos",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Acos(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"acosh",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Acosh(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"asin",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Asin(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"asinh",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Asinh(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"atan",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Atan(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"atanh",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Atanh(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"atan2",
			[]types.Type{
				types.NumType{},
			},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				y := float64(inputValue.(states.NumValue))
				x := float64(argumentValues[0].(states.NumValue))
				return states.NumValue(math.Atan2(y, x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"cbrt",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Cbrt(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"ceil",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Ceil(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"clz32",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := uint32(inputValue.(states.NumValue))
				return states.NumValue(bits.LeadingZeros32(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"cos",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Cos(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"cosh",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Cos(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"exp",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Exp(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"expm1",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Expm1(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"floor",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Floor(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"fround",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(float32(x)), nil
			},
		),
		expressions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"hypot",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				v := inputValue.(*states.ArrValue)
				x := make([]float64, 0)
				for v != nil {
					x = append(x, float64(v.Head.(states.NumValue)))
					var err error
					v, err = v.GetTail()
					if err != nil {
						return nil, err
					}
				}
				hypot := varhypot.Hypot(x...)
				return states.NumValue(float64(hypot)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"imul",
			[]types.Type{
				types.NumType{},
			},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := int32(int64(inputValue.(states.NumValue)))
				y := int32(int64(argumentValues[0].(states.NumValue)))
				return states.NumValue(x * y), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"log",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"log1p",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log1p(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"log10",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log10(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"log2",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log2(x)), nil
			},
		),
		// TODO max, min with key function
		expressions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"max",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				v := inputValue.(*states.ArrValue)
				max := math.Inf(-1)
				for v != nil {
					max = math.Max(max, float64(v.Head.(states.NumValue)))
					var err error
					v, err = v.GetTail()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(max), nil
			},
		),
		expressions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"min",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				v := inputValue.(*states.ArrValue)
				min := math.Inf(1)
				for v != nil {
					min = math.Min(min, float64(v.Head.(states.NumValue)))
					var err error
					v, err = v.GetTail()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(min), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"**",
			[]types.Type{
				types.NumType{},
			},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				y := float64(argumentValues[0].(states.NumValue))
				return states.NumValue(math.Pow(x, y)), nil
			},
		),
		expressions.SimpleFuncer(
			types.AnyType{},
			"random", // TODO integer, choice
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(rand.Float64()), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"round",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Round(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"sign",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				if math.Signbit(x) {
					if x == 0 {
						return states.NumValue(math.Copysign(0, -1)), nil
					} else {
						return states.NumValue(-1), nil
					}
				} else {
					if x == 0 {
						return states.NumValue(0), nil
					} else {
						return states.NumValue(1), nil
					}
				}
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"sin",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Sin(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"sinh",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Sinh(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"sqrt",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Sqrt(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"tan",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Tan(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"tanh",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Tanh(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NumType{},
			"trunc",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Trunc(x)), nil
			},
		),
	})
}
