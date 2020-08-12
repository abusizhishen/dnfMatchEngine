package main

import (
	"fmt"
	engine2 "github.com/abusizhishen/expressMatchEngine/src/engine"
)

func main() {
	// int match
	engine,err := engine2.New("[(age>{1}|0)]")
	if err != nil{
		panic(err)
	}
	fmt.Println(engine.Match(map[string]string{"age":"1"}))
	// output false,nil

	fmt.Println(engine.Match(map[string]string{"age":"2"}))
	// output true,nil

	//string match
	engine,err = engine2.New("[(word={hello}|3)]")
	if err != nil{
		panic(err)
	}
	fmt.Println(engine.Match(map[string]string{"word":"hell"}))
	// output false,nil

	fmt.Println(engine.Match(map[string]string{"word":"hello"}))
	// output true,nil
}
