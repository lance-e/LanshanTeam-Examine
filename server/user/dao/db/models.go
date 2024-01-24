package db

import (
	"LanshanTeam-Examine/server/user/pkg/utils"
)

type UserInfo struct {
	Username     string `gorm:"username" json:"username"`
	Password     string `gorm:"password" json:"password"`
	PhoneNumber  int    `gorm:"phone_number" json:"phone_number"`
	Email        string `gorm:"email" json:"email"`
	IsGithubUser bool   `gorm:"is_github_user" json:"is_github_user"`
}

func (u UserInfo) TableName() string {
	return "user_info"
}

// Get 数据库操作：进行查询，what是字段名，value是值，user是查询到的地方
func (u *UserInfo) Get(what, value string, user *UserInfo) error {
	err := DB.Model(u).Where(what+"= ? ", value).First(&user).Error
	if err != nil {
		utils.UserLogger.Error("GET: " + err.Error())
		return err
	}
	return nil
}
func (u *UserInfo) Create() error {
	err := DB.Model(&UserInfo{}).Create(u).Error
	if err != nil {
		utils.UserLogger.Error("CREATE: " + err.Error())
		return err
	}
	return nil
}
func (u *UserInfo) Update() {

}
func (u *UserInfo) Delete() {

}
