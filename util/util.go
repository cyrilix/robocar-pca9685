package util

// MapRange Linear mapping between two ranges of values
func MapRange(x, xmin, xmax, ymin, ymax float64) int {
	Xrange := xmax - xmin
	Yrange := ymax - ymin
	XYratio := Xrange / Yrange

	y := (x-xmin)/XYratio + ymin

	return int(y)
}
