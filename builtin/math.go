package builtin

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"strconv"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/golang-variadic-hypot/varhypot"
)

var MathFuncers = []shapes.Funcer{
	shapes.SimpleFuncer(
		"Adds two numbers.",
		types.Num{},
		"the first summand",
		"+",
		[]*params.Param{
			params.SimpleParam("b", "the second summand", types.Num{}),
		},
		types.Num{},
		"the sum",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.NumValue(inputNum + argumentNum), nil
		},
		[]shapes.Example{
			{"1 +1", "Num", "2", nil},
		},
	),
	shapes.SimpleFuncer(
		"Subtracts a number from another.",
		types.Num{},
		"the minuend",
		"-",
		[]*params.Param{
			params.SimpleParam("b", "the subtrahend", types.Num{}),
		}, types.Num{},
		"the difference",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.NumValue(inputNum - argumentNum), nil
		},
		[]shapes.Example{
			{"5 -3", "Num", "2", nil},
		},
	),

	shapes.SimpleFuncer(
		"Returns the additive inverse of a number.",
		types.Any{},
		"any value (is ignored)",
		"-",
		[]*params.Param{
			params.SimpleParam("n", "a number", types.Num{}),
		},
		types.Num{},
		"the additive inverse (opposite number) of n",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(argumentValues[0].(states.NumValue))
			if math.Signbit(x) {
				return states.NumValue(math.Copysign(x, 1)), nil
			} else {
				return states.NumValue(math.Copysign(x, -1)), nil
			}
		},
		[]shapes.Example{
			{"-1", "Num", "-1", nil},
			{"-(-2.0)", "Num", "2", nil},
			{"-inf", "Num", "-inf", nil},
			{"-nan", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Multiplies two numbers.",
		types.Num{},
		"the first factor",
		"*",
		[]*params.Param{
			params.SimpleParam("b", "the second factor", types.Num{}),
		},
		types.Num{},
		"the product",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.NumValue(inputNum * argumentNum), nil
		},
		[]shapes.Example{
			{"2 *3", "Num", "6", nil},
		},
	),
	shapes.SimpleFuncer(
		"Divides a number by another.",
		types.Num{},
		"the dividend",
		"/",
		[]*params.Param{
			params.SimpleParam("b", "the divisor", types.Num{}),
		},
		types.Num{},
		"the quotient",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.NumValue(inputNum / argumentNum), nil
		},
		[]shapes.Example{
			{"3 /2", "Num", "1.5", nil},
		},
	),
	shapes.SimpleFuncer(
		"Remainder",
		types.Num{},
		"the dividend",
		"%",
		[]*params.Param{
			params.SimpleParam("b", "the divisor", types.Num{}),
		},
		types.Num{},
		"the remainder of integer division (rounded towards zero)",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.NumValue(math.Mod(float64(inputNum), float64(argumentNum))), nil
		},
		[]shapes.Example{
			{"3 %2", "Num", "1", nil},
			{"-8.5 %3", "Num", "-2.5", nil},
		},
	),
	shapes.SimpleFuncer(
		"Less than",
		types.Num{},
		"a number",
		"<",
		[]*params.Param{
			params.SimpleParam("b", "another number", types.Num{}),
		},
		types.Bool{},
		"true iff the input is smaller than b",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.BoolValue(inputNum < argumentNum), nil
		},
		[]shapes.Example{
			{"2 <1", "Bool", "false", nil},
			{"-inf <inf", "Bool", "true", nil},
			{"0 <0", "Bool", "false", nil},
		},
	),
	shapes.SimpleFuncer(
		"Greater than",
		types.Num{},
		"a number",
		">",
		[]*params.Param{
			params.SimpleParam("b", "another number", types.Num{}),
		},
		types.Bool{},
		"true iff the input is greater than b",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.BoolValue(inputNum > argumentNum), nil
		},
		[]shapes.Example{
			{"2 >1", "Bool", "true", nil},
			{"-inf >inf", "Bool", "false", nil},
			{"0 >0", "Bool", "false", nil},
		},
	),
	shapes.SimpleFuncer(
		"Less than or equal to",
		types.Num{},
		"a number",
		"<=",
		[]*params.Param{
			params.SimpleParam("b", "another number", types.Num{}),
		},
		types.Bool{},
		"true iff the input is less than or equal to b",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.BoolValue(inputNum <= argumentNum), nil
		},
		[]shapes.Example{
			{"2 <=1", "Bool", "false", nil},
			{"-inf <=inf", "Bool", "true", nil},
			{"0 <=0", "Bool", "true", nil},
		},
	),
	shapes.SimpleFuncer(
		"Greater than or equal to",
		types.Num{},
		"a number",
		">=",
		[]*params.Param{
			params.SimpleParam("b", "another number", types.Num{}),
		},
		types.Bool{},
		"true iff the input is greater than or equal to b",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputNum := inputValue.(states.NumValue)
			argumentNum := argumentValues[0].(states.NumValue)
			return states.BoolValue(inputNum >= argumentNum), nil
		},
		[]shapes.Example{
			{"2 >=1", "Bool", "true", nil},
			{"-inf >=inf", "Bool", "false", nil},
			{"0 >=0", "Bool", "true", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the sum of several numbers.",
		types.NewArr(types.Num{}),
		"an array of numbers",
		"sum",
		nil,
		types.Num{},
		"their sum",
		"",
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
		[]shapes.Example{
			{"[1, 2, 3, 4] sum", "Num", "10", nil},
			{"[] sum", "Num", "0", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the arithmetic mean (average) of several numbers.",
		types.NewArr(types.Num{}),
		"an array of numbers",
		"mean",
		nil,
		types.Num{},
		"their mean",
		"",
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
		[]shapes.Example{
			{"[2, 3, 5, 7] mean", "Num", "4.25", nil},
			{"[1.25] mean", "Num", "1.25", nil},
			{"[] mean", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the special number value representing positive infinity.",
		types.Any{},
		"any value (is ignored)",
		"inf",
		nil,
		types.Num{},
		"positive infinity",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Inf(1)), nil
		},
		[]shapes.Example{
			{"inf", "Num", "inf", nil},
			{"-inf", "Num", "-inf", nil},
			{"inf +inf", "Num", "inf", nil},
			{"inf -inf", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Retuns the special number value representing “not a number”.",
		types.Any{},
		"any value (is ignored)",
		"nan",
		nil,
		types.Num{},
		"not a number",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.NaN()), nil
		},
		[]shapes.Example{
			{"nan", "Num", "nan", nil},
			{"nan ==2", "Bool", "false", nil},
			{"nan ==nan", "Bool", "false", nil},
			{"-nan", "Num", "nan", nil},
		},
	),

	shapes.SimpleFuncer("", types.Num{}, "", "isFinite", nil, types.Bool{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.BoolValue(!math.IsInf(x, 0)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "isNaN", nil, types.Bool{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.BoolValue(math.IsNaN(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "epsilon", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.Nextafter(1, 2) - 1), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "largestSafeInteger", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(9007199254740991), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "largestNum", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.MaxFloat64), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "smallestSafeInteger", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(-9007199254740991), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "smallestPositiveNum", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.SmallestNonzeroFloat64), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "isInteger", nil, types.Bool{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.BoolValue(x == float64(int(x))), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "isSafeInteger", nil, types.Bool{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.BoolValue(x >= -9007199254740991 && x <= 9007199254740991), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "toBase", []*params.Param{
		params.SimpleParam("base", "", types.Num{}),
	}, types.Str{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		xInt := int64(x)
		if x != float64(xInt) {
			return nil, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(inputValue),
				errors.Message("base conversion for non-integers not yet supported"),
			)
		}
		radix := float64(argumentValues[0].(states.NumValue))
		radixInt := int(radix)
		if radix != float64(radixInt) || radixInt < 2 || radixInt > 36 {
			return nil, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(argumentValues[0]),
				errors.Message("radix must be an integer between 2 and 36 (inclusive)"),
			)
		}
		return states.StrValue(strconv.FormatInt(xInt, radixInt)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "toExponential", nil, types.Str{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.StrValue(fmt.Sprintf("%e", x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "toExponential", []*params.Param{
		params.SimpleParam("precision", "", types.Num{}),
	}, types.Str{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		f := int(argumentValues[0].(states.NumValue))
		format := "%." + strconv.Itoa(f) + "e"
		return states.StrValue(fmt.Sprintf(format, x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "toFixed", []*params.Param{
		params.SimpleParam("precision", "", types.Num{}),
	}, types.Str{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		f := int(argumentValues[0].(states.NumValue))
		format := "%." + strconv.Itoa(f) + "f"
		return states.StrValue(fmt.Sprintf(format, x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "e", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.E), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "ln2", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.Ln2), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "ln10", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.Ln10), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "log2e", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.Log2E), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "log10e", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.Log10E), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "pi", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.Pi), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "sqrt1_2", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(0.7071067811865476), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "sqrt2", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(math.Sqrt2), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "abs", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Abs(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "acos", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Acos(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "acosh", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Acosh(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "asin", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Asin(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "asinh", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Asinh(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "atan", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Atan(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "atanh", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Atanh(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "atan2", []*params.Param{
		params.SimpleParam("x", "", types.Num{}),
	}, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		y := float64(inputValue.(states.NumValue))
		x := float64(argumentValues[0].(states.NumValue))
		return states.NumValue(math.Atan2(y, x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "cbrt", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Cbrt(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "ceil", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Ceil(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "clz32", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := uint32(inputValue.(states.NumValue))
		return states.NumValue(bits.LeadingZeros32(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "cos", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Cos(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "cosh", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Cos(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "exp", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Exp(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "expm1", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Expm1(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "floor", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Floor(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "fround", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(float32(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.NewArr(types.Num{}), "", "hypot", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
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
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "imul", []*params.Param{
		params.SimpleParam("y", "", types.Num{}),
	}, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := int32(int64(inputValue.(states.NumValue)))
		y := int32(int64(argumentValues[0].(states.NumValue)))
		return states.NumValue(x * y), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "log", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Log(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "log1p", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Log1p(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "log10", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Log10(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "log2", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Log2(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.NewArr(types.Num{}), "", "max", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
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
	}, nil),

	shapes.SimpleFuncer("", types.NewArr(types.Num{}), "", "min", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
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
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "**", []*params.Param{
		params.SimpleParam("y", "", types.Num{}),
	}, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		y := float64(argumentValues[0].(states.NumValue))
		return states.NumValue(math.Pow(x, y)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "random", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		return states.NumValue(rand.Float64()), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "round", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Round(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "sign", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
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
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "sin", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Sin(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "sinh", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Sinh(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "sqrt", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Sqrt(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "tan", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Tan(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "tanh", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Tanh(x)), nil
	}, nil),

	shapes.SimpleFuncer("", types.Num{}, "", "trunc", nil, types.Num{}, "", "", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
		x := float64(inputValue.(states.NumValue))
		return states.NumValue(math.Trunc(x)), nil
	}, nil),
}
