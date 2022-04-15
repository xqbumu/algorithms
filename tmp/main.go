package main

import (
	"fmt"
	"log"
	"math/big"
	"sync"
	"unicode/utf8"
)

func main() {
	num := &big.Rat{}

	num.SetString("3")
	num.Neg(num)

	log.Panicln(num.IsInt())

	log.Println(utf8.DecodeRune([]byte("1+1")))

	out := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			out <- i
		}
		close(out)
	}()
	go func() {
		defer wg.Done()
		for i := range out {
			fmt.Println(i)
		}
	}()
	wg.Wait()
}
