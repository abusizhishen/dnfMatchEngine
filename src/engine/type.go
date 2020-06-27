package engine

type engine interface {
	Parse(string2 string)
	Match(map[string]string)bool
}

type node engine
type collection engine
type token engine