package handle

import (
	"LanshanTeam-Examine/caller/pkg/consts"
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/userModule"
	"LanshanTeam-Examine/caller/rpc/userModule/pb"
	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) {
	username, ok := c.Get("username")
	utils.ClientLogger.Debug("username is : " + username.(string))
	if !ok {
		c.JSON(400, gin.H{
			"code":    consts.TokenInvalid,
			"message": "token invalid",
			"error":   "token invalid",
		})
		return
	}
	resp, err := userModule.UserClient.HomePage(c, &pb.HomePageReq{
		Username: username.(string),
	})
	if err != nil {
		c.JSON(400, gin.H{
			"code":    consts.GetUserAllInformationFailed,
			"message": "serve unavailable",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":         consts.GetUserAllInformationSuccess,
		"username":     resp.GetUsername(),
		"phone_number": resp.GetPhoneNumber(),
		"email":        resp.GetEmail(),
		"score":        resp.GetScore(),
		"error":        nil,
	})
}
