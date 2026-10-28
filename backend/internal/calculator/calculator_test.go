package calculator_test

import (
	"errors"
	"math"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestApply(t *testing.T) {
	negZero := math.Copysign(0, -1)

	tests := []struct {
		name     string
		op       string
		operands []float64
		want     float64
		wantErr  error
	}{
		{name: "add integers", op: "add", operands: []float64{2, 3}, want: 5},
		{name: "add negatives", op: "add", operands: []float64{-2, -3}, want: -5},
		{name: "add keeps float64 precision", op: "add", operands: []float64{0.1, 0.2}, want: 0.30000000000000004},
		{name: "add with zero", op: "add", operands: []float64{0, 7}, want: 7},

		{name: "subtract", op: "subtract", operands: []float64{10, 4}, want: 6},
		{name: "subtract operand order", op: "subtract", operands: []float64{4, 10}, want: -6},

		{name: "multiply", op: "multiply", operands: []float64{3, 4}, want: 12},
		{name: "multiply decimals", op: "multiply", operands: []float64{-2, 0.5}, want: -1},
		{name: "multiply normalizes negative zero", op: "multiply", operands: []float64{-1, 0}, want: 0},
		{name: "multiply overflow", op: "multiply", operands: []float64{1e308, 10}, wantErr: calculator.ErrResultOutOfRange},

		{name: "divide", op: "divide", operands: []float64{10, 4}, want: 2.5},
		{name: "divide operand order", op: "divide", operands: []float64{4, 10}, want: 0.4},
		{name: "divide by zero", op: "divide", operands: []float64{1, 0}, wantErr: calculator.ErrDivisionByZero},
		{name: "divide by negative zero", op: "divide", operands: []float64{1, negZero}, wantErr: calculator.ErrDivisionByZero},
		{name: "divide zero by zero", op: "divide", operands: []float64{0, 0}, wantErr: calculator.ErrDivisionByZero},

		{name: "power", op: "power", operands: []float64{2, 10}, want: 1024},
		{name: "power operand order", op: "power", operands: []float64{10, 2}, want: 100},
		{name: "power negative exponent", op: "power", operands: []float64{2, -1}, want: 0.5},
		{name: "power zero to the zero", op: "power", operands: []float64{0, 0}, want: 1},
		{name: "power fractional exponent of negative base", op: "power", operands: []float64{-8, 1.0 / 3}, wantErr: calculator.ErrInvalidOperand},
		{name: "power overflow", op: "power", operands: []float64{10, 400}, wantErr: calculator.ErrResultOutOfRange},
		{name: "power zero to a negative exponent", op: "power", operands: []float64{0, -1}, wantErr: calculator.ErrResultOutOfRange},

		{name: "sqrt", op: "sqrt", operands: []float64{16}, want: 4},
		{name: "sqrt irrational", op: "sqrt", operands: []float64{2}, want: math.Sqrt2},
		{name: "sqrt of zero", op: "sqrt", operands: []float64{0}, want: 0},
		{name: "sqrt of negative zero", op: "sqrt", operands: []float64{negZero}, want: 0},
		{name: "sqrt of negative", op: "sqrt", operands: []float64{-1}, wantErr: calculator.ErrInvalidOperand},

		{name: "percentage", op: "percentage", operands: []float64{15, 200}, want: 30},
		{name: "percentage is exact for common values", op: "percentage", operands: []float64{7, 100}, want: 7},
		{name: "percentage operand order", op: "percentage", operands: []float64{200, 15}, want: 30},
		{name: "percentage divides first when the product overflows", op: "percentage", operands: []float64{1e300, 1e10}, want: 1e308},
		{name: "percentage overflow", op: "percentage", operands: []float64{1e308, 1e10}, wantErr: calculator.ErrResultOutOfRange},

		{name: "NaN operand", op: "add", operands: []float64{math.NaN(), 1}, wantErr: calculator.ErrInvalidOperand},
		{name: "infinite operand", op: "add", operands: []float64{1, math.Inf(1)}, wantErr: calculator.ErrInvalidOperand},
		{name: "too few operands", op: "add", operands: []float64{1}, wantErr: calculator.ErrWrongArity},
		{name: "too many operands", op: "sqrt", operands: []float64{4, 2}, wantErr: calculator.ErrWrongArity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, ok := calculator.Lookup(tt.op)
			if !ok {
				t.Fatalf("Lookup(%q) not found", tt.op)
			}

			got, err := op.Apply(tt.operands...)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Apply(%v) error = %v, want %v", tt.operands, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Apply(%v) unexpected error: %v", tt.operands, err)
			}
			if got != tt.want || math.Signbit(got) != math.Signbit(tt.want) {
				t.Errorf("Apply(%v) = %v, want %v", tt.operands, got, tt.want)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	for _, name := range []string{"add", "subtract", "multiply", "divide", "power", "sqrt", "percentage"} {
		op, ok := calculator.Lookup(name)
		if !ok || op.Name != name {
			t.Errorf("Lookup(%q) = %+v, %v; want registered operation", name, op, ok)
		}
	}

	for _, name := range []string{"", "modulo", "ADD"} {
		if _, ok := calculator.Lookup(name); ok {
			t.Errorf("Lookup(%q) found an operation, want none", name)
		}
	}
}

func TestArity(t *testing.T) {
	want := map[string]int{"add": 2, "sqrt": 1, "percentage": 2}
	for name, arity := range want {
		op, _ := calculator.Lookup(name)
		if op.Arity != arity {
			t.Errorf("%s arity = %d, want %d", name, op.Arity, arity)
		}
	}
}
