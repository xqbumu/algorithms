package main

import (
	"fmt"
	"reflect"
	"runtime"

	"golang.org/x/sync/singleflight"
)

type Controller struct{}

func (c Controller) GetAll() {
}

func main() {
	c := Controller{}
	fn := c.GetAll
	str := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	fmt.Println(str) // Output: GetAll
}

func exampleSingleflight() {
	g := new(singleflight.Group)

	block := make(chan struct{})
	res1c := g.DoChan("key", func() (interface{}, error) {
		<-block
		return "func 1", nil
	})
	res2c := g.DoChan("key", func() (interface{}, error) {
		<-block
		return "func 2", nil
	})
	close(block)

	res1 := <-res1c
	res2 := <-res2c

	// Results are shared by functions executed with duplicate keys.
	fmt.Println("Shared:", res2.Shared)
	// Only the first function is executed: it is registered and started with "key",
	// and doesn't complete before the second funtion is registered with a duplicate key.
	fmt.Println("Equal results:", res1.Val.(string) == res2.Val.(string))
	fmt.Println("Result:", res1.Val)
}
