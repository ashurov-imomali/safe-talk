package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"safe_talk/pkg/logger"
	"safe_talk/pkg/models"
)

type Repository struct {
	l logger.Logger
	p *gorm.DB
}

func NewRepos(p *gorm.DB, l logger.Logger) Repository {
	return Repository{l: l, p: p}
}

func (r *Repository) GetUserByLogin(login string) (*models.User, bool, error) {
	var user models.User
	//return &user, r
	err := r.p.Table("users").Where("login = ?", login).First(&user).Error
	if err != nil {
		return &user, errors.Is(err, gorm.ErrRecordNotFound), err
	}
	return &user, false, nil
}

func (r *Repository) AddUser(u *models.User) error {
	return r.p.Create(&u).Error
}

func (r *Repository) UpdateUserPassword(userId uuid.UUID, password string) error {
	return r.p.Table("users").Where("id = ?", userId).Update("password", password).Error
}

func (r *Repository) AddMessage(message models.SMessage) error {
	return r.p.Create(&message).Error
}

func (r *Repository) GetUserMessages(chatId string) ([]models.SMessage, error) {
	var nMessages []models.SMessage
	return nMessages, r.p.Select("m.*").Table("messages m").
		Joins("join chats c on c.id = m.chat_id and c.id = ?", chatId).
		Order("m.created_at desc").Scan(&nMessages).Error
}

func (r *Repository) GetUserChat(userId string) ([]models.Chat, error) {
	//select c.id, u.login, c.last_message from chats c
	//	join users2chats u2c on c.id = u2c.chat_id and u2c.user_id = ''
	//	join users u on u2c.user_id = u.id;
	var result []models.Chat
	return result, r.p.Select("c.id, u.login, c.last_message").Table("chats c").
		Joins("join users2chats u2c on c.id = u2c.chat_id and u2c.user_id = ?", userId).
		Joins("join users u on u2c.user_id = u.id").
		Scan(&result).Error
}

func (r *Repository) CreateChat(c models.NChat) (uuid.UUID, error) {
	return c.ID, r.p.Create(&c).Error
}

func (r *Repository) AddUsers2Chat(m models.User2Chats) error {
	return r.p.Create(&m).Error
}

func (r *Repository) GetChatUsers(chatId, userId string) (models.User, error) {
	//
	//select * from users u
	//	join users2chats u2c on u.id = u2c.user_id
	//	join chats c on c.id = u2c.chat_id and c.id = ?
	var user models.User
	return user, r.p.Select("u.*").Table("users u").
		Joins("join users2chats u2c on u.id = u2c.user_id and u.id != ?", userId).
		Joins("join chats c on c.id = u2c.chat_id and c.id = ?", chatId).First(&user).Error
}

func (r *Repository) GetUsersByLogin(login string) ([]models.User, error) {
	var users []models.User
	return users, r.p.Where("login ilike ?", fmt.Sprintf(`%%%s%%`, login)).Find(&users).Error
}

func (r *Repository) UpdateMessage(id, text string) error {
	return r.p.Table("messages").Where("id = ?", id).UpdateColumn("text", text).Error
}

func (r *Repository) DeleteMessage(id string) error {
	return r.p.Delete(&models.SMessage{}, id).Error
}
