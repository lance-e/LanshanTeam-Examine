package gameModule

import (
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/gameModule/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GameClient pb.GameClient

	GameConn *grpc.ClientConn
)

func GameRPC() {

	var opts []grpc.DialOption
	var err error
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	GameConn, err = grpc.Dial("localhost:10002", opts...)

	if err != nil {
		utils.ClientLogger.Debug(err.Error())
		fmt.Println("TEST CONNECTION...........ERROR")
	}
	GameClient = pb.NewGameClient(GameConn)

}
