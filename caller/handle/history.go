package handle

import (
	"LanshanTeam-Examine/caller/pkg/consts"
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/gameModule"
	"LanshanTeam-Examine/caller/rpc/gameModule/pb"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func ShowHistory(c *gin.Context) {
	host := c.PostForm("room_host")
	resp, err := gameModule.GameClient.ShowSteps(c, &pb.ShowStepsReq{RoomHost: host})
	if err != nil {
		utils.ClientLogger.Error("show steps failed,error:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "can't show the history",
			"error":   err.Error(),
		})
		return
	}
	var steps string
	for _, v := range resp.GetAllStep() {
		steps += strings.Join([]string{(*v).Player, ":(", strconv.Itoa(int((*v).Row)), " , ", strconv.Itoa(int((*v).Column)), ")<- "}, "")
	}
	c.JSON(200, gin.H{
		"code":  consts.ShowHistorySuccess,
		"Steps": steps,
		"error": nil,
	})
	return
}
