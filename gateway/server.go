package main

import (
	auth "gateway/auth"
	"gateway/biz"

	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

type GetUserRequest struct {
	AuthKey string `json:"authKey"`
	UserId  string `json:"userId"`

	WithSqlInject bool `json:"withSqlInject"`
}

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

func getUsers(c *gin.Context) {
	var req GetUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	usersJSON, err := biz.GetUsers(req.AuthKey, req.UserId, req.WithSqlInject)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	} else {
		c.JSON(200, usersJSON)
	}
}

func main() {
	port := flag.Int("port", 6433, "Port number")
	flag.Parse()
	auth.Init()
	biz.Init()
	defer auth.Close()
	defer biz.Close()
	r := gin.Default()
	r.GET("/auth", getAuthKey)
	r.POST("/get-users", getUsers)
	r.Run(fmt.Sprintf(":%d", *port))
}
