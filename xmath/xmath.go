package xmath

import "math"

func Equal(a, b, epsilon float64) bool {
	// 1. NaN：任何比较都应失败
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}

	// 2. Inf：只在符号相同的情况下相等
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}

	// 3. epsilon 必须是非负的（明确语义）
	if epsilon < 0 {
		return false
	}

	// 4. 有限数的绝对误差比较（IEEE 友好）
	return math.Abs(a-b) <= epsilon
}
