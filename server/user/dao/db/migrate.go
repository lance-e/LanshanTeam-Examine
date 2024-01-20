package db

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() error {
	return DB.AutoMigrate(&UserInfo{})
	//if err != nil {
	//	utils.UserLogger.Panic("couldn't migrate the user table")
	//	return
	//}
}
