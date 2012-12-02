package skiplist

import "math"

const (
	defaultMaxElems = 65536
	p               = 0.25
)

var maxLevel int

func init() {
	maxLevel = expectedLevels(p, defaultMaxElems)
}

func expectedLevels(p float64, maxElems uint) int {
	return int(math.Ceil(math.Log(float64(maxElems)) / math.Log(1/p)))
}
