package xmath

import (
	"math"
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		b       float64
		epsilon float64
		want    bool
	}{
		{
			name:    "exact_equal_zero",
			a:       0,
			b:       0,
			epsilon: 0,
			want:    true,
		},
		{
			name:    "exact_equal_non_zero",
			a:       1.2345,
			b:       1.2345,
			epsilon: 0,
			want:    true,
		},
		{
			name:    "within_epsilon",
			a:       1.0,
			b:       1.0000001,
			epsilon: 1e-6,
			want:    true,
		},
		{
			name:    "equal_to_epsilon_boundary",
			a:       1.0,
			b:       1.000001,
			epsilon: 0.000001,
			want:    true,
		},
		{
			name:    "outside_epsilon",
			a:       1.0,
			b:       1.00001,
			epsilon: 1e-6,
			want:    false,
		},
		{
			name:    "negative_numbers_within_epsilon",
			a:       -1.0,
			b:       -1.0000005,
			epsilon: 1e-6,
			want:    true,
		},
		{
			name:    "mixed_sign_within_epsilon",
			a:       -1e-9,
			b:       1e-9,
			epsilon: 3e-9,
			want:    true,
		},
		{
			name:    "epsilon_negative_should_work",
			a:       1.0,
			b:       1.0000001,
			epsilon: -1e-6,
			want:    false,
		},
		{
			name:    "epsilon_zero_not_equal",
			a:       1.0,
			b:       1.0 + math.SmallestNonzeroFloat64,
			epsilon: 0,
			want:    true,
		},
		{
			name:    "very_large_numbers_within_epsilon",
			a:       1e308,
			b:       1e308 + 1e292,
			epsilon: 1e293,
			want:    true,
		},
		{
			name:    "inf_equal",
			a:       math.Inf(1),
			b:       math.Inf(1),
			epsilon: 0,
			want:    true,
		},
		{
			name:    "inf_not_equal",
			a:       math.Inf(1),
			b:       math.Inf(-1),
			epsilon: math.MaxFloat64,
			want:    false,
		},
		{
			name:    "nan_always_false",
			a:       math.NaN(),
			b:       math.NaN(),
			epsilon: 1,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Equal(tt.a, tt.b, tt.epsilon)
			if got != tt.want {
				t.Fatalf(
					"Equal(%v, %v, %v) = %v, want %v",
					tt.a, tt.b, tt.epsilon, got, tt.want,
				)
			}
		})
	}
}
