package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"strings"
	"net/http"
)

func subscribe() {
	conn := pool.Get()
	defer conn.Close()
	connSub := pool.Get()
	defer connSub.Close()
	sub := redis.PubSubConn{Conn: connSub}
	sub.Subscribe("__keyevent@0__:expired")
	for {
		switch v := sub.Receive().(type) {
		case redis.Message:
			key := strings.Replace(string(v.Data), "str:", "hm:", 1)
			res, err := redis.Values(conn.Do("HMGET", key, "callbackUrl", "data"))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			var callbackUrl string
			var data string
			if _, err := redis.Scan(res, &callbackUrl, &data); err != nil {
				fmt.Println(err.Error())
				continue
			}
			if _, err := http.Post(callbackUrl, "application/json", strings.NewReader(data)); err != nil {
				fmt.Println(err.Error())
				continue
			}
			conn.Do("DEL", key)
		case error:
			fmt.Println(v)
		}
	}
}
