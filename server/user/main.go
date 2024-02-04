package main

import (
	"LanshanTeam-Examine/server/user/dao/Init"
	"LanshanTeam-Examine/server/user/handle"
	"LanshanTeam-Examine/server/user/pb"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"LanshanTeam-Examine/server/user/serveRegister"
	"google.golang.org/grpc"
	"net"
)

func init() {
	Init.InitMysql()
	Init.InitRedis()
}
func main() {

	listener, err := net.Listen("tcp", "localhost:10001")
	if err != nil {
		utils.UserLogger.Panic("ERROR:" + err.Error())
	}
	server := grpc.NewServer()
	//注册服务
	pb.RegisterUserServer(server, &handle.UserServer{})
	register, err := serveRegister.NewEtcdRegister("localhost:2379")
	if err != nil {
		utils.UserLogger.Debug("can't init a new service register")
		return
	}
	err = register.ServiceRegister("user", "localhost:10001", 30)
	if err != nil {
		utils.UserLogger.Debug("service register failed")
		return
	}
	defer register.Close()
	//启动服务
	if err := server.Serve(listener); err != nil {
		utils.UserLogger.Panic("ERROR" + err.Error())
	}

}
