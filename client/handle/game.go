package handle

import (
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/ws"
	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
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
	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()
	if err != nil {
		utils.ClientLogger.Error("ws error:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.StartGameFailed,
			"message": "start game failed",
			"error":   err.Error(),
		})
		return
	}
	//新建一个用户连接
	user1 := ws.NewUserConn(username.(string), conn)
	go user1.SendGameResp() //并发获取游戏响应

	//创建房间
	room := user1.NewRoom()

	go room.Start() //启动房间进程，接收两个用户之间的消息
	select {
	case info := <-user1.GameLogicChannel:
		user1.Conn.WriteJSON(info)
	}

}
func Join(c *gin.Context) {
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
	user1 := c.PostForm("user1")
	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()
	if err != nil {
		utils.ClientLogger.Error("ws error:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.StartGameFailed,
			"message": "start game failed",
			"error":   err.Error(),
		})
		return
	}
	//新建一个用户连接
	user2 := ws.NewUserConn(username.(string), conn)
	go user2.SendGameResp() //并发获取游戏响应

	//加入目标房间
	targetRoom, ok := ws.AllRoom.Rooms[user1]
	if !ok {
		utils.ClientLogger.Debug("this room not exists")
		c.JSON(400, gin.H{
			"code":    consts.NotFoundTargetGameRoom,
			"message": "not found target game room",
			"error":   "room not found",
		})
		return
	}
	err = user2.JoinRoom(targetRoom)
	if err != nil {
		utils.ClientLogger.Debug("can't join the room")
	}
	//向房间发送进入房间信息
	targetRoom.GameLogicChannel <- &ws.GameLogic{
		Sender:  user2,
		Message: user2.Username + " is join this room",
	}

}
func Ready(c *gin.Context) {
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
	targetUser, ok := ws.AllUserConn.Users[username.(string)]
	if !ok {
		utils.ClientLogger.Error("ready function not found in AllUserConn")
		c.JSON(400, gin.H{
			"code":    consts.ReadyToPlayGameFailed,
			"message": "ready failed",
			"error":   "can't find in connection pool",
		})
		return
	}
	targetUser.IsReadyToPlay = true
	utils.ClientLogger.Debug("ready success")
	c.JSON(400, gin.H{
		"code":    consts.ReadyToPlayGameSuccess,
		"message": "ready success",
		"error":   nil,
	})
}
