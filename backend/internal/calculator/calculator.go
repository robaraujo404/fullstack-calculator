// Package calculator implements the arithmetic operations exposed by the API.
// It knows nothing about HTTP or JSON: callers look up an operation by name
// and apply it to float64 operands.
package calculator

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrDivisionByZero   = errors.New("division by zero is undefined")
	ErrInvalidOperand   = errors.New("operand is outside the operation's domain")
	ErrResultOutOfRange = errors.New("result is out of range")
	ErrWrongArity       = errors.New("wrong number of operands")
)

// Operation is a named arithmetic operation that takes Arity operands.
type Operation struct {
	Name  string
	Arity int
	fn    func(x []float64) (float64, error)
}

var operations = map[string]Operation{
	"add":        {Name: "add", Arity: 2, fn: add},
	"subtract":   {Name: "subtract", Arity: 2, fn: subtract},
	"multiply":   {Name: "multiply", Arity: 2, fn: multiply},
	"divide":     {Name: "divide", Arity: 2, fn: divide},
	"power":      {Name: "power", Arity: 2, fn: power},
	"sqrt":       {Name: "sqrt", Arity: 1, fn: sqrt},
	"percentage": {Name: "percentage", Arity: 2, fn: percentage},
}

// Lookup returns the operation registered under name. Names are case-sensitive.
func Lookup(name string) (Operation, bool) {
	op, ok := operations[name]
	return op, ok
}

// Apply runs the operation and guarantees a finite result.
//
// JSON cannot represent NaN or ±Inf, so they never leave this function:
// NaN means the operands were outside the operation's domain (e.g. the square
// root of a negative number) and ±Inf means the result overflowed float64.
func (o Operation) Apply(operands ...float64) (float64, error) {
	if len(operands) != o.Arity {
		return 0, fmt.Errorf("%w: %s takes %d, got %d", ErrWrongArity, o.Name, o.Arity, len(operands))
	}
	for _, x := range operands {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return 0, ErrInvalidOperand
		}
	}

	result, err := o.fn(operands)
	switch {
	case err != nil:
		return 0, err
	case math.IsNaN(result):
		return 0, ErrInvalidOperand
	case math.IsInf(result, 0):
		return 0, ErrResultOutOfRange
	case result == 0:
		return 0, nil // normalizes -0 to 0
	}
	return result, nil
}

func add(x []float64) (float64, error)      { return x[0] + x[1], nil }
func subtract(x []float64) (float64, error) { return x[0] - x[1], nil }
func multiply(x []float64) (float64, error) { return x[0] * x[1], nil }
func power(x []float64) (float64, error)    { return math.Pow(x[0], x[1]), nil }
func sqrt(x []float64) (float64, error)     { return math.Sqrt(x[0]), nil }

// divide checks the divisor explicitly: in IEEE-754, x/0 is ±Inf or NaN,
// which Apply would otherwise report as a range or domain error.
func divide(x []float64) (float64, error) {
	if x[1] == 0 { // also true for -0
		return 0, ErrDivisionByZero
	}
	return x[0] / x[1], nil
}

// percentage returns a% of b. Multiplying before dividing keeps common cases
// exact (7% of 100 is 7, whereas 7/100*100 is 7.000000000000001). When the
// product overflows, dividing first still finds any result that fits in a
// float64.
func percentage(x []float64) (float64, error) {
	if p := x[0] * x[1] / 100; !math.IsInf(p, 0) {
		return p, nil
	}
	return x[0] / 100 * x[1], nil
}
