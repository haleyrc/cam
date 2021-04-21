package math

func ClampLower(val, bound int) int {
	if val < bound {
		return bound
	}
	return val
}

func ClampUpper(val, bound int) int {
	if val > bound {
		return bound
	}
	return val
}
