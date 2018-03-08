package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"time"
)

type Message struct {
	Action      string `json:"action" binding:"required"`
	UUID        string `json:"uuid" binding:"required"`
	CallbackUrl string `json:"callbackUrl" binding:"required"`
	Data        string `json:"data" binding:"required"`
	Timestamp   int64 `json:"timestamp"`
	Override    bool   `json:"override"`
}

func getKey(action string, uuid string) string {
	return fmt.Sprintf("str:%s:%s", action, uuid)
}

func getHMKey(action string, uuid string) string {
	return fmt.Sprintf("hm:%s:%s", action, uuid)
}

func publish(ctx *gin.Context) interface{} {
	conn := pool.Get()
	defer conn.Close()
	var message Message
	if err := ctx.BindJSON(&message); err == nil {
		key := getKey(message.Action, message.UUID)
		hmKey := getHMKey(message.Action, message.UUID)
		if !message.Override {
			if res, err := conn.Do("GET", key); err != nil {
				return err.Error()
			} else if res != nil {
				return "can't override current schedule"
			}
		}
		var ttl int64 = 1
		if message.Timestamp != 0 {
			t := time.Now().Unix()
			if message.Timestamp > t {
				ttl = message.Timestamp - t
			}
		}
		if err := conn.Send("SET", key, "1", "EX", ttl); err != nil {
			return err.Error()
		}
		if err := conn.Send("HMSET", hmKey, "callbackUrl", message.CallbackUrl, "data", message.Data); err != nil {
			return err.Error()
		}
		if err := conn.Flush(); err != nil {
			return err.Error()
		}
		return nil
	} else {
		return err.Error()
	}
}
