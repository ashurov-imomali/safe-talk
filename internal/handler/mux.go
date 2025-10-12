package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes(h Handler) http.Handler {
	//gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(gin.Recovery(), h.Logger())
	{
		e.GET("/ping", h.ping)
		e.POST("/sign-up", h.registration)
		e.POST("/sign-in", h.signIn)
		e.POST("/reset-password", h.resetPassword)
	}
	chats := e.Group("", gin.Recovery(), h.auth(), h.Logger())
	{
		chats.POST("/chat", h.createChat)
		chats.GET("/chat-history", h.getChatHistory)
		chats.GET("/user-chats", h.getUserChats)
	}
	rtConnection := e.Group("/connection", gin.Recovery(), h.auth())
	rtConnection.GET("", h.getConn)

	return e
}
