package util

import "testing"

func TestRangeMapping(t *testing.T) {
	cases := []struct {
		x, xmin, xmax, ymin, ymax float64
		expected                  int
	}{
		// Test positive
		{-100, -100, 100, 0, 1000, 0},
		{0, -100, 100, 0, 1000, 500},
		{100, -100, 100, 0, 1000, 1000},

		// Test negative
		{0, 0, 100, 0, 1000, 0},
		{50, 0, 100, 0, 1000, 500},
		{100, 0, 100, 0, 1000, 1000},

		// Reverse
		{0, 100, 0, 0, 1000, 1000},
		{50, 100, 0, 0, 1000, 500},
		{100, 100, 0, 0, 1000, 0},
	}

	for _, c := range cases {
		val := MapRange(c.x, c.xmin, c.xmax, c.ymin, c.ymax)
		if val != c.expected {
			t.Errorf("MapRange(%.0f, %.0f, %.0f, %.0f, %.0f): %d, wants %d", c.x, c.xmin, c.xmax, c.ymin, c.ymax, val, c.expected)
		}
	}
}
