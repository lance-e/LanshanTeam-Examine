package handle

import (
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/userModule"
	"LanshanTeam-Examine/client/rpc/userModule/pb"
	"github.com/gin-gonic/gin"
)

type userinfo struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	PhoneNumber int    `form:"phone_number" json:"phone_number"`
	Email       string `form:"phone_number" json:"email"`
}

func RegisterByPassword(c *gin.Context) {
	var user userinfo

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

func Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": userModule.UserConn.GetState().String(),
		"code":   200,
	})
}
