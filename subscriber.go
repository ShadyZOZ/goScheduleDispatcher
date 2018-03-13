package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"bytes"
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
			if callbackUrl == "" {
				fmt.Println("no callbackUrl")
				continue
			}
			action := strings.Split(key, ":")[1]
			postData := map[string]string{"action": action, "data": data}
			jsonData, _ := json.Marshal(postData)
			if _, err := http.Post(callbackUrl, "application/json", bytes.NewBuffer(jsonData)); err != nil {
				fmt.Println(err.Error())
				continue
			}
			//resp, err := http.Post(callbackUrl + "?token=oTGmBUNIspCPcXJZxQih1ig1", "application/json", strings.NewReader(data))
			//body, _ := ioutil.ReadAll(resp.Body)
			//fmt.Println(string(body))
			//if err != nil {
			//	fmt.Println(err.Error())
			//	continue
			//}
			conn.Do("DEL", key)
		case error:
			fmt.Println("error:", v)
		}
	}
}
