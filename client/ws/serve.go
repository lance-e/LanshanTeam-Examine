package ws

import (
	"LanshanTeam-Examine/client/pkg/utils"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

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
	MessageChannel   chan *Message   `json:"message_channel,omitempty"`
}
type GameRoom struct {
	User1            *UserConn       `json:"user_1,omitempty"`
	User2            *UserConn       `json:"user_2,omitempty"`
	TurnUser         *UserConn       `json:"turn_user,omitempty"`
	ChessBoard       [10][10]int64   `json:"chess_board,omitempty"`
	GameLogicChannel chan *GameLogic `json:"game_logic_channel,omitempty"`
	MessageChannel   chan *Message   `json:"message_channel,omitempty"`
}
type allRoom struct {
	Rooms    map[string]*GameRoom
	RoomName chan string
}

// 单例，保存所有创建的房间
var AllRoom = allRoom{
	Rooms: make(map[string]*GameRoom),
}

type allUserConn struct {
	Users map[string]*UserConn
}

var AllUserConn = allUserConn{
	Users: make(map[string]*UserConn),
}

// 客户端应该传其中一个的json
type AllInfo struct {
	*Message
	*GameLogic
}

type Message struct {
	Sender  string `json:"sender,omitempty"`
	Content string `json:"content,omitempty"`
}

type GameLogic struct {
	Player string `json:"player,omitempty"`
	Row    int64  `json:"row,omitempty"`
	Column int64  `json:"column,omitempty"`
}

// 下棋和消息的请求
func (u *UserConn) GameReq(g *GameRoom) error {
	var all AllInfo
	var err error
	for {
		err = u.Conn.ReadJSON(&all)

		if err != nil {
			utils.ClientLogger.Error("unmarshal failed")
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if all.Message != nil {
			g.MessageChannel <- all.Message
		}
		if all.GameLogic != nil {
			g.GameLogicChannel <- all.GameLogic
		}

	}
}

// 下棋的响应
func (u *UserConn) GameLogicResp() {
	for g := range u.GameLogicChannel {
		_ = u.Conn.WriteJSON(g)
	}
}

// 消息的响应
func (u *UserConn) MessageResp() {
	for m := range u.MessageChannel {
		_ = u.Conn.WriteJSON(m)
	}
}

// 新建一个用户连接
func NewUserConn(name string, conn *websocket.Conn) *UserConn {
	//mutex.Lock()
	//defer mutex.Unlock()
	u := &UserConn{
		Username:         name,
		Conn:             conn,
		GameLogicChannel: make(chan *GameLogic),
		MessageChannel:   make(chan *Message),
	}
	AllUserConn.Users[name] = u
	log.Println("New Userconn : ", u)
	return u
}

func (u *UserConn) Close() {
	mutex.Lock()
	defer mutex.Unlock()
	delete(AllRoom.Rooms, u.Username)
}

// 创建房间
func (u *UserConn) NewRoom() *GameRoom {
	//mutex.Lock()
	//defer mutex.Unlock()
	r := &GameRoom{
		User1:            u,
		User2:            NewUserConn("", &websocket.Conn{}),
		TurnUser:         u,
		GameLogicChannel: make(chan *GameLogic),
		MessageChannel:   make(chan *Message),
	}
	AllRoom.Rooms[u.Username] = r
	log.Println("New room :", r)
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
			log.Println("==========> logic !!!!")
			//不是该玩家的回合
			if logic.Player != g.TurnUser.Username {
				AllUserConn.Users[logic.Player].MessageChannel <- &Message{
					Sender:  "room",
					Content: "not your round",
				}
			}
			//是该玩家的回合
			if logic.Player == g.User1.Username {
				g.ChessBoard[logic.Row][logic.Column] = 1
				g.IsWin()
				g.TurnUser = g.User2

			} else if logic.Player == g.User2.Username {
				g.ChessBoard[logic.Row][logic.Column] = 2
				g.IsWin()
				g.TurnUser = g.User1
			}

		case msg := <-g.MessageChannel:
			log.Println("==============>message !!!")
			//聊天信息
			utils.ClientLogger.Debug("the message is :" + msg.Content)
			if msg.Sender == g.User1.Username {
				g.User2.MessageChannel <- msg
			}
			if msg.Sender == g.User2.Username {
				g.User1.MessageChannel <- msg
			}
		}
	}
}
func (g *GameRoom) Close() {
	mutex.Lock()
	defer mutex.Unlock()
	delete(AllRoom.Rooms, g.User1.Username)
}
func (g *GameRoom) IsWin() bool {
	//游戏输赢判断
	return false
}
