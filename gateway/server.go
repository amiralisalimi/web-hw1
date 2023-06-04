package main

import (
	auth "gateway/auth"

	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func getAuthKey(c *gin.Context) {
	messageId, err := auth.SendPGRequest(0)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	key, _, err := auth.SendDHParamsRequest(messageId + 1)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"key": key})
}

func main() {
	port := flag.Int("port", 6433, "Port number")
	flag.Parse()
	r := gin.Default()
	r.GET("/auth", getAuthKey)
	r.Run(fmt.Sprintf(":%d", *port))
}
