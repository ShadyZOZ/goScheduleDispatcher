package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func runServer() {
	fmt.Println("start server")
	r := gin.Default()
	r.POST("/schedule", func(ctx *gin.Context) {
		err := publish(ctx)
		if err == nil {
			ctx.Status(201)
		} else {
			ctx.AbortWithStatusJSON(400, gin.H{"error": err})
		}
	})
	r.Run(":5001")
}
