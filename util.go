package main

func bound(i, max int) int {
	// loop backwards
	if i > max {
		return 0
	}

	// loop forwards
	if i < 0 {
		return max
	}

	return i
}
