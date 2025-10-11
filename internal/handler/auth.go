package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"safe_talk/pkg/models"
)

func (h *Handler) registration(c *gin.Context) {
	var user models.AuthData
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"message": "Ошибка парсинга проверьте теги!"})
		return
	}

	status, msg := h.u.SignUp(user)
	c.JSON(status, gin.H{"message": msg})
}

func (h *Handler) signIn(c *gin.Context) {
	var data models.AuthData
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"message": "Ошибка парсинга проверьте теги!"})
		return
	}

	status, message := h.u.SignIn(data)
	c.JSON(status, gin.H{"message": message})
}

func (h *Handler) resetPassword(c *gin.Context) {
	var data models.AuthData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"message": "Ошибка парсинга проверьте теги!"})
		return
	}

	status, message := h.u.ResetPassword(data)

	c.JSON(status, gin.H{"message": message})

}

func (h *Handler) getChatHistory(c *gin.Context) {
	//token := c.GetHeader("Authorization")
	//if token == "" {
	//	c.JSON(401, gin.H{"message": "Не авторизован"})
	//	return
	//}

	userId := c.Query("user_id")

	messages, err := h.u.GetNewMessages(userId)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, messages)

}

var clients map[string]*websocket.Conn

func init() {
	clients = make(map[string]*websocket.Conn)
}

func (h *Handler) getConn(c *gin.Context) {
	conn, err := h.ws.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()
	//userID := c.Value("user_uuid").(string)
	userID := c.Query("user_uuid")
	clients[userID] = conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			delete(clients, userID)
			break
		}
		var message models.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		if err := h.u.AddMessage(models.SMessage{Text: message.Message, FromUser: userID, ToUser: message.UserID}); err != nil {
			c.JSON(500, gin.H{"message": "Что то с БД"})
			return
		}

		if toConn, ok := clients[message.UserID]; ok {
			b, _ := json.MarshalIndent(message, " ", "")
			toConn.WriteMessage(websocket.TextMessage, b)
		}
	}
}
