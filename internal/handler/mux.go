package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes(h Handler) http.Handler {
	//gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	//e.Use(gin.Recovery(), h.Logger())

	e.GET("/ping", h.ping)
	e.POST("/sign-up", h.registration)
	e.POST("/sign-in", h.signIn)
	e.POST("/reset-password", h.resetPassword)
	e.GET("/chat-history", h.getChatHistory)
	e.GET("/connection", h.getConn)
	return e
}
