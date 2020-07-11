# dnfMatchEngine

```
自定义dnf表达式匹配
通过自定义dnf表达式生成匹配引擎，对匹配引擎输入相关标签数据，进行一定的计算即可实现匹配
支持扩展自定义数据类型实现匹配

```

```go
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
    
    // output true <nil>
}

```