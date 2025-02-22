package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kod/kod"
)

func main() {
	if err := kod.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

// app is the main component of the application. kod.Run creates
// it and passes it to serve.
type app struct {
	kod.Implements[kod.Main]
}

// serve is called by kod.Run and contains the body of the application.
func serve(context.Context, *app) error {
	fmt.Println("Hello")
	return nil
}
