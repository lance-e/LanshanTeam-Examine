package handle

import (
	"LanshanTeam-Examine/caller/api/middleware"
	"LanshanTeam-Examine/caller/model"
	"LanshanTeam-Examine/caller/pkg/consts"
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/userModule"
	"LanshanTeam-Examine/caller/rpc/userModule/pb"
	"strconv"

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
		Username:     user.Username,
		Password:     user.Password,
		IsGithubUser: false,
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
		"error":   nil,
		"token":   token,
	})

}
func LoginByPhoneNumber(c *gin.Context) {
	var params model.JustForLoginByPhoneNumber
	if err := c.ShouldBind(&params); err != nil {
		utils.ClientLogger.Error("ERROR: " + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.LackParams,
			"message": "you should enter your information completely,你应该完善你的信息",
			"error":   err.Error(),
		})
		return
	}
	if err := VerifyCode(params.Code, params.PhoneNumber); err != nil {
		utils.ClientLogger.Error("verifyCode in login by phone number ,code wrong")
		c.JSON(400, gin.H{
			"code":    consts.CodeWrong,
			"message": "code wrong ,验证码错误",
			"error":   err.Error(),
		})
		return
	}
	utils.ClientLogger.Info("verify success")

	num, _ := strconv.Atoi(params.PhoneNumber)
	resp, err := userModule.UserClient.Login(c, &pb.LoginReq{
		Username:    params.Username,
		PhoneNumber: int64(num),
	})
	if err != nil {
		utils.ClientLogger.Error("login by phone number failed")
		c.JSON(400, gin.H{
			"code":    consts.PhoneNumberUnavailable,
			"message": resp.GetMessage(),
			"error":   err.Error(),
		})
		return
	} else {
		utils.ClientLogger.Error("login by phone number success")

		token, err := middleware.GetToken(&model.Userinfo{
			Username: params.Username,
		})
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
			"message": resp.GetMessage(),
			"error":   nil,
			"token":   token,
		})

	}

}
