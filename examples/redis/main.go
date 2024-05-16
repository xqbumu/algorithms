package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	options := &redis.Options{
		Addr:         "redis.infra.orb.local:6379",
		DB:           0,
		Password:     "",
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
	cli := redis.NewClient(options)

	ret := cli.Get("test").String()

	log.Println(ret)
}
