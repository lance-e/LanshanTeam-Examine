package main

import (
	"LanshanTeam-Examine/caller/api/router"
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/gameModule"
	"LanshanTeam-Examine/caller/rpc/userModule"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	go userModule.UserRPC()
	go gameModule.GameRPC()
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

	//优雅地关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer userModule.UserConn.Close()
	defer gameModule.GameConn.Close()
	if err := server.Shutdown(ctx); err != nil {
		utils.ClientLogger.Info("server shutdown...")
	}
	utils.ClientLogger.Info("server closed...")
}
