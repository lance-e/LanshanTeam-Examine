package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
)

var salt = "LanShanTeam examine :D"

func Encrypt(password string) string {

	pw, _ := hex.DecodeString(password + salt)
	hash := sha256.Sum256(pw)
	return hex.EncodeToString(hash[:])
}

func Compare(hashedPassword string, enteredPassword string) error {
	enteredPassword = enteredPassword + salt

	hashed := Encrypt(enteredPassword)
	log.Println(enteredPassword)
	log.Println(hashedPassword)
	if hashed != hashedPassword {
		return errors.New("hashedPassword is not the hash of the given password ")
	}
	return nil
}
