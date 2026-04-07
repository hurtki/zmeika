package app

func InGaps(c cord, size int) bool {
	if c.x < 0 || c.x >= size {
		return false
	}
	if c.y < 0 || c.y >= size {
		return false
	}
	return true
}
