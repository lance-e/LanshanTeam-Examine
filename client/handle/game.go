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
	go user1.GameLogicResp() //并发获取游戏响应
	go user1.MessageResp()
	//创建房间
	room := user1.NewRoom()
	//for loop ，发送游戏请求
	err = user1.GameReq(room)

	if err == nil {
		utils.ClientLogger.Info("Connection close ,game over....")
		user1.Conn.Close()
		room.Close()
		user1.Close()
		c.JSON(200, gin.H{
			"code":    consts.GameOver,
			"message": "game over ,游戏结束",
			"error":   nil,
		})
		return
	} else {
		utils.ClientLogger.Info("Connection close ,these are some error that can't connect....")
		user1.Conn.Close()
		room.Close()
		user1.Close()
		c.JSON(400, gin.H{
			"code":    consts.GameOver,
			"message": "connection error",
			"error":   err.Error(),
		})
		return
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
	go user2.GameLogicResp() //并发获取游戏响应
	go user2.MessageResp()
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
	targetRoom.MessageChannel <- &ws.Message{
		Sender:  user2,
		Content: user2.Username + " is join this room",
	}

	err = user2.GameReq(targetRoom)

	if err == nil {
		utils.ClientLogger.Info("Connection close ,game over....")
		user2.Conn.Close()
		user2.Close()
		c.JSON(200, gin.H{
			"code":    consts.GameOver,
			"message": "game over ,游戏结束",
			"error":   nil,
		})
		return
	} else {
		utils.ClientLogger.Info("Connection close ,these are some error that can't connect....")
		user2.Conn.Close()
		user2.Close()
		c.JSON(400, gin.H{
			"code":    consts.GameOver,
			"message": "connection error",
			"error":   err.Error(),
		})
		return
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
