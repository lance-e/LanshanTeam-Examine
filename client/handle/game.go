package handle

import (
	"LanshanTeam-Examine/client/api/middleware"
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

func Create(c *gin.Context) {
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
	//====
	value := c.Request.Header.Get("Authorization")
	utils.ClientLogger.Debug("TOKEN : " + value)
	tokenstr := strings.SplitN(value, " ", 2)
	if tokenstr[0] != "Bearer" {
		utils.ClientLogger.Debug("JWT格式不正确")
		conn.WriteMessage(websocket.TextMessage, []byte("JWT格式不正确"))
		return
	}
	if tokenstr[1] == "" {
		utils.ClientLogger.Debug("JWT为空")
		conn.WriteMessage(websocket.TextMessage, []byte("JWT为空"))
		return
	}
	cliam, err := middleware.ParseJWT(tokenstr[1])
	if err != nil {
		utils.ClientLogger.Debug("解析失败")
		conn.WriteMessage(websocket.TextMessage, []byte("解析失败"))
		return
	} else if cliam.ExpiresAt < time.Now().Unix() {
		utils.ClientLogger.Debug("token 超时")
		conn.WriteMessage(websocket.TextMessage, []byte("token 超时"))
		return
	}
	utils.ClientLogger.Debug("NAME: " + cliam.Username + " coming")
	c.Set("username", cliam.Username)
	//====
	username, ok := c.Get("username")
	log.Println(username.(string))
	utils.ClientLogger.Debug("username is : " + username.(string))
	if !ok {
		utils.ClientLogger.Error("token invalid")
		conn.WriteMessage(websocket.TextMessage, []byte("token invalid"))
		return
	}
	//
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

		utils.ClientLogger.Info("game over ,游戏结束")
		conn.WriteMessage(websocket.TextMessage, []byte("game over ,游戏结束"))
		return
	} else {
		utils.ClientLogger.Info("Connection close ,these are some error that can't connect....")
		user1.Conn.Close()
		room.Close()
		user1.Close()
		utils.ClientLogger.Error("connection error")
		conn.WriteMessage(websocket.TextMessage, []byte("connection error"))
		return
	}
}
func Join(c *gin.Context) {

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

	value := c.Request.Header.Get("Authorization")
	utils.ClientLogger.Debug("TOKEN : " + value)
	tokenstr := strings.SplitN(value, " ", 2)
	if tokenstr[0] != "Bearer" {
		utils.ClientLogger.Debug("JWT格式不正确")
		conn.WriteMessage(websocket.TextMessage, []byte("JWT格式不正确"))
		return
	}
	if tokenstr[1] == "" {
		utils.ClientLogger.Debug("JWT为空")
		conn.WriteMessage(websocket.TextMessage, []byte("JWT为空"))
		return
	}
	cliam, err := middleware.ParseJWT(tokenstr[1])
	if err != nil {
		utils.ClientLogger.Debug("解析失败")
		conn.WriteMessage(websocket.TextMessage, []byte("解析失败"))
		return
	} else if cliam.ExpiresAt < time.Now().Unix() {
		utils.ClientLogger.Debug("token 超时")
		conn.WriteMessage(websocket.TextMessage, []byte("token 超时"))
		return
	}
	utils.ClientLogger.Debug("NAME: " + cliam.Username + " coming")
	c.Set("username", cliam.Username)
	//
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
	//
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
func tool(ctx *gin.Context) {
	value := ctx.Request.Header.Get("Authorization")
	utils.ClientLogger.Debug("TOKEN : " + value)
	tokenstr := strings.SplitN(value, " ", 2)
	if tokenstr[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "JWT 格式不正确",
		})
		ctx.Abort()
		return
	}
	if tokenstr[1] == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "JWT 为空",
		})
		ctx.Abort()
		return
	}
	cliam, err := middleware.ParseJWT(tokenstr[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "解析失败",
		})
		ctx.Abort()
		return
	} else if cliam.ExpiresAt < time.Now().Unix() {
		ctx.JSON(400, gin.H{
			"error": "token 超时",
		})
		ctx.Abort()
		return
	}

	ctx.Set("username", cliam.Username)
}
