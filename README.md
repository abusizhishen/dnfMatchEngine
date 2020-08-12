# dnfMatchEngine

```
自定义dnf表达式匹配
通过自定义dnf表达式生成匹配引擎，对匹配引擎输入表示式，自顶向下解析，生成为可执行程序，
通过对程序输入相应的数据，计算即可实现匹配逻辑
支持扩展自定义数据类型实现匹配

```

```go
// var Identification = map[int]string{
//	And: "^",
//	Or:  "∨",
//	Not: "┐",
//
//	Gt: ">",
//	GtOrEq: "≥",
//	Eq: "=",
//	NotEq: "≠",
//	LtOrEq: "≤",
//	Lt: "<",
//	In: "∈",
//	NotIn: "∉",
//
//	Int: "int",
//	Float64: "float64",
//	String: "string",
//	Time: "time",
//}
```

```go
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
```

```go
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


```