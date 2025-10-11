package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math"
	"net/http"
	"os"
	"safe_talk/internal/usecase"
	"safe_talk/pkg/logger"
	"time"
)

type Handler struct {
	u  usecase.UseCase
	l  logger.Logger
	ws websocket.Upgrader
}

func New(u usecase.UseCase, l logger.Logger) Handler {
	return Handler{u: u, l: l, ws: websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {
		return true
	}}}
}

func (h *Handler) ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func (h *Handler) Logger() gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		//h.l.KVLog("hostname", hostname)
		//h.l.KVLog("latency", latency)
		//h.l.KVLog("clientIP", clientIP)
		//h.l.KVLog("method", c.Request.Method)
		//h.l.KVLog("referrer", referer)
		//h.l.KVLog("dataLength", dataLength)
		//h.l.KVLog("statusCode", statusCode)
		//h.l.KVLog("clientUserAgent", clientUserAgent)
		if len(c.Errors) > 0 {
			h.l.Error(errors.New(c.Errors.ByType(gin.ErrorTypePrivate).String()), "ERROR")
		} else {
			format := "[ClientIP]: %s | [HostName]: %s | [Time]: %s | [Method]: %s | [Path]: %s | [StatusCode]: %d | [DataLength]: %d | [ClientUserAgent]: %s | [Latency]: %d | [Referer] : %s"
			msg := fmt.Sprintf(format, clientIP, hostname, time.Now().Format("02/Jan/2006:15:04:05 +5"), c.Request.Method, path, statusCode, dataLength, clientUserAgent, latency, referer)
			if statusCode > 499 {
				h.l.Error(errors.New(msg), "ERROR")
			} else if statusCode > 399 {
				h.l.Warn(msg)
			} else {
				h.l.Info(msg)
			}
		}
	}
}
