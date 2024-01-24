package model

import "time"

type Userinfo struct {
	Username     string `form:"username" json:"username" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	PhoneNumber  string `form:"phone_number" json:"phone_number"`
	Email        string `form:"phone_number" json:"email"`
	IsGithubUser bool   `form:"is_github_user" json:"is_github_user"`
	Code         int64  `form:"code" json:"code"`
}

// 验证码信息
type CodeInfo struct {
	Code     int64     `form:"code" json:"code"`
	ExpireAt time.Time `form:"expire_at" json:"expire_at"`
}

//type CodeInfo struct {
//	Code     int64     `form:"code" json:"code"`
//	ExpireAt time.Time `form:"expire_at" json:"expire_at"`
//	RequestNumber string `form:"request_number" json:"request_number"`
//}

type Config struct {
	GithubConfig `mapstructure:"github" json:"github"`
	AliyunConfig `mapstructure:"aliyun" json:"aliyun"`
}

type GithubConfig struct {
	ClientId      string `mapstructure:"client_id" json:"client_id"`
	ClientSecrets string `mapstructure:"client_secrets" json:"client_secrets"`
}
type AliyunConfig struct {
	AccessKeyId     string `mapstructure:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret"`
}
