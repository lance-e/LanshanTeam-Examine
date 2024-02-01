package ws

import (
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/gameModule"
	"LanshanTeam-Examine/client/rpc/gameModule/pb"
	"LanshanTeam-Examine/client/rpc/userModule"
	pb2 "LanshanTeam-Examine/client/rpc/userModule/pb"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 消息格式{"sender":"lance","content":"this is message","player":"longxu","row":12,"column":13}
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
	ChessBoard       [10][11]int64   `json:"chess_board,omitempty"`
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
			//判断用户是否准备
			if AllUserConn.Users[logic.Player].IsReadyToPlay {
				//不是该玩家的回合
				if logic.Player != g.TurnUser.Username {
					AllUserConn.Users[logic.Player].MessageChannel <- &Message{
						Sender:  "room",
						Content: "not your round",
					}
				} else {
					//是该玩家的回合
					if logic.Row >= 10 || logic.Row < 0 || logic.Column >= 10 || logic.Column < 0 {
						AllUserConn.Users[logic.Player].MessageChannel <- &Message{
							Sender:  "room",
							Content: "bad operation ",
						}
					} else if g.ChessBoard[logic.Row][logic.Column] != 0 {
						AllUserConn.Users[logic.Player].MessageChannel <- &Message{
							Sender:  "room",
							Content: "there was already set",
						}
					} else {
						if logic.Player == g.User1.Username {
							g.ChessBoard[logic.Row][logic.Column] = 1
							_, err := gameModule.GameClient.Save(context.Background(), &pb.SaveReq{
								RoomHost: g.User1.Username,
								Player:   logic.Player,
								Row:      logic.Row,
								Column:   logic.Column,
							})
							if err != nil {
								utils.ClientLogger.Debug("can't save the step,error:" + err.Error())
							}
							if g.IsWin(1) {
								utils.ClientLogger.Debug("user1 has win")
								g.User1.MessageChannel <- &Message{
									Sender:  "room",
									Content: logic.Player + " has win !!!",
								}
								g.User2.MessageChannel <- &Message{
									Sender:  "room",
									Content: logic.Player + " has win !!!",
								}
								g.ChessBoard = [10][11]int64{}
								addScore(logic.Player)

							} else {
								g.TurnUser = g.User2
							}
						} else if logic.Player == g.User2.Username {
							g.ChessBoard[logic.Row][logic.Column] = 2
							_, err := gameModule.GameClient.Save(context.Background(), &pb.SaveReq{
								RoomHost: g.User1.Username,
								Player:   logic.Player,
								Row:      logic.Row,
								Column:   logic.Column,
							})
							if err != nil {
								utils.ClientLogger.Debug("can't save the step,error:" + err.Error())
							}
							if g.IsWin(2) {
								utils.ClientLogger.Debug("user2 has win")
								g.User1.MessageChannel <- &Message{
									Sender:  "room",
									Content: logic.Player + " has win !!!",
								}
								g.User2.MessageChannel <- &Message{
									Sender:  "room",
									Content: logic.Player + "user2 has win !!!",
								}
								g.ChessBoard = [10][11]int64{}
								addScore(logic.Player)
							}
							g.TurnUser = g.User1
						}
					}
				}
			} else {
				AllUserConn.Users[logic.Player].MessageChannel <- &Message{
					Sender:  "room",
					Content: "you should ready for game first",
				}
			}
			//每次游戏逻辑发送过来之后，房间会发送目前棋局情况
			g.User1.MessageChannel <- &Message{
				Sender:  "room",
				Content: g.ShowBoard(),
			}
			g.User2.MessageChannel <- &Message{
				Sender:  "room",
				Content: g.ShowBoard(),
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
func (g *GameRoom) IsWin(value int64) bool {
	//游戏输赢判断
	//竖直
	for i := 0; i < 10; i++ {
		for j := 0; j < 6; j++ {
			if g.ChessBoard[i][j] == value {
				if g.ChessBoard[i][j] == g.ChessBoard[i][j+1] && g.ChessBoard[i][j+1] == g.ChessBoard[i][j+2] &&
					g.ChessBoard[i][j+2] == g.ChessBoard[i][j+3] && g.ChessBoard[i][j+3] == g.ChessBoard[i][j+4] {
					return true
				}
			}
		}
	}
	//水平
	for i := 0; i < 6; i++ {
		for j := 0; j < 10; j++ {
			if g.ChessBoard[i][j] == value {
				if g.ChessBoard[i][j] == g.ChessBoard[i+1][j] && g.ChessBoard[i+1][j] == g.ChessBoard[i+2][j] &&
					g.ChessBoard[i+2][j] == g.ChessBoard[i+3][j] && g.ChessBoard[i+3][j] == g.ChessBoard[i+4][j] {
					return true
				}
			}
		}
	}
	//反斜
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			if g.ChessBoard[i][j] == value {
				if g.ChessBoard[i][j] == g.ChessBoard[i+1][j+1] && g.ChessBoard[i+1][j+1] == g.ChessBoard[i+2][j+2] &&
					g.ChessBoard[i+2][j+2] == g.ChessBoard[i+3][j+3] && g.ChessBoard[i+3][j+3] == g.ChessBoard[i+4][j+4] {
					return true
				}
			}
		}

	}
	for i := 4; i < 10; i++ {
		for j := 0; j < 6; j++ {
			if g.ChessBoard[i][j] == value {
				if g.ChessBoard[i][j] == g.ChessBoard[i-1][j+1] && g.ChessBoard[i-1][j+1] == g.ChessBoard[i-2][j+2] &&
					g.ChessBoard[i-2][j+2] == g.ChessBoard[i-3][j+3] && g.ChessBoard[i-3][j+3] == g.ChessBoard[i-4][j+4] {
					return true
				}
			}
		}
	}
	return false
}
func (g *GameRoom) ShowBoard() string {
	var board string
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			board += strconv.Itoa(int(g.ChessBoard[i][j]))
		}
		board += "  \n"
	}
	return board
}
func addScore(username string) {
	if username[(len(username)-8):] != "(github)" {
		_, err := userModule.UserClient.AddScore(context.Background(), &pb2.AddScoreReq{
			Username:     username,
			IsGithubName: false,
		})
		if err != nil {
			utils.ClientLogger.Debug("addScore rpc request failed,error:" + err.Error())
		}
	} else {
		_, err := userModule.UserClient.AddScore(context.Background(), &pb2.AddScoreReq{
			Username:     username,
			IsGithubName: false,
		})
		if err != nil {
			utils.ClientLogger.Debug("addScore rpc request failed,error:" + err.Error())
		}
	}
}
