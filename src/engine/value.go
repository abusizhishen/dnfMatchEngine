package engine

import (
	json "github.com/json-iterator/go"
	"time"
)

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
}

var TypeMap = map[string]int{}

func init() {
	for idx, k := range TypeArr {
		TypeMap[k] = idx
	}
}

type Value struct {
	Type       string             `json:"type"`
	Origin     string             `json:"origin"`
	Float64    float64            `json:"float_64"`
	Int        int                `json:"int"`
	Int64      int64              `json:"int64"`
	String     string             `json:"string"`
	Time       time.Time          `json:"time"`
	TagValue   []string           `json:"tag_value"`
	Float64Arr map[float64]string `json:"-"`
	IntArr     map[int]string     `json:"int_arr"`
	Int64Arr   map[int64]string   `json:"int64_arr"`
	StringArr  map[string]string  `json:"string_arr"`
}

func (v *Value) MarshalJSON() ([]byte, error) {
	var floatArr []float64
	for idx, _ := range v.Float64Arr {
		floatArr = append(floatArr, idx)
	}
	return json.Marshal(
		struct {
			Type       string            `json:"type"`
			Origin     string            `json:"origin"`
			Float64    float64           `json:"float_64"`
			Int        int               `json:"int"`
			Int64      int64             `json:"int64"`
			String     string            `json:"string"`
			Time       time.Time         `json:"time"`
			TagValue   []string          `json:"tag_value"`
			Float64Arr []float64         `json:"float_64_arr"`
			IntArr     map[int]string    `json:"int_arr"`
			Int64Arr   map[int64]string  `json:"int64_arr"`
			StringArr  map[string]string `json:"string_arr"`
		}{
			v.Type,
			v.Origin,
			v.Float64,
			v.Int,
			v.Int64,
			v.String,
			v.Time,
			v.TagValue,
			floatArr,
			v.IntArr,
			v.Int64Arr,
			v.StringArr,
		})
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var floatArr []float64
	for idx, _ := range v.Float64Arr {
		floatArr = append(floatArr, idx)
	}
	type V struct {
		Type       string            `json:"type"`
		Origin     string            `json:"interface"`
		Float64    float64           `json:"float_64"`
		Int        int               `json:"int"`
		Int64      int64             `json:"int64"`
		String     string            `json:"string"`
		Time       time.Time         `json:"time"`
		TagValue   []string          `json:"tag_value"`
		Float64Arr []float64         `json:"float_64_arr"`
		IntArr     map[int]string    `json:"int_arr"`
		Int64Arr   map[int64]string  `json:"int64_arr"`
		StringArr  map[string]string `json:"string_arr"`
	}

	var vv V
	err := json.Unmarshal(data, &vv)
	if err != nil {
		return err
	}

	v.Type = vv.Type
	v.Origin = vv.Origin
	v.Float64 = vv.Float64
	v.Int = vv.Int
	v.Int64 = vv.Int64
	v.String = vv.String
	v.Time = vv.Time
	v.TagValue = vv.TagValue
	floatArr = vv.Float64Arr
	v.IntArr = vv.IntArr
	v.Int64Arr = vv.Int64Arr
	v.StringArr = vv.StringArr

	return nil
}

func (v *Value) Init() {
	v.Float64Arr = make(map[float64]string)
	v.IntArr = make(map[int]string)
	v.Int64Arr = make(map[int64]string)
	v.StringArr = make(map[string]string)
}

var layout = "2006-01-02 15:04:05"
