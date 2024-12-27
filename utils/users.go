package utils

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/clients/responses"
	"github.com/gin-gonic/gin"
	"log"
)

func SetUserInContext(c *gin.Context, user *responses.UserInfo) {
	log.Println("Setting user in context")
	c.Set("UserInfo", user)
}

func GetUserInfoFromContext(c *gin.Context) *responses.UserInfo {
	log.Println("Getting user info from context")
	userInfo, _ := c.Get("UserInfo")

	user, _ := userInfo.(*responses.UserInfo)

	return user
}
