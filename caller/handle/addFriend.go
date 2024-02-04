package handle

import (
	"LanshanTeam-Examine/caller/pkg/consts"
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/userModule"
	"LanshanTeam-Examine/caller/rpc/userModule/pb"
	"github.com/gin-gonic/gin"
)

func SendFriendRequest(c *gin.Context) {
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
	receiver := c.PostForm("receiver")
	_, err := userModule.UserClient.AddFriend(c, &pb.AddFriendReq{
		Sender:    username.(string),
		Receiver:  receiver,
		IsRequest: true,
	})
	if err != nil {
		if err.Error() == "you are friends" {
			c.JSON(400, gin.H{
				"code":    consts.ServeUnavailable,
				"message": "you are friends",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "serve unavailable",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    consts.AddFriendRequestSuccess,
		"message": "add friend request send success",
		"error":   nil,
	})

}
func ReceiveFriendRequest(c *gin.Context) {
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
	sender := c.PostForm("sender")
	_, err := userModule.UserClient.AddFriend(c, &pb.AddFriendReq{
		Sender:    sender,
		Receiver:  username.(string),
		IsRequest: false,
	})
	if err != nil {
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "serve unavailable",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    consts.CreateFriendSuccess,
		"message": "create friend success",
		"error":   nil,
	})
}
