package ws

import (
	"LanshanTeam-Examine/client/pkg/utils"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"sync"
	"time"
)

//var AllRoom = make(map[string]*GameRoom)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Minute,
}
var mutex sync.Mutex

type UserConn struct {
	Username         string          `json:"username,omitempty"`
	Conn             *websocket.Conn `json:"conn,omitempty"`
	IsReadyToPlay    bool            `json:"is_ready_to_play,omitempty"`
	GameLogicChannel chan *GameLogic `json:"game_logic_channel,omitempty"`
}
type GameRoom struct {
	User1            *UserConn       `json:"user_1,omitempty"`
	User2            *UserConn       `json:"user_2,omitempty"`
	TurnUser         *UserConn       `json:"turn_user,omitempty"`
	ChessBoard       [15][15]int64   `json:"chess_board,omitempty"`
	GameLogicChannel chan *GameLogic `json:"game_logic_channel,omitempty"`
}
type allRoom struct {
	Rooms    map[string]*GameRoom
	RoomName chan string
}

// 单例，保存所有创建的房间
var AllRoom = allRoom{
	Rooms: make(map[string]*GameRoom),
	//RoomName: make(chan string),
}

type allUserConn struct {
	Users map[string]*UserConn
}

var AllUserConn = allUserConn{
	Users: make(map[string]*UserConn),
}

//type Message struct {
//	Sender  *UserConn
//	Content string
//}

type GameLogic struct {
	Sender  *UserConn `json:"sender,omitempty"`
	Row     int64     `json:"row,omitempty"`
	Column  int64     `json:"column,omitempty"`
	Message string    `json:"message,omitempty"`
}

//func (u *UserConn) SendMessage() {
//
//}

// 下棋的请求
func (u *UserConn) SendGameReq(g *GameRoom) error {
	var logic *GameLogic //
	var err error
	for {
		err = u.Conn.ReadJSON(logic)

		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		g.GameLogicChannel <- logic
	}
}

// 下棋的响应
func (u *UserConn) SendGameResp() {
	for g := range u.GameLogicChannel {
		_ = u.Conn.WriteJSON(g)
	}
}

// 新建一个用户连接
func NewUserConn(name string, conn *websocket.Conn) *UserConn {
	mutex.Lock()
	defer mutex.Unlock()
	u := &UserConn{
		Username:         name,
		Conn:             conn,
		GameLogicChannel: make(chan *GameLogic),
	}
	AllUserConn.Users[name] = u
	return u
}

// 创建房间
func (u *UserConn) NewRoom() *GameRoom {
	mutex.Lock()
	defer mutex.Unlock()
	r := &GameRoom{
		User1:            u,
		User2:            NewUserConn("", &websocket.Conn{}),
		TurnUser:         u,
		GameLogicChannel: make(chan *GameLogic),
	}
	AllRoom.Rooms[u.Username] = r
	return r
}

// 进入房间
func (u *UserConn) JoinRoom(g *GameRoom) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	g.User2 = u
	defer func() {
		if r := recover(); r != nil {
			utils.ClientLogger.Debug("JoinRoom happen panic,and already recovered")
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()
	return nil
}
func (g *GameRoom) Start() {
	for {
		select {

		case logic := <-g.GameLogicChannel:

			if logic.Sender == g.User1 {
				g.User2.GameLogicChannel <- logic
				defer func() {
					if r := recover(); r != nil {
						utils.ClientLogger.Debug("can't send logic to user2")
						return
					}
				}()
			}
			if logic.Sender == g.User2 {
				g.User1.GameLogicChannel <- logic
				defer func() {
					if r := recover(); r != nil {
						utils.ClientLogger.Debug("can't send logic to user2")
						return
					}
				}()
			}
			utils.ClientLogger.Debug("the sender is not this room")
		}
	}
}
func (g *GameRoom) Close() {
	mutex.Lock()
	defer mutex.Unlock()

}
