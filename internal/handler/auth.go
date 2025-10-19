package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"safe_talk/pkg/models"
	"safe_talk/pkg/utils"
	"strings"
)

func (h *Handler) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		userID, err := utils.JWTConfirm(strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			h.l.Errorf("Ошибка при извлечении данных с токена. ОШИБКА: [%v]", err)
			c.AbortWithStatus(401)
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

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
	chatID := c.Query("chat_id")

	messages, err := h.u.GetNewMessages(chatID)
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
	anyUserId, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"message": "не авторизован"})
		return
	}
	userID := anyUserId.(string)
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
		if message.Message == "" {
			c.JSON(400, gin.H{"message": "Поле сообщения не может быть пустой!"})
			return
		}
		ToUser, err := h.u.AddMessage(models.SMessage{Text: message.Message, FromUser: userID, ChatId: message.ChatId})
		if err != nil {
			c.JSON(500, gin.H{"message": "Что то с БД"})
			return
		}
		message.FromUser = userID
		if toConn, ok := clients[ToUser]; ok {
			b, _ := json.MarshalIndent(message, " ", "")
			toConn.WriteMessage(websocket.TextMessage, b)
		}
	}
}

func (h *Handler) getUserChats(c *gin.Context) {
	value, find := c.Get("user_id")
	if !find {
		c.JSON(401, gin.H{"message": "Не авторизован ((("})
		return
	}
	userID := value.(string)
	chats, err := h.u.GetUserChats(userID)
	if err != nil {
		h.l.Errorf("Ошибка при получении данных с БД. ОШИБКА [%v]", err)
		c.JSON(500, gin.H{"message": "Ошибка при обращение к БД"})
		return
	}

	c.JSON(200, chats)
}

func (h *Handler) createChat(c *gin.Context) {
	chat := struct {
		UserIds []uuid.UUID `json:"user_ids"`
	}{}

	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(400, gin.H{"message": "Проверьте данные"})
		return
	}

	chatId, status, err := h.u.CreateChat(chat.UserIds)
	if err != nil {
		c.JSON(status, gin.H{"message": "Ошибка в БД или существующий чат"})
		return
	}
	c.JSON(status, gin.H{"chat_id": chatId})
}

func (h *Handler) getUserByLogin(c *gin.Context) {
	value := c.Query("login")
	users, err := h.u.GetUsersByLogin(value)
	if err != nil {
		h.l.Errorf("Ошибка при получении данных с БД. ОШИБКА [%v]", err)
		c.JSON(500, gin.H{"message": "Ошибка на стороне сервера :(("})
		return
	}

	c.JSON(200, users)
}

func (h *Handler) updateMessage(c *gin.Context) {
	updMEssage := struct {
		Id    interface{}
		NText string
	}{}

	if err := c.ShouldBindJSON(&updMEssage); err != nil {
		c.JSON(400, gin.H{"message": "Не корректные данные"})
		return
	}

	if err := h.u.UpdateMessage(updMEssage.Id, updMEssage.NText); err != nil {
		c.JSON(500, gin.H{"message": "Внутренная ошибка :))"})
		h.l.Errorf("Ошибка при обоашении к БД. ОШИБКА %v", err)
		return
	}
	c.JSON(200, gin.H{"message": "Успешно обновлено :)"})
}

func (h *Handler) deleteMessage(c *gin.Context) {
	id := c.Query("message_id")
	if err := h.u.DeleteMessage(id); err != nil {
		c.JSON(500, gin.H{"message": "Внутренная ошибка :))"})
		h.l.Errorf("Ошибка при обоашении к БД. ОШИБКА %v", err)
		return
	}
	c.JSON(200, gin.H{"message": "Успешно удалено :)"})
}

func (h *Handler) sendFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": "Ошибка на стороне клиента (*_*)"})
		return
	}
	defer file.Close()
	value, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"message": "Не авторизован"})
		return
	}
	if err := h.u.SaveFileToServer(value.(string), header.Filename, file); err != nil {
		c.JSON(500, gin.H{"message": "Ошибка при отправки файла"})
		return
	}

	c.JSON(200, gin.H{"message": "Успешно доставленно"})
}
