package engine

var TypeArr = []string{
	"int",
	"float64",
	"int64",
	"string",
	"time",
	"[]int",
	"[]float64",
	"[]string",
	"[]int64",
	"[]mixed", // [int or string or float or time... or mixed]
}

var TypeMap = map[string]int{}

func init() {
	for idx, k := range TypeArr {
		TypeMap[k] = idx
	}
}
