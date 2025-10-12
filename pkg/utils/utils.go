package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"regexp"
	"time"
)

func CheckPhoneNum(phone *string) bool {
	n := regexp.MustCompile(`^\d{9}$`)
	t := regexp.MustCompile(`^\+\d{12}$`)
	if n.MatchString(*phone) {
		*phone = "+992" + *phone
		return true
	}
	return t.MatchString(*phone)
}

func GetSha256Hash(a ...interface{}) string {
	var data string
	for _, i := range a {
		data += fmt.Sprintf("%v", i)
	}

	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func CheckLogin(login string) bool {
	if len(login) < 2 {
		return false
	}
	r := regexp.MustCompile("^[a-zA-Z0-9_]{3,16}$")
	return r.MatchString(login)
}

func GenerateJWT(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(50 * time.Minute).Unix() //todo 50  minut !!!!!
	claims["user_uuid"] = userId
	return token.SignedString([]byte("secretKey"))
}

func JWTConfirm(token string) (string, error) {
	var claims jwt.MapClaims
	parseT, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secretKey"), nil
	})
	if err != nil || !parseT.Valid {
		sprintf := fmt.Sprintf("invalid token. Err: %v", err)
		return "", errors.New(sprintf)
	}
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", errors.New("token is expired")
		}
	} else {
		return "", errors.New("token is expired")
	}
	return claims["user_uuid"].(string), nil
}

func CheckPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	return hasDigit && hasSpecial && hasLower && hasUpper
}
