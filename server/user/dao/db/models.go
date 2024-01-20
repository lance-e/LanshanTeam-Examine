package db

type UserInfo struct {
	Username    string `gorm:"username" json:"username"`
	Password    string `gorm:"password" json:"password"`
	PhoneNumber int    `gorm:"phone_number" json:"phone_number"`
	Email       string `gorm:"email" json:"email"`
	GitHubName  string `gorm:"git_hub_name" json:"git_hub_name"`
}

func (u UserInfo) TableName() string {
	return "user_info"
}
