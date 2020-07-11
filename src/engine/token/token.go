package token

const (
	//relation
	And = iota
	Or
	Not

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

)

var Identification = map[int]string{
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
}
