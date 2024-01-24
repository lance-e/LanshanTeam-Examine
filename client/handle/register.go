package handle

import (
	"LanshanTeam-Examine/client/model"
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/userModule"
	"LanshanTeam-Examine/client/rpc/userModule/pb"
	"github.com/gin-gonic/gin"
	"strconv"
)

func RegisterByPassword(c *gin.Context) {
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

	resp, err := userModule.UserClient.Register(c, &pb.RegisterReq{
		Username: user.Username,
		Password: user.Password,
	})

	if err != nil {
		utils.ClientLogger.Error("ERROR:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.RegisterFailed,
			"message": "happen error when create user,创建用户时出错",
			"error":   "服务器端出错" + err.Error(),
		})
		return
	}
	if !resp.GetFlag() && err == nil {
		utils.ClientLogger.Info(resp.GetMessage())
		c.JSON(400, gin.H{
			"code":    consts.UserAlreadyExist,
			"message": "user already exist, 用户已存在",
			"error":   resp.GetMessage(),
		})
		return
	}
	utils.ClientLogger.Info(resp.GetMessage())
	c.JSON(200, gin.H{
		"code":    consts.RegisterSuccess,
		"message": "user create success,用户创建成功",
		"error":   "",
	})
}
func RegisterByPhoneNumber(c *gin.Context) {
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

	if err := VerifyCode(user.Code, user.PhoneNumber); err != nil {
		utils.ClientLogger.Error("verifyCode in register by phone number ,code wrong")
		c.JSON(400, gin.H{
			"code":    consts.CodeWrong,
			"message": "code wrong ,验证码错误",
			"error":   err.Error(),
		})
		return
	}
	utils.ClientLogger.Info("verify success")
	number, _ := strconv.Atoi(user.PhoneNumber)
	resp, err := userModule.UserClient.Register(c, &pb.RegisterReq{
		Username:    user.Username,
		Password:    user.Password,
		PhoneNumber: int64(number),
	})
	if err != nil {
		utils.ClientLogger.Error("ERROR:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.RegisterFailed,
			"message": "happen error when create user,创建用户时出错",
			"error":   "服务器端出错" + err.Error(),
		})
		return
	}
	if !resp.GetFlag() && err == nil {
		utils.ClientLogger.Info(resp.GetMessage())
		c.JSON(400, gin.H{
			"code":    consts.UserAlreadyExist,
			"message": "user already exist, 用户已存在",
			"error":   resp.GetMessage(),
		})
		return
	}
	utils.ClientLogger.Info(resp.GetMessage())
	c.JSON(200, gin.H{
		"code":    consts.RegisterSuccess,
		"message": "user create success,用户创建成功",
		"error":   "",
	})

}
