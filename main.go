package main

import (
	_ "net/http/pprof"
	"log"
	"net/http"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync"
)

var (
	pool      *redis.Pool
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wg.Wait()
}
