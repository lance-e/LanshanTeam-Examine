package router

import (
	"LanshanTeam-Examine/client/api/middleware"
	"LanshanTeam-Examine/client/handle"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.Cors()) //解决跨域问题

	engine.GET("/", handle.HomePage)
	user := engine.Group("/user")
	{

		user.POST("/register/byPassword", handle.RegisterByPassword)
		//user.POST("/register/byPhoneNumber", handle.RegisterByPhoneNumber)

		user.POST("/login/byPassword", handle.Login)

		user.GET("/githubRegisterAndLogin", handle.GithubRegisterAndLogin)
		user.GET("/githubCallback", handle.GithubCallback)

		user.Use(middleware.JWT())

	}

	return engine
}
