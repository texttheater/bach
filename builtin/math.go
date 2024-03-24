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
		"Returns infinity.",
		types.Any{},
		"any value (is ignored)",
		"inf",
		nil,
		types.Num{},
		"the special number value representing positive infinity",
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
		"Returns NaN.",
		types.Any{},
		"any value (is ignored)",
		"nan",
		nil,
		types.Num{},
		"the special number value representing “not a number”",
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
	shapes.SimpleFuncer(
		"Checks whether a number is finite.",
		types.Num{},
		"a number",
		"isFinite",
		nil,
		types.Bool{},
		"true unless the input is positive or negative infinity",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.BoolValue(!math.IsInf(x, 0)), nil
		},
		[]shapes.Example{
			{"1024 isFinite", "Bool", "true", nil},
			{"-1024 isFinite", "Bool", "true", nil},
			{"inf isFinite", "Bool", "false", nil},
			{"-inf isFinite", "Bool", "false", nil},
			{"nan isFinite", "Bool", "true", nil},
		},
	),
	shapes.SimpleFuncer(
		"Checks whether a number is NaN.",
		types.Num{},
		"a number",
		"isNaN",
		nil,
		types.Bool{},
		"true iff the input is NaN",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.BoolValue(math.IsNaN(x)), nil
		},
		[]shapes.Example{
			{"1024 isNaN", "Bool", "false", nil},
			{"nan isNaN", "Bool", "true", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the floating point epsilon.",
		types.Any{},
		"any value (is ignored)",
		"epsilon",
		nil,
		types.Num{},
		"the difference between 1 and the smallest floating point number greater than 1",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Nextafter(1, 2) - 1), nil
		},
		[]shapes.Example{
			{"epsilon", "Num", "2.220446049250313e-16", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the largest safe integer.",
		types.Any{},
		"any value (is ignored)",
		"largestSafeInteger",
		nil,
		types.Num{},
		"the largest integer that can be represented as an IEEE-754 double precision number and cannot be the result of rounding another number to fit IEEE-754",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(9007199254740991), nil
		},
		[]shapes.Example{
			{"largestSafeInteger +0", "Num", "9007199254740991", nil},
			{"largestSafeInteger +1", "Num", "9007199254740992", nil},
			{"largestSafeInteger +2", "Num", "9007199254740992", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the largest representable number.",
		types.Any{},
		"any value (is ignored)",
		"largestNum",
		nil,
		types.Num{},
		"the largest number representable as an IEEE-754 double precision number",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.MaxFloat64), nil
		},
		[]shapes.Example{
			{"largestNum", "Num", "1.7976931348623157e+308", nil},
			{"largestNum +largestNum", "Num", "inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the smallest safe integer.",
		types.Any{},
		"any value (is ignored)",
		"smallestSafeInteger",
		nil,
		types.Num{},
		"the smallest integer that can be represented as an IEEE-754 double precision number and cannot be the result of rounding another integer to fit the IEEE-754 double precision representation",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(-9007199254740991), nil
		},
		[]shapes.Example{
			{"smallestSafeInteger -0", "Num", "-9007199254740991", nil},
			{"smallestSafeInteger -1", "Num", "-9007199254740992", nil},
			{"smallestSafeInteger -2", "Num", "-9007199254740992", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the smallest representable positive number.",
		types.Any{},
		"any value (is ignored)",
		"smallestPositiveNum",
		nil,
		types.Num{},
		"the smallest representable positive number",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.SmallestNonzeroFloat64), nil
		},
		[]shapes.Example{
			{"smallestPositiveNum", "Num", "5e-324", nil},
		},
	),
	shapes.SimpleFuncer(
		"Checks whether a number is integer.",
		types.Num{},
		"a number",
		"isInteger",
		nil,
		types.Bool{},
		"true iff the input represents a whole number",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.BoolValue(x == float64(int(x))), nil
		},
		[]shapes.Example{
			{"1.0 isInteger", "Bool", "true", nil},
			{"1.1 isInteger", "Bool", "false", nil},
		},
	),
	shapes.SimpleFuncer(
		"Checks whether a number is a safe integer.",
		types.Num{},
		"a number",
		"isSafeInteger",
		nil,
		types.Bool{},
		"true iff the input is an integer and cannot be the result of rounding another integer to fit the IEEE-754 double precision representation",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.BoolValue(x == float64(int(x)) && x >= -9007199254740991 && x <= 9007199254740991), nil
		},
		[]shapes.Example{
			{"-9007199254740992 isSafeInteger", "Bool", "false", nil},
			{"-9007199254740991 isSafeInteger", "Bool", "true", nil},
			{"0 isSafeInteger", "Bool", "true", nil},
			{"0.1 isSafeInteger", "Bool", "false", nil},
			{"9007199254740991 isSafeInteger", "Bool", "true", nil},
			{"9007199254740992 isSafeInteger", "Bool", "false", nil},
		},
	),
	shapes.SimpleFuncer(
		"Converts an integer to a specified base.",
		types.Num{},
		"an integer",
		"toBase",
		[]*params.Param{
			params.SimpleParam("base", "an integer between 2 and 36 (inclusive)", types.Num{}),
		},
		types.Str{},
		"a string representation of the input in the specified base",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
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
		},
		[]shapes.Example{
			{"233 toBase(16)", "Str", `"e9"`, nil},
			{"11 toBase(16)", "Str", `"b"`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Converts a number to exponential notation.",
		types.Num{},
		"a number",
		"toExponential",
		nil,
		types.Str{},
		"a string representation of the input in exponential notation with 6 digits after the decimal point",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.StrValue(fmt.Sprintf("%e", x)), nil
		},
		[]shapes.Example{
			{"77.1234 toExponential", "Str", `"7.712340e+01"`, nil},
			{"77 toExponential", "Str", `"7.700000e+01"`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Converts a number to exponential notation.",
		types.Num{},
		"a number",
		"toExponential",
		[]*params.Param{
			params.SimpleParam("precision", "the number digits after the decimal point", types.Num{}),
		},
		types.Str{},
		"a string representation of the input in exponential notation with the specified number of digits after the decimal point",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			f := int(argumentValues[0].(states.NumValue))
			format := "%." + strconv.Itoa(f) + "e"
			return states.StrValue(fmt.Sprintf(format, x)), nil
		},
		[]shapes.Example{
			{"77.1234 toExponential(4)", "Str", `"7.7123e+01"`, nil},
			{"77.1234 toExponential(2)", "Str", `"7.71e+01"`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Converts a number to fixed-point notation.",
		types.Num{},
		"a number",
		"toFixed",
		[]*params.Param{
			params.SimpleParam("precision", "the number digits after the decimal point", types.Num{}),
		},
		types.Str{},
		"a rounded string representation of the input with the specified number of digits after the decimal point",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			f := int(argumentValues[0].(states.NumValue))
			format := "%." + strconv.Itoa(f) + "f"
			return states.StrValue(fmt.Sprintf(format, x)), nil
		},
		[]shapes.Example{
			{"123.456 toFixed(2)", "Str", `"123.46"`, nil},
			{"0.004 toFixed(2)", "Str", `"0.00"`, nil},
			{"1.23e+5 toFixed(2)", "Str", `"123000.00"`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns Euler's number.",
		types.Any{},
		"any value (is ignored)",
		"e",
		nil,
		types.Num{},
		"an approximation of Euler's number",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.E), nil
		},
		[]shapes.Example{
			{"e", "Num", "2.718281828459045", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the natural logarithm of 2.",
		types.Any{},
		"any value (is ignored)",
		"ln2",
		nil,
		types.Num{},
		"the approximate natural logarithm of 2",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Ln2), nil
		},
		[]shapes.Example{
			{"ln2", "Num", "0.6931471805599453", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the natural logarithm of 10.",
		types.Any{},
		"any value (is ignored)",
		"ln10",
		nil,
		types.Num{},
		"the approximate natural logarithm of 10",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Ln10), nil
		},
		[]shapes.Example{
			{"ln10", "Num", "2.302585092994046", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the base-2 logarithm of e",
		types.Any{},
		"any value (is ignored)",
		"log2e",
		nil,
		types.Num{},
		"the approximate base-2 logarithm of Euler's number",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Log2E), nil
		},
		[]shapes.Example{
			{"log2e", "Num", "1.4426950408889634", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the base-10 logarithm of e",
		types.Any{},
		"any value (is ignored)",
		"log10e",
		nil,
		types.Num{},
		"the approximate base-10 logarithm of Euler's number",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Log10E), nil
		},
		[]shapes.Example{
			{"log10e", "Num", "0.4342944819032518", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns pi.",
		types.Any{},
		"any value (is ignored)",
		"pi",
		nil,
		types.Num{},
		"an approximation of pi",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Pi), nil
		},
		[]shapes.Example{
			{"pi", "Num", "3.141592653589793", nil},
			{"for Num def radiusToCircumference Num as *2 *pi ok 10 radiusToCircumference", "Num", "62.83185307179586", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the square root of 1/2.",
		types.Any{},
		"any value (is ignored)",
		"sqrt1_2",
		nil,
		types.Num{},
		"the approximate square root of 1/2",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(0.7071067811865476), nil
		},
		[]shapes.Example{
			{"sqrt1_2", "Num", "0.7071067811865476", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the square root of 2.",
		types.Any{},
		"any value (is ignored)",
		"sqrt2",
		nil,
		types.Num{},
		"the approximate square root of 2",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(math.Sqrt2), nil
		},
		[]shapes.Example{
			{"sqrt2", "Num", "1.4142135623730951", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the absolute value.",
		types.Num{},
		"a number",
		"abs",
		nil,
		types.Num{},
		"the absolute value of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Abs(x)), nil
		},
		[]shapes.Example{
			{"3 -5 abs", "Num", "2", nil},
			{"5 -3 abs", "Num", "2", nil},
			{"1.23456 -7.89012 abs", "Num", "6.6555599999999995", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the inverse cosine.",
		types.Num{},
		"a number in the interval [-1, 1]",
		"acos",
		nil,
		types.Num{},
		"the inverse cosine (in radians) of the input, or `nan` if the input is invalid",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Acos(x)), nil
		},
		[]shapes.Example{
			{"-2 acos", "Num", "nan", nil},
			{"-1 acos", "Num", "3.141592653589793", nil},
			{"0 acos", "Num", "1.5707963267948966", nil},
			{"1 acos", "Num", "0", nil},
			{"1.1 acos", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the inverse hyperbolic cosine.",
		types.Num{},
		"a number greater than or equal to 1",
		"acosh",
		nil,
		types.Num{},
		"the inverse hyperoblic cosine of the input, or `nan` if the input is invalid",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Acosh(x)), nil
		},
		[]shapes.Example{
			{"0.9 acosh", "Num", "nan", nil},
			{"1 acosh", "Num", "0", nil},
			{"10 acosh", "Num", "2.993222846126381", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the inverse sine.",
		types.Num{},
		"a number in the interval [-1, 1]",
		"asin",
		nil,
		types.Num{},
		"the inverse sine (in radians) of the input, of nan if the input is invalid",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Asin(x)), nil
		},
		[]shapes.Example{
			{"-2 asin", "Num", "nan", nil},
			{"-1 asin", "Num", "-1.5707963267948966", nil},
			{"0 asin", "Num", "0", nil},
			{"1 asin", "Num", "1.5707963267948966", nil},
			{"1.1 asin", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the inverse hyperbolic sine.",
		types.Num{},
		"a number",
		"asinh",
		nil,
		types.Num{},
		"the inverse hyperbolic sine of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Asinh(x)), nil
		},
		[]shapes.Example{
			{"-1 asinh", "Num", "-0.881373587019543", nil},
			{"0 asinh", "Num", "0", nil},
			{"1 asinh", "Num", "0.881373587019543", nil},
			{"2 asinh", "Num", "1.4436354751788103", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the inverse tangent.",
		types.Num{},
		"a number",
		"atan",
		nil,
		types.Num{},
		"the inverse tangent (in radians) of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Atan(x)), nil
		},
		[]shapes.Example{
			{"-10 atan", "Num", "-1.4711276743037345", nil},
			{"-1 atan", "Num", "-0.7853981633974483", nil},
			{"0 atan", "Num", "0", nil},
			{"1 atan", "Num", "0.7853981633974483", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the angle in the plane.",
		types.Num{},
		"a number y (y-coordinate)",
		"atan2",
		[]*params.Param{
			params.SimpleParam("x", "a number (x-coordinate)", types.Num{}),
		},
		types.Num{},
		"the angle in the plane (in radians) between the positive x-axis and the ray from (0, 0) to (x, y)",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			y := float64(inputValue.(states.NumValue))
			x := float64(argumentValues[0].(states.NumValue))
			return states.NumValue(math.Atan2(y, x)), nil
		},
		[]shapes.Example{
			{"5 atan2(5)", "Num", "0.7853981633974483", nil},
			{"10 atan2(10)", "Num", "0.7853981633974483", nil},
			{"10 atan2(0)", "Num", "1.5707963267948966", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the inverse hyperbolic tangent.",
		types.Num{},
		"a number in the interval [-1, 1]",
		"atanh",
		nil,
		types.Num{},
		"the inverse hyperbolic tangent of the input, or `nan` if the input is invalid",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Atanh(x)), nil
		},
		[]shapes.Example{
			{"-2 atanh", "Num", "nan", nil},
			{"-1 atanh", "Num", "-inf", nil},
			{"0 atanh", "Num", "0", nil},
			{"0.5 atanh", "Num", "0.5493061443340548", nil},
			{"1 atanh", "Num", "inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the cube root.",
		types.Num{},
		"a number",
		"cbrt",
		nil,
		types.Num{},
		"the cube root of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Cbrt(x)), nil
		},
		[]shapes.Example{
			{"-1 cbrt", "Num", "-1", nil},
			{"1 cbrt", "Num", "1", nil},
			{"inf cbrt", "Num", "inf", nil},
			{"64 cbrt", "Num", "4", nil},
		},
	),
	shapes.SimpleFuncer(
		"Rounds a number up.",
		types.Num{},
		"a number",
		"ceil",
		nil,
		types.Num{},
		"the smallest integer greater than or equal to the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Ceil(x)), nil
		},
		[]shapes.Example{
			{".95 ceil", "Num", "1", nil},
			{"4 ceil", "Num", "4", nil},
			{"7.004 ceil", "Num", "8", nil},
			{"-7.004 ceil", "Num", "-7", nil},
		},
	),
	shapes.SimpleFuncer(
		"Count leading zeros.",
		types.Num{},
		"a number (is truncated to integer)",
		"clz32",
		nil,
		types.Num{},
		"the number of leading zero bits in the 32-bit binary representation of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := uint32(inputValue.(states.NumValue))
			return states.NumValue(bits.LeadingZeros32(x)), nil
		},
		[]shapes.Example{
			{"-inf clz32", "Num", "32", nil},
			{"-4 clz32", "Num", "0", nil},
			{"-1 clz32", "Num", "0", nil},
			{"0 clz32", "Num", "32", nil},
			{"0.5 clz32", "Num", "32", nil},
			{"1 clz32", "Num", "31", nil},
			{"1.1 clz32", "Num", "31", nil},
			{"4 clz32", "Num", "29", nil},
			{"4.7 clz32", "Num", "29", nil},
			{"1000 clz32", "Num", "22", nil},
			{"inf clz32", "Num", "32", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the cosine.",
		types.Num{},
		"an angle in radians",
		"cos",
		nil,
		types.Num{},
		"the cosine of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Cos(x)), nil
		},
		[]shapes.Example{
			{"-inf cos", "Num", "nan", nil},
			{"-0 cos", "Num", "1", nil},
			{"0 cos", "Num", "1", nil},
			{"1 cos", "Num", "0.5403023058681398", nil},
			{"pi cos", "Num", "-1", nil},
			{"pi *2 cos", "Num", "1", nil},
			{"inf cos", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the hyperbolic cosine.",
		types.Num{},
		"a number",
		"cosh",
		nil,
		types.Num{},
		"the hyperbolic cosine of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Cosh(x)), nil
		},
		[]shapes.Example{
			{"0 cosh", "Num", "1", nil},
			{"1 cosh", "Num", "1.5430806348152437", nil},
			{"-1 cosh", "Num", "1.5430806348152437", nil},
			{"2 cosh", "Num", "3.7621956910836314", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the exponential function.",
		types.Num{},
		"the exponent",
		"exp",
		nil,
		types.Num{},
		"e (Euler's number) raised to the exponent",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Exp(x)), nil
		},
		[]shapes.Example{
			{"-inf exp", "Num", "0", nil},
			{"-1 exp", "Num", "0.36787944117144233", nil},
			{"0 exp", "Num", "1", nil},
			{"1 exp", "Num", "2.718281828459045", nil},
			{"inf exp", "Num", "inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the exponential function, subtracted by one.",
		types.Num{},
		"the exponent",
		"expm1",
		nil,
		types.Num{},
		"",
		"e (Euler's number) raised to the exponent, minus 1",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Expm1(x)), nil
		},
		[]shapes.Example{
			{"-inf expm1", "Num", "-1", nil},
			{"-1 expm1", "Num", "-0.6321205588285577", nil},
			{"-0 expm1", "Num", "-0", nil},
			{"0 expm1", "Num", "0", nil},
			{"1 expm1", "Num", "1.718281828459045", nil},
			{"inf expm1", "Num", "inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Rounds down.",
		types.Num{},
		"a number",
		"floor",
		nil,
		types.Num{},
		"the largest integer less than or equal to the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Floor(x)), nil
		},
		[]shapes.Example{
			{"5.95 floor", "Num", "5", nil},
			{"5.05 floor", "Num", "5", nil},
			{"5 floor", "Num", "5", nil},
			{"-5.05 floor", "Num", "-6", nil},
		},
	),
	shapes.SimpleFuncer(
		"Rounds to 32-bit precision.",
		types.Num{},
		"a number",
		"fround",
		nil,
		types.Num{},
		"the nearest 32-bit single precision float representation",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(float32(x)), nil
		},
		[]shapes.Example{
			{"5.5 fround", "Num", "5.5", nil},
			{"5.05 fround", "Num", "5.050000190734863", nil},
			{"5 fround", "Num", "5", nil},
			{"-5.05 fround", "Num", "-5.050000190734863", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the square root of the sum of squares",
		types.NewArr(types.Num{}),
		"an array of numbers",
		"hypot",
		nil,
		types.Num{},
		"the square root of the sum of the squares of the input numbers",
		"",
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
		[]shapes.Example{
			{"[3, 4] hypot", "Num", "5", nil},
			{"[5, 12] hypot", "Num", "13", nil},
			{"[3, 4, 5] hypot", "Num", "7.0710678118654755", nil},
			{"[-5] hypot", "Num", "5", nil},
		},
	),
	shapes.SimpleFuncer(
		"32-bit multiplication",
		types.Num{},
		"the first factor",
		"imul",
		[]*params.Param{
			params.SimpleParam("y", "the second factor", types.Num{}),
		},
		types.Num{},
		"the product of the 32-bit versions (cf. fround) of the factors",
		"", func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := int32(int64(inputValue.(states.NumValue)))
			y := int32(int64(argumentValues[0].(states.NumValue)))
			return states.NumValue(x * y), nil
		},
		[]shapes.Example{
			{"3 imul(4)", "Num", "12", nil},
			{"-5 imul(12)", "Num", "-60", nil},
			{`"ffffffff" parseInt(16) imul(5)`, "Num", "-5", nil},
			{`"fffffffe" parseInt(16) imul(5)`, "Num", "-10", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the natural logarithm.",
		types.Num{},
		"a number",
		"log",
		nil,
		types.Num{},
		"the natural (base e) logarithm of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Log(x)), nil
		},
		[]shapes.Example{
			{"-1 log", "Num", "nan", nil},
			{"-0 log", "Num", "-inf", nil},
			{"0 log", "Num", "-inf", nil},
			{"1 log", "Num", "0", nil},
			{"10 log", "Num", "2.302585092994046", nil},
			{"inf log", "Num", "inf", nil},
			{"8 log /(2 log)", "Num", "3", nil},
			{"625 log /(5 log)", "Num", "4", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the base 10 logarithm.",
		types.Num{},
		"a number",
		"log10",
		nil,
		types.Num{},
		"the base 10 logarithm of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Log10(x)), nil
		},
		[]shapes.Example{
			{"-2 log10", "Num", "nan", nil},
			{"-0 log10", "Num", "-inf", nil},
			{"0 log10", "Num", "-inf", nil},
			{"1 log10", "Num", "0", nil},
			{"2 log10", "Num", "0.3010299956639812", nil},
			{"100000 log10", "Num", "5", nil},
			{"inf log10", "Num", "inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the natural logarithm of x + 1.",
		types.Num{},
		"a number (x)",
		"log1p",
		nil,
		types.Num{},
		"",
		"the natural (base e) logarithm of (x + 1)",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Log1p(x)), nil
		},
		[]shapes.Example{
			{"1 log1p", "Num", "0.6931471805599453", nil},
			{"0 log1p", "Num", "0", nil},
			{"-1 log1p", "Num", "-inf", nil},
			{"-2 log1p", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the base 2 logarithm.",
		types.Num{},
		"a number",
		"log2",
		nil,
		types.Num{},
		"the base 2 logarithm of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Log2(x)), nil
		},
		[]shapes.Example{
			{"3 log2", "Num", "1.5849625007211563", nil},
			{"2 log2", "Num", "1", nil},
			{"1 log2", "Num", "0", nil},
			{"0 log2", "Num", "-inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Finds the largest number",
		types.NewArr(types.Num{}),
		"an array of numbers",
		"max",
		nil,
		types.Num{},
		"the largest number in the input, or -inf if the input is empty",
		"",
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
		[]shapes.Example{
			{"[1, 3, 2] max", "Num", "3", nil},
			{"[-1, -3, -2] max", "Num", "-1", nil},
			{"[] max", "Num", "-inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Finds the smallest number",
		types.NewArr(types.Num{}),
		"an array of numbers",
		"min",
		nil,
		types.Num{},
		"the smallest number in the input, or inf if the input is empty",
		"",
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
		[]shapes.Example{
			{"[1, 3, 2] min", "Num", "1", nil},
			{"[-1, -3, -2] min", "Num", "-3", nil},
			{"[] min", "Num", "inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes powers (exponentiation).",
		types.Num{},
		"the base",
		"**",
		[]*params.Param{
			params.SimpleParam("y", "the exponent", types.Num{}),
		},
		types.Num{},
		"the base taken to the y-th power",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			y := float64(argumentValues[0].(states.NumValue))
			return states.NumValue(math.Pow(x, y)), nil
		},
		[]shapes.Example{
			{"7 **3", "Num", "343", nil},
			{"4 **.5", "Num", "2", nil},
			{"7 **(-2)", "Num", "0.02040816326530612", nil},
			{"-7 **0.5", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns a random number between 0 and 1.",
		types.Any{},
		"any value (is ignored)",
		"random",
		nil,
		types.Num{},
		"",
		"a floating-point, pseudo-random number n with 0 <= n < 1 and approximately uniform distribution over that range",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NumValue(rand.Float64()), nil
		},
		[]shapes.Example{
			{"[null] repeat(1000) each(random) each(>=0) all", "Bool", "true", nil},
			{"[null] repeat(1000) each(random) each(<1) all", "Bool", "true", nil},
		},
	),
	// TODO a funcer that returns a random integer in a specified interval
	shapes.SimpleFuncer(
		"Rounds a number to the nearest integer.",
		types.Num{},
		"a number",
		"round",
		nil,
		types.Num{},
		"the nearest integer, or away from zero if there's two nearest integers",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Round(x)), nil
		},
		[]shapes.Example{
			{"0.9 round", "Num", "1", nil},
			{"5.95 round", "Num", "6", nil},
			{"5.5 round", "Num", "6", nil},
			{"5.05 round", "Num", "5", nil},
			{"-5.05 round", "Num", "-5", nil},
			{"-5.5 round", "Num", "-6", nil},
			{"-5.95 round", "Num", "-6", nil},
		},
	),
	shapes.SimpleFuncer(
		"Determines the sign of a number.",
		types.Num{},
		"a number",
		"sign",
		nil,
		types.Num{},
		"1 if the input is positive, -1 if negative, 0 if it's 0, and -0 if it's -0",
		"",
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
		[]shapes.Example{
			{"3 sign", "Num", "1", nil},
			{"0 sign", "Num", "0", nil},
			{"-0 sign", "Num", "-0", nil},
			{"-3", "Num", "-3", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the sine.",
		types.Num{},
		"an angle in radians",
		"sin",
		nil,
		types.Num{},
		"the sine of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Sin(x)), nil
		},
		[]shapes.Example{
			{"-inf sin", "Num", "nan", nil},
			{"-0 sin", "Num", "-0", nil},
			{"0 sin", "Num", "0", nil},
			{"1 sin", "Num", "0.8414709848078965", nil},
			{"pi /2 sin", "Num", "1", nil},
			{"inf sin", "Num", "nan", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the hyperbolic sine.",
		types.Num{},
		"a number",
		"sinh",
		nil,
		types.Num{},
		"the hyperbolic sine of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Sinh(x)), nil
		},
		[]shapes.Example{
			{"0 sinh", "Num", "0", nil},
			{"1 sinh", "Num", "1.1752011936438014", nil},
			{"-1 sinh", "Num", "-1.1752011936438014", nil},
			{"2 sinh", "Num", "3.626860407847019", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes square roots.",
		types.Num{},
		"a number",
		"sqrt",
		nil,
		types.Num{},
		"the square root of the input, or nan if it's negative",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Sqrt(x)), nil
		},
		[]shapes.Example{
			{"-1 sqrt", "Num", "nan", nil},
			{"-0 sqrt", "Num", "-0", nil},
			{"0 sqrt", "Num", "0", nil},
			{"1 sqrt", "Num", "1", nil},
			{"2 sqrt", "Num", "1.4142135623730951", nil},
			{"9 sqrt", "Num", "3", nil},
			{"inf sqrt", "Num", "inf", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the tangent.",
		types.Num{},
		"an angle in radians",
		"tan",
		nil,
		types.Num{},
		"the tangent of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Tan(x)), nil
		},
		[]shapes.Example{
			{"0 *pi /180 tan", "Num", "0", nil},
			{"45 *pi /180 tan", "Num", "1", nil},
			{"90 *pi /180 tan", "Num", "16331239353195392", nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes the hyperbolic tangent.",
		types.Num{},
		"a number",
		"tanh",
		nil,
		types.Num{},
		"the hyperbolic tangent of the input",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Tanh(x)), nil
		},
		[]shapes.Example{
			{"-1 tanh", "Num", "-0.7615941559557649", nil},
			{"0 tanh", "Num", "0", nil},
			{"inf tanh", "Num", "1", nil},
			{"1 tanh", "Num", "0.7615941559557649", nil},
		},
	),
	shapes.SimpleFuncer(
		"Rounds towards zero.",
		types.Num{},
		"a number",
		"trunc",
		nil,
		types.Num{},
		"the input without fractional digits",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			x := float64(inputValue.(states.NumValue))
			return states.NumValue(math.Trunc(x)), nil
		},
		[]shapes.Example{
			{"13.37 trunc", "Num", "13", nil},
			{"42.84 trunc", "Num", "42", nil},
			{"0.123 trunc", "Num", "0", nil},
			{"-0.123 trunc", "Num", "-0", nil},
		},
	),
}
