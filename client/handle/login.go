package handle

import (
	"LanshanTeam-Examine/client/api/middleware"
	"LanshanTeam-Examine/client/model"
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/userModule"
	"LanshanTeam-Examine/client/rpc/userModule/pb"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user model.Userinfo

	if err := c.ShouldBind(&user); err != nil {
		utils.ClientLogger.Error("ERROR: " + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.LackParams,
			"message": "you should enter your information completely,你应该完善你的信息",
			"error":   err.Error(),
		})
		return
	}
	resp, err := userModule.UserClient.Login(c, &pb.LoginReq{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		if err.Error() == "wrong password" {
			c.JSON(400, gin.H{
				"code":    consts.LoginPasswordWrong,
				"message": "wrong password,密码错误",
				"error":   err.Error(),
			})
			return
		} else {
			c.JSON(400, gin.H{
				"code":    consts.LoginFailed,
				"message": "can't get user information from server,无法从服务器获取信息",
				"error":   err.Error(),
			})
			return
		}
	} else if !resp.GetFlag() {
		c.JSON(400, gin.H{
			"code":    consts.UserNotFound,
			"message": "user not found ,用户不存在",
			"error":   "user not found",
		})
		return
	}
	token, err := middleware.GetToken(&user)
	if err != nil {
		utils.ClientLogger.Error("generate token failed,ERROR: " + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.GenerateTokenFailed,
			"message": "generate token failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    consts.LoginSuccess,
		"message": "login success,用户登陆成功",
		"error":   "",
		"token":   token,
	})

}
