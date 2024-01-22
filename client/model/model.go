package model

type Userinfo struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	PhoneNumber int    `form:"phone_number" json:"phone_number"`
	Email       string `form:"phone_number" json:"email"`
}
