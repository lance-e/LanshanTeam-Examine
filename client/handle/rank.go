package handle

import (
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/userModule"
	"LanshanTeam-Examine/client/rpc/userModule/pb"
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
		"No1":   rank[0],
		"No2":   rank[1],
		"No3":   rank[2],
		"No4":   rank[3],
		"No5":   rank[4],
		"No6":   rank[5],
		"No7":   rank[6],
		"No8":   rank[7],
		"No9":   rank[8],
		"No10":  rank[9],
		"error": nil,
	})
}
