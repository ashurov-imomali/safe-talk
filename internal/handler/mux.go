package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes(h Handler) http.Handler {
	//gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(gin.Recovery(), h.Logger(), h.cors())
	{
		e.GET("/ping", h.ping)
		e.POST("/sign-up", h.registration)
		e.POST("/sign-in", h.signIn)
		e.POST("/reset-password", h.resetPassword)
	}
	chats := e.Group("", gin.Recovery(), h.auth(), h.cors())
	{
		chats.POST("/chat", h.createChat)
		chats.GET("/chat-history", h.getChatHistory)
		chats.GET("/user-chats", h.getUserChats)
		chats.GET("/user", h.getUserByLogin)
		chats.PUT("/message", h.updateMessage)
		chats.DELETE("/message", h.deleteMessage)
		chats.POST("/file", h.sendFile)
		chats.GET("/file", h.getFile)
	}

	rtConnection := e.Group("/connection", gin.Recovery(), h.auth(), h.cors())
	rtConnection.GET("", h.getConn)

	return e
}
