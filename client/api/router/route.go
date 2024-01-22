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
		register := user.Group("/register")
		{
			register.POST("/byPassword", handle.RegisterByPassword)
			register.POST("/byPhoneNumber", handle.RegisterByPhoneNumber)
			//register.POST("/byEmail", handle.RegisterByEmail)
		}

		login := user.Group("/login")
		{
			login.POST("/byPassword", handle.Login)
		}

		user.Use(middleware.JWT())

	}
	engine.POST()
	return engine
}
