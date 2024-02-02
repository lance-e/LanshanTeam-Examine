package db

import (
	"LanshanTeam-Examine/server/game/utils"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() error {
	err := DB.AutoMigrate(&GameSteps{})
	if err != nil {
		utils.GameLogger.Error("couldn't migrate the user table")
		return err
	}
	return nil
}
