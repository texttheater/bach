package builtin

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"strconv"
	"time"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/golang-variadic-hypot/varhypot"
)

func initMath() {
	rand.Seed(time.Now().UnixNano())
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.FuncerDefinition{
		// for Num +Num Num
		expressions.SimpleFuncer(
			types.Num{},
			"+",
			[]types.Type{types.Num{}},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum + argumentNum), nil
			},
		),
		// for Num -Num Num
		expressions.SimpleFuncer(
			types.Num{},
			"-",
			[]types.Type{types.Num{}},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum - argumentNum), nil
			},
		),
		// for Any -Num Num
		expressions.SimpleFuncer(
			types.Any{},
			"-",
			[]types.Type{types.Num{}},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(argumentValues[0].(states.NumValue))
				if math.Signbit(x) {
					return states.NumValue(math.Copysign(x, 1)), nil
				} else {
					return states.NumValue(math.Copysign(x, -1)), nil
				}
			},
		),
		// for Num *Num Num
		expressions.SimpleFuncer(
			types.Num{},
			"*",
			[]types.Type{types.Num{}},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum * argumentNum), nil
			},
		),
		// for Num /Num Num
		expressions.SimpleFuncer(
			types.Num{},
			"/",
			[]types.Type{types.Num{}},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum / argumentNum), nil
			},
		),
		// for Num %Num Num
		expressions.SimpleFuncer(
			types.Num{},
			"%",
			[]types.Type{types.Num{}},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(math.Mod(float64(inputNum), float64(argumentNum))), nil
			},
		),
		// for Num <Num Bool
		expressions.SimpleFuncer(
			types.Num{},
			"<",
			[]types.Type{types.Num{}},
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum < argumentNum), nil
			},
		),
		// for Num >Num Bool
		expressions.SimpleFuncer(
			types.Num{},
			">",
			[]types.Type{types.Num{}},
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum > argumentNum), nil
			},
		),
		// for Num <=Num Bool
		expressions.SimpleFuncer(
			types.Num{},
			"<=",
			[]types.Type{types.Num{}},
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum <= argumentNum), nil
			},
		),
		// for Num >=Num Bool
		expressions.SimpleFuncer(
			types.Num{},
			">=",
			[]types.Type{types.Num{}},
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum >= argumentNum), nil
			},
		),
		// for Arr<Num> sum Num
		expressions.SimpleFuncer(
			types.NewArr(types.Num{}),
			"sum",
			nil,
			types.Num{},
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
		// for Arr<Num> mean Num
		expressions.SimpleFuncer(
			types.NewArr(types.Num{}),
			"mean",
			nil,
			types.Num{},
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
		// for Any inf Num
		expressions.SimpleFuncer(
			types.Any{},
			"inf",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Inf(1)), nil
			},
		),
		// for Any nan Num
		expressions.SimpleFuncer(
			types.Any{},
			"nan",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.NaN()), nil
			},
		),
		// for Num isFinite Bool
		expressions.SimpleFuncer(
			types.Num{},
			"isFinite",
			nil,
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(!math.IsInf(x, 0)), nil
			},
		),
		// for Num isNaN Bool
		expressions.SimpleFuncer(
			types.Num{},
			"isNaN",
			nil,
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(math.IsNaN(x)), nil
			},
		),
		// for Any epsilon Num
		expressions.SimpleFuncer(
			types.Any{},
			"epsilon",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Nextafter(1, 2) - 1), nil
			},
		),
		// for Any largestSafeInteger Num
		expressions.SimpleFuncer(
			types.Any{},
			"largestSafeInteger",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(9007199254740991), nil
			},
		),
		// for Any largestNum Num
		expressions.SimpleFuncer(
			types.Any{},
			"largestNum",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.MaxFloat64), nil
			},
		),
		// for Any smallestSafeInteger Num
		expressions.SimpleFuncer(
			types.Any{},
			"smallestSafeInteger",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(-9007199254740991), nil
			},
		),
		// for Any smallestPositiveNum Num
		expressions.SimpleFuncer(
			types.Any{},
			"smallestPositiveNum",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.SmallestNonzeroFloat64), nil
			},
		),
		// for Num isInteger Bool
		expressions.SimpleFuncer(
			types.Num{},
			"isInteger",
			nil,
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(x == float64(int(x))), nil
			},
		),
		// for Num isSafeInteger Bool
		expressions.SimpleFuncer(
			types.Num{},
			"isSafeInteger",
			nil,
			types.Bool{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.BoolValue(x >= -9007199254740991 && x <= 9007199254740991), nil
			},
		),
		// for Num toBase Str
		expressions.SimpleFuncer(
			types.Num{},
			"toBase",
			[]types.Type{types.Num{}},
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				xInt := int64(x)
				if x != float64(xInt) {
					return nil, errors.ValueError(
						errors.Code(errors.UnexpectedValue),
						errors.GotValue(inputValue),
						errors.Message("base conversion for non-integers not yet supported"),
						// TODO add pos
					)
				}
				radix := float64(argumentValues[0].(states.NumValue))
				radixInt := int(radix)
				if radix != float64(radixInt) || radixInt < 2 || radixInt > 36 {
					return nil, errors.ValueError(
						errors.Code(errors.UnexpectedValue),
						errors.GotValue(argumentValues[0]),
						errors.Message("radix must be an integer between 2 and 36 (inclusive)"),
						// TODO add pos
					)
				}
				return states.StrValue(strconv.FormatInt(xInt, radixInt)), nil
			},
		),
		// for Num toExponential Str
		expressions.SimpleFuncer(
			types.Num{},
			"toExponential",
			nil,
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.StrValue(fmt.Sprintf("%e", x)), nil
			},
		),
		// for Num toExponential(Num) Str
		expressions.SimpleFuncer(
			types.Num{},
			"toExponential",
			[]types.Type{types.Num{}},
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				f := int(argumentValues[0].(states.NumValue))
				format := "%." + strconv.Itoa(f) + "e"
				return states.StrValue(fmt.Sprintf(format, x)), nil
			},
		),
		// for Num toFixed(Num) Str
		expressions.SimpleFuncer(
			types.Num{},
			"toFixed",
			[]types.Type{types.Num{}},
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				f := int(argumentValues[0].(states.NumValue))
				format := "%." + strconv.Itoa(f) + "f"
				return states.StrValue(fmt.Sprintf(format, x)), nil
			},
		),
		// for Any e Num
		expressions.SimpleFuncer(
			types.Any{},
			"e",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.E), nil
			},
		),
		// for Any ln2 Num
		expressions.SimpleFuncer(
			types.Any{},
			"ln2",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Ln2), nil
			},
		),
		// for Any ln10 Num
		expressions.SimpleFuncer(
			types.Any{},
			"ln10",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Ln10), nil
			},
		),
		// for Any log2e Num
		expressions.SimpleFuncer(
			types.Any{},
			"log2e",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Log2E), nil
			},
		),
		// for Any log10e Num
		expressions.SimpleFuncer(
			types.Any{},
			"log10e",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Log10E), nil
			},
		),
		// for Any pi Num
		expressions.SimpleFuncer(
			types.Any{},
			"pi",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Pi), nil
			},
		),
		// for Any sqrt1_2 Num
		expressions.SimpleFuncer(
			types.Any{},
			"sqrt1_2",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(0.7071067811865476), nil
			},
		),
		// for Any sqrt2 Num
		expressions.SimpleFuncer(
			types.Any{},
			"sqrt2",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Sqrt2), nil
			},
		),
		// for Num abs Num
		expressions.SimpleFuncer(
			types.Num{},
			"abs",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Abs(x)), nil
			},
		),
		// for Num acos Num
		expressions.SimpleFuncer(
			types.Num{},
			"acos",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Acos(x)), nil
			},
		),
		// for Num acosh Num
		expressions.SimpleFuncer(
			types.Num{},
			"acosh",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Acosh(x)), nil
			},
		),
		// for Num asin Num
		expressions.SimpleFuncer(
			types.Num{},
			"asin",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Asin(x)), nil
			},
		),
		// for Num asinh Num
		expressions.SimpleFuncer(
			types.Num{},
			"asinh",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Asinh(x)), nil
			},
		),
		// for Num atan Num
		expressions.SimpleFuncer(
			types.Num{},
			"atan",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Atan(x)), nil
			},
		),
		// for Num atanh Num
		expressions.SimpleFuncer(
			types.Num{},
			"atanh",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Atanh(x)), nil
			},
		),
		// for Num atan2(Num) Num
		expressions.SimpleFuncer(
			types.Num{},
			"atan2",
			[]types.Type{
				types.Num{},
			},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				y := float64(inputValue.(states.NumValue))
				x := float64(argumentValues[0].(states.NumValue))
				return states.NumValue(math.Atan2(y, x)), nil
			},
		),
		// for Num cbrt Num
		expressions.SimpleFuncer(
			types.Num{},
			"cbrt",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Cbrt(x)), nil
			},
		),
		// for Num ceil Num
		expressions.SimpleFuncer(
			types.Num{},
			"ceil",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Ceil(x)), nil
			},
		),
		// for Num clz32 Num
		expressions.SimpleFuncer(
			types.Num{},
			"clz32",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := uint32(inputValue.(states.NumValue))
				return states.NumValue(bits.LeadingZeros32(x)), nil
			},
		),
		// for Num cos Num
		expressions.SimpleFuncer(
			types.Num{},
			"cos",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Cos(x)), nil
			},
		),
		// for Num cosh Num
		expressions.SimpleFuncer(
			types.Num{},
			"cosh",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Cos(x)), nil
			},
		),
		// for Num exp Num
		expressions.SimpleFuncer(
			types.Num{},
			"exp",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Exp(x)), nil
			},
		),
		// for Num expm1 Num
		expressions.SimpleFuncer(
			types.Num{},
			"expm1",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Expm1(x)), nil
			},
		),
		// for Num floor Num
		expressions.SimpleFuncer(
			types.Num{},
			"floor",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Floor(x)), nil
			},
		),
		// for Num fround Num
		expressions.SimpleFuncer(
			types.Num{},
			"fround",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(float32(x)), nil
			},
		),
		// for Arr<Num> hypot Num
		expressions.SimpleFuncer(
			types.NewArr(types.Num{}),
			"hypot",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				v := inputValue.(*states.ArrValue)
				x := make([]float64, 0)
				for v != nil {
					x = append(x, float64(v.Head.(states.NumValue)))
					var err error
					v, err = v.Tail.EvalArr()
					if err != nil {
						return nil, err
					}
				}
				hypot := varhypot.Hypot(x...)
				return states.NumValue(float64(hypot)), nil
			},
		),
		// for Num imul(Num) Num
		expressions.SimpleFuncer(
			types.Num{},
			"imul",
			[]types.Type{
				types.Num{},
			},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := int32(int64(inputValue.(states.NumValue)))
				y := int32(int64(argumentValues[0].(states.NumValue)))
				return states.NumValue(x * y), nil
			},
		),
		// for Num log Num
		expressions.SimpleFuncer(
			types.Num{},
			"log",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log(x)), nil
			},
		),
		// for Num log1p Num
		expressions.SimpleFuncer(
			types.Num{},
			"log1p",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log1p(x)), nil
			},
		),
		// for Num log10 Num
		expressions.SimpleFuncer(
			types.Num{},
			"log10",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log10(x)), nil
			},
		),
		// for Num log2 Num
		expressions.SimpleFuncer(
			types.Num{},
			"log2",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Log2(x)), nil
			},
		),
		expressions.SimpleFuncer(
			types.NewArr(types.Num{}),
			"max",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				v := inputValue.(*states.ArrValue)
				max := math.Inf(-1)
				for v != nil {
					max = math.Max(max, float64(v.Head.(states.NumValue)))
					var err error
					v, err = v.Tail.EvalArr()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(max), nil
			},
		),
		expressions.SimpleFuncer(
			types.NewArr(types.Num{}),
			"min",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				v := inputValue.(*states.ArrValue)
				min := math.Inf(1)
				for v != nil {
					min = math.Min(min, float64(v.Head.(states.NumValue)))
					var err error
					v, err = v.Tail.EvalArr()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(min), nil
			},
		),
		// for Num **Num Num
		expressions.SimpleFuncer(
			types.Num{},
			"**",
			[]types.Type{
				types.Num{},
			},
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				y := float64(argumentValues[0].(states.NumValue))
				return states.NumValue(math.Pow(x, y)), nil
			},
		),
		// for Any random Num
		expressions.SimpleFuncer(
			types.Any{},
			"random", // TODO integer, choice
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(rand.Float64()), nil
			},
		),
		// for Num round Num
		expressions.SimpleFuncer(
			types.Num{},
			"round",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Round(x)), nil
			},
		),
		// for Num sign Num
		expressions.SimpleFuncer(
			types.Num{},
			"sign",
			nil,
			types.Num{},
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
		// for Num sin Num
		expressions.SimpleFuncer(
			types.Num{},
			"sin",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Sin(x)), nil
			},
		),
		// for Num sinh Num
		expressions.SimpleFuncer(
			types.Num{},
			"sinh",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Sinh(x)), nil
			},
		),
		// for Num sqrt Num
		expressions.SimpleFuncer(
			types.Num{},
			"sqrt",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Sqrt(x)), nil
			},
		),
		// for Num tan Num
		expressions.SimpleFuncer(
			types.Num{},
			"tan",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Tan(x)), nil
			},
		),
		// for Num tanh Num
		expressions.SimpleFuncer(
			types.Num{},
			"tanh",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Tanh(x)), nil
			},
		),
		// for Num trunc Num
		expressions.SimpleFuncer(
			types.Num{},
			"trunc",
			nil,
			types.Num{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				x := float64(inputValue.(states.NumValue))
				return states.NumValue(math.Trunc(x)), nil
			},
		),
	})
}
