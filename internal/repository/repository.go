package repository

import (
	"errors"
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

func (r *Repository) GetUserMessages(userId string) ([]models.SMessage, error) {
	var nMessages []models.SMessage
	return nMessages, r.p.Where("to_user = ?", userId).Find(&nMessages).Error
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
