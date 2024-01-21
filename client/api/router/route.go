package router

import (
	"LanshanTeam-Examine/client/api/middleware"
	"LanshanTeam-Examine/client/handle"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.Cors()) //解决跨域问题

	engine.GET("/state", handle.Get)
	user := engine.Group("/user")
	{
		register := user.Group("/register")
		{
			register.POST("/byPassword", handle.RegisterByPassword)
		}

	}

	return engine
}
