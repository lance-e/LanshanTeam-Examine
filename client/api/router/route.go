package router

import (
	"LanshanTeam-Examine/client/api/middleware"
	"LanshanTeam-Examine/client/handle"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.Cors()) //解决跨域问题

	user := engine.Group("/user")
	{

		user.POST("/register/byPassword", handle.RegisterByPassword)
		user.POST("/register/byPhoneNumber", handle.RegisterByPhoneNumber)

		user.POST("/login/byPassword", handle.Login)
		user.POST("/login/byPhoneNumber", handle.LoginByPhoneNumber)

		//验证码
		user.POST("/sendCode", handle.SendCode)

		//github第三方注册登陆
		user.GET("/githubRegisterAndLogin", handle.GithubRegisterAndLogin)
		user.GET("/githubCallback", handle.GithubCallback)

		user.Use(middleware.JWT())
		user.GET("/information", handle.HomePage)
		user.POST("/addFriend", handle.SendFriendRequest)

	}

	return engine
}
