package db

import (
	"LanshanTeam-Examine/server/user/pkg/utils"
	"gorm.io/gorm"
	"log"
)

type UserInfo struct {
	gorm.Model
	Username     string `gorm:"username" json:"username"`
	Password     string `gorm:"password" json:"password"`
	PhoneNumber  int    `gorm:"phone_number" json:"phone_number"`
	Email        string `gorm:"email" json:"email"`
	IsGithubUser bool   `gorm:"is_github_user" json:"is_github_user"`
	Score        int    `gorm:"score" json:"score"`
}

func (u UserInfo) TableName() string {
	return "user_info"
}

type FriendShip struct {
	gorm.Model
	Sender   string `gorm:"sender" json:"sender"`
	Receiver string `gorm:"receiver" json:"receiver"`
	Status   bool   `gorm:"status" json:"status"`
}

func (f *FriendShip) TableName() string {
	return "friend_ship"
}

// Get 数据库操作：进行查询，what是字段名，value是值，user是查询到的地方
func (u *UserInfo) Get(what, value string, user *UserInfo) error {
	err := DB.Model(&UserInfo{}).Where(what+"= ? ", value).First(&user).Error
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
func (u *UserInfo) AddScore() error {
	log.Println("first:", *u)
	err := u.Get("username", u.Username, u)
	if err != nil {
		utils.UserLogger.Error("can't get this user's information")
		return err
	}
	log.Println("second", *u)
	err = DB.Model(&UserInfo{}).Where("username = ?", u.Username).Update("score", u.Score+5).Error
	if err != nil {
		utils.UserLogger.Debug("can't add score")
		return err
	}
	log.Println("third", *u)
	return nil
}
func (u *UserInfo) Delete() {

}

func (f *FriendShip) IsRequestAlreadyExists() bool {
	err := DB.Model(&FriendShip{}).Where("sender = ? && receiver = ? ", f.Sender, f.Receiver).First(nil).Error
	if err == nil {
		utils.UserLogger.Info("request already exists or database is error")
		return true
	}
	return false
}
func (f *FriendShip) CheckState() bool {
	var flag bool
	err := DB.Model(&FriendShip{}).Select("status").Where("sender = ? && receiver = ?", f.Sender, f.Receiver).Find(&flag).Error
	if err != nil {
		err = DB.Model(&FriendShip{}).Select("status").Where("sender = ? && receiver = ?", f.Receiver, f.Sender).Find(&flag).Error
		if err != nil {
			utils.UserLogger.Error("not any connection ")
			return false
		}
		return flag
	}
	return flag
}
func (f *FriendShip) Create() error {
	err := DB.Model(&FriendShip{}).Create(f).Error
	if err != nil {
		utils.UserLogger.Error("CREATE: " + err.Error())
		return err
	}
	return nil
}
func (f *FriendShip) Update() error {
	err := DB.Model(&FriendShip{}).Where("sender = ? && receiver = ? ", f.Sender, f.Receiver).Update("status", true).Error
	if err != nil {
		utils.UserLogger.Error("UPDATE Friendship failed")
		return err
	}
	return nil
}
func (f *FriendShip) Delete() {

}
