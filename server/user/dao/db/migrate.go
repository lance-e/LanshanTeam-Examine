package db

import (
	"LanshanTeam-Examine/server/user/pkg/utils"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() error {
	err := DB.AutoMigrate(&UserInfo{}, &FriendShip{})
	if err != nil {
		utils.UserLogger.Error("couldn't migrate the user table")
		return err
	}
	return nil
}
