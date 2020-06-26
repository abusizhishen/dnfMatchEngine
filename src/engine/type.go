package engine

type engine interface {
	Parse(string2 string)
	Match(map[string]string)bool
}

type node interface {
	Parse(string2 string)
	Match(map[string]string)bool
}

type collection interface {
	Parse(string2 string)
	Match(map[string]string)bool
}

type token interface {
	Parse(string2 string)
	Match(map[string]string)bool
}