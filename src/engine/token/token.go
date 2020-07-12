package token

const (
	//relation
	And = iota
	Or
	Not
	DefaultRelation

	//symble
	Gt
	GtOrEq
	Eq
	NotEq
	LtOrEq
	Lt
	In
	NotIn

	//data type
	Int
	Float64
	String
	Time
	IntArr
	Float64Arr
	StringArr
	TimeArr

)

type Token int
var Tokens = [...]string{
	And: "^",
	Or:  "∨",
	Not: "┐",

	Gt: ">",
	GtOrEq: "≥",
	Eq: "=",
	NotEq: "≠",
	LtOrEq: "≤",
	Lt: "<",
	In: "∈",
	NotIn: "∉",

	Int: "int",
	Float64: "float64",
	String: "string",
	Time: "time",
	IntArr: "intArr",
	Float64Arr: "float64Arr",
	StringArr: "stringArr",
	TimeArr: "timeArr",
}

var Map map[string]int

func init() {
	for id, s := range Tokens{
		Map[s] = id
	}
}
