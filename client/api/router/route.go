package router

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	engine := gin.Default()

	user := engine.Group("/user")
	{
		user.POST("/register/byPassword") //未实名
		user.POST("/register/byPhoneNumber")
		user.POST("/register/byEmail")
		user.POST("/register/byGitHub") //OAuth2.0注册

	}

	return engine
}
