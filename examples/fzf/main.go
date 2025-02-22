package main

import (
	"fmt"
	"os"
	"time"

	fzf "github.com/junegunn/fzf/src"
)

func main() {
	inputChan := make(chan string)
	go func() {
		for _, s := range []string{"a", "b", "c"} {
			time.Sleep(time.Millisecond * 300)
			inputChan <- s
		}
		close(inputChan)
	}()

	outputChan := make(chan string)
	go func() {
		for s := range outputChan {
			fmt.Println("Got: " + s)
		}
	}()

	exit := func(code int, err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(code)
	}

	// Build fzf.Options
	options, err := fzf.ParseOptions(
		true, // whether to load defaults ($FZF_DEFAULT_OPTS_FILE and $FZF_DEFAULT_OPTS)
		[]string{"--multi", "--reverse", "--border", "--height=40%"},
	)
	if err != nil {
		exit(fzf.ExitError, err)
	}

	// Set up input and output channels
	options.Input = inputChan
	options.Output = outputChan

	// Run fzf
	code, err := fzf.Run(options)
	exit(code, err)
}
