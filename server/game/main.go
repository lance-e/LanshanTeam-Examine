package main

import (
	"LanshanTeam-Examine/server/game/dao/Init"
	"LanshanTeam-Examine/server/game/handle"
	"LanshanTeam-Examine/server/game/pb"
	"LanshanTeam-Examine/server/game/serveRegister"
	"LanshanTeam-Examine/server/game/utils"
	"google.golang.org/grpc"
	"net"
)

func init() {
	Init.InitMysql()
	Init.InitRedis()
}

func main() {
	listener, err := net.Listen("tcp", "localhost:10002")
	if err != nil {
		utils.GameLogger.Panic("ERROR:" + err.Error())
	}
	server := grpc.NewServer()
	//注册服务
	pb.RegisterGameServer(server, &handle.GameServer{})
	register, err := serveRegister.NewEtcdRegister("localhost:2379")
	if err != nil {
		utils.GameLogger.Debug("can't init a new service register")
		return
	}
	err = register.ServiceRegister("game", "localhost:10002", 30)
	if err != nil {
		utils.GameLogger.Debug("service register failed")
		return
	}
	defer register.Close()
	//启动服务
	if err := server.Serve(listener); err != nil {
		utils.GameLogger.Panic("ERROR" + err.Error())
	}

}
