package testutil_test

import (
	"testing"
	"time"

	"github.com/vnnyx/employee-management/pkg/optional"
	"github.com/vnnyx/employee-management/pkg/testutil"
)

type testCase struct {
	name     string
	x        interface{}
	y        interface{}
	expected bool
}

func TestEqualVerbose(t *testing.T) {
	duration := 5 * time.Second

	tests := []testCase{
		// Bool
		{
			name:     "Equal optional.Bool",
			x:        optional.NewBool(true),
			y:        optional.NewBool(true),
			expected: true,
		},
		{
			name:     "Not equal optional.Bool",
			x:        optional.NewBool(true),
			y:        optional.NewBool(false),
			expected: false,
		},
		{
			name:     "Empty vs non-empty optional.Bool",
			x:        optional.NewBool(),
			y:        optional.NewBool(true),
			expected: false,
		},
		{
			name:     "Both empty optional.Bool",
			x:        optional.NewBool(),
			y:        optional.NewBool(),
			expected: true,
		},

		// Int32
		{
			name:     "Equal optional.Int32",
			x:        optional.NewInt32(123),
			y:        optional.NewInt32(123),
			expected: true,
		},
		{
			name:     "Not equal optional.Int32",
			x:        optional.NewInt32(123),
			y:        optional.NewInt32(321),
			expected: false,
		},
		{
			name:     "Empty vs non-empty optional.Int32",
			x:        optional.NewInt32(),
			y:        optional.NewInt32(123),
			expected: false,
		},

		// Int64
		{
			name:     "Equal optional.Int64",
			x:        optional.NewInt64(9999999999),
			y:        optional.NewInt64(9999999999),
			expected: true,
		},
		{
			name:     "Not equal optional.Int64",
			x:        optional.NewInt64(1),
			y:        optional.NewInt64(2),
			expected: false,
		},
		{
			name:     "Both empty optional.Int64",
			x:        optional.NewInt64(),
			y:        optional.NewInt64(),
			expected: true,
		},

		// String
		{
			name:     "Equal optional.String",
			x:        optional.NewString("golang"),
			y:        optional.NewString("golang"),
			expected: true,
		},
		{
			name:     "Not equal optional.String",
			x:        optional.NewString("go"),
			y:        optional.NewString("lang"),
			expected: false,
		},
		{
			name:     "Empty optional.String",
			x:        optional.NewString(),
			y:        optional.NewString("test"),
			expected: false,
		},

		// Float64
		{
			name:     "Equal optional.Float64",
			x:        optional.NewFloat64(3.14),
			y:        optional.NewFloat64(3.14),
			expected: true,
		},
		{
			name:     "Not equal optional.Float64",
			x:        optional.NewFloat64(3.14),
			y:        optional.NewFloat64(2.71),
			expected: false,
		},

		// Float32
		{
			name:     "Equal optional.Float32",
			x:        optional.NewFloat32(1.23),
			y:        optional.NewFloat32(1.23),
			expected: true,
		},
		{
			name:     "Not equal optional.Float32",
			x:        optional.NewFloat32(1.23),
			y:        optional.NewFloat32(4.56),
			expected: false,
		},

		// Duration
		{
			name:     "Equal optional.Duration",
			x:        optional.NewDuration(duration),
			y:        optional.NewDuration(duration),
			expected: true,
		},
		{
			name:     "Not equal optional.Duration",
			x:        optional.NewDuration(duration),
			y:        optional.NewDuration(duration + time.Second),
			expected: false,
		},

		{
			name:     "Equal Option[int32]",
			x:        optional.Option[int32]{}.SetAndReturn(123),
			y:        optional.Option[int32]{}.SetAndReturn(123),
			expected: true,
		},
		{
			name:     "Not equal Option[int32]",
			x:        optional.Option[int32]{}.SetAndReturn(123),
			y:        optional.Option[int32]{}.SetAndReturn(321),
			expected: false,
		},
		{
			name:     "Equal Option[string]",
			x:        optional.Option[string]{}.SetAndReturn("go"),
			y:        optional.Option[string]{}.SetAndReturn("go"),
			expected: true,
		},
		{
			name:     "Empty vs set Option[bool]",
			x:        optional.Option[bool]{}, // empty
			y:        optional.Option[bool]{}.SetAndReturn(true),
			expected: false,
		},
		{
			name:     "Both empty Option[float64]",
			x:        optional.Option[float64]{}, // empty
			y:        optional.Option[float64]{}, // empty
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := testutil.EqualVerbose(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("EqualVerbose(%v, %v) = %v; expected %v", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}
