package models

import "github.com/google/uuid"

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	KeyWord  string `json:"key_word"`
}

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Login    string
	Password string
	Keyword  string
}

type Message struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type SMessage struct {
	Id       int `gorm:"primaryKey"`
	Text     string
	FromUser string
	ToUser   string
}

func (SMessage) TableName() string {
	return "messages"
}

//
//create table messages(
//id serial primary key,
//text text,
//--     chat_id uuid references chats
//from_user uuid references users,
//to_user uuid references users
//);
