package router

import (
	"LanshanTeam-Examine/caller/api/middleware"
	"LanshanTeam-Examine/caller/handle"
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
		user.POST("/acceptAddFriend", handle.ReceiveFriendRequest)

		user.GET("/rank", handle.ShowRank)
	}
	game := engine.Group("/game")
	{
		game.GET("/createRoom", handle.Create)
		game.GET("/joinRoom", handle.Join)
		game.Use(middleware.JWT())
		game.GET("/ready", handle.Ready)
		game.POST("/showHistory", handle.ShowHistory)
	}
	return engine
}
