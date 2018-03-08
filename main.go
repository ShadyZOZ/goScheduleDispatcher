package main

import (
	"fmt"
	"sync"
	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
	redisAddr = "localhost:6379"
)

func init() {
	pool = newPool(redisAddr)
}

func main() {
	fmt.Println("go schedule dispatcher")
	// avoid `subscribe()` and `runServer()` blocking process, so send them to goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go subscribe()
	go runServer()
	wg.Wait()
}
