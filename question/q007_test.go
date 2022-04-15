package question

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var value int32

func setValue(delta int32) {
	for {
		v := value
		if atomic.CompareAndSwapInt32(&value, v, (v + delta)) {
			break
		}
	}
}

func Test_setValue(t *testing.T) {
	type args struct {
		delta int32
	}
	tests := []struct {
		name string
		args args
	}{
		{``, args{1}},
		{``, args{10}},
		{``, args{2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setValue(tt.args.delta)
		})
	}
}

func Test_sub8(t *testing.T) {
	abc := make(chan int, 1000)
	for i := 0; i < 10; i++ {
		abc <- i
	}
	go func() {
		for a := range abc {
			fmt.Println("a: ", a)
		}
	}()
	close(abc)
	fmt.Println("close")
	time.Sleep(time.Second * 5)
}

func Test_sub9(t *testing.T) {
	type Student struct {
		name string
	}
	m := map[string]Student{"people": {"zhoujielun"}}
	s := m["people"]
	s.name = "wuyanzu"
	log.Println(m, s)
}

func Test_sub11(t *testing.T) {
	var i byte
	go func() {
		for i = 0; i <= 255; i++ {
			log.Println(i)
		}
	}()
	fmt.Println("Dropping mic")
	// Yield execution to force executing other goroutines
	runtime.Gosched()
	runtime.GC()
	fmt.Println("Done")
}
