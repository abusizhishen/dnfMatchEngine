package engine

type engine interface {
	Parse(string2 string)error
	Init()error
	Match(map[string]string)(bool,error)
}


type node interface {
	Parse(string2 string)error
	Init()error
	Match(map[string]string)(bool,error)
}


type array interface {
	Parse(string2 string)error
	Init()error
	Match(map[string]string)(bool,error)
}

type item interface {
	Parse(string2 string)error
	Init()error
	Match(map[string]string)(bool,error)
}

type token int