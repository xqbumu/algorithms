package main

import (
	"algorithms/tmp"
	"log"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
)

func main() {
	fp, err := tmp.Content.Open("index.js")
	if err != nil {
		panic(err)
	}

	p, err := parser.ParseFile(nil, "", fp, 0)
	if err != nil {
		panic(err)
	}

	log.Println(p)

	vm := otto.New()
	script, err := vm.Compile("", fp)
	if err != nil {
		panic(err)
	}
	log.Println(script.String())
	vm.Run(`
			abc = 2 + 2;
			console.log("The value of abc is " + abc); // 4
	`)
}
