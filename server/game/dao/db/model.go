package db

import (
	"LanshanTeam-Examine/server/game/utils"
	"gorm.io/gorm"
)

type GameSteps struct {
	gorm.Model
	RoomHost string `gorm:"room_host" json:"room_host,omitempty"`
	Player   string `gorm:"player" json:"player,omitempty"`
	Row      int64  `gorm:"row" json:"row,omitempty"`
	Column   int64  `gorm:"column" json:"column,omitempty"`
}

func (g *GameSteps) TableName() string {
	return "game_steps"
}
func (g *GameSteps) Create() error {
	err := DB.Model(&GameSteps{}).Create(g).Error
	if err != nil {
		utils.GameLogger.Error("gameSteps create failed , error :" + err.Error())
		return err
	}
	return nil
}
func (g *GameSteps) Get() ([]GameSteps, error) {
	var info = make([]GameSteps, 0, 10)
	err := DB.Model(&GameSteps{}).Order("created_at desc").Where("room_host = ?", g.RoomHost).Find(&info).Error
	if err != nil {
		utils.GameLogger.Error("gameSteps get failed , error :" + err.Error())
		return nil, err
	}
	utils.GameLogger.Debug("GET get success")
	return info, nil
}
