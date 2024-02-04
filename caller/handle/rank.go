package handle

import (
	"LanshanTeam-Examine/caller/pkg/consts"
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/userModule"
	"LanshanTeam-Examine/caller/rpc/userModule/pb"
	"github.com/gin-gonic/gin"
)

func ShowRank(c *gin.Context) {
	resp, err := userModule.UserClient.Rank(c, &pb.RankReq{})
	if err != nil {
		utils.ClientLogger.Debug("get rank failed")
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "get rank failed",
			"error":   err.Error(),
		})
		return
	}
	rank := resp.GetRank()
	c.JSON(200, gin.H{
		"code":  consts.GetRankSuccess,
		"rank":  rank,
		"error": nil,
	})
}
