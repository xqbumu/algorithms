package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wulijun/go-php-serialize/phpserialize"
)

func main() {
	content, err := os.ReadFile("/Users/bumu/Downloads/serialize.txt")
	if err != nil {
		panic(err)
	}
	result, err := phpserialize.Decode(string(content))
	if err != nil {
		panic(err)
	}
	buf, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}
