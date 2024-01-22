package utils

import (
	"golang.org/x/crypto/bcrypt"
)

var salt = "LanShanTeam examine :D"

func Encrypt(password string) (string, error) {
	password = password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		UserLogger.Error("ERROR: " + err.Error())
		return "", err
	}
	return string(hashedPassword), nil
}

func Compare(hashedPassword string, enteredPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword))

}
