package usecase

import (
	"encoding/json"
	"github.com/google/uuid"
	"safe_talk/pkg/models"
	"safe_talk/pkg/utils"
)

func (u *UseCase) SignUp(data models.AuthData) (int, string) {
	_, notFound, err := u.r.GetUserByLogin(data.Login)
	if err != nil && !notFound {
		u.l.Errorf("Ошибка при извлечении данных с БД. Ошибка %v", err)
		return 500, "Ошибка при связи с БД. Повторте попытку позже :("
	}

	if !notFound {
		return 400, "существуюший пользователь. Выберите другой логин"
	}

	if !utils.CheckLogin(data.Login) || !utils.CheckPassword(data.Password) {
		return 400, "Пароль или логин не соответвует требованиям"
	}

	hash := utils.GetSha256Hash(data.Password)

	if err := u.r.AddUser(&models.User{Login: data.Login, Password: hash, Keyword: data.KeyWord}); err != nil {
		u.l.Errorf("Ошибка при добавлении пользователья в БД. Ошибка: %v", err)
		return 500, "Ошибка при связи с БД. Повторте попытку позже :("
	}

	return 200, "Успешная регистрация. Поздравляю :)"
}

func (u *UseCase) SignIn(data models.AuthData) (int, string) {
	user, notFound, err := u.r.GetUserByLogin(data.Login)
	if err != nil && notFound {
		return 404, "Такого пользователя не существует :("
	}

	if err != nil {
		u.l.Errorf("Ошибка при извлечении user из БД. Ошибка: %v", err)
		return 500, "Ошибка при обращение в БД"
	}

	hash := utils.GetSha256Hash(data.Password)

	if hash != user.Password {
		return 401, "Не корректный логин или пароль"
	}

	jwt, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		u.l.Errorf("Ошибка при генерации токена. Ошибка:%v", err)
		return 500, "Ошибка при генерации токена"
	}

	res := struct {
		Token  string
		UserId string
	}{Token: jwt, UserId: user.ID.String()}

	temp, _ := json.MarshalIndent(res, " ", "")
	return 200, string(temp)
}

func (u *UseCase) ResetPassword(data models.AuthData) (int, string) {
	user, notFound, err := u.r.GetUserByLogin(data.Login)
	if err != nil && notFound {
		return 404, "Такого пользователя не существует :("
	}

	if err != nil {
		u.l.Errorf("Ошибка при извлечении user из БД. Ошибка: %v", err)
		return 500, "Ошибка при обращение в БД :("
	}

	hash := utils.GetSha256Hash(user.Keyword)

	if hash != data.KeyWord {
		return 400, "Не правильное ключевое слово :("
	}

	if err := u.r.UpdateUserPassword(user.ID, data.Password); err != nil {
		u.l.Errorf("Ошибка при обновлении пароля. Ошибка: %v", err)
		return 500, "Ошибка при обрашении к БД :("
	}

	return 200, "Успешно. Можете войти в свой аккаунт :)"
}

func (u *UseCase) GetNewMessages(userId string) ([]models.SMessage, error) {
	return u.r.GetUserMessages(userId)
}

func (u *UseCase) AddMessage(msg models.SMessage) error {
	return u.r.AddMessage(msg)
}

func (u *UseCase) GetUserChats(userId string) ([]models.Chat, error) {
	return u.r.GetUserChat(userId)
}

func (u *UseCase) CreateChat(userIDs []uuid.UUID) error {
	chatId, err := u.r.CreateChat(models.NChat{IsActive: true})
	if err != nil {
		u.l.Errorf("Ошибка при создании нового чата. ОШИБКА %v", err)
		return err
	}

	for _, id := range userIDs {
		user2Chats := models.User2Chats{
			UserId: id,
			ChatId: chatId,
		}
		if err := u.r.AddUsers2Chat(user2Chats); err != nil {
			u.l.Errorf("Ошибка при добавлении пользователья в чат. ОШИБКА %v", err)
			return err
		}
	}
	return nil

}
