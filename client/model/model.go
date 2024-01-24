package model

type Userinfo struct {
	Username     string `form:"username" json:"username" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	PhoneNumber  int    `form:"phone_number" json:"phone_number"`
	Email        string `form:"phone_number" json:"email"`
	IsGithubUser bool   `form:"is_github_user" json:"is_github_user"`
}

type Config struct {
	GithubConfig `mapstructure:"github" json:"github"`
}

type GithubConfig struct {
	ClientId      string `mapstructure:"client_id" json:"client_id"`
	ClientSecrets string `mapstructure:"client_secrets" json:"client_secrets"`
}
