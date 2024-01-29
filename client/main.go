package main

import (
	"LanshanTeam-Examine/client/api/router"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/userModule"
	"LanshanTeam-Examine/client/ws"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	userModule.UserRPC()
}
func main() {
	//启动所有的房间协程，接收两个用户之间的消息
	for _, v := range ws.AllRoom.Rooms {
		go v.Start()
	}

	engine := router.NewRouter()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.ClientLogger.Panic("listen error: " + err.Error())
		}
	}()
	//fmt.Println("======", userModule.UserConn.GetState())
	//userModule.UserRPC() //仅使用一次连接
	//fmt.Println("======", userModule.UserConn.GetState())
	//go func() {
	//	for {
	//		fmt.Println("==============RPC=============")
	//
	//		userModule.UserRPC()
	//
	//		state := userModule.UserConn.GetState()
	//
	//	}
	//
	//}()

	defer userModule.UserConn.Close()
	//优雅地关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		utils.ClientLogger.Info("server shutdown...")
	}
	utils.ClientLogger.Info("server closed...")
}
