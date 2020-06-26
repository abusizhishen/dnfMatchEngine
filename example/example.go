package main

import (
	engine2 "github.com/abusizhishen/expressMatchEngine/src/engine"
	"log"
)

func main() {
	engine,err := engine2.New("┐[(age>{1}|0)]")
	if err != nil{
		panic(err)
	}
	match := engine.Match(map[string]string{"age":"1"})
	log.Println(match)
}
